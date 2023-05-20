package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type mensagem struct {
	tipo  int    // tipo da mensagem para fazer o controle do que fazer (eleição, confirmacao da eleicao)
	corpo [3]int // conteudo da mensagem para colocar os ids (usar um tamanho ocmpativel com o numero de processos no anel)
}

var (
	chans = []chan mensagem{ // vetor de canias para formar o anel de eleicao - chan[0], chan[1] and chan[2] ...
		make(chan mensagem),
		make(chan mensagem),
		make(chan mensagem),
		make(chan mensagem),
	}
	controle = make(chan int)
	wg       sync.WaitGroup // wg is used to wait for the program to finish
)

func ElectionControler(in chan int) {
	defer wg.Done() //Só executada no final da função

	var temp mensagem
	rand.Seed(time.Now().UnixNano())
	randomErrorNode := rand.Intn(3)
	randomErrorChannel := rand.Intn(3) //Ele da erro ao ser usado como o temp
	// comandos para o anel iciam aqui

	//fmt.Printf("consome %d %d\n", randomErrorNode, randomErrorChannel)
	// mudar o processo 0 - canal de entrada 3 - para falho (defini mensagem tipo 2 pra isto)

	temp.tipo = randomErrorNode //1 e 4 e 5 se forem atribuidos dão erro
	//Aloca Valores dos Aneis que não forem os escolhidos-------------------
	switch randomErrorChannel {
	case 0:
		{
			temp.tipo = 99 //Mensagem generica
			chans[1] <- temp
			fmt.Printf("%d\n", <-in) //Controle de Confirmação
			chans[2] <- temp
			fmt.Printf("%d\n", <-in) //Controle de Confirmação
			chans[3] <- temp
			fmt.Printf("%d\n", <-in) //Controle de Confirmação

		}
	case 1:
		{
			temp.tipo = 99 //Mensagem generica
			chans[0] <- temp
			fmt.Printf("%d\n", <-in) //Controle de Confirmação
			chans[2] <- temp
			fmt.Printf("%d\n", <-in) //Controle de Confirmação
			chans[3] <- temp
			fmt.Printf("%d\n", <-in) //Controle de Confirmação
		}
	case 2:
		{
			temp.tipo = 99 //Mensagem generica
			chans[0] <- temp
			fmt.Printf("%d\n", <-in) //Controle de Confirmação
			chans[1] <- temp
			fmt.Printf("%d\n", <-in) //Controle de Confirmação
			chans[3] <- temp
			fmt.Printf("%d\n", <-in) //Controle de Confirmação
		}
	case 3:
		{
			temp.tipo = 99 //Mensagem generica
			chans[0] <- temp
			fmt.Printf("%d\n", <-in) //Controle de Confirmação
			chans[1] <- temp
			fmt.Printf("%d\n", <-in) //Controle de Confirmação
			chans[2] <- temp
			fmt.Printf("%d\n", <-in) //Controle de Confirmação
		}
	default:
		{
			fmt.Printf("Nodo não existe \n")
		}
	}
	//Não posso deixar de atribuir valor pros Channels
	chans[randomErrorChannel] <- temp
	fmt.Printf("Nodo: %d Recebeu: %d\n", randomErrorChannel, randomErrorNode)

	fmt.Printf("Controle: confirmação %d\n", <-in) // receber e imprimir confirmação

	// mudar o processo 1 - canal de entrada 0 - para falho (defini mensagem tipo 2 pra isto)

	/*temp.tipo = randomErrorNode
	chans[0] <- temp
	fmt.Printf("Controle: mudar o processo 1 para falho\n")
	fmt.Printf("Controle: confirmação %d\n", <-in) // receber e imprimir confirmação

	// matar os outrs processos com mensagens não conhecidas (só pra cosumir a leitura)

	temp.tipo = randomErrorNode //Não podem enviar para o mesmo tipo de mensagem sem o imprimir confirmação
	chans[1] <- temp
	fmt.Printf("Controle: confirmação %d\n", <-in) // receber e imprimir confirmação //Precisa dele para não dar deadlock
	temp.tipo = randomErrorNode
	chans[2] <- temp
	fmt.Printf("Controle: confirmação %d\n", <-in) // receber e imprimir confirmação
	*/

	fmt.Println("\n   Processo controlador concluído\n")
}

func ElectionStage(TaskId int, in chan mensagem, out chan mensagem, leader int) {
	defer wg.Done()

	// variaveis locais que indicam se este processo é o lider e se esta ativo

	var actualLeader int
	var bFailed bool = false // todos inciam sem falha

	actualLeader = leader // indicação do lider veio por parâmatro

	temp := <-in // ler mensagem
	//Corpo é o vetor de votação e precisa inclunir todos valores de preocessos que pertencem ao anel
	fmt.Printf("%2d: recebi mensagem %d, [ %d, %d, %d ]\n", TaskId, temp.tipo, temp.corpo[0], temp.corpo[1], temp.corpo[2])
	switch temp.tipo {
	case 1:
		{
			fmt.Printf("Caso 1 \n")

			controle <- -5
		}
	case 2:
		{
			//fmt.Printf("%2d: recebi mensagem %d, [ %d, %d, %d ]\n", TaskId, temp.tipo, temp.corpo[0], temp.corpo[1], temp.corpo[2])
			bFailed = true
			fmt.Printf("%2d: falho %v \n", TaskId, bFailed)
			fmt.Printf("%2d: lider atual %d\n", TaskId, actualLeader)
			controle <- -5
		}
	case 3:
		{

			//fmt.Printf("%2d: recebi mensagem %d, [ %d, %d, %d ]\n", TaskId, temp.tipo, temp.corpo[0], temp.corpo[1], temp.corpo[2])
			bFailed = false
			fmt.Printf("%2d: falho %v \n", TaskId, bFailed)
			fmt.Printf("%2d: lider atual %d\n", TaskId, actualLeader)
			controle <- -5
		}
	case 4:
		{
			//fmt.Printf("%2d: recebi mensagem %d, [ %d, %d, %d ]\n", TaskId, temp.tipo, temp.corpo[0], temp.corpo[1], temp.corpo[2])
			fmt.Printf("Caso 4 \n")
			controle <- -5
		}
	case 5:
		{
			//fmt.Printf("%2d: recebi mensagem %d, [ %d, %d, %d ]\n", TaskId, temp.tipo, temp.corpo[0], temp.corpo[1], temp.corpo[2])
			fmt.Printf("Caso 5 \n")
			controle <- -5
		}
	default:
		{
			//fmt.Printf("%2d: não conheço este tipo de mensagem\n", TaskId)
			//fmt.Printf("%2d: lider atual %d\n", TaskId, actualLeader)
			controle <- 99

		}
	}

	//fmt.Printf("%2d: terminei \n", TaskId)
}

func main() {

	wg.Add(5) // Add a count of four, one for each goroutine

	// criar os processo do anel de eleicao

	go ElectionStage(0, chans[3], chans[0], 0) // este é o lider
	go ElectionStage(1, chans[0], chans[1], 0) // não é lider, é o processo 0
	go ElectionStage(2, chans[1], chans[2], 0) // não é lider, é o processo 0
	go ElectionStage(3, chans[2], chans[3], 0) // não é lider, é o processo 0

	fmt.Println("\n   Anel de processos criado")

	// criar o processo controlador

	go ElectionControler(controle)

	fmt.Println("\n   Processo controlador criado\n")

	wg.Wait() // Wait for the goroutines to finish\
}
