package spec

type PoolExecutionSummary struct {
	Id 					[]byte // pubkey
	StartingEpoch 		uint64 // a.k.a previous epoch
	EndEpoch 			uint64 //
	EndBalance			uint64 // taken from the BeaconChain
	Performance			Performance
}
