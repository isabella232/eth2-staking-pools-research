import participant
import config
import threading
import logging
import crypto

last_logged_epoch = 0
participants = []

def main():
    global participants

    pool_pk = {}
    for p_idx in range(1, config.NUMBER_OF_POOLS + 1):

        logging.debug("Pool %d with %d participants (via DKG)", p_idx, config.POOL_SIZE)
        ids = range(1 + (p_idx - 1) * config.POOL_SIZE,
                    ((p_idx - 1) * config.POOL_SIZE) + config.POOL_SIZE + 1)
        dkg = crypto.DKG(config.POOL_THRESHOLD - 1, ids) # following Shamir's secret sharing, degree is threshold - 1
        dkg.run()
        sks = dkg.calculate_participants_sks()
        logging.debug("     Group sk: %s", dkg.group_sk())
        logging.debug("     Group pk: %s", crypto.readable_pk(dkg.group_pk()).hex())

        for i in sks:
            p = participant.Participant(i, sks[i])
            participants.append(p)

        pool_pk[p_idx] = crypto.readable_pk(dkg.group_pk())

    # connect all participants together and update them with groups
    logging.debug("connecting participants to each-other")
    for i in range(len(participants)):
        for p_idx in pool_pk:
            participants[i].node.state.save_pool_info(p_idx, pool_pk[p_idx])
        # connect nodes to each other
        for j in range(i+1, len(participants)):
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

    epoch = node.state.get_epoch(last_logged_epoch)
    pools = epoch.pools_info
    # """
    #     Epoch stats
    # """
    logging.debug("\n\n----------------EPOCH %d Summary ----------------\n", last_logged_epoch)
    logging.debug("Pools for epoch %d: %s", last_logged_epoch, pools)
    for p in participants:
        shares = epoch.participant_shares(p.id)
        logging.debug("P(%d) shares received: %d", p.id, len(shares))

    for p_id in pools:
        sigs = epoch.aggregated_sig_for_pool(p_id)
        if sigs is not None:
            logging.debug("pool %d %s sig verified: %s",
                              p_id,
                              sigs["ids"],
                              sigs["is_verified"],
                          )
        else:
            logging.debug("pool %d no sigs found", p_id)
    logging.debug("\n\n-------------------------------------------------\n")

    last_logged_epoch = last_logged_epoch + 1

    # run again
    run_continously(node)


from  py_ecc.bls.g2_primatives import G1_to_pubkey,pubkey_to_G1
import time
def benchmark_dkg():
    num_of_part = 10
    logging.debug("creating %d participants via DKG", num_of_part)
    ids = range(1, num_of_part)

    start = time.time()
    dkg = crypto.DKG(config.POOL_THRESHOLD, ids)
    dkg.run()
    sks = dkg.calculate_participants_sks()
    logging.debug("     Group sk:       %s", dkg.group_sk())
    logging.debug("     Group sig:      %s", crypto.readable_sig(crypto.sign_with_sk(dkg.group_sk(),config.TEST_EPOCH_MSG)).hex())
    end = time.time()
    logging.debug("dkg: %f sec", (end-start))
    logging.debug("     Group pk:       %s", crypto.readable_pk(dkg.group_pk()).hex())
    logging.debug("     real Group pk:  %s", crypto.readable_pk(crypto.pk_from_sk(dkg.group_sk())).hex())

    sigs = {}
    pks = {}
    for sk in sks:
        sig = crypto.sign_with_sk(sks[sk], config.TEST_EPOCH_MSG)
        pk = crypto.pk_from_sk(sks[sk])

        sigs[sk] = sig
        pks[sk] = pk


    # reconstruct sig and pk
    start=time.time()
    recon_pk = crypto.reconstruct_pk(pks)
    recon_sig = crypto.reconstruct_group_sig(sigs)
    end = time.time()
    logging.debug("reconstruct sk/pk: %f sec", (end - start))

    recon_pk = crypto.readable_pk(recon_pk)
    logging.debug("pk after reconstruction: %s", recon_pk.hex())
    recon_sig = crypto.readable_sig(recon_sig)
    logging.debug("sig after reconstruction: %s", recon_sig.hex())


    ## redistribuite
    start = time.time()
    re_distro_shares = {}
    re_distro_comm = {}
    for p_indx in sks:
        sk = sks[p_indx]

        redistro = crypto.Redistribuition(config.POOL_THRESHOLD -1, sk, ids)  # following Shamir's secret sharing, degree is threshold - 1
        shares, commitment = redistro.generate_shares()
        for i in ids:
            if i not in re_distro_shares:
                re_distro_shares[i] = {}
            re_distro_shares[i][p_indx] = shares[i]
        re_distro_comm[p_indx] = commitment

    sk_per_id = {}
    pk_per_id = {}
    for idx in ids:
        sk_per_id[idx] = crypto.reconstruct_sk(re_distro_shares[idx])
        pk_per_id[idx] = crypto.pk_from_sk(sk_per_id[idx])

    group_pk_after_redistro = G1_to_pubkey(crypto.reconstruct_pk(pk_per_id))

    logging.debug("     Group pk after redistro: %s", group_pk_after_redistro.hex())
    end = time.time()
    logging.debug("redistro: %f sec", (end - start))


if __name__ == '__main__':
    logging.basicConfig(format='%(asctime)s-%(levelname)s-%(message)s', level=logging.DEBUG)
    main()
    # benchmark_dkg()
