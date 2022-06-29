package hw06pipelineexecution

import "sync"

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := in

	wg := sync.WaitGroup{}
	wg.Add(len(stages))

	for _, stage := range stages {
		stage := stage
		out = stage.Execute(done, out)
	}

	go func() {
		wg.Wait()
	}()

	return out
}

func (stage Stage) Execute(done, in In) <-chan interface{} {
	outputStream := make(chan interface{})
	go func() {
		defer close(outputStream)

		for s := range stage(in) {
			select {
			case <-done:
				return
			default:
				outputStream <- s
			}
		}
	}()

	return outputStream
}
