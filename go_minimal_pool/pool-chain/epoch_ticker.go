package pool_chain

import (
	"time"
)

type EpochTicker struct {
	ticker *time.Ticker
	interval time.Duration
	number int
	tickerChan chan int
	done chan bool
}

func NewEpochTicker(interval time.Duration) *EpochTicker {
	return &EpochTicker{
		interval: interval,
		number:    0,
		tickerChan: make(chan int),
		done: make(chan bool),
	}
}

func (t *EpochTicker) Start () {
	t.ticker = time.NewTicker(t.interval)
	go func() {
		// send first epoch now
		t.tickerChan <- t.number
		t.number += 1
		for {
			select {
			case <-t.done:
				return
			case _ = <-t.ticker.C:
				t.tickerChan <- t.number
				t.number += 1
			}
		}
	}()

}

func (t *EpochTicker) Stop () {
	t.ticker.Stop()
	t.done <- true
}

func (t *EpochTicker) C() <-chan int  {
	return t.tickerChan
}