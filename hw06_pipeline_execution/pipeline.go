package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	_ = done
	dataChan := in
	for _, stage := range stages {
		dataChan = stage(dataChan)
	}
	return dataChan
}
