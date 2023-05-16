package main

import (
	"sync"
	"fmt"
	"math/rand"
	"time"
  )

type nodo struct {
	id  int  // Id do Nodo
	act bool true // Se esta ativo ou não
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

func ElectionControler(){
	defer wg.Done()
	rand.Seed(time.Now().UnixNano())
	int randomErrorNode <- rand.Intn(5)

}