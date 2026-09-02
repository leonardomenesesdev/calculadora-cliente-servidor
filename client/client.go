package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	host := flag.String("host", "127.0.0.1", "endereco do servidor")
	porta := flag.String("porta", "9000", "porta do servidor")
	flag.Parse()

	endereco := net.JoinHostPort(*host, *porta) //junta host e porta (host:porta)

	conn, err := net.Dial("tcp", endereco) //tenta estabelecerr a conexao tcp com o server
	if err != nil {
		fmt.Println("Nao foi possivel conectar em", endereco, "->", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Printf("Conectado ao servidor %s\n\n", endereco)

	scanner := bufio.NewScanner(os.Stdin)
	buffer := make([]byte, 1024)

	for {
		fmt.Print("calc> ")
		//scanner retorna false se o usuario encerra a entrada
		if !scanner.Scan() {
			fmt.Println()
			conn.Write([]byte("SAIR"))
			break
		}

		comando := strings.TrimSpace(scanner.Text())
		if comando == "" {
			continue
		}

		if _, err := conn.Write([]byte(comando)); err != nil {
			fmt.Println("Erro ao enviar comando:", err)
			break
		}

		if strings.ToUpper(comando) == "SAIR" {
			break
		}

		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("\nConexao encerrada pelo servidor.")
			return
		}

		fmt.Println(strings.TrimSpace(string(buffer[:n])))
		fmt.Println()
	}

	fmt.Println("Conexao encerrada.")
}
