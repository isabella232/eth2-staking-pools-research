import participant
from crypto import generate_sk

POOL_SIZE = 3
POOL_THRESHOLD = 2
NUM_OF_PARTICIPANTS = POOL_SIZE * 3
EPOCH_TIME = 6 # seconds

def main():
    participants = []

    print("creating ",NUM_OF_PARTICIPANTS," participants")
    for i in range(NUM_OF_PARTICIPANTS):
        p = participant.Participant(i+1)
        sk = generate_sk()
        participants.append(p)
        print("     participant ",i, " initializing with secret: ",sk)
        p.generate_polynomial(sk,POOL_THRESHOLD)

    # connect all participants together
    print("connecting participants to each-other")
    for i in range(len(participants)):
        for j in range(i+1,len(participants)):
            print("     connecting participants ",i," to ",j)
            participants[i].node.connect(participants[j].node)
            print("     connecting participants ",j," to ",i)
            participants[j].node.connect(participants[i].node)

    # subscribe all nodes to topics
    for t in ["shares_for_pool"]:
        for p in participants:
            p.node.subscribe_to_topic(p.id,t)

    # distribuite shares
    print("distribuiting shares")
    for i in range(len(participants)):
        p = participants[i]
        shares = p.distribuite_shares([participants[i].id for i in range(len(participants))])
        p.node.broadcast_shares(p.id,shares,"shares_for_pool")
        print("     participant ",p.id," shared shares: ",shares)

    # reconstruct individual secrets
    for p in participants:
        p.reconstruct_group_secret()

    # p = participant.Participant()
    # p.generate_polynomial(1111,20)
    # print(p.distribuite_shares([1,2,3,4,5,6,7,8,9]))
    # p.node.subscribe_to_topic(p.uuid,"hello")

if __name__ == '__main__':
    main()
