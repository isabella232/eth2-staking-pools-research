import uuid

class Message:
    def __init__(self, type, data, sender_id):
        self.id = uuid.uuid4()
        self.type = type
        self.data = data
        self.sender_id = sender_id
