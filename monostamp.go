package monostamp

import (
	"sync/atomic"
	"time"
)

type Monostamp struct {
	clock   func() int64
	lastTS  atomic.Int64
	onDrift func(int64, int64)
}

func New(clock func() int64, startTS int64, onDrift func(int64, int64)) *Monostamp {
	m := Monostamp{
		clock:   clock,
		lastTS:  atomic.Int64{},
		onDrift: onDrift,
	}
	m.lastTS.Store(startTS)
	return &m
}

func (m *Monostamp) Next() int64 {
	for {
		lastTS := m.lastTS.Load()
		ts := m.clock()
		if ts <= lastTS {
			if m.onDrift != nil {
				m.onDrift(lastTS+1, ts)
			}
			ts = lastTS + 1
		}
		if m.lastTS.CompareAndSwap(lastTS, ts) {
			return ts
		}
	}
}

func UnixNano() int64 {
	return time.Now().UnixNano()
}

func UnixMicro() int64 {
	return time.Now().UnixMicro()
}

func UnixMilli() int64 {
	return time.Now().UnixMilli()
}

func Unix() int64 {
	return time.Now().Unix()
}

type DriftReporter struct {
	threshold int64
	interval  int64
	callback  func(int64, int64)
	lastRep   atomic.Int64
}

func NewDriftReporter(
	threshold int64, interval int64, callback func(int64, int64),
) *DriftReporter {
	return &DriftReporter{
		threshold: threshold,
		interval:  interval,
		callback:  callback,
		lastRep:   atomic.Int64{},
	}
}

func (d *DriftReporter) Report(generatedTS int64, clockTS int64) {
	if generatedTS-clockTS < d.threshold {
		return
	}
	lastRepVal := d.lastRep.Load()
	if clockTS < lastRepVal+d.interval {
		return
	}
	if !d.lastRep.CompareAndSwap(lastRepVal, clockTS) {
		return
	}
	d.callback(generatedTS, clockTS)
}
