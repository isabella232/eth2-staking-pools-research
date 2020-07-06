import random
from py_ecc.optimized_bls12_381 import curve_order,add as ec_add, multiply as ec_mul,G1
from py_ecc.bls.g2_primatives import G1_to_pubkey, signature_to_G2, G2_to_signature
import milagro_bls_binding as milagro_bls
from hashlib import sha256
import config
import math

bls = milagro_bls
order = curve_order

def generate_sk():
    return random.getrandbits(config.KEY_SIZE_BITS) % order

"""
    given a set of shares, this will reconstruct a participant's group secret 
"""
def reconstruct_sk(shares):
    l = LagrangeInterpolation(shares, order)
    return l.evaluate()

"""
    given a set of shares, this will reconstruct the group's public key 
"""
def reconstruct_pk(shares):
    l = ECLagrangeInterpolation(shares, order)
    ev = l.evaluate()
    return ev

def pk_from_sk(sk):
    return ec_mul(G1, sk)

def readable_pk(optimized_pk):
    return G1_to_pubkey(optimized_pk)

def sign_with_sk(sk, msg):
    return signature_to_G2(bls.Sign(sk.to_bytes(32, config.ENDIANNESS), msg))

def readable_sig(optimized_sig):
    return G2_to_signature(optimized_sig)

def reconstruct_group_sig(shares):
    l = ECLagrangeInterpolation(shares, order)
    ev = l.evaluate()

    return ev

def aggregate_sigs(sigs):
    return bls.Aggregate(sigs)

def aggregate_pks(pks):
    return bls._AggregatePKs(pks)

def verify_aggregated_sigs(pks, message, sig):
    return bls.FastAggregateVerify(pks,message,sig)

def verify_sig(pk, message, sig):
    return bls.Verify(pk, message,sig)

def hash(x: bytes) -> bytes:
    return sha256(x).digest()

class Polynomial:
    def __init__(self,secret, degree, mod):
        self.coefficients = []
        self.secret = secret
        self.degree = degree
        self.mod = mod

    def generate_random(self):
        self.coefficients.append(self.secret)  # important it's in index 0, see valuate
        self.coefficients.extend([generate_sk() for _ in range(0, self.degree)])

    def evaluate(self, point):
        return sum([self.coefficients[i] * (point ** i) for i in range(len(self.coefficients))]) % self.mod

    def coefficients_commitment(self):
        return [pk_from_sk(self.coefficients[i]) for i in range(len(self.coefficients))]

class LagrangeInterpolation:
    def __init__(self, shares, mod):
        self.shares = shares
        self.mod = mod

    def evaluate(self):
        return self.sum_func([self.mul_func(self.shares[i], self.lagrange_coefficients(i)) for i in self.shares]) % self.mod

    def lagrange_coefficients(self, i, x_0=0):
        num = 1
        den = 1
        for j in self.shares:
            if i != j:
                num = (num * (x_0 - j))
                den = (den * (i - j))
        return ((num % self.mod) * self.mod_inverse(den, self.mod)) % self.mod

    def mod_inverse(self, b, m):
        g = math.gcd(b, m)
        if (g != 1):
            raise AssertionError("Inverse doesn't exist")
        else:
            # If b and m are relatively prime,
            # then modulo inverse is b^(m-2) mode m
            return pow(b, m - 2, m)

    def sum_func(self, lst):
        if len(lst) == 0:
            return 0

        ret = lst[0]
        for i in range(1, len(lst)):
            ret = self.add_func(ret, lst[i])
        return ret

    def add_func(self, a, b):
        return a+b

    def mul_func(self, a, b):
        return a*b

class ECLagrangeInterpolation(LagrangeInterpolation):
    def evaluate(self):
        mull = []
        for i in self.shares:
            m = self.mul_func(self.shares[i], self.lagrange_coefficients(i))
            mull.append(m)

        return self.sum_func(mull)

    def add_func(self, a, b):
        return ec_add(a, b)

    def mul_func(self, a, b):
        return ec_mul(a, b)


class Redistribuition:
    def __init__(self, threshold, sk, participant_indexes):
        self.threshold = threshold
        self.indexes = participant_indexes
        self.sk = sk

    def generate_shares(self):
        poly = Polynomial(self.sk, self.threshold, order)
        poly.generate_random()
        commitment = poly.coefficients_commitment()

        shares = {}
        for p_idx in self.indexes:
            s = poly.evaluate(p_idx)
            shares[p_idx] = s

        return shares, commitment


class DKG:
    def __init__(self, threshold, participant_indexes):
        self.threshold = threshold
        self.shares = {}
        self.commitments = {}
        self.indexes = participant_indexes

    def run(self):
        for idx in self.indexes:
            poly_sk = generate_sk()
            poly = Polynomial(poly_sk, self.threshold, order)
            poly.generate_random()
            commitment = poly.coefficients_commitment()

            shares = {}
            for p_idx in self.indexes:
                s = poly.evaluate(p_idx)
                if p_idx not in shares:
                    shares[p_idx] = []
                shares[p_idx] = s

            self.add_participant(idx, shares, commitment)

    def add_participant(self, id, shares, commitment):
        if len(shares) < self.threshold:
            raise AssertionError("not enough shares, required: ", len(shares))
        if id in self.commitments:
            raise AssertionError("% already distributed shares", id)

        for s in shares:
            p_id = s
            share = shares[s]
            if p_id in self.shares:
                self.shares[p_id].append(share)
            else:
                self.shares[p_id] = [share]
        self.commitments[id] = [c for c in commitment]

    def calculate_participants_sks(self):
        if len(self.shares) < self.threshold:
            raise AssertionError("not enough participants, requires: %d", self.threshold)

        group_secrets = {}
        for p in self.shares:
            p_shares = self.shares[p]
            sk = sum([s for s in p_shares]) % order
            group_secrets[p] = sk

        return group_secrets

    def group_sk(self):
        sks = self.calculate_participants_sks()
        l = LagrangeInterpolation(sks, order)
        return  l.evaluate()

    def group_pk(self):
        sks = self.calculate_participants_sks()
        pks = {}
        for i in sks:
            pks[i] = pk_from_sk(sks[i])
        return reconstruct_pk(pks)

