package core

type Config struct {
	PoolExecutorsNumber 			uint64
	DKGParticipantsNumber 			uint64
	BaseEth2DutyReward  			uint64
	DKGReward 						uint64
}

func TestConfig() *Config {
	return &Config{
		PoolExecutorsNumber: 		128,
		DKGParticipantsNumber: 		128,
		BaseEth2DutyReward:  		100,
		DKGReward:					1000,
	}
}
