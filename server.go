package main

import (
	"flag"
	"fmt"
	"net"
	"strconv"
	"strings"
)

// Mensagens de erro do protocolo.
const (
	erroFormato = "ERRO: formato invalido (use: OPERACAO NUM1 NUM2)"
	erroComando = "ERRO: comando desconhecido"
	erroDivZero = "ERRO: divisao por zero"
)

func main() {
	porta := flag.String("porta", "9000", "porta TCP em que o servidor escuta")
	flag.Parse()

	endereco := ":" + *porta

	listener, err := net.Listen("tcp", endereco)
	if err != nil {
		fmt.Println("Erro ao abrir a porta", *porta, "->", err)
		return
	}
	defer listener.Close()

	fmt.Println("Servidor de calculo escutando na porta", *porta)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Erro ao aceitar conexao:", err)
			continue
		}

		fmt.Printf("[CONECTOU] %s\n", conn.RemoteAddr())

		go tratarConexao(conn)
	}
}

func tratarConexao(conn net.Conn) {
	endereco := conn.RemoteAddr().String()

	defer func() {
		conn.Close()
		fmt.Printf("[DESCONECTOU] %s\n", endereco)
	}()

	buffer := make([]byte, 1024)

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			// EOF ou conexao perdida: encerra o atendimento deste cliente.
			return
		}

		mensagem := strings.TrimSpace(string(buffer[:n]))
		fmt.Printf("[%s] recebido: %q\n", endereco, mensagem)

		if strings.ToUpper(mensagem) == "SAIR" {
			return
		}

		resposta := processar(mensagem)

		if _, err := conn.Write([]byte(resposta + "\n")); err != nil {
			return
		}
	}
}

// processar interpreta uma mensagem do protocolo e devolve a resposta.
func processar(mensagem string) string {
	campos := strings.Fields(mensagem)
	if len(campos) == 0 {
		return erroFormato
	}

	operacao := strings.ToUpper(campos[0])

	// Primeiro validamos o comando: "HELLO" e comando desconhecido,
	// enquanto "MUL 3" e um comando valido com formato invalido.
	switch operacao {
	case "SOMA", "SUB", "MUL", "DIV":
	default:
		return erroComando
	}

	if len(campos) != 3 {
		return erroFormato
	}

	a, err1 := strconv.Atoi(campos[1])
	b, err2 := strconv.Atoi(campos[2])
	if err1 != nil || err2 != nil {
		return erroFormato
	}

	switch operacao {
	case "SOMA":
		return fmt.Sprintf("RESULTADO %d", a+b)
	case "SUB":
		return fmt.Sprintf("RESULTADO %d", a-b)
	case "MUL":
		return fmt.Sprintf("RESULTADO %d", a*b)
	case "DIV":
		if b == 0 {
			return erroDivZero
		}
		return fmt.Sprintf("RESULTADO %d", a/b)
	}

	return erroComando
}
