import crypto
import pool_node
import config
import logging
import threading
from  py_ecc.bls.g2_primatives import pubkey_to_G1

class Participant:
    def __init__(self,id, key):
        self.current_polynomial = None
        self.id = id
        self.node = pool_node.PoolNode(self.id,self.new_msg)
        self.incoming_shares = []
        self.collected_sigs = []
        self.key = key
        self.epoch_transition_lock = threading.Lock()

    def reconstruct_group_sk(self):
        shares = {}
        for m in self.incoming_shares:
            if m["epoch"] == self.node.state.epoch:
                shares[m["from_p_id"]] = m["share"]
        self.key = crypto.reconstruct_sk(shares)

    def reconstruct_group_pk(self):
        last_epoch = self.node.state.epoch

        pks = {}
        for s in self.collected_sigs:
            if s["epoch"] == last_epoch:
                pks[s["from_p_id"]] = pubkey_to_G1(s["pk"])
        return crypto.reconstruct_pk(pks)

    def sign_epoch_msg(self,msg):
        return self.sign(msg)

    def sign(self,message):
        return crypto.sign_with_sk(self.key, message)

    def end_epoch(self):
        with self.epoch_transition_lock:
            # remove last round's shares
            last_epoch = self.node.state.epoch
            last_epoch_shares = []
            next_epoch_shares = []
            for i in range(len(self.incoming_shares)):
                s = self.incoming_shares[i]
                if s["epoch"] == last_epoch:
                    last_epoch_shares.append(s)
                else:
                    next_epoch_shares.append(s)

            # save
            self.node.state.save_participant_shares(last_epoch_shares, last_epoch, self.id)
            self.incoming_shares = next_epoch_shares

            # reconstruct group pk
            group_pk = self.reconstruct_group_pk()

            # verify sigs and save them
            sigs = []
            pks = []
            for s in self.collected_sigs:
                if s["epoch"] == last_epoch:
                    sigs.append(s["sig"])
                    pks.append(s["pk"])
            agg_sigs = crypto.aggregate_sigs(sigs)
            agg_pks = crypto.aggregate_pks(pks)
            is_verified = crypto.verify_sig(agg_pks, config.TEST_EPOCH_MSG, agg_sigs)

            logging.debug("collected sig(p %d): %d",self.id, len(self.collected_sigs))
            logging.debug("group pk: %s", group_pk.hex())


            # self.node.state.save_aggregated_sig(agg_sigs, pks, is_verified, last_epoch)
            self.collected_sigs = []
            logging.debug("participant %d epoch end", self.id)

    def mid_epoch(self):
        # reconstruct my share and group's public key
        self.reconstruct_group_sk()

        # broadcast my sig
        sig = self.sign(config.TEST_EPOCH_MSG)
        pk = crypto.pk_from_sk(self.key)
        self.node.broadcast_sig(
            self.id,
            sig,
            pk,
            self.node.current_epoch_pool_assignment(self.id)
        )
        self.collected_sigs.append({  # TODO - find a better way to store own share
            "from_p_id": self.id,
            "sig": sig,
            "pk": pk,
            "pool_id": self.node.current_epoch_pool_assignment(self.id),
            "epoch": self.node.state.epoch,
        })
        logging.debug("participant %d epoch mid", self.id)

    def start_epoch(self):
        pool_asssignment = self.node.current_epoch_pool_assignment(self.id)
        pool_participants = self.node.current_epoch_pool_participants(pool_asssignment)
        redistro_polynomial = crypto.Redistribuition(config.POOL_THRESHOLD-1, self.key, pool_participants) # following Shamir's secret sharing, degree is threshold - 1
        shares_to_distrb, commitments = redistro_polynomial.generate_shares()
        self.node.broadcast_shares(
            self.id,
            shares_to_distrb,
            commitments,
            pool_asssignment
        )

        # add own share to self
        self.incoming_shares.append({  # TODO - find a better way to store own share
            "epoch": self.node.state.epoch,
            "share": shares_to_distrb[self.id],
            "commitments": commitments,
            "from_p_id": self.id,
            "p": self.id,
            "pool_id": pool_asssignment,
        })
        logging.debug("participant %d epoch start", self.id)

    def new_msg(self,msg):
        if msg.type == config.MSG_SHARE_DISTRO:
            with self.epoch_transition_lock:
                if self.node.current_epoch_pool_assignment(self.id) == msg.data["pool_id"] and msg.data["p"] == self.id:
                    if msg.data not in self.incoming_shares:
                        self.incoming_shares.append(msg.data)
        if msg.type == config.MSG_NEW_EPOCH:
            self.start_epoch()
        if msg.type == config.MSG_END_EPOCH:
            self.end_epoch()
        if msg.type == config.MSG_MID_EPOCH:
            self.mid_epoch()
        if msg.type == config.MSG_EPOCH_SIG:
            with self.epoch_transition_lock:
                if self.node.current_epoch_pool_assignment(self.id) == msg.data["pool_id"]:
                    if msg.data not in self.collected_sigs:
                        self.collected_sigs.append(msg.data)

