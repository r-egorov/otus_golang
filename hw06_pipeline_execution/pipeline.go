package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	inFirstStage := redirectInputToFirstStage(in, done)
	outLastStage := launchStages(inFirstStage, stages)
	outPipeline := redirectOutputFromLastStage(outLastStage, done)

	return outPipeline
}

// Creates a `dest` channel, where redirects values from `source` to,
// stops if `done` is closed, closes `dest` when done.
func redirectInputToFirstStage(source In, done In) In {
	dest := make(Bi)

	go func() {
		defer close(dest)
		for inValue := range source {
			select {
			case dest <- inValue:
			case <-done:
				return
			}
		}
	}()

	return dest
}

// Launches all the stages,
// takes the input channel of the first stage as a parameter,
// returns the output channel of the last stage.
func launchStages(inStage In, stages []Stage) Out {
	var outLastStage Out
	for _, stage := range stages {
		outLastStage = stage(inStage)
		inStage = outLastStage
	}
	return outLastStage
}

// Redirects results to the pipeline's output channel,
// stops if `done` is closed,
// closes pipeline's out when done.
func redirectOutputFromLastStage(outLastStage In, done In) Out {
	out := make(Bi)

	go func() {
		defer close(out)
		for {
			select {
			case result, more := <-outLastStage:
				if !more {
					return
				}
				select {
				case out <- result:
				case <-done:
					return
				}
			case <-done:
				return
			}
		}
	}()

	return out
}
