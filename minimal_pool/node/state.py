import config
import random
import threading
import crypto

class Epoch:
    def __init__(self, number, seed, agg_sigs = {}, shares = {}):
        self.number = number
        self.seed = seed
        self.pools_info = self._calculate_pools()
        self.agg_sigs = agg_sigs
        self.shares = shares

    def save_participant_shares(self, shares, p_id):
        self.shares[p_id] = shares

    def participant_shares(self, p_id):
        if p_id in self.shares:
            return self.shares[p_id]
        return []

    def save_aggregated_sig(self, pool_id, sig, ids, is_verified):
        self.agg_sigs[pool_id] = {
            "sig": sig.hex(),
            "ids": ids,
            "is_verified": bool(is_verified),
        }

    def aggregated_sig_for_pool(self, pool_id):
        if pool_id not in self.agg_sigs:
            return None
        return self.agg_sigs[pool_id]

    def pool_id_for_participant(self, index):
        lst = list(range(1, config.NUM_OF_PARTICIPANTS+1))  # indexes must run from 1
        rnd = random.Random(self.seed)
        rnd.shuffle(lst)
        return lst[index-1] % config.NUMBER_OF_POOLS + 1  # indexes must run from 1

    def pool_participants_by_id(self, pool_id):
        pools = self._calculate_pools()
        return pools[pool_id]

    def _calculate_pools(self):
        pools = {}
        for i in range(1, config.NUM_OF_PARTICIPANTS + 1):  # indexes must run from 1
            pool_id = self.pool_id_for_participant(i)
            if pool_id in pools:
                pools[pool_id].append(i)
            else:
                pools[pool_id] = [i]
        return pools

class State:
    def __init__(self, seed):
        self.seed = seed
        self.epoch = config.STARTING_EPOCH
        self.pool_info = {}
        # self.pool_per_epoch = {}
        # self.shares_per_epoch = {}
        # self.aggregated_sig = {}

        # locks
        self.epochs_lock = threading.Lock()

        self.epochs = {}

    def _increase_epoch(self):
        self.epoch += 1

    def _mix_seed(self):
        epoch_number_bytes = self.epoch.to_bytes(32, config.ENDIANNESS)
        mixer = int.from_bytes(crypto.hash(epoch_number_bytes), config.ENDIANNESS)
        self.seed = (self.seed * mixer) % config.KEY_SIZE_BITS

    """
        will create a new epoch for current epoch number
    """
    def new_poch(self):
        with self.epochs_lock:
            self._increase_epoch()
            self._mix_seed()

            e = Epoch(
                self.epoch,
                self.seed
            )
            self.epochs[self.epoch] = e
            return e

    def save_epoch(self, epoch):
        with self.epochs_lock:
            self.epochs[epoch.number] = epoch

    def get_epoch(self, epoch_number):
        with self.epochs_lock:
            return self.epochs[epoch_number]

    def save_pool_info(self, pool_id, pk):
        self.pool_info[pool_id] = {
            "pk": pk
        }

    def pool_info_by_id(self, pool_id):
        if pool_id not in self.pool_info:
            raise AssertionError("%d pool id does not exist", pool_id)
        return self.pool_info[pool_id]



    # def save_pool_participants(self,pools,epoch):
    #     self.pool_per_epoch[epoch] = pools
    #
    # def pool_participants_for_epoch(self, epoch):
    #     if epoch not in self.pool_per_epoch:
    #         raise AssertionError("%d epoch could not be found", epoch)
    #     return self.pool_per_epoch[epoch]

    # def save_epoch_sig(self, sig, pk, is_verified, epoch):
    #     self.aggregated_sig[epoch] = {
    #         "sig": sig.hex(),
    #         "pks": pk.hex(),
    #         "is_verified": bool(is_verified),
    #     }
    #
    # def aggregated_sig_for_epoch(self,epoch):
    #     if epoch not in self.aggregated_sig:
    #         raise AssertionError("%d epoch could not be found", epoch)
    #     return self.aggregated_sig[epoch]


    # def save_participant_shares(self,shares,epoch,p_id):
    #     self.shares_per_epoch[epoch] = {p_id:shares}
    #
    # def participant_shares_for_epoch(self,epoch,p_id):
    #     if epoch in self.shares_per_epoch:
    #         return self.shares_per_epoch[epoch][p_id]
    #     return []