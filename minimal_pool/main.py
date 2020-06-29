import participant
import config
import threading
import logging
import crypto
from node import node

last_logged_epoch = 0
participants = []

def main():
    global participants

    n = node.PoolNode(-1, None)
    epoch_0 = n.state.get_epoch(0)

    pool_pk = {}
    for p_idx in range(1, config.NUMBER_OF_POOLS + 1):
        ids = epoch_0.pool_participants_by_id(p_idx)
        logging.debug("Pool %d participants: %s (via DKG)", p_idx, ids)

        dkg = crypto.DKG(config.POOL_THRESHOLD - 1, ids)  # following Shamir's secret sharing, degree is threshold - 1
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
    threading.Timer(config.EPOCH_TIME + 3, log_end_of_round, args=[node]).start()

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


from  py_ecc.bls.g2_primatives import G1_to_pubkey
import time

def recon():
    epoch = 1
    shares = {
        6: {4: 14718138269560737216416227068691712518063415736892723469991188442463872446793, 5: 19473304902709436137684775681958310937706886046016930823718672551874933384376, 6: 24228471535858135058953324295224909357350356355141138177446156661285994321959},
        3: {4: 19107830316373382133392187649707847744762743540741444818782689361709666871276, 5: 6512133818541930846602739912419488497525745439281473109838990550931696893286, 6: 46352312495836670039261032683317095087979299838349139223498950440092308099809},
        2: {1: 24823305088538841062236704057955487963920130642512542631945894719330085420101, 2: 3787328452096883045094978748880236782414559497208749947976983348366856371279, 3: 35187226990781115507400993947990951438599540852432595086611730677342208506970},
        1: {4: 29207243774599538527733816913975530521349546213513835131527351551096198890153, 5: 37425137065026029827672967458656577724541870831190510130709247095962047926719, 6: 45643030355452521127612118003337624927734195448867185129891142640827896963285},
        5: {1: 32271963851158130763083596796342000256868737635685158382787658212760448963158, 2: 5799072670918972458102592684613643042276270285224173671479449834578545812583, 3: 31762056665806004632569329081071251665374355435290826782774900156335223846521},
        4: {1: 46976724845799864745806785285087711627229943594023290574803958741247294647441, 2: 22025160703965873593295113339722949945780332434501161374037851491789056001248, 3: 49509471737258072920231181902544154102021273775506669995875402942269398539568}
    }

    ordered_shares = {}
    for from_p_idx in shares:
        srs = shares[from_p_idx]
        for dest_p_idx in srs:
            if dest_p_idx not in ordered_shares:
                ordered_shares[dest_p_idx] = {}
            ordered_shares[dest_p_idx][from_p_idx] = srs[dest_p_idx]

    pool_1 = [4,5,6]
    sks_1 = {}
    for p_idx in pool_1:
        sk = crypto.reconstruct_sk(ordered_shares[p_idx])
        sks_1[p_idx] = sk

    sk_1 = crypto.reconstruct_sk(sks_1)
    logging.debug("sk 1: %d", sk_1)

    pool_2 = [1,2,3]
    sks_2 = {}
    for p_idx in pool_2:
        sk = crypto.reconstruct_sk(ordered_shares[p_idx])
        sks_2[p_idx] = sk
    sk_2 = crypto.reconstruct_sk(sks_2)
    logging.debug("sk 2: %d", sk_2)


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
    # recon()
    # benchmark_dkg()
