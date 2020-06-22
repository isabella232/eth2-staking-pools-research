from py_ecc.bls import G2ProofOfPossession as bls
import random
from crypto import Polynomial

KEY_SIZE_BITS = 256

class Participant:
    def __init__(self):
        self.current_polynomial = None

    def distribuite_shares(self,share_indexes):
        if self.current_polynomial == None:
            raise AssertionError('set a polynomial before distribuiting shares')
        return [self.current_polynomial.evaluate(i) for i in share_indexes]

    def generate_polynomial(self,secret,degree):
        self.current_polynomial = Polynomial(secret,degree)
        self.current_polynomial.generate_random()


        # d = self.current_polynomial.evaluate(0)

        # rnd = random.getrandbits(KEY_SIZE_BITS)
        # pk = bls.SkToPk(rnd)
        # message = b'\xab' * 32
        # sig = bls.Sign(rnd,message)
        # ver = bls.Verify(pk,message,sig)
        #
        # rnd2 = random.getrandbits(KEY_SIZE_BITS)
        # pk2 = bls.SkToPk(rnd2)
        # sig2 = bls.Sign(rnd2, message)
        # ver2 = bls.Verify(pk2, message, sig2)
        #
        #
        # agg = bls.Aggregate([sig,sig2])
        # agg_ver = bls.FastAggregateVerify([pk,pk2],message,agg)

        print("hello")