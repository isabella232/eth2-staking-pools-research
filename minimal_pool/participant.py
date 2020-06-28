import crypto
from node import node
import config
import logging
import threading

class Participant:
    def __init__(self, id, key):
        self.current_polynomial = None
        self.id = id
        self.node = node.PoolNode(self.id, self.new_msg)
        self.incoming_shares = []
        self.collected_sigs = []
        self.key = key

        # locks
        self.incoming_shares_lock = threading.Lock()
        self.collected_sigs_lock = threading.Lock()
        self.key_lock = threading.Lock()

    def reconstruct_group_sk(self, epoch):
        with self.incoming_shares_lock:
            shares = {}
            for m in self.incoming_shares:
                if m["epoch"].number == epoch.number:
                    shares[m["from_p_id"]] = m["share"]
            with self.key_lock:
                self.key = crypto.reconstruct_sk(shares)

    # TODO - remove not needed as we save pool pk by id
    # def reconstruct_group_pk(self):
    #     last_epoch = self.node.state.epoch
    #
    #     pks = {}
    #     for s in self.collected_sigs:
    #         if s["epoch"] == last_epoch:
    #             pks[s["from_p_id"]] = s["pk"]
    #     return crypto.reconstruct_pk(pks)

    def reconstruct_group_sig(self):
        with self.collected_sigs_lock:
            last_epoch = self.node.state.epoch
            sigs = {}
            for s in self.collected_sigs:
                if s["epoch"].number == last_epoch:
                    sigs[s["from_p_id"]] = s["sig"]
            return crypto.reconstruct_group_sig(sigs)

    def sign_epoch_msg(self,msg):
        return self.sign(msg)

    def sign(self,message):
        return crypto.sign_with_sk(self.key, message)

    def end_epoch(self, epoch):
        # remove last round's shares
        last_epoch = self.node.state.epoch
        my_pool_id = epoch.pool_id_for_participant(self.id)

        last_epoch_shares = []
        next_epoch_shares = []
        with self.incoming_shares_lock:
            for i in range(len(self.incoming_shares)):
                s = self.incoming_shares[i]
                if s["epoch"] == last_epoch:
                    last_epoch_shares.append(s)
                else:
                    next_epoch_shares.append(s)
            self.incoming_shares = next_epoch_shares

        # save
        epoch.save_participant_shares(last_epoch_shares, self.id)
        self.node.state.save_epoch(epoch)

        # reconstruct group pk and sig
        # group_pk = self.node.state.pool_info_by_id(my_pool_id)["pk"]
        # group_sig = self.reconstruct_group_sig()

        # verify sigs and save them
        # is_verified = crypto.verify_sig(group_pk, config.TEST_EPOCH_MSG, group_sig)
        # self.node.state.save_epoch_sig(group_sig, group_pk, is_verified, last_epoch)
        with self.collected_sigs_lock:
            self.collected_sigs = []

        logging.debug("participant %d epoch %d end", self.id, epoch.number)


    def mid_epoch(self, epoch):
        # reconstruct my share and group's public key
        self.reconstruct_group_sk(epoch)

        # broadcast my sig
        sig = self.sign(config.TEST_EPOCH_MSG)
        pk = crypto.pk_from_sk(self.key)
        self.node.broadcast_sig(
            epoch,
            self.id,
            sig,
            pk,
            epoch.pool_id_for_participant(self.id)
        )
        with self.collected_sigs_lock:
            self.collected_sigs.append({  # TODO - find a better way to store own share
                "from_p_id": self.id,
                "sig": sig,
                "pk": pk,
                "pool_id": epoch.pool_id_for_participant(self.id),
                "epoch": epoch,
            })
        logging.debug("participant %d epoch %d mid", self.id, epoch.number)

    def start_epoch(self, epoch):
        pool_asssignment = epoch.pool_id_for_participant(self.id)
        pool_participants = epoch.pool_participants_by_id(pool_asssignment)
        with self.key_lock:
            redistro_polynomial = crypto.Redistribuition(config.POOL_THRESHOLD-1, self.key, pool_participants) # following Shamir's secret sharing, degree is threshold - 1
            shares_to_distrb, commitments = redistro_polynomial.generate_shares()
        self.node.broadcast_shares(
            epoch,
            self.id,
            shares_to_distrb,
            commitments,
            pool_asssignment
        )

        # add own share to self
        with self.incoming_shares_lock:
            self.incoming_shares.append({  # TODO - find a better way to store own share
                "epoch": epoch,
                "share": shares_to_distrb[self.id],
                "commitments": commitments,
                "from_p_id": self.id,
                "p": self.id,
                "pool_id": pool_asssignment,
            })
        logging.debug("participant %d epoch %d start", self.id, self.node.state.epoch)

    def new_msg(self,msg):
        if msg.type == config.MSG_SHARE_DISTRO:
            e = msg.data["epoch"]
            if e.pool_id_for_participant(self.id) == msg.data["pool_id"] and msg.data["p"] == self.id:
                with self.incoming_shares_lock:
                    if msg.data not in self.incoming_shares:
                        self.incoming_shares.append(msg.data)
        if msg.type == config.MSG_NEW_EPOCH:
            threading.Thread(target=self.start_epoch, args=[msg.data["epoch"]], daemon=True).start()
        if msg.type == config.MSG_END_EPOCH:
            threading.Thread(target=self.end_epoch, args=[msg.data["epoch"]], daemon=True).start()
        if msg.type == config.MSG_MID_EPOCH:
            threading.Thread(target=self.mid_epoch, args=[msg.data["epoch"]], daemon=True).start()
        if msg.type == config.MSG_EPOCH_SIG:
            e = msg.data["epoch"]
            if e.pool_id_for_participant(self.id) == msg.data["pool_id"]:
                with self.collected_sigs_lock:
                    if msg.data not in self.collected_sigs:
                        self.collected_sigs.append(msg.data)


