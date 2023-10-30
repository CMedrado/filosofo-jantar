package main

import (
	"fmt"
	"sync"
	"time"
)

// Definindo a estrutura do Filósofo
type Filosofo struct {
	id       int         // Identificador único do filósofo
	garfoEsq *sync.Mutex // Garfo à esquerda do filósofo (Mutex para garantir exclusão mútua)
	garfoDir *sync.Mutex // Garfo à direita do filósofo (Mutex para garantir exclusão mútua)
	monitor  *sync.Mutex // Mutex para controle da região crítica onde o filósofo decide comer ou pensar
	pensando bool        // Indica se o filósofo está pensando (true) ou comendo (false)
	done     chan bool   // Canal para sinalizar quando o filósofo deve parar sua execução
}

// Função para criar um novo Filósofo com os parâmetros fornecidos
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

// Função para simular o estado de pensamento do Filósofo
func (f *Filosofo) pensar() {
	fmt.Printf("Filósofo %d está pensando.\n", f.id)
	time.Sleep(time.Duration(f.id+1) * time.Second) // Simula o tempo que o filósofo passa pensando
}

// Função para simular o estado de comer do Filósofo
func (f *Filosofo) comer() {
	f.monitor.Lock() // Entra na região crítica para decidir se deve comer

	if f.pensando {
		// Tenta pegar os garfos (Mutexes) para comer
		f.garfoEsq.Lock()
		f.garfoDir.Lock()

		// O Filósofo está comendo
		fmt.Printf("Filósofo %d está comendo.\n", f.id)
		time.Sleep(time.Duration(f.id+1) * time.Second) // Simula o tempo que o filósofo passa comendo

		// Libera os garfos após comer
		f.garfoEsq.Unlock()
		f.garfoDir.Unlock()

		// O Filósofo terminou de comer e volta a pensar
		fmt.Printf("Filósofo %d terminou de comer.\n", f.id)
		f.pensando = true
	}

	f.monitor.Unlock() // Sai da região crítica
}

// Função para iniciar a rotina do Filósofo
func (f *Filosofo) iniciar() {
	for {
		select {
		case <-f.done: // Se recebe um sinal para parar, a goroutine termina
			return
		default:
			f.pensar() // O Filósofo pensa
			f.comer()  // O Filósofo tenta comer
		}
	}
}

func main() {
	var wg sync.WaitGroup

	// Inicializa os garfos e o monitor como Mutexes
	garfos := make([]*sync.Mutex, 5)
	monitor := &sync.Mutex{}
	done := make(chan bool) // Canal para sinalizar quando os filósofos devem parar

	// Inicializa os garfos como Mutexes
	for i := 0; i < 5; i++ {
		garfos[i] = &sync.Mutex{}
	}

	// Inicializa os filósofos
	filosofos := make([]*Filosofo, 5)
	for i := 0; i < 5; i++ {
		// Cada filósofo é inicializado com seu identificador, os Mutexes dos garfos à esquerda e à direita,
		// o monitor (Mutex para a região crítica) e o canal done para sinalizar quando eles devem parar
		filosofos[i] = NovoFilosofo(i+1, garfos[i], garfos[(i+1)%5], monitor, done)
	}

	// Adiciona as goroutines para cada Filósofo à WaitGroup
	wg.Add(5)
	for i := 0; i < 5; i++ {
		go func(f *Filosofo) {
			defer wg.Done()
			f.iniciar() // Inicia a rotina do Filósofo
		}(filosofos[i])
	}

	// Aguarde um pouco antes de encerrar as goroutines
	time.Sleep(10 * time.Second)
	close(done) // Sinaliza para os filósofos que eles devem parar

	wg.Wait() // Aguarde as goroutines terminarem
}
