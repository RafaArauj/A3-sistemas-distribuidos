package main

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"
)

type Requisicao struct {
	Operacao string `json:"operacao"`
	Texto    string `json:"texto"`
}

type Resposta struct {
	Operacao  string `json:"operacao"`
	Maiusculas string `json:"maiusculas"`
	Invertido string `json:"invertido"`
}

func inverterTexto(texto string) string {
	letras := []rune(texto)

	for i, j := 0, len(letras)-1; i < j; i, j = i+1, j-1 {
		letras[i], letras[j] = letras[j], letras[i]
	}

	return string(letras)
}

func executorServidor2(requisicao Requisicao, canal chan Resposta) {
	resposta := Resposta{
		Operacao:   requisicao.Operacao,
		Maiusculas: strings.ToUpper(requisicao.Texto),
		Invertido: inverterTexto(requisicao.Texto),
	}

	canal <- resposta
}

func main() {
	listener, err := net.Listen("tcp", ":8082")
	if err != nil {
		fmt.Println("Erro ao iniciar servidor:", err)
		return
	}

	defer listener.Close()

	fmt.Println("Servidor 2 iniciado na porta 8082...")
	fmt.Println("Aguardando conexão...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Erro ao aceitar conexão:", err)
			continue
		}

		fmt.Println("Cliente conectado!")

		buffer := make([]byte, 4096)

		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Erro ao receber mensagem:", err)
			conn.Close()
			continue
		}

		var requisicao Requisicao

		err = json.Unmarshal(buffer[:n], &requisicao)
		if err != nil {
			fmt.Println("Erro ao interpretar JSON:", err)
			conn.Close()
			continue
		}

		fmt.Println("Texto recebido:", requisicao.Texto)

		canal := make(chan Resposta)

go executorServidor2(requisicao, canal)

resposta := <-canal

dados, err := json.Marshal(resposta)
		if err != nil {
			fmt.Println("Erro ao criar JSON:", err)
			conn.Close()
			continue
		}

		conn.Write(dados)

		conn.Close()
	}
}