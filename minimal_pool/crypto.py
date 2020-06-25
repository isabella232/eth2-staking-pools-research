import random
from py_ecc.bls import G2ProofOfPossession as py_ecc_bls
from py_ecc.optimized_bls12_381 import curve_order
import milagro_bls_binding as milagro_bls
from hashlib import sha256
import time
import logging
import config

bls = milagro_bls
order = 10#curve_order

def generate_sk():
    return random.getrandbits(config.KEY_SIZE_BITS) % order

def reconstruct_secret(shares):
    return sum([i for i in shares]) % order

def pk_from_sk(sk):
    return bls.SkToPk(sk.to_bytes(32,config.ENDIANNESS))

def sign_with_sk(sk, msg):
    return bls.Sign(sk.to_bytes(32,config.ENDIANNESS),msg)

def aggregate_signatures(sigs):
    return bls.Aggregate(sigs)

def verify_aggregated_sigs(pks, message, sig):
    return bls.FastAggregateVerify(pks,message,sig)

def verify_sig(pk, message, sig):
    return bls.Verify(pk,message,sig)

def hash(x: bytes) -> bytes:
    return sha256(x).digest()


def test():
    logging.debug("started")
    priv = []
    pub = []
    sigs = []
    message = b'\xab' * 32
    for i in range(0,100):

        k1 = generate_sk()#.to_bytes(32,config.ENDIANNESS)
        k1_pub = pk_from_sk(k1)

        k1_sig = sign_with_sk(k1,message)

        priv.append(k1)
        pub.append(k1_pub)
        sigs.append(k1_sig)
        
    start = time.time()
    agg_sig = aggregate_signatures(sigs)
    end = time.time()
    logging.debug("aggregation time %f", (start - end))

    start = time.time()
    res = verify_aggregated_sigs(pub,message,agg_sig)
    end = time.time()
    logging.debug("FastAggregateVerify time %f with res %s", (start - end),res)

class Polynomial:
    def __init__(self,secret,degree):
        self.coefficients = []
        self.secret = secret
        self.degree = degree

    def generate_random(self):
        self.coefficients.append(self.secret) # important it's in index 0, see valuatee
        self.coefficients.extend([generate_sk() for _ in range(0, self.degree)])


    def evaluate(self,point):
        return sum([self.coefficients[i] * (point ** i) for i in range(len(self.coefficients))]) % order

