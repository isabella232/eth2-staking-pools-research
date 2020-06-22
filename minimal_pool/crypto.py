import random

KEY_SIZE_BITS = 256

class Polynomial:
    def __init__(self,secret,degree):
        self.coefficients = []
        self.secret = secret
        self.degree = degree

    def generate_random(self):
        self.coefficients.append(self.secret)
        self.coefficients.extend([random.getrandbits(KEY_SIZE_BITS) for _ in range(1, self.degree)])

    def evaluate(self,point):
        return sum([self.coefficients[i] * (point ** i) for i in range(len(self.coefficients))])