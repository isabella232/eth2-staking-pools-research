package src

type Config struct {
	PoolExecutorsNumber uint64
	BaseEth2DutyReward  uint64
}

func TestConfig() *Config {
	return &Config{
		PoolExecutorsNumber: 128,
		BaseEth2DutyReward:  100,
	}
}
