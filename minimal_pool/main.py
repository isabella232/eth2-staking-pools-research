import participant
from crypto import generate_sk
import config
import threading
import logging

last_logged_epoch = 0
participants = []

def main():
    global participants

    logging.debug("creating %d participants",config.NUM_OF_PARTICIPANTS)
    for i in range(config.NUM_OF_PARTICIPANTS):
        p = participant.Participant(i+1) # id can't be 0 as it's the secret
        sk = generate_sk()
        participants.append(p)
        p.secret = sk

    # connect all participants together
    logging.debug("connecting participants to each-other")
    for i in range(len(participants)):
        # participants[i].node.connect(participants[i].node) # connect to self to propogate msg, could be optimized
        # connect nodes to eachother
        for j in range(i+1,len(participants)):
            participants[i].node.connect(participants[j].node)
            participants[j].node.connect(participants[i].node)

    # subscribe all nodes to topics
    for t in [config.MSG_SHARE_DISTRO]:
        for p in participants:
            p.node.subscribe_to_topic(p.id,t)

    # distribuite shares
    # print("distribuiting shares")
    # for i in range(len(participants)):
    #     p = participants[i]
    #     shares = p.distribuite_shares([participant_index_to_distro_index(participants[i].id) for i in range(len(participants))])
    #     p.node.broadcast_shares(p.id,shares,"shares_for_pool")
    #     print("     participant ",p.id," shared shares: ",shares)
    #
    # # reconstruct individual secrets
    # for p in participants:
    #     p.reconstruct_group_secret()
    #
    # # sign
    # message = b'\xab' * 32
    # sigs = []
    # for p in participants:
    #     sigs.append(p.sign(message))

    # aggregate and verify
    # aggregated = aggregate_signatures(sigs)
    # pks = [p.pub_group_key() for p in participants]
    # is_verified = verify_aggregated_sigs(pks,message,aggregated)
    #
    # print("verified aggregated sig: " + str(is_verified))

    # start epoch execution
    [threading.Thread(target=p.node.execute_epoch(), args=(p.id,), daemon=True).start() for p in participants]

    # start epoch logging
    logging.debug("start epoch logging")
    run_continously(participants[1].node)

def run_continously(node):
    threading.Timer(config.EPOCH_TIME+1, log_end_of_round,args=[node]).start()

def log_end_of_round(node):
    global last_logged_epoch

    pools = node.state.pool_participants_for_epoch(last_logged_epoch)
    # """
    #     Epoch stats
    # """
    logging.debug("\n\n----------------EPOCH %d Summary ----------------\n",last_logged_epoch)
    logging.debug("Pools for epoch %d: %s", last_logged_epoch, pools)
    for p in participants:
        shares = p.node.state.participant_shares_for_epoch(last_logged_epoch,p.id)
        logging.debug("P(%d) shares recieved: %d",p.id,len(shares))
    logging.debug("\n\n-------------------------------------------------\n")

    last_logged_epoch = last_logged_epoch + 1

    # run again
    run_continously(node)



if __name__ == '__main__':
    logging.basicConfig(format='%(asctime)s-%(levelname)s-%(message)s',level=logging.DEBUG)
    main()
