import random
import py_ecc
from py_ecc.bls import G2ProofOfPossession as bls

KEY_SIZE_BITS = 256


def generate_sk():
    return random.getrandbits(KEY_SIZE_BITS)

def reconstruct_secret(shares):
    return sum([i for i in shares]) % py_ecc.optimized_bls12_381.curve_order

def pubkey_from_sk(sk):
    return bls.SkToPk(sk)

def sign(sk,msg):
    return bls.Sign(sk,msg)

def aggregate_signatures(sigs):
    return bls.Aggregate(sigs)

def verify_aggregated_sigs(pks, message, sig):
    return bls.FastAggregateVerify(pks,message,sig)

class Polynomial:
    def __init__(self,secret,degree):
        self.coefficients = []
        self.secret = secret
        self.degree = degree

    def generate_random(self):
        self.coefficients.append(self.secret)
        self.coefficients.extend([generate_sk() for _ in range(1, self.degree)])

    def evaluate(self,point):
        return sum([self.coefficients[i] * (point ** i) for i in range(len(self.coefficients))])

