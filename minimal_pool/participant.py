from crypto import Polynomial,reconstruct_secret,sign,pubkey_from_sk
import pool_node
import config
import logging

class Participant:
    def __init__(self,id):
        self.current_polynomial = None
        self.id = id
        self.node = pool_node.PoolNode(self.id,self.new_msg)
        self.next_round_shares = []
        self.secret = -1

    def reconstruct_group_secret(self):
        if self.id == 1:
            logging.debug("mid round")
        # self.secret = reconstruct_secret(self.next_round_shares)

    def sign(self,message):
        return sign(self.secret, message)

    def pub_group_key(self):
        return pubkey_from_sk(self.secret)

    def reset_for_epoch(self):
        # save
        self.node.state.save_participant_shares(self.next_round_shares,self.node.state.epoch-1,self.id)

        # # prepare for next epoch
        # if self.node.state.epoch == 1: # important for initial share distro
        #     return
        self.next_round_shares = []

    def new_msg(self,msg):
        if msg.type == config.MSG_SHARE_DISTRO:
            if self.node.current_epoch_pool_assignment(self.id) == msg.data["pool_id"] and msg.data["p"] == self.id:
                if msg.data["share"] not in self.next_round_shares:
                    self.next_round_shares.append(msg.data)
                    logging.debug("epoch(%d), share from %d to %d (msg id %s)",msg.data["epoch"],msg.data["from_p_id"],self.id,msg.id)
        if msg.type == config.MSG_NEW_EPOCH:
            if self.id == 3:
                logging.debug("new round")

            pool_asssignment = self.node.current_epoch_pool_assignment(self.id)
            pool_participants = self.node.current_epoch_pool_participants(pool_asssignment)
            redistro_polynomial = Polynomial(self.secret, config.POOL_THRESHOLD)
            redistro_polynomial.generate_random()
            shares_to_distrb = [[p_id,redistro_polynomial.evaluate(p_id)] for p_id in pool_participants]
            self.node.broadcast_shares(self.id,shares_to_distrb,pool_asssignment)

            # add own share to self
            # self.next_round_shares.append(redistro_polynomial.evaluate(self.id))
        if msg.type == config.MSG_END_EPOCH:
            if self.id == 3:
                logging.debug("end round")
            self.reset_for_epoch()
        if msg.type == config.MSG_MID_EPOCH:
            self.reconstruct_group_secret()
            # if self.id == 1:
            #     logging.debug("shares for next epoch: %d",len(self.next_round_shares))


