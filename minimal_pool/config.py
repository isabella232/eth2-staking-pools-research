POOL_SIZE = 4
NUMBER_OF_POOLS = 2
POOL_THRESHOLD = 3
NUM_OF_PARTICIPANTS = NUMBER_OF_POOLS * POOL_SIZE
EPOCH_TIME = 6  # seconds
# used to deterministically simulate a random source for the network
GENESIS_SEED = 95637185274827421279466020522819461017051869823440212610124288004910264493064
STARTING_EPOCH = -1  # TODO - this is set to -1 because when execute_epoch (Node) is called it increases it by +1.

MSG_SHARE_DISTRO        = "share_distro"
MSG_EPOCH_SIG           = "epoch_sig"
MSG_NEW_EPOCH           = "new_epoch"
MSG_END_EPOCH           = "end_epoch"
MSG_MID_EPOCH           = "mid_epoch"

ENDIANNESS = "big"
KEY_SIZE_BYTES = 32
KEY_SIZE_BITS = KEY_SIZE_BYTES * 8

TEST_EPOCH_MSG = b'\xab' * 32