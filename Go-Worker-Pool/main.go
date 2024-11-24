package main

import (
	"fmt"
	"sync"
	"time"
)

type Task struct {
	Id int
}

func (obj *Task) Process() {
	fmt.Println("Processing Task with ID :", obj.Id)
	time.Sleep(2 * time.Second)
}

type WorkerPool struct {
	Tasks       []Task
	concurrency int
	taskChan    chan Task
	wg          sync.WaitGroup
}

func (obj *WorkerPool) worker() {
	for task := range obj.taskChan {
		task.Process()
		obj.wg.Done()
	}
}

func (obj *WorkerPool) Run() {

	obj.taskChan = make(chan Task, len(obj.Tasks))
	for i := 0; i < obj.concurrency; i++ {
		go obj.worker()
	}

	obj.wg.Add(len(obj.Tasks))
	for _, task := range obj.Tasks {
		obj.taskChan <- task
	}

	close(obj.taskChan)
	obj.wg.Wait()

}

func main() {
	tasks := make([]Task, 20)
	for i := 0; i < 20; i++ {
		tasks[i] = Task{Id: i + 1}
	}

	workerPool := WorkerPool{
		Tasks:       tasks,
		concurrency: 3,
	}
	workerPool.Run()
}
