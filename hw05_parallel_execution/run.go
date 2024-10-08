package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

type workerPool struct {
	workerCount int
	errLimit    int
	errCount    int

	tasksCh  chan Task
	cancelCh chan struct{}

	wg  sync.WaitGroup
	mut sync.Mutex
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
		errLimit:    m,
		cancelCh:    make(chan struct{}),
		tasksCh:     make(chan Task),
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
	s.wg.Add(1)
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
	s.wg.Add(s.workerCount)
	// создаем N воркеров
	for i := 1; i <= s.workerCount; i++ {
		go s.worker()
	}
}

func (s *workerPool) wait() bool {
	// ждем завершения всех воркеров
	s.wg.Wait()

	return s.errCount < s.errLimit
}

func (s *workerPool) worker() {
	defer s.wg.Done()

	// читаем канал задач пока канал открыт
	// когда у продюсера закончатся задачи - он закроет канал tasksCh и цикл завершится
	for task := range s.tasksCh {
		if task() != nil {
			s.notifyError()
		}
	}
}

func (s *workerPool) notifyError() {
	s.mut.Lock()
	defer s.mut.Unlock()

	if s.errCount >= s.errLimit {
		return
	}

	s.errCount++

	if s.errCount >= s.errLimit {
		close(s.cancelCh)
	}
}
