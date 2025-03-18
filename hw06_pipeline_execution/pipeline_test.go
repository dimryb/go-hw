package hw06pipelineexecution

import (
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const (
	sleepPerStage     = time.Millisecond * 100
	fastSleepPerStage = time.Microsecond * 100
	fault             = sleepPerStage / 2
)

var isFullTesting = true

func TestPipeline(t *testing.T) {
	// Stage generator
	g := func(_ string, f func(v interface{}) interface{}) Stage {
		return func(in In) Out {
			out := make(Bi)
			go func() {
				defer close(out)
				for v := range in {
					time.Sleep(sleepPerStage)
					out <- f(v)
				}
			}()
			return out
		}
	}

	stages := []Stage{
		g("Dummy", func(v interface{}) interface{} { return v }),
		g("Multiplier (* 2)", func(v interface{}) interface{} { return v.(int) * 2 }),
		g("Adder (+ 100)", func(v interface{}) interface{} { return v.(int) + 100 }),
		g("Stringifier", func(v interface{}) interface{} { return strconv.Itoa(v.(int)) }),
	}

	t.Run("simple case", func(t *testing.T) {
		in := make(Bi)
		data := []int{1, 2, 3, 4, 5}

		go func() {
			for _, v := range data {
				in <- v
			}
			close(in)
		}()

		result := make([]string, 0, 10)
		start := time.Now()
		for s := range ExecutePipeline(in, nil, stages...) {
			result = append(result, s.(string))
		}
		elapsed := time.Since(start)

		require.Equal(t, []string{"102", "104", "106", "108", "110"}, result)
		require.Less(t,
			int64(elapsed),
			// ~0.8s for processing 5 values in 4 stages (100ms every) concurrently
			int64(sleepPerStage)*int64(len(stages)+len(data)-1)+int64(fault))
	})

	t.Run("done case", func(t *testing.T) {
		in := make(Bi)
		done := make(Bi)
		data := []int{1, 2, 3, 4, 5}

		// Abort after 200ms
		abortDur := sleepPerStage * 2
		go func() {
			<-time.After(abortDur)
			close(done)
		}()

		go func() {
			for _, v := range data {
				in <- v
			}
			close(in)
		}()

		result := make([]string, 0, 10)
		start := time.Now()
		for s := range ExecutePipeline(in, done, stages...) {
			result = append(result, s.(string))
		}
		elapsed := time.Since(start)

		require.Len(t, result, 0)
		require.Less(t, int64(elapsed), int64(abortDur)+int64(fault))
	})

	t.Run("empty input", func(t *testing.T) {
		in := make(Bi)
		close(in)

		result := make([]string, 0, 10)
		for s := range ExecutePipeline(in, nil, stages...) {
			result = append(result, s.(string))
		}

		require.Len(t, result, 0)
	})

	t.Run("partial processing with done", func(t *testing.T) {
		in := make(Bi)
		done := make(Bi)
		data := []int{1, 2, 3, 4, 5}

		abortDur := sleepPerStage * 3
		go func() {
			<-time.After(abortDur)
			close(done)
		}()

		go func() {
			for _, v := range data {
				in <- v
			}
			close(in)
		}()

		result := make([]string, 0, 10)
		for s := range ExecutePipeline(in, done, stages...) {
			result = append(result, s.(string))
		}

		require.Equal(t, len(result), 0)
	})

	t.Run("cancel before processing", func(t *testing.T) {
		in := make(Bi)
		done := make(Bi)
		close(done)

		go func() {
			defer close(in)
			for i := 1; i <= 5; i++ {
				in <- i
			}
		}()

		result := make([]string, 0, 10)
		for s := range ExecutePipeline(in, done, stages...) {
			result = append(result, s.(string))
		}

		require.Len(t, result, 0)
	})
}

func TestAllStageStop(t *testing.T) {
	if !isFullTesting {
		return
	}
	wg := sync.WaitGroup{}
	// Stage generator
	g := func(_ string, f func(v interface{}) interface{}) Stage {
		return func(in In) Out {
			out := make(Bi)
			wg.Add(1)
			go func() {
				defer wg.Done()
				defer close(out)
				for v := range in {
					time.Sleep(sleepPerStage)
					out <- f(v)
				}
			}()
			return out
		}
	}

	stages := []Stage{
		g("Dummy", func(v interface{}) interface{} { return v }),
		g("Multiplier (* 2)", func(v interface{}) interface{} { return v.(int) * 2 }),
		g("Adder (+ 100)", func(v interface{}) interface{} { return v.(int) + 100 }),
		g("Stringifier", func(v interface{}) interface{} { return strconv.Itoa(v.(int)) }),
	}

	t.Run("done case", func(t *testing.T) {
		in := make(Bi)
		done := make(Bi)
		data := []int{1, 2, 3, 4, 5}

		// Abort after 200ms
		abortDur := sleepPerStage * 2
		go func() {
			<-time.After(abortDur)
			close(done)
		}()

		go func() {
			for _, v := range data {
				in <- v
			}
			close(in)
		}()

		result := make([]string, 0, 10)
		for s := range ExecutePipeline(in, done, stages...) {
			result = append(result, s.(string))
		}
		wg.Wait()

		require.Len(t, result, 0)
	})
}

func TestPipelineFastStage(t *testing.T) {
	// Stage generator
	g := func(_ string, f func(v interface{}) interface{}) Stage {
		return func(in In) Out {
			out := make(Bi)
			go func() {
				defer close(out)
				for v := range in {
					time.Sleep(fastSleepPerStage)
					out <- f(v)
				}
			}()
			return out
		}
	}

	// Data generator
	dg := func(ln int, offset int) []int {
		data := make([]int, ln)
		for i := range data {
			data[i] = i + 1 + offset
		}
		return data
	}

	stages := []Stage{
		g("Dummy", func(v interface{}) interface{} { return v }),
		g("Multiplier (* 2)", func(v interface{}) interface{} { return v.(int) * 2 }),
		g("Adder (+ 100)", func(v interface{}) interface{} { return v.(int) + 100 }),
		g("Stringifier", func(v interface{}) interface{} { return strconv.Itoa(v.(int)) }),
	}

	// Result generator
	rg := func(data []int) []string {
		expectResult := make([]string, len(data))
		for i := range expectResult {
			val := data[i]*2 + 100
			expectResult[i] = strconv.Itoa(val)
		}
		return expectResult
	}

	t.Run("large data set", func(t *testing.T) {
		in := make(Bi)
		data := dg(1000, 0)

		go func() {
			for _, v := range data {
				in <- v
			}
			close(in)
		}()

		result := make([]string, 0, len(data))
		start := time.Now()
		for s := range ExecutePipeline(in, nil, stages...) {
			result = append(result, s.(string))
		}
		elapsed := time.Since(start)

		require.Len(t, result, len(data))
		require.Less(t,
			int64(elapsed),
			int64(sleepPerStage)*int64(len(stages)+len(data)-1)+int64(fault))
		require.Equal(t, rg(data), result)
	})

	t.Run("no stages", func(t *testing.T) {
		in := make(Bi)

		data := dg(100, 0)

		go func() {
			for _, v := range data {
				in <- v
			}
			close(in)
		}()

		result := make([]int, 0, len(data))
		for v := range ExecutePipeline(in, nil) {
			result = append(result, v.(int))
		}

		require.Equal(t, data, result)
	})

	t.Run("multiple pipelines", func(t *testing.T) {
		in1 := make(Bi)
		in2 := make(Bi)
		data1 := dg(100, 0)
		data2 := dg(50, 1000)

		go func() {
			for _, v := range data1 {
				in1 <- v
			}
			close(in1)
		}()
		go func() {
			for _, v := range data2 {
				in2 <- v
			}
			close(in2)
		}()

		result1 := make([]string, 0, len(data1))
		result2 := make([]string, 0, len(data2))

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()
			for s := range ExecutePipeline(in1, nil, stages...) {
				result1 = append(result1, s.(string))
			}
		}()
		go func() {
			defer wg.Done()
			for s := range ExecutePipeline(in2, nil, stages...) {
				result2 = append(result2, s.(string))
			}
		}()

		wg.Wait()

		require.Equal(t, rg(data1), result1)
		require.Equal(t, rg(data2), result2)
	})
}
