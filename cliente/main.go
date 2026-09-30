package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"sync"
)

type Requisicao struct {
	Operacao string `json:"operacao"`
	Texto    string `json:"texto"`
}

type RespostaServidor1 struct {
	Operacao   string `json:"operacao"`
	Palavras   int    `json:"palavras"`
	Caracteres int    `json:"caracteres"`
	Vogais     int    `json:"vogais"`
}

type RespostaServidor2 struct {
	Operacao  string `json:"operacao"`
	Maiusculas string `json:"maiusculas"`
	Invertido string `json:"invertido"`
}

func consultarServidor1(dados []byte, wg *sync.WaitGroup) {
	defer wg.Done()

	conn, err := net.Dial("tcp", "localhost:8081")
	if err != nil {
		fmt.Println("Erro ao conectar ao Servidor 1:", err)
		return
	}

	defer conn.Close()

	conn.Write(dados)

	buffer := make([]byte, 4096)

	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Erro ao receber do Servidor 1:", err)
		return
	}

	var resposta RespostaServidor1

	err = json.Unmarshal(buffer[:n], &resposta)
	if err != nil {
		fmt.Println("Erro ao interpretar resposta do Servidor 1:", err)
		return
	}

	fmt.Println("\n--- SERVIDOR 1 ---")
	fmt.Println("Palavras:", resposta.Palavras)
	fmt.Println("Caracteres:", resposta.Caracteres)
	fmt.Println("Vogais:", resposta.Vogais)
}

func consultarServidor2(dados []byte, wg *sync.WaitGroup) {
	defer wg.Done()

	conn, err := net.Dial("tcp", "localhost:8082")
	if err != nil {
		fmt.Println("Erro ao conectar ao Servidor 2:", err)
		return
	}

	defer conn.Close()

	conn.Write(dados)

	buffer := make([]byte, 4096)

	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Erro ao receber do Servidor 2:", err)
		return
	}

	var resposta RespostaServidor2

	err = json.Unmarshal(buffer[:n], &resposta)
	if err != nil {
		fmt.Println("Erro ao interpretar resposta do Servidor 2:", err)
		return
	}

	fmt.Println("\n--- SERVIDOR 2 ---")
	fmt.Println("Maiúsculas:", resposta.Maiusculas)
	fmt.Println("Invertido:", resposta.Invertido)
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Digite um texto: ")

	texto, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Erro ao ler o texto:", err)
		return
	}

	requisicao := Requisicao{
		Operacao: "processar_texto",
		Texto:    texto,
	}

	dados, err := json.Marshal(requisicao)
	if err != nil {
		fmt.Println("Erro ao criar JSON:", err)
		return
	}

	fmt.Println("\nJSON enviado:")
	fmt.Println(string(dados))

	var wg sync.WaitGroup

	wg.Add(2)

	go consultarServidor1(dados, &wg)
	go consultarServidor2(dados, &wg)

	wg.Wait()

	fmt.Println("\nProcessamento concluído!")
}