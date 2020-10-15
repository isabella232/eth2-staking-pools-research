package shared

import (
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/core"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/shared/params"
	"github.com/prysmaticlabs/go-bitfield"
)

/**
	def is_slashable_attestation_data(data_1: AttestationData, data_2: AttestationData) -> bool:
		"""
		Check if ``data_1`` and ``data_2`` are slashable according to Casper FFG rules.
		"""
		return (
			# Double vote
			(data_1 != data_2 and data_1.target.epoch == data_2.target.epoch) or
			# Surround vote
			(data_1.source.epoch < data_2.source.epoch and data_2.target.epoch < data_1.target.epoch)
		)
 */
func IsSlashableAttestationData (att1 *core.AttestationData, att2 *core.AttestationData) bool {
	return (!core.AttestationDataEqual(att1, att2) && att1.Target.Epoch == att2.Target.Epoch) ||
		(att1.Source.Epoch < att2.Source.Epoch && att2.Target.Epoch < att2.Target.Epoch)
}

/**
def is_valid_indexed_attestation(state: BeaconState, indexed_attestation: IndexedAttestation) -> bool:
    """
    Check if ``indexed_attestation`` is not empty, has sorted and unique indices and has a valid aggregate signature.
    """
    # Verify indices are sorted and unique
    indices = indexed_attestation.attesting_indices
    if len(indices) == 0 or not indices == sorted(set(indices)):
        return False
    # Verify aggregate signature
    pubkeys = [state.validators[i].pubkey for i in indices]
    domain = get_domain(state, DOMAIN_BEACON_ATTESTER, indexed_attestation.data.target.epoch)
    signing_root = compute_signing_root(indexed_attestation.data, domain)
    return bls.FastAggregateVerify(pubkeys, signing_root, indexed_attestation.signature)
// TODO - is_valid_indexed_attestation
 */
func IsValidIndexedAttestation(state *core.State) bool {
	return false
}

/**
def compute_committee(indices: Sequence[ValidatorIndex],
                      seed: Bytes32,
                      index: uint64,
                      count: uint64) -> Sequence[ValidatorIndex]:
    """
    Return the committee corresponding to ``indices``, ``seed``, ``index``, and committee ``count``.
    """
    start = (len(indices) * index) // count
    end = (len(indices) * uint64(index + 1)) // count
    return [indices[compute_shuffled_index(uint64(i), uint64(len(indices)), seed)] for i in range(start, end)]
 */
func ComputeCommittee(indices []uint64, seed []byte, index uint64, count uint64) ([]uint64, error) {
	start := uint64(len(indices)) * index / count
	end := uint64(len(indices)) * uint64(index + 1) / count

	ret := []uint64{}
	for i := start ; i < end ; i++ {
		idx, err := computeShuffledIndex(i, uint64(len(indices)), SliceToByte32(seed), true,10) // TODO - shuffle round via config
		if err != nil {
			return []uint64{}, err
		}

		ret = append(ret, idx)
	}
	return ret, nil
}

/**
def get_committee_count_per_slot(state: BeaconState, epoch: Epoch) -> uint64:
    """
    Return the number of committees in each slot for the given ``epoch``.
    """
    return max(uint64(1), min(
        MAX_COMMITTEES_PER_SLOT,
        uint64(len(get_active_validator_indices(state, epoch))) // SLOTS_PER_EPOCH // TARGET_COMMITTEE_SIZE,
    ))
 */
func GetCommitteeCountPerSlot(state *core.State, slot uint64) uint64 {
	epoch := ComputeEpochAtSlot(slot)
	bps := GetActiveBlockProducers(state, epoch)
	committeePerSlot := uint64(len(bps)) / params.ChainConfig.SlotsInEpoch / params.ChainConfig.MinAttestationCommitteeSize

	if committeePerSlot > params.ChainConfig.MaxCommitteesPerSlot {
		return params.ChainConfig.MaxCommitteesPerSlot
	}
	if committeePerSlot == 0 {
		return 1
	}
	return committeePerSlot
}

/**
def get_beacon_committee(state: BeaconState, slot: Slot, index: CommitteeIndex) -> Sequence[ValidatorIndex]:
    """
    Return the beacon committee at ``slot`` for ``index``.
    """
    epoch = compute_epoch_at_slot(slot)
    committees_per_slot = get_committee_count_per_slot(state, epoch)
    return compute_committee(
        indices=get_active_validator_indices(state, epoch),
        seed=get_seed(state, epoch, DOMAIN_BEACON_ATTESTER),
        index=(slot % SLOTS_PER_EPOCH) * committees_per_slot + index,
        count=committees_per_slot * SLOTS_PER_EPOCH,
    )
 */
func GetAttestationCommittee(state *core.State, slot uint64, index uint64) ([]uint64, error) {
	epoch := ComputeEpochAtSlot(slot)
	committeesPerSlot := GetCommitteeCountPerSlot(state, slot)
	seed := GetSeed(state, epoch, params.ChainConfig.DomainBeaconAttester)
	return ComputeCommittee(
			GetActiveBlockProducers(state, epoch),
			seed[:],
			(slot & params.ChainConfig.SlotsInEpoch) * committeesPerSlot + index,
			committeesPerSlot * params.ChainConfig.SlotsInEpoch,
		)
}


// Vault committee is a randomly selected committee of BPs that are chosen to generate the pool's keys via DKG
//
// Pool committee is chosen randomly by shuffling a seed + category (pool %d committee)
// The previous epoch's seed is used to choose the DKG committee as the current one (the block's epoch)
func GetVaultCommittee(state *core.State, poolId uint64, epoch uint64) ([]uint64,error) {
	// TODO - handle integer overflow
	seed, err := GetEpochSeed(state, epoch - 1) // we always use the seed from previous epoch
	if err != nil {
		return []uint64{}, err
	}

	vault, err := ComputeCommittee(
		GetActiveBlockProducers(state, epoch),
		seed,
		poolId,
		params.ChainConfig.VaultSize)

	//shuffled, err := shuffleActiveBPs(
	//	GetActiveBlockProducers(state, epoch),
	//	SliceToByte32(seed),
	//	[]byte(fmt.Sprintf("pool %d committee", poolId)),
	//)
	if err != nil {
		return nil, err
	}
	return vault, nil
}

/**
def get_attesting_indices(state: BeaconState,
                          data: AttestationData,
                          bits: Bitlist[MAX_VALIDATORS_PER_COMMITTEE]) -> Set[ValidatorIndex]:
    """
    Return the set of attesting indices corresponding to ``data`` and ``bits``.
    """
    committee = get_beacon_committee(state, data.slot, data.index)
    return set(index for i, index in enumerate(committee) if bits[i])
 */
func GetAttestingIndices(state *core.State, data *core.AttestationData, bits bitfield.Bitlist) ([]uint64, error) {
	committee, err := GetAttestationCommittee(state, data.Slot, data.CommitteeIndex)
	if err != nil {
		return nil, err
	}
	ret := []uint64{}
	for i := range bits {
		if bits.BitAt(uint64(i)) {
			ret = append(ret, committee[i])
		}
	}
	return ret, nil
}