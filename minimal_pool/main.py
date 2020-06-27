import participant
import config
import threading
import logging
import crypto

last_logged_epoch = 0
participants = []

def main():
    global participants

    logging.debug("creating a %d participants pool via DKG",config.NUM_OF_PARTICIPANTS)
    ids = range(1, config.NUM_OF_PARTICIPANTS+1)
    dkg = crypto.DKG(config.POOL_THRESHOLD - 1, ids) # following Shamir's secret sharing, degree is threshold - 1
    dkg.run()
    sks = dkg.calculate_participants_sks()
    logging.debug("     Group sk: %s", dkg.group_sk())
    logging.debug("     Group pk: %s", dkg.group_pk().hex())

    for i in sks:
        p = participant.Participant(i, sks[i])
        p.node.state.save_pool_info(1, dkg.group_pk())
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
        logging.debug("P(%d) sig verified: %s",
                          p.id,
                          sigs["is_verified"],
                      )
    logging.debug("\n\n-------------------------------------------------\n")

    last_logged_epoch = last_logged_epoch + 1

    # run again
    run_continously(node)


from  py_ecc.bls.g2_primatives import G1_to_pubkey,pubkey_to_G1
if __name__ == '__main__':
    logging.basicConfig(format='%(asctime)s-%(levelname)s-%(message)s',level=logging.DEBUG)
    main()

    # logging.debug("creating %d participants via DKG", config.NUM_OF_PARTICIPANTS)
    # ids = range(1, config.NUM_OF_PARTICIPANTS + 1)
    # dkg = crypto.DKG(config.POOL_THRESHOLD, ids)
    # dkg.run()
    # sks = dkg.calculate_participants_sks()
    # logging.debug("     Group sk: %s", dkg.group_sk())
    # logging.debug("     Group pk:       %s", dkg.group_pk().hex())
    # logging.debug("     real Group pk:  %s", crypto.pk_from_sk(dkg.group_sk()).hex())
    #
    # sigs = []
    # pks = []
    # for sk in sks:
    #     sig = crypto.sign_with_sk(sks[sk],config.TEST_EPOCH_MSG)
    #     pk = crypto.pk_from_sk(sks[sk])
    #     is1 = crypto.verify_sig(
    #         pk,
    #         config.TEST_EPOCH_MSG,
    #         sig)
    #     logging.debug("verified %d with group pk %s", sk,is1)
    #     sigs.append(sig)
    #     pks.append(pk)
    #
    # agg = crypto.aggregate_sigs(sigs)
    # is1 = crypto.verify_aggregated_sigs(
    #     pks,
    #     config.TEST_EPOCH_MSG,
    #     agg)
    # logging.debug("verified with multi pk %s", is1)
    #
    # agg_pks = crypto.aggregate_pks(pks)
    # is1 = crypto.verify_sig(
    #     agg_pks,
    #     config.TEST_EPOCH_MSG,
    #     agg)
    # logging.debug("verified with group sk/pk %s", is1)
    #
    #
    #
    #
    #
    # ## redistribuite
    # re_distro_shares = {}
    # re_distro_comm = {}
    # for p_indx in sks:
    #     sk = sks[p_indx]
    #
    #     redistro = crypto.Redistribuition(config.POOL_THRESHOLD -1, sk, ids) # following Shamir's secret sharing, degree is threshold - 1
    #     shares, commitment = redistro.generate_shares()
    #     for idx in ids:
    #         if p_indx not in re_distro_shares:
    #             re_distro_shares[p_indx] = {}
    #         re_distro_shares[p_indx][idx] = shares[idx]
    #     re_distro_comm[p_indx] = commitment
    #
    # sk_per_id = {}
    # pk_per_id = {}
    # for idx in ids:
    #     sk_per_id[idx] = crypto.reconstruct_sk(re_distro_shares[idx])
    #     pk_per_id[idx] = crypto._optimized_pk_from_sk(sk_per_id[idx])
    #
    # group_pk_after_redistro = G1_to_pubkey(crypto.reconstruct_pk(pk_per_id))
    #
    # logging.debug("     Group pk after redistro: %s", group_pk_after_redistro.hex())
