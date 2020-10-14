package shared

import (
	"fmt"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/core"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/shared/params"
)

/**
	def is_active_validator(validator: Validator, epoch: Epoch) -> bool:
		"""
		Check if ``validator`` is active.
		"""
		return validator.activation_epoch <= epoch < validator.exit_epoch
 */
func IsActiveBP(bp *core.BlockProducer, epoch uint64) bool {
	return bp.ActivationEpoch <= epoch && epoch < bp.ExitEpoch
}

/**
	def is_eligible_for_activation_queue(validator: Validator) -> bool:
		"""
		Check if ``validator`` is eligible to be placed into the activation queue.
		"""
		return (
			validator.activation_eligibility_epoch == FAR_FUTURE_EPOCH
			and validator.effective_balance == MAX_EFFECTIVE_BALANCE
		)
*/
func IsEligibleForActivationQueue(bp *core.BlockProducer) bool {
	return bp.ActivationEligibilityEpoch == params.ChainConfig.FarFutureEpoch && bp.Stake == params.ChainConfig.MaxEffectiveBalance
}

/**
	def is_eligible_for_activation(state: BeaconState, validator: Validator) -> bool:
		"""
		Check if ``validator`` is eligible for activation.
		"""
		return (
			# Placement in queue is finalized
			validator.activation_eligibility_epoch <= state.finalized_checkpoint.epoch
			# Has not yet been activated
			and validator.activation_epoch == FAR_FUTURE_EPOCH
		)
 */
func IsEligibleForActivation(state *core.State, bp *core.BlockProducer) bool {
	return bp.ActivationEligibilityEpoch <= state.FinalizedCheckpoint.Epoch && // Placement in queue is finalized
					bp.ActivationEpoch == params.ChainConfig.FarFutureEpoch // Has not yet been activated
}

/**
	def is_slashable_validator(validator: Validator, epoch: Epoch) -> bool:
		"""
		Check if ``validator`` is slashable.
		"""
		return (not validator.slashed) and (validator.activation_epoch <= epoch < validator.withdrawable_epoch)
 */
func IsSlashableBp(bp *core.BlockProducer, epoch uint64) bool {
	return !bp.Slashed && (bp.ActivationEpoch <= epoch && epoch < bp.WithdrawableEpoch)
}

/**
def compute_proposer_index(state: BeaconState, indices: Sequence[ValidatorIndex], seed: Bytes32) -> ValidatorIndex:
    """
    Return from ``indices`` a random index sampled by effective balance.
    """
    assert len(indices) > 0
    MAX_RANDOM_BYTE = 2**8 - 1
    i = uint64(0)
    total = uint64(len(indices))
    while True:
        candidate_index = indices[compute_shuffled_index(i % total, total, seed)]
        random_byte = hash(seed + uint_to_bytes(uint64(i // 32)))[i % 32]
        effective_balance = state.validators[candidate_index].effective_balance
        if effective_balance * MAX_RANDOM_BYTE >= MAX_EFFECTIVE_BALANCE * random_byte:
            return candidate_index
        i += 1
 */
func ComputeProposerIndex(state *core.State, indices []uint64, seed []byte) (uint64, error) {
	if len(indices) == 0 {
		return 0, fmt.Errorf("couldn't compute proposer, indices list empty")
	}
	maxRandomByte := uint64(2^8-1)
	i := uint64(0)
	total := uint64(len(indices))
	for {
		idx, err := computeShuffledIndex(i % total, total, SliceToByte32(seed), true,10) // TODO - shuffle round via config
		if err != nil {
			return 0, err
		}

		candidateIndex := indices[idx]
		b := append(seed[:], Bytes8(i / 32)...)
		randomByte := Hash(b)[i%32]

		bp := GetBlockProducer(state, candidateIndex)
		if bp == nil {
			return 0, fmt.Errorf("could not find shuffled BP index %d", candidateIndex)
		}
		stake := bp.Stake

		if stake * maxRandomByte >= params.ChainConfig.MaxEffectiveBalance * uint64(randomByte) {
			return candidateIndex, nil
		}
	}
}

/**
def compute_activation_exit_epoch(epoch: Epoch) -> Epoch:
    """
    Return the epoch during which validator activations and exits initiated in ``epoch`` take effect.
    """
    return Epoch(epoch + 1 + MAX_SEED_LOOKAHEAD)
 */
func ComputeActivationExitEpoch(epoch uint64) uint64 {
	return epoch + 1 + params.ChainConfig.MaxSeedLookahead
}

/**
def get_active_validator_indices(state: BeaconState, epoch: Epoch) -> Sequence[ValidatorIndex]:
    """
    Return the sequence of active validator indices at ``epoch``.
    """
    return [ValidatorIndex(i) for i, v in enumerate(state.validators) if is_active_validator(v, epoch)]
 */
func GetActiveBlockProducers(state *core.State, epoch uint64) []uint64 {
	var activeBps []uint64
	for _, bp := range state.BlockProducers {
		if IsActiveBP(bp, epoch) {
			activeBps = append(activeBps, bp.GetId())
		}
	}
	return activeBps
}

/**
def get_validator_churn_limit(state: BeaconState) -> uint64:
    """
    Return the validator churn limit for the current epoch.
    """
    active_validator_indices = get_active_validator_indices(state, get_current_epoch(state))
    return max(MIN_PER_EPOCH_CHURN_LIMIT, uint64(len(active_validator_indices)) // CHURN_LIMIT_QUOTIENT)
 */
func GetValidatorChurnLimit(state *core.State) uint64 {
	activeBPs := GetActiveBlockProducers(state, GetCurrentEpoch(state))
	churLimit := uint64(len(activeBPs)) / params.ChainConfig.ChurnLimitQuotient
	if churLimit < params.ChainConfig.MinPerEpochChurnLimit {
		churLimit = params.ChainConfig.MinPerEpochChurnLimit
	}
	return churLimit
}

/**
def get_beacon_proposer_index(state: BeaconState) -> ValidatorIndex:
    """
    Return the beacon proposer index at the current slot.
    """
    epoch = get_current_epoch(state)
    seed = hash(get_seed(state, epoch, DOMAIN_BEACON_PROPOSER) + uint_to_bytes(state.slot))
    indices = get_active_validator_indices(state, epoch)
    return compute_proposer_index(state, indices, seed)
 */
func GetBlockProposerIndex(state *core.State) (uint64, error) {
	epoch := GetCurrentEpoch(state)
	seed := GetSeed(state, epoch, params.ChainConfig.DomainBeaconProposer)
	SeedWithSlot := append(seed[:], Bytes8(state.CurrentSlot)...)
	hash := Hash(SeedWithSlot)

	bps := GetActiveBlockProducers(state, epoch)
	return ComputeProposerIndex(state, bps, hash[:])
}