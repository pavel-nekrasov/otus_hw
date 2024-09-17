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
	return pipeLine.Run()
}

func newPipeline(in In, done In, stages []Stage) *pipeline {
	return &pipeline{
		in:     in,
		done:   done,
		stages: stages,
	}
}

func (p *pipeline) Run() Out {
	ch := p.in
	for _, stage := range p.stages {
		ch = stage(p.multiplex(ch))
	}

	return ch
}

func (p *pipeline) multiplex(in In) Out {
	s := outputAdapter{
		out: make(Bi),
	}

	// тут тонкость в том, что stage не проверяет внутри себя done поэтому после закрытия done,
	// если сразу выйти из горутины
	// горутина объявленная в предыдущем stage может зависнуть пытаясь писать в канал который никто не читает.
	// Поэтому при закрытии done закрываем out, выставляем флаг, но из цикла не выходим до тех пор,
	// пока предыдущий stage не закроет свой канал записи. Фактическив момент закрытия done
	// мы продолжвем читать из in пока он не закроется и просто отбрасываем приходящие оттуда данные
	go func() {
		defer s.close()
		for {
			select {
			case v, ok := <-in:
				if !ok {
					return
				}
				s.write(v)
			case <-p.done:
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
