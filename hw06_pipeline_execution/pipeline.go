package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(inPipeline In, done In, stages ...Stage) Out {
	outPipeline := make(Bi)

	inFirstStage, outLastStage := launchStages(stages)

	go redirectInput(inPipeline, inFirstStage, done)
	go redirectLastStageResult(outLastStage, outPipeline, done)

	return outPipeline
}

// Launches all the stages.
// Returns a bidirectional channel, used as input of the first stage
// and an output channel of the last stage.
func launchStages(stages []Stage) (Bi, Out) {
	inFirstStage := make(Bi)

	var inStage In = inFirstStage
	var outLastStage Out
	for _, stage := range stages {
		outLastStage = stage(inStage)
		inStage = outLastStage
	}
	return inFirstStage, outLastStage
}

// Consumes values from `source`, redirects them to `dest`,
// stops if `done` is closed,
// closes `dest` when done.
func redirectInput(source In, dest Bi, done In) {
	defer close(dest)

	for inValue := range source {
		select {
		case dest <- inValue:
		case <-done:
			return
		}
	}
}

// Redirects results to the pipeline's output channel,
// stops if `done` is closed,
// closes pipeline's out when done.
func redirectLastStageResult(outLastStage In, outPipeline Bi, done In) {
	defer close(outPipeline)

	for {
		select {
		case result, more := <-outLastStage:
			if more {
				select {
				case outPipeline <- result:
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
}
