package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

type pipeline struct {
	stages []Stage
	in     In
	done   In
}

type outputAdapter struct {
	out    Bi
	closed bool
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	pipeLine := newPipeline(in, done, stages)
	return pipeLine.Run(in, done)
}

func newPipeline(in In, done In, stages []Stage) *pipeline {
	return &pipeline{
		in:     in,
		done:   done,
		stages: stages,
	}
}

func (p *pipeline) Run(in In, done In) Out {
	ch := in
	for _, stage := range p.stages {
		ch = stage(p.multiplex(ch, done))
	}

	return ch
}

func (p *pipeline) multiplex(in In, done In) Out {
	s := outputAdapter{
		out: make(Bi),
	}

	go func() {
		defer s.close()
		for {
			select {
			case v, ok := <-in:
				if !ok {
					return
				}
				s.write(v)
			case <-done:
				s.close()
				continue
			}
		}
	}()

	return s.out
}

func (s *outputAdapter) write(v interface{}) {
	if !s.closed {
		s.out <- v
	}
}

func (s *outputAdapter) close() {
	if !s.closed {
		s.closed = true
		close(s.out)
	}
}
