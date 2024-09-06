package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

type workerPool struct {
	workerCount int

	tasksCh  chan Task
	cancelCh chan struct{}
	errsCh   chan struct{}

	success bool

	wg sync.WaitGroup
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	return newPool(n, m).Execute(tasks)
}

func newPool(n, m int) *workerPool {
	// если кол-во ошибок меньше либо равно 1 то выходим при любой ошибке
	if m < 1 {
		m = 1
	}
	return &workerPool{
		workerCount: n,
		cancelCh:    make(chan struct{}),
		tasksCh:     make(chan Task, n),
		errsCh:      make(chan struct{}, m-1),
		success:     true,
	}
}

func (s *workerPool) Execute(tasks []Task) error {
	s.runConsumers()

	s.runProducer(tasks)

	if !s.wait() {
		return ErrErrorsLimitExceeded
	}
	return nil
}

func (s *workerPool) runProducer(tasks []Task) {
	go func() {
		defer s.wg.Done()
		defer close(s.tasksCh)

		// читаем задачи и если сигнала отмены нет - пишем в канал задач
		for _, task := range tasks {
			select {
			case <-s.cancelCh:
				return
			case s.tasksCh <- task:
			}
		}
	}()
}

func (s *workerPool) runConsumers() {
	s.wg.Add(s.workerCount + 1)
	// создаем N воркеров
	for i := 1; i <= s.workerCount; i++ {
		go s.worker()
	}
}

func (s *workerPool) wait() bool {
	defer close(s.errsCh)
	// ждем завершения всех воркеров
	s.wg.Wait()

	return s.success
}

func (s *workerPool) worker() {
	defer s.wg.Done()

	// читаем канал задач пока канал открыт
	// когда у продюсера закончатся задачи - он закроет канал tasksCh и цикл завершится
	for task := range s.tasksCh {
		select {
		// перед выполнением проверяем - не превышено ли число ошибок другими воркерами
		case <-s.cancelCh:
			return
		default:
			err := task()
			if err != nil {
				select {
				case <-s.cancelCh:
					return
				// пишем в канал ошибок если лимит ошибок не превышен
				case s.errsCh <- struct{}{}:
				// сюда попадаем если канал ошибок заполнен, т.е. в первый раз превышен лимит ошибок
				default:
					// отменяем выполнение
					s.success = false
					close(s.cancelCh)
					return
				}
			}
		}
	}
}
