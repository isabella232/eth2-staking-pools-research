import crypto
from node import node
import config
import logging
import threading

class Participant:
    def __init__(self, id, key):
        self.current_polynomial = None
        self.id = id
        self.node = node.PoolNode(self.id, self.handle_messages)
        self.incoming_shares = []
        self.collected_sigs = []
        self.key = key

        # locks
        self.epoch_lock = threading.Lock() # will lock at the beginning of each epoch, releases at end
        self.incoming_shares_lock = threading.Lock()
        self.collected_sigs_lock = threading.Lock()
        self.key_lock = threading.Lock()

    def reconstruct_group_sk(self, epoch):
        with self.incoming_shares_lock:
            shares = {}
            ids = []
            for m in self.incoming_shares:
                if m["epoch"].number == epoch.number:
                    shares[m["from_p_id"]] = m["share"]
                    ids.append(m["from_p_id"])
            with self.key_lock:
                logging.debug("p %d, epoch %d, reconstruct_sk from: %s", self.id, epoch.number, ids)
                self.key = crypto.reconstruct_sk(shares)

    # TODO - remove not needed as we save pool pk by id
    def reconstruct_group_pk(self, for_epoch):
        with self.collected_sigs_lock:
            pks = {}
            ids = []
            for s in self.collected_sigs:
                if s["epoch"].number == for_epoch:
                    pks[s["from_p_id"]] = s["pk"]
                    ids.append(s["from_p_id"])
            return ids, crypto.reconstruct_pk(pks)

    def reconstruct_group_sig(self, for_epoch):
        with self.collected_sigs_lock:
            sigs = {}
            ids = []
            for s in self.collected_sigs:
                if s["epoch"].number == for_epoch:
                    sigs[s["from_p_id"]] = s["sig"]
                    ids.append(s["from_p_id"])

            if len(ids) == 0:
                return [], None
            return ids, crypto.reconstruct_group_sig(sigs)

    def sign_epoch_msg(self, msg):
        return self.sign(msg)

    def sign(self,message):
        return crypto.sign_with_sk(self.key, message)

    def end_epoch(self, epoch):
        # remove last round's shares
        # TODO - move incoming_shares to mid_epoch?
        last_epoch = self.node.state.epoch
        my_pool_id = epoch.pool_id_for_participant(self.id)

        last_epoch_shares = []
        next_epoch_shares = []
        with self.incoming_shares_lock:
            for i in range(len(self.incoming_shares)):
                s = self.incoming_shares[i]
                if s["epoch"].number == last_epoch:
                    last_epoch_shares.append(s)
                else:
                    next_epoch_shares.append(s)
            self.incoming_shares = next_epoch_shares

        # save
        epoch.save_participant_shares(last_epoch_shares, self.id)
        self.node.state.save_epoch(epoch)

        # reconstruct group pk and sig
        group_pk = self.node.state.pool_info_by_id(my_pool_id)["pk"]
        ids, group_sig = self.reconstruct_group_sig(epoch.number)
        if group_sig is not None:
            group_sig = crypto.readable_sig(group_sig)

            ids, pk = self.reconstruct_group_pk(epoch.number)
            logging.debug("p %d pk from: %s VS %s, group pk: %s", self.id, ids, epoch.pool_participants_by_id(my_pool_id),
                      crypto.readable_pk(pk).hex())


            # verify sigs and save them
            is_verified = crypto.verify_sig(group_pk, config.TEST_EPOCH_MSG, group_sig)
            epoch.save_aggregated_sig(my_pool_id, group_sig, ids, is_verified)
            self.node.state.save_epoch(epoch)

        with self.collected_sigs_lock:
            self.collected_sigs = []

        logging.debug("participant %d epoch %d end", self.id, epoch.number)

    def mid_epoch(self, epoch):
        my_pool_id = epoch.pool_id_for_participant(self.id)

        # reconstruct my group share
        if epoch.number != 0:  # TODO - change from special case
            self.reconstruct_group_sk(epoch)

        # broadcast my sig
        sig = self.sign(config.TEST_EPOCH_MSG)
        pk = crypto.pk_from_sk(self.key)

        self.node.broadcast_sig(
            epoch,
            self.id,
            sig,
            pk,
            my_pool_id
        )
        with self.collected_sigs_lock:
            self.collected_sigs.append({  # TODO - find a better way to store own share
                "from_p_id": self.id,
                "sig": sig,
                "pk": pk,
                "pool_id": my_pool_id,
                "epoch": epoch,
            })
        logging.debug("participant %d epoch %d mid", self.id, epoch.number)

    """ 
        At the start of every epoch we distribute shares to the next, deterministically chosen, the shares for this
        participant's key
    """
    def start_epoch(self, epoch):
        next_epoch = self.node.state.get_epoch(epoch.number + 1)
        my_pool_id = epoch.pool_id_for_participant(self.id)
        pool_participants = next_epoch.pool_participants_by_id(
            epoch.share_distribuition_target(my_pool_id)
        )

        with self.key_lock:
            redistro_polynomial = crypto.Redistribuition(config.POOL_THRESHOLD - 1, self.key, pool_participants)  # following Shamir's secret sharing, degree is threshold - 1
            shares_to_distrb, commitments = redistro_polynomial.generate_shares()

        logging.debug("p %d, epoch %d, shares to: %s", self.id, epoch.number, pool_participants)

        self.node.broadcast_shares(
            next_epoch,
            self.id,
            shares_to_distrb,
            commitments,
            my_pool_id
        )

        # add own share to self
        if self.id in pool_participants:
            with self.incoming_shares_lock:
                self.incoming_shares.append({  # TODO - find a better way to store own share
                    "epoch": next_epoch,
                    "share": shares_to_distrb[self.id],
                    "commitments": commitments,
                    "from_p_id": self.id,
                    "p": self.id,
                    "pool_id": my_pool_id,
                })


        logging.debug("participant %d epoch %d start", self.id, epoch.number)

    def handle_messages(self, msg):
        if msg.type == config.MSG_SHARE_DISTRO:
            if msg.data["p"] == self.id:
                with self.incoming_shares_lock:
                    if msg.data not in self.incoming_shares:
                        self.incoming_shares.append(msg.data)
        if msg.type == config.MSG_NEW_EPOCH:
            self.epoch_lock.acquire()
            threading.Thread(target=self.start_epoch, args=[msg.data["epoch"]], daemon=True).start()
        if msg.type == config.MSG_END_EPOCH:
            def wait_and_release(t, participant):
                t.join()
                participant.epoch_lock.release()
            t = threading.Thread(target=self.end_epoch, args=[msg.data["epoch"]], daemon=True)
            t.start()
            threading.Thread(target=wait_and_release, args=[t, self], daemon=True).start()
        if msg.type == config.MSG_MID_EPOCH:
            threading.Thread(target=self.mid_epoch, args=[msg.data["epoch"]], daemon=True).start()
        if msg.type == config.MSG_EPOCH_SIG:
            e = msg.data["epoch"]
            my_pool_in_epoch = e.pool_id_for_participant(self.id)
            if my_pool_in_epoch == msg.data["pool_id"]:  # add to sigs only if i'm in the relevant pool
                with self.collected_sigs_lock:
                    if msg.data not in self.collected_sigs:
                        self.collected_sigs.append(msg.data)



