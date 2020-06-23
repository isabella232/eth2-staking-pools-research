import uuid
import config
import threading
from crypto import hash,KEY_SIZE_BITS
import random

class State:
    def __init__(self,seed):
        self.seed = seed
        self.epoch = 0

    def increase_epoch(self):
        self.epoch += 1

    def mix_seed(self):
        self.seed = (self.seed * int.from_bytes(hash(self.epoch.to_bytes(length=32,byteorder=config.ENDIANNESS)),byteorder=config.ENDIANNESS)) % KEY_SIZE_BITS

class Message:
    def __init__(self, type, data, sender_id,topic):
        self.id = uuid.uuid4()
        self.type = type
        self.data = data
        self.sender_id = sender_id
        self.topic = topic

"""
This is a simple pool node that emulates a node in a distributed network 
"""
class PoolNode:
    def __init__(self,id,subscriber):
        self.id = id
        self.peers = []
        self.subscriber = subscriber
        self.topics = {}
        self.state = State(config.GENESIS_SEED)
        self.known_messages = []

    def execute_round(self):
        self.state.increase_epoch()
        self.state.mix_seed()
        self.send_to_subscriberr(Message(
                config.MSG_NEW_EPOCH,
                {"epoch":self.state.epoch},
                self.id,
                None
            )
        )

        # run it again
        threading.Timer(config.EPOCH_TIME, self.execute_round).start()

    def current_epoch_pool_assignment(self, index):
        lst = list(range(config.NUM_OF_PARTICIPANTS))
        rnd = random.Random(self.state.seed)
        rnd.shuffle(lst)
        return lst[index]

    def current_epoch_pools(self):
        pools = {}
        for i in range(config.NUM_OF_PARTICIPANTS):
            pool_id = self.current_epoch_pool_assignment(i) % config.NUMBER_OF_POOLS
            if pool_id in pools:
                pools[pool_id].append(i)
            else:
                pools[pool_id] = [i]
        return pools

    """
        Networking
    """

    def connect(self, node):
        if node.id != self.id:
            self.peers.append(node)

    def disconnect(self,node):
        self.peers.remove(node)

    def recieve(self,msg):
        # do not handle known messages
        if msg.id in self.known_messages:
            return
        else:
            self.known_messages.append(msg.id)

        if msg.type == config.MSG_SHARE_DISTRO:
            if msg.data["p"] == self.id: # send only for specific participant
                self.send_to_subscriberr(msg)

        # let other nodes know
        self.send(msg)

    def send(self,msg):
        # let other nodes know
        for p in self.peers:
            if msg.topic in self.topics and p.id in self.topics[msg.topic]:
                p.recieve(msg)

    def subscribe_to_topic(self,sender_id,topic):
        if topic in self.topics and sender_id in self.topics[topic]:
            return

        if topic in self.topics:
            if sender_id not in self.topics[topic]:
                self.topics[topic].append(sender_id)
        else:
            self.topics[topic] = [sender_id]

        # let other nodes know
        for p in self.peers:
            p.subscribe_to_topic(sender_id,topic)

    def remove_from_topic(self,sender_id,topic):
        if topic in self.topics and sender_id in self.topics[topic]:
            self.topics[topic].remove(sender_id)

        # let other nodes know
        for p in self.peers:
            p.remove_from_topic(sender_id, topic)

    def send_to_subscriberr(self,msg):
        if self.subscriber != None:
            self.subscriber(msg)

    """
        messages
    """
    def broadcast_shares(self,sender_id,shares,pool_id):
        for i in range(len(shares)):
            msg = Message("share_distro",{"p":(i+1),"share":shares[i]},sender_id,pool_id) # p is the participant's index and we assume shares are ordered
            self.send(msg)