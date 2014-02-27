package glowflake

import (
	"fmt"
	"sync"
	"time"
)

const (
	WorkerIdBits     = 5
	DatacenterIdBits = 5
	SequenceBits     = 12
	MaxWorkerId      = -1 ^ (-1 << WorkerIdBits)
	MaxDatacenterId  = -1 ^ (-1 << DatacenterIdBits)
	MaxSequence      = -1 ^ (-1 << SequenceBits)
)

var (
	twepoch int64 = time.Date(2013, time.September, 4, 0, 30, 0, 0, time.UTC).UnixNano() / 1e6
)

type GlowFlake struct {
	workerId      int8
	datacenterId  int8
	sequence      int16
	lastTimestamp int64
	lock          sync.Mutex
}

func (gf *GlowFlake) NextId() (int64, error) {
	gf.lock.Lock()
	defer gf.lock.Unlock()

	ts := timeGen()

	if ts < gf.lastTimestamp {
		return 0, fmt.Errorf("Invalid timestamp: %v - precedes %v", ts, gf)
	}

	if ts == gf.lastTimestamp {
		gf.sequence = (gf.sequence + 1) & MaxSequence
		if gf.sequence == 0 {
			ts = tilNextMillis(ts)
		}
	} else {
		gf.sequence = 0
	}

	gf.lastTimestamp = ts
	return (gf.lastTimestamp-twepoch)<<22 |
		int64(gf.datacenterId)<<17 |
		int64(gf.workerId)<<12 |
		int64(gf.sequence), nil
}

func NewGlowFlake(workerId int8, datacenterId int8) (*GlowFlake, error) {
	if workerId < 0 || workerId > MaxWorkerId {
		return nil, fmt.Errorf("Worker id %v is invalid", workerId)
	}
	if datacenterId < 0 || datacenterId > MaxDatacenterId {
		return nil, fmt.Errorf("DatacenterId id %v is invalid", datacenterId)
	}
	return &GlowFlake{workerId: workerId, datacenterId: datacenterId}, nil
}

func timeGen() int64 {
	return time.Now().UnixNano() / 1e6
}

func tilNextMillis(lastTimestamp int64) int64 {
	i := timeGen()
	for i <= lastTimestamp {
		i = timeGen()
	}
	return i
}
