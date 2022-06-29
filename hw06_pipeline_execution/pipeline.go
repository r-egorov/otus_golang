package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := make(Bi)

	// Make an input channel for the first stage
	// so that we can close it if need to
	inFirstStage := make(Bi)

	// Startup the stages,
	// save input channel of the first stage
	// save output channel of the last stage (going to be returned)
	var inNextStage In = inFirstStage
	var outLastStage Out
	for _, stage := range stages {
		outLastStage = stage(inNextStage)
		inNextStage = outLastStage
	}

	// start sending values to the pipeline
	go func() {
		defer close(inFirstStage)

		for inValue := range in {
			select {
			case inFirstStage <- inValue:
			case <-done:
				return
			}
		}
	}()

	// start sending results to the output channel
	// break if done is closed
	go func() {
		defer close(out)

		for {
			select {
			case result, more := <-outLastStage:
				if more {
					select {
					case out <- result:
					case <-done:
						return
					}
				} else {
					return
				}
			case <-done:
				return
			}
		}
	}()

	return out
}
