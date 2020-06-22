import uuid

class State:
    def __init__(self):
        self.d = 1

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
        self.state = State()
        self.known_messages = []

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

        if msg.type == "share_distro":
            if msg.data["p"] == self.id: # send only for specific participant
                self.subscriber(msg)

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

    """
        messages
    """
    def broadcast_shares(self,sender_id,shares,pool_id):
        for i in range(len(shares)):
            msg = Message("share_distro",{"p":(i+1),"share":shares[i]},sender_id,pool_id) # p is the participant's index and we assume shares are ordered
            self.send(msg)