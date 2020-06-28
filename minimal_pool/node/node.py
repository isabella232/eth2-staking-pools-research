from node.state import State
from node.message import Message
import config
import threading
import random


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

    def chain_round(self):
        # mark mid round
        threading.Timer(config.EPOCH_TIME / 2, self.mid_round_mark).start()
        # run it again
        threading.Timer(config.EPOCH_TIME, self.execute_epoch).start()

    def execute_epoch(self):
        # save epoch
        # self.state.save_pool_participants(self.current_epoch_pools(), self.state.epoch)

        # end epoch only if not first
        if self.state.epoch != config.STARTING_EPOCH:
            e = self.state.get_epoch(self.state.epoch)
            self.send_to_subscriber(Message(
                    config.MSG_END_EPOCH,
                    {"epoch": e},
                    self.id,
                )
            )

        # start new epoch
        e = self.state.new_poch()
        self.send_to_subscriber(Message(
                config.MSG_NEW_EPOCH,
                {"epoch": e},
                self.id,
            )
        )

        self.chain_round()

    def mid_round_mark(self):
        e = self.state.get_epoch(self.state.epoch)
        self.send_to_subscriber(Message(
                config.MSG_MID_EPOCH,
                {"epoch": e},
                self.id,
            )
        )

    # def current_epoch_pool_assignment(self, index):
    #     lst = list(range(1, config.NUM_OF_PARTICIPANTS+1)) # indexes must run from 1
    #     rnd = random.Random(self.state.seed)
    #     rnd.shuffle(lst)
    #     return lst[index-1] % config.NUMBER_OF_POOLS + 1 # indexes must run from 1
    #
    # def current_epoch_pool_participants(self, pool_id):
    #     pools = self.current_epoch_pools()
    #     return pools[pool_id]
    #
    # def current_epoch_pools(self):
    #     pools = {}
    #     for i in range(1, config.NUM_OF_PARTICIPANTS+1): # indexes must run from 1
    #         pool_id = self.current_epoch_pool_assignment(i)
    #         if pool_id in pools:
    #             pools[pool_id].append(i)
    #         else:
    #             pools[pool_id] = [i]
    #     return pools

    """
        Networking
    """

    def connect(self, node):
        self.peers.append(node)


    def disconnect(self,node):
        self.peers.remove(node)

    def recieve(self,msg):
        # do not handle known messages
        if msg.id in self.known_messages:
            return
        else:
            self.known_messages.append(msg.id)

        if msg.type == config.MSG_SHARE_DISTRO or msg.type == config.MSG_EPOCH_SIG:
            self.send_to_subscriber(msg)

        # TODO - let other nodes know
        #self.send(msg)

    def send(self,msg):
        for p in self.peers:
            if msg.type in self.topics and p.id in self.topics[msg.type]:
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

    def send_to_subscriber(self, msg):
        if self.subscriber != None:
            self.subscriber(msg)

    """
        messages
    """
    def broadcast_shares(self, epoch, sender_id, shares, commitments, pool_id):
        for p_indx in shares:
            s = shares[p_indx]
            msg = Message(
                config.MSG_SHARE_DISTRO,
                {
                    "from_p_id": sender_id,
                    "p": p_indx,
                    "share": s,
                    "commitments": commitments,
                    "pool_id": pool_id,
                    "epoch": epoch,
                 },
                sender_id
            ) # p is the participant's index and we assume shares are ordered
            self.send(msg)

    def broadcast_sig(self, epoch, sender_id, sig, pk, pool_id):
        msg = Message(
            config.MSG_EPOCH_SIG,
            {
                "from_p_id": sender_id,
                "sig": sig,
                "pk": pk,
                "pool_id": pool_id,
                "epoch": epoch,
            },
            sender_id
        )
        self.send(msg)