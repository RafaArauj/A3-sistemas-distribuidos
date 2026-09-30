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
	Palavras  int    `json:"palavras"`
	Caracteres int   `json:"caracteres"`
	Vogais    int    `json:"vogais"`
}

func contarVogais(texto string) int {
	vogais := "aeiouAEIOU"
	contador := 0

	for _, letra := range texto {
		if strings.ContainsRune(vogais, letra) {
			contador++
		}
	}

	return contador
}

func executorServidor1(requisicao Requisicao, canal chan Resposta) {
	resposta := Resposta{
		Operacao:   requisicao.Operacao,
		Palavras:   len(strings.Fields(requisicao.Texto)),
		Caracteres: len([]rune(requisicao.Texto)),
		Vogais:     contarVogais(requisicao.Texto),
	}

	canal <- resposta
}

func main() {
	listener, err := net.Listen("tcp", ":8081")
	if err != nil {
		fmt.Println("Erro ao iniciar servidor:", err)
		return
	}

	defer listener.Close()

	fmt.Println("Servidor 1 iniciado na porta 8081...")
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

go executorServidor1(requisicao, canal)

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