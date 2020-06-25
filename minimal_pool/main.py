import participant
import config
import threading
import logging
import crypto

last_logged_epoch = 0
participants = []

def main():
    global participants

    logging.debug("creating %d participants via DKG",config.NUM_OF_PARTICIPANTS)
    ids = range(1,config.NUM_OF_PARTICIPANTS)
    dkg = crypto.DKG(config.POOL_THRESHOLD, ids)
    dkg.run()
    sks = dkg.calculate_group_sk()
    for i in sks:
        p = participant.Participant(i,sks[i])
        participants.append(p)

    # connect all participants together
    logging.debug("connecting participants to each-other")
    for i in range(len(participants)):
        # connect nodes to eachother
        for j in range(i+1,len(participants)):
            participants[i].node.connect(participants[j].node)
            participants[j].node.connect(participants[i].node)

    # subscribe all nodes to topics
    for t in [config.MSG_SHARE_DISTRO,config.MSG_EPOCH_SIG]:
        for p in participants:
            p.node.subscribe_to_topic(p.id,t)

    # start epoch execution
    [threading.Thread(target=p.node.execute_epoch, daemon=True).start() for p in participants]

    # start epoch logging
    run_continously(participants[1].node)

def run_continously(node):
    threading.Timer(config.EPOCH_TIME + 3, log_end_of_round,args=[node]).start()

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
        logging.debug("P(%d) shares received: %d",p.id,len(shares))

        sigs = p.node.state.aggregated_sig_for_epoch(last_logged_epoch)
        logging.debug("P(%d) sig verified: %s, sig count: %d",p.id,sigs["is_verified"],len(sigs["pks"]))
    logging.debug("\n\n-------------------------------------------------\n")

    last_logged_epoch = last_logged_epoch + 1

    # run again
    run_continously(node)



if __name__ == '__main__':
    logging.basicConfig(format='%(asctime)s-%(levelname)s-%(message)s',level=logging.DEBUG)
    main()