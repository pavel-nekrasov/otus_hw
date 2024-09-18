package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	ch := createInputChannel(in, done)
	for index, stage := range stages {
		var inputCh In
		if index > 0 {
			inputCh = createInputChannel(ch, done)
		} else {
			inputCh = ch
		}
		ch = stage(inputCh)
	}
	return ch
}

func createInputChannel(in In, done In) Out {
	out := make(Bi)

	// тут тонкость в том, что stage не проверяет внутри себя done поэтому после закрытия done,
	// если сразу выйти из горутины
	// горутина объявленная в предыдущем stage может зависнуть пытаясь писать в канал который никто не читает.
	// Поэтому при выходе через defer закрываем out,
	// но продолжаем читать из in пока он не закроется и просто отбрасываем приходящие оттуда данные
	go func() {
		defer func() {
			close(out)
			//revive:disable:empty-block
			for range in {
			}
			//revive:enable:empty-block
		}()
		for {
			select {
			case v, ok := <-in:
				if !ok {
					return
				}
				out <- v
			case <-done:
				return
			}
		}
	}()

	return out
}
