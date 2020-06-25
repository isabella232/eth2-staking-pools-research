import random
from py_ecc.bls import G2ProofOfPossession as py_ecc_bls
from py_ecc.optimized_bls12_381 import curve_order
import milagro_bls_binding as milagro_bls
from hashlib import sha256
import logging
import config
import math

bls = milagro_bls
order = 251#curve_order

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

class Polynomial:
    def __init__(self,secret,degree):
        self.coefficients = []
        self.secret = secret
        self.degree = degree

    def generate_random(self):
        self.coefficients.append(self.secret) # important it's in index 0, see valuate
        self.coefficients.extend([generate_sk() for _ in range(0, self.degree)])

    def evaluate(self,point):
        return sum([self.coefficients[i] * (point ** i) for i in range(len(self.coefficients))]) % order

    def coefficients_commitment(self):
        return [pk_from_sk(self.coefficients[i]) for i in range(len(self.coefficients))]

class LagrangeInterpolation:
    def __init__(self, shares):
        self.shares = shares

    def evaluate(self):
        return sum([self.shares[i] * self.lagrange_coefficients(i) for i in self.shares]) % order

    def lagrange_coefficients(self, i, x_0=0):
        num = 1
        den = 1
        for j in self.shares:
            if i != j:
                num = (num * (x_0 - j))
                den = (den * (i - j))
        return ((num % order) * self.mod_inverse(den, order)) % order

    def mod_inverse(self, b, m):
        g = math.gcd(b, m)
        if (g != 1):
            raise AssertionError("Inverse doesn't exist")
        else:
            # If b and m are relatively prime,
            # then modulo inverse is b^(m-2) mode m
            return pow(b, m - 2, m)

class DKG:
    def __init__(self, threshold, participant_indexes):
        self.threshold = threshold
        self.shares = {}
        self.commitments = {}
        self.indexes = participant_indexes

    def run(self):
        for idx in self.indexes:
            poly_sk = generate_sk()
            poly = Polynomial(poly_sk, self.threshold-1) # following Shamir's secret sharing, degree is threshold - 1
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
        self.commitments[id] = commitment

    def calculate_group_sk(self):
        if len(self.shares) < self.threshold:
            raise AssertionError("not enough participants, requires: %d", self.threshold)

        group_secrets = {}
        for p in self.shares:
            p_shares = self.shares[p]
            sk = sum([s for s in p_shares]) % order
            group_secrets[p] = sk

        return group_secrets

