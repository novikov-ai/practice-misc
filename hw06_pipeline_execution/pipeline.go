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
		stage := stage

		select {
		case <-done:
			return nil
		default:
			out = executor(done, stage(out))
		}
	}

	return out
}

func executor(done, input In) Out {
	outStream := make(chan interface{})

	go func() {
		defer close(outStream)

		for value := range input {
			select {
			case <-done:
				return
			case outStream <- value:
			}
		}
	}()

	return outStream
}
