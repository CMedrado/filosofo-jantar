package main

import (
	"fmt"
	"sync"
	"time"
)

type Filosofo struct {
	id       int
	garfoEsq *sync.Mutex
	garfoDir *sync.Mutex
	monitor  *sync.Mutex
	pensando bool
	done     chan bool
}

func NovoFilosofo(id int, garfoEsq, garfoDir *sync.Mutex, monitor *sync.Mutex, done chan bool) *Filosofo {
	return &Filosofo{
		id:       id,
		garfoEsq: garfoEsq,
		garfoDir: garfoDir,
		monitor:  monitor,
		pensando: true,
		done:     done,
	}
}

func (f *Filosofo) pensar() {
	fmt.Printf("Filósofo %d está pensando.\n", f.id)
	time.Sleep(time.Duration(f.id+1) * time.Second)
}

func (f *Filosofo) comer() {
	f.monitor.Lock()

	if f.pensando {
		f.garfoEsq.Lock()
		f.garfoDir.Lock()

		fmt.Printf("Filósofo %d está comendo.\n", f.id)
		time.Sleep(time.Duration(f.id+1) * time.Second)

		f.garfoEsq.Unlock()
		f.garfoDir.Unlock()

		fmt.Printf("Filósofo %d terminou de comer.\n", f.id)
		f.pensando = true
	}

	f.monitor.Unlock()
}

func (f *Filosofo) iniciar() {
	for {
		select {
		case <-f.done:
			return
		default:
			f.pensar()
			f.comer()
		}
	}
}

func main() {
	var wg sync.WaitGroup

	garfos := make([]*sync.Mutex, 5)
	monitor := &sync.Mutex{}
	done := make(chan bool)

	for i := 0; i < 5; i++ {
		garfos[i] = &sync.Mutex{}
	}

	filosofos := make([]*Filosofo, 5)
	for i := 0; i < 5; i++ {
		filosofos[i] = NovoFilosofo(i+1, garfos[i], garfos[(i+1)%5], monitor, done)
	}

	wg.Add(5)
	for i := 0; i < 5; i++ {
		go func(f *Filosofo) {
			defer wg.Done()
			f.iniciar()
		}(filosofos[i])
	}

	time.Sleep(10 * time.Second)
	close(done)

	wg.Wait()
}
