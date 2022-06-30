package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := in

	for _, stage := range stages {
		out = Execute(done, stage(out))
	}

	return out
}

func Execute(done, in In) <-chan interface{} {
	outputStream := make(chan interface{})
	go func() {
		defer close(outputStream)

		for {
			select {
			case <-done:
				return
			case value, open := <-in:
				if !open {
					return
				}
				outputStream <- value
			}
		}
	}()

	return outputStream
}
