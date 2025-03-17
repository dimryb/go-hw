package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func runStage(in In, done In, stage Stage) Out {
	out := make(Bi)

	go func() {
		defer close(out)
		stageOut := stage(in)

		for {
			select {
			case <-done:
				return
			case v, ok := <-stageOut:
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
