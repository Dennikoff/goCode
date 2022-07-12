package main

import (
	"fmt"
	"sync"
)

type Pool struct {
	mu    sync.Mutex
	size  int
	tasks chan Task
	kill  chan struct{}
	wg    sync.WaitGroup
}

type Task struct {
	text string
}

func (t Task) Execute() {
	fmt.Println(t.text)
}

func NewPool(size int) (pool *Pool) {
	pool = &Pool{
		tasks: make(chan Task, 128),
		kill:  make(chan struct{}),
	}
	pool.Resize(size)
	return pool
}

func (p *Pool) Resize(n int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	for p.size < n {
		p.size++
		p.wg.Add(1)
		go p.work()
	}
	for p.size > n {
		p.size--
		p.kill <- struct{}{}
	}
}

func (p *Pool) work() {
	defer p.wg.Done()
	for {
		select {
		case task, ok := <-p.tasks:
			if !ok {
				return
			}
			task.Execute()
		case <-p.kill:
			return
		}
	}
}

func (p *Pool) Close() {
	close(p.tasks)
}

func (p *Pool) Exec(text []string) {
	for i := range text {
		task := Task{text[i]}
		p.tasks <- task
	}
}

func (p *Pool) Wait() {
	p.wg.Wait()
}

func main() {
	pool := NewPool(6)
	pool.Exec([]string{"kek", "lol1", "lol2", "lol3", "lol4", "lol5", "lol6", "lol7", "lol8"})
	pool.Close()
	pool.Wait()
}
