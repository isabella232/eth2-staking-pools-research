package shared

import (
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/core"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/shared/params"
)

/**
def get_total_balance(state: BeaconState, indices: Set[ValidatorIndex]) -> Gwei:
    """
    Return the combined effective balance of the ``indices``.
    ``EFFECTIVE_BALANCE_INCREMENT`` Gwei minimum to avoid divisions by zero.
    Math safe up to ~10B ETH, afterwhich this overflows uint64.
    """
    return Gwei(max(EFFECTIVE_BALANCE_INCREMENT, sum([state.validators[index].effective_balance for index in indices])))
 */
func GetTotalStake(state *core.State, indices []uint64) uint64 {
	sum := uint64(0)
	for _, index := range indices {
		bp := GetBlockProducer(state, index)
		if bp != nil {
			sum += bp.Stake
		}
	}

	if sum < params.ChainConfig.EffectiveBalanceIncrement {
		return params.ChainConfig.EffectiveBalanceIncrement
	}
	return sum
}

/**
def get_total_active_balance(state: BeaconState) -> Gwei:
    """
    Return the combined effective balance of the active validators.
    Note: ``get_total_balance`` returns ``EFFECTIVE_BALANCE_INCREMENT`` Gwei minimum to avoid divisions by zero.
    """
    return get_total_balance(state, set(get_active_validator_indices(state, get_current_epoch(state))))
 */
func GetTotalActiveStake(state *core.State) uint64 {
	indices := GetActiveBlockProducers(state, GetCurrentEpoch(state))
	return GetTotalStake(state, indices)
}