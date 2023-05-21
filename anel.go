package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type mensagem struct {
	tipo  int    // tipo da mensagem para fazer o controle do que fazer (eleição, confirmacao da eleicao)
	corpo [4]int // conteudo da mensagem para colocar os ids (usar um tamanho ocmpativel com o numero de processos no anel)
}

var (
	chans = []chan mensagem{ // vetor de canias para formar o anel de eleicao - chan[0], chan[1] and chan[2] ...
		make(chan mensagem, 1),
		make(chan mensagem, 1),
		make(chan mensagem, 1),
		make(chan mensagem, 1),
	}
	controle   = make(chan int)
	wg         sync.WaitGroup // wg is used to wait for the program to finish
	letHimKnow bool           = false
	mutex      sync.Mutex     // mutex for synchronization
)

func ElectionControler(in chan int, firstLeader int) {
	defer wg.Done() //Só executada no final da função

	rand.Seed(time.Now().UnixNano())

	var actualLeader int
	var msg mensagem
	var warnChannel int

	actualLeader = firstLeader
	fmt.Println(actualLeader)

	msg.tipo = -1

	chans[0] <- msg

	fmt.Println("Controle: envio de mensagem vazia para todos os nós (inicialização)")

	for {
		msg.tipo = 2
		chans[actualLeader] <- msg

		warnChannel = generateRandomValue(len(chans), actualLeader)

		msg.tipo = 1
		chans[warnChannel] <- msg

		oldLeader := actualLeader
		fmt.Println("Controle: Esperando pelo novo lider")
		actualLeader = <-in
		fmt.Printf("Controle: Novo lider definido: Processo %d\n", actualLeader)

		// Aguardar 10 segundos
		time.Sleep(5 * time.Second)

		msg.tipo = 3
		msg.corpo[0] = actualLeader
		chans[oldLeader] <- msg
	}
}

func generateRandomValue(max int, avoid int) int {
	randomValue := rand.Intn(max)
	if randomValue == avoid {
		randomValue = (randomValue + 1) % max
	}
	return randomValue
}

func getMax(arr [4]int) int {
	if len(arr) == 0 {
		return 0 // Retorna um valor padrão caso o array esteja vazio
	}

	max := arr[0] // Assumimos que o primeiro elemento é o máximo inicialmente

	for _, element := range arr {
		if element > max {
			max = element // Atualiza o valor máximo se encontrar um elemento maior
		}
	}

	return max
}

func contains(arr [4]int, value int) bool {
	for _, element := range arr {
		if element == value {
			return true
		}
	}
	return false
}

func ElectionStage(TaskId int, in chan mensagem, out chan mensagem, leader int) {
	defer wg.Done()

	// variáveis locais que indicam se este processo é o líder e se está ativo
	var temp mensagem
	var actualLeader int
	var bFailed bool = false // todos iniciam sem falha

	actualLeader = leader // indicação do líder veio por parâmetro

	for {
		temp = <-in
		fmt.Printf("%2d: recebi mensagem %d, [ %d, %d, %d, %d ]\n", TaskId, temp.tipo, temp.corpo[0], temp.corpo[1], temp.corpo[2], temp.corpo[3])

		switch temp.tipo {
		case 1:
			//msg = temp
			fmt.Printf("%2d: Fui alertado de que houve falha\n", TaskId)
			temp.tipo = 4
			temp.corpo = [4]int{-1, -1, -1, -1}
			temp.corpo[TaskId] = TaskId
			fmt.Printf("%2d: ka\n", TaskId)
			out <- temp
			fmt.Printf("%2d: buum\n", TaskId)
			fmt.Printf("%2d: NOVA ELEIÇÃO DISPARADA\n", TaskId)

		case 2:
			bFailed = true
			fmt.Printf("%2d: falhei\n", TaskId)
			fmt.Printf("%2d: líder atual %d\n", TaskId, actualLeader)

		case 3:
			bFailed = false
			fmt.Printf("%2d: voltei a ativa, mas estou meio perdido\n", TaskId)
			fmt.Printf("%2d: líder atual %d\n", TaskId, actualLeader)
			if letHimKnow {
				fmt.Printf("%2d: descobri que o lider agora é %d\n", TaskId, temp.corpo[0])
				actualLeader = temp.corpo[0]
			}

		case 4:
			//msg = temp
			if !bFailed {
				fmt.Printf("%2d: Elegendo\n", TaskId)
				if contains(temp.corpo, TaskId) {
					eleito := getMax(temp.corpo)
					temp.tipo = 5
					temp.corpo = [4]int{-1, -1, -1, -1}
					temp.corpo[0] = eleito
					actualLeader = eleito

					fmt.Printf("%2d: Processo %d é o novo lider\n", TaskId, eleito)

					out <- temp
				} else {
					temp.corpo[TaskId] = TaskId
					out <- temp
				}
			} else {
				fmt.Printf("%2d: To morto cara, nao posso eleger\n", TaskId)
				out <- temp
			}

		case 5:
			//msg = temp
			fmt.Printf("%2d: bFailed: %t.\n", TaskId, bFailed)
			if !bFailed {
				if actualLeader == temp.corpo[0] {
					fmt.Printf("%2d: Eleição finalizada, todos cientes do novo lider %d.\n", TaskId, actualLeader)
					controle <- temp.corpo[0]
					temp.tipo = -1
					temp.corpo = [4]int{-1, -1, -1, -1}

					out <- temp
				} else {
					actualLeader = temp.corpo[0]

					fmt.Printf("%2d: Processo %d é o novo lider\n", TaskId, actualLeader)

					out <- temp
				}
			} else {
				out <- temp
			}

		case -1:
			fmt.Printf("%2d: idle\n", TaskId)
			temp.tipo = -1
			out <- temp

		default:
			fmt.Printf("%2d: não conheço este tipo de mensagem\n", TaskId)
			fmt.Printf("%2d: líder atual %d\n", TaskId, actualLeader)
			temp.tipo = -1
			out <- temp
		}
		time.Sleep(1000 * time.Millisecond)

		// fmt.Printf("%2d: terminei\n", TaskId)
	}
}

func main() {

	wg.Add(5) // Incrementamos o contador de WaitGroup para incluir o ElectionControler

	// Criar o processo controlador
	go ElectionControler(controle, 0)
	fmt.Println("\n   Processo controlador criado\n")

	// Esperar pela inicialização do processo controlador

	// Criar os processos do anel de eleição
	go ElectionStage(0, chans[0], chans[1], 0) // Este é o líder
	go ElectionStage(1, chans[1], chans[2], 0) // Não é líder, é o processo 0
	go ElectionStage(2, chans[2], chans[3], 0) // Não é líder, é o processo 0
	go ElectionStage(3, chans[3], chans[0], 0) // Não é líder, é o processo 0

	fmt.Println("\n   Anel de processos criado")

	wg.Wait() // Aguardar o término das goroutines
}
