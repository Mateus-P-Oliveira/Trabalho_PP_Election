package main

import (
	"math/rand"
	"sync"
	"time"
)

type nodo struct {
	id int // Id do Nodo
	//act bool  // Se esta ativo ou não
}

var (
	channels = []chan nodo{ // vetor de canias para formar o anel de eleicao - chan[0], chan[1] and chan[2] ...
		make(chan nodo),
		make(chan nodo),
		make(chan nodo),
		make(chan nodo),
	}
	controls = make(chan int)
	wg       sync.WaitGroup // wg is used to wait for the program to finish
)

func ElectionControler() {
	defer wg.Done()

	rand.Seed(time.Now().UnixNano())
	randomErrorNode := rand.Intn(5)

}

func RingCreation() {
	defer wg.Done()

	// variaveis locais que indicam se este processo é o lider e se esta ativo

	var actualLeader int
	var bFailed bool = false // todos inciam sem falha

	actualLeader = leader // indicação do lider veio por parâmatro

}

func main() {
	wg.Add(5)

}
