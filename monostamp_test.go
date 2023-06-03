package monostamp_test

import (
	"sort"
	"sync"
	"testing"
	"time"

	"github.com/Flamefork/monostamp"
	"github.com/stretchr/testify/assert"
)

func TestNextIfClockIsIncreasing(t *testing.T) {
	clValues := []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	tsValues := []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	m := monostamp.New(clockWithValues(clValues), 0, nil)

	for i := 0; i < 10; i++ {
		assert.Equal(t, tsValues[i], m.Next())
	}
}

func TestNextIfClockIsStuck(t *testing.T) {
	clValues := []int64{1, 1, 1, 1, 1, 1, 1, 1, 1, 1}
	tsValues := []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	m := monostamp.New(clockWithValues(clValues), 0, nil)

	for i := 0; i < 10; i++ {
		assert.Equal(t, tsValues[i], m.Next())
	}
}

func TestNextIfClockIsRecovered(t *testing.T) {
	clValues := []int64{1, 1, 1, 1, 1, 1, 1, 9, 9, 9}
	tsValues := []int64{1, 2, 3, 4, 5, 6, 7, 9, 10, 11}
	m := monostamp.New(clockWithValues(clValues), 0, nil)

	for i := 0; i < 10; i++ {
		assert.Equal(t, tsValues[i], m.Next())
	}
}

func TestNextWithStartTS(t *testing.T) {
	clValues := []int64{1, 1, 1, 1, 1, 19, 19, 19, 19, 19}
	tsValues := []int64{3, 4, 5, 6, 7, 19, 20, 21, 22, 23}
	m := monostamp.New(clockWithValues(clValues), 2, nil)

	for i := 0; i < 10; i++ {
		assert.Equal(t, tsValues[i], m.Next())
	}
}

func TestWithProvidedClock(t *testing.T) {
	{
		m := monostamp.New(monostamp.Unix, 0, nil)
		assert.InEpsilon(t, time.Now().Unix(), m.Next(), 0.001)
	}
	{
		m := monostamp.New(monostamp.UnixMilli, 0, nil)
		assert.InEpsilon(t, time.Now().UnixMilli(), m.Next(), 0.001)
	}
	{
		m := monostamp.New(monostamp.UnixMicro, 0, nil)
		assert.InEpsilon(t, time.Now().UnixMicro(), m.Next(), 0.001)
	}
	{
		m := monostamp.New(monostamp.UnixNano, 0, nil)
		assert.InEpsilon(t, time.Now().UnixNano(), m.Next(), 0.001)
	}
}

func TestNextFromGoroutines(t *testing.T) {
	clock := func() int64 { return 1 }
	m := monostamp.New(clock, 0, nil)
	tsValues := make([]int64, 1000)
	wg := sync.WaitGroup{}
	for i := 0; i < 1000; i++ {
		i := i
		wg.Add(1)
		go func() {
			tsValues[i] = m.Next()
			wg.Done()
		}()
	}
	wg.Wait()
	sort.Slice(tsValues, func(i, j int) bool { return tsValues[i] < tsValues[j] })
	for i := 0; i < 1000; i++ {
		assert.Equal(t, int64(i+1), tsValues[i])
	}
}

func TestOnDrift(t *testing.T) {
	clValues := []int64{1, 1, 3, 3, 5, 5, 5, 8, 9, 10}
	tsValues := []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	drValues := []int64{0, 1, 0, 1, 0, 1, 2, 0, 0, 0}

	var drift int64
	onDrift := func(genTS int64, clockTS int64) {
		drift = genTS - clockTS
	}
	m := monostamp.New(clockWithValues(clValues), 0, onDrift)

	for i := 0; i < 10; i++ {
		drift = 0
		assert.Equal(t, tsValues[i], m.Next())
		assert.Equal(t, drValues[i], drift)
	}
}

func TestDriftReporterBasics(t *testing.T) {
	var drift int
	onDrift := func(genTS int64, clockTS int64) {
		drift = int(genTS - clockTS)
	}
	dr := monostamp.NewDriftReporter(1, 1, onDrift)

	drift = -1
	dr.Report(1, 1)
	assert.Equal(t, -1, drift)

	drift = -1
	dr.Report(5, 1)
	assert.Equal(t, 4, drift)

	drift = -1
	dr.Report(6, 6)
	assert.Equal(t, -1, drift)
}

func TestDriftReporterThreshold(t *testing.T) {
	var drift int
	onDrift := func(genTS int64, clockTS int64) {
		drift = int(genTS - clockTS)
	}
	dr := monostamp.NewDriftReporter(2, 1, onDrift)

	drift = -1
	dr.Report(1, 1)
	assert.Equal(t, -1, drift)

	drift = -1
	dr.Report(5, 1)
	assert.Equal(t, 4, drift)

	drift = -1
	dr.Report(7, 6)
	assert.Equal(t, -1, drift)
}

func TestDriftReporterInterval(t *testing.T) {
	var drift int
	onDrift := func(genTS int64, clockTS int64) {
		drift = int(genTS - clockTS)
	}
	dr := monostamp.NewDriftReporter(1, 2, onDrift)

	drift = -1
	dr.Report(10, 10)
	assert.Equal(t, -1, drift)

	drift = -1
	dr.Report(15, 11)
	assert.Equal(t, 4, drift)

	drift = -1
	dr.Report(16, 12)
	assert.Equal(t, -1, drift)

	drift = -1
	dr.Report(17, 13)
	assert.Equal(t, 4, drift)
}

func clockWithValues(values []int64) func() int64 {
	i := -1
	return func() int64 {
		i++
		return values[i]
	}
}
