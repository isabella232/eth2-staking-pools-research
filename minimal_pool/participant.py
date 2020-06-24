from crypto import Polynomial,reconstruct_secret,sign,pubkey_from_sk
import pool_node
import config
import logging
import threading

class Participant:
    def __init__(self,id):
        self.current_polynomial = None
        self.id = id
        self.node = pool_node.PoolNode(self.id,self.new_msg)
        self.incoming_shares = []
        self.secret = -1
        self.epoch_transition_lock = threading.Lock()

    def reconstruct_group_secret(self):
        1+1
        # self.secret = reconstruct_secret(self.next_round_shares)

    def sign(self,message):
        return sign(self.secret, message)

    def pub_group_key(self):
        return pubkey_from_sk(self.secret)

    def end_epoch(self):
        with self.epoch_transition_lock:
            # isolate last round's shares
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

    def start_epoch(self):
        pool_asssignment = self.node.current_epoch_pool_assignment(self.id)
        pool_participants = self.node.current_epoch_pool_participants(pool_asssignment)
        redistro_polynomial = Polynomial(self.secret, config.POOL_THRESHOLD)
        redistro_polynomial.generate_random()
        shares_to_distrb = [[p_id, redistro_polynomial.evaluate(p_id)] for p_id in pool_participants]
        self.node.broadcast_shares(self.id, shares_to_distrb, pool_asssignment)

        # add own share to self
        self.incoming_shares.append({  # TODO - find a better way to store own share
            "epoch": self.node.state.epoch,
            "share": redistro_polynomial.evaluate(self.id),
            "from_p_id": self.id,
            "p": self.id,
            "pool_id": pool_asssignment,
        })


    def new_msg(self,msg):
        if msg.type == config.MSG_SHARE_DISTRO:
            with self.epoch_transition_lock:
                if self.node.current_epoch_pool_assignment(self.id) == msg.data["pool_id"] and msg.data["p"] == self.id:
                    if msg.data["share"] not in self.incoming_shares:
                        self.incoming_shares.append(msg.data)
        if msg.type == config.MSG_NEW_EPOCH:
            self.start_epoch()
            logging.debug("participant %d epoch start", self.id)
        if msg.type == config.MSG_END_EPOCH:
            self.end_epoch()
            logging.debug("participant %d epoch end",self.id)
        if msg.type == config.MSG_MID_EPOCH:
            self.reconstruct_group_secret()
            logging.debug("participant %d epoch mid", self.id)


