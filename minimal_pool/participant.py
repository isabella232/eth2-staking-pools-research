from crypto import Polynomial,reconstruct_secret,sign,pubkey_from_sk
import pool_node

KEY_SIZE_BITS = 256

class Participant:
    def __init__(self,id):
        self.current_polynomial = None
        self.id = id
        self.node = pool_node.PoolNode(self.id,self.new_msg)
        self.round_shares = []
        self.group_secret = -1

    def distribuite_shares(self,share_indexes):
        if self.current_polynomial == None:
            raise AssertionError('set a polynomial before distribuiting shares')
        return [self.current_polynomial.evaluate(i) for i in share_indexes]

    def reconstruct_group_secret(self):
        self.group_secret = reconstruct_secret(self.round_shares)

    def sign(self,message):
        return sign(self.group_secret,message)

    def pub_group_key(self):
        return pubkey_from_sk(self.group_secret)

    def generate_polynomial(self,secret,degree):
        self.current_polynomial = Polynomial(secret,degree)
        self.current_polynomial.generate_random()

    def new_msg(self,msg):
        if msg.type == "share_distro":
            self.round_shares.append(msg.data["share"])
