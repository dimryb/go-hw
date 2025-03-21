package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func runStage(in In, done In, stage Stage) Out {
	out := make(Bi)
	buf := make(Bi)

	stageOut := stage(in)

	go func() {
		defer close(buf)

		for v := range stageOut {
			select {
			case buf <- v:
			case <-done:
				for val := range stageOut {
					_ = val // suppress linter
				}
				return
			}
		}
	}()

	go func() {
		defer close(out)

		for {
			select {
			case <-done:
				return
			case v, ok := <-buf:
				if !ok {
					return
				}
				out <- v
			}
		}
	}()
	return out
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	dataChan := in
	for _, stage := range stages {
		dataChan = runStage(dataChan, done, stage)
	}
	return dataChan
}
