# Calculadora Remota — Sockets TCP em Go

Atividade prática de Redes de Computadores. Cliente e servidor conversam por um
protocolo textual simples sobre TCP.

## Como executar

Abra **dois terminais** na pasta do projeto.

Terminal 1 — servidor:

    go run server.go

Terminal 2 — cliente:

    go run client.go

Flags opcionais:

    go run server.go -porta 9000
    go run client.go -host 127.0.0.1 -porta 9000

> Os dois arquivos são `package main`, então rode sempre `go run server.go` /
> `go run client.go` (nomeando o arquivo). Um `go run .` tentaria compilar os
> dois juntos e reclamaria de dois `func main`.

## Protocolo

| Envio        | Resposta                                        |
|--------------|-------------------------------------------------|
| `SOMA 10 20` | `RESULTADO 30`                                  |
| `SUB 50 8`   | `RESULTADO 42`                                  |
| `MUL 6 7`    | `RESULTADO 42`                                  |
| `DIV 20 4`   | `RESULTADO 5`                                   |
| `DIV 10 0`   | `ERRO: divisao por zero`                        |
| `HELLO`      | `ERRO: comando desconhecido`                    |
| `MUL 3`      | `ERRO: formato invalido (use: OPERACAO NUM1 NUM2)` |
| `SAIR`       | (servidor fecha a conexão)                      |

## Roteiro de teste para a correção

1. Suba o servidor e conecte um cliente — o servidor imprime `[CONECTOU] 127.0.0.1:xxxxx`.
2. Rode a sequência: `SOMA 10 20`, `SUB 50 8`, `MUL 6 7`, `DIV 20 4`.
3. Erros: `DIV 10 0`, `HELLO`, `MUL 3`, `SOMA a b`.
4. **Múltiplos clientes**: abra um terceiro terminal com outro `go run client.go`
   e mande comandos nos dois ao mesmo tempo — ambos são atendidos.
5. Digite `SAIR` em um deles: o cliente encerra, o servidor loga `[DESCONECTOU]`
   e continua atendendo o outro cliente normalmente.

## Detalhes da implementação

- **Servidor**: `net.Listen` + laço de `Accept`; cada conexão vai para
  `go tratarConexao(conn)`, que lê em laço com um buffer `[]byte` até receber
  `SAIR` ou o cliente cair (`Read` retorna erro/EOF). O `defer` fecha a conexão
  e imprime o log de desconexão.
- **Parsing**: `strings.Fields` separa os campos; o comando é validado *antes*
  da contagem de argumentos, para que `HELLO` seja "comando desconhecido" e
  `MUL 3` seja "formato invalido". `strconv.Atoi` valida os operandos.
- **Cliente**: `net.Dial`, `bufio.Scanner` em `os.Stdin` para o prompt `calc> `,
  envia com `Write` e lê a resposta com `Read` no mesmo buffer.
- Comandos são aceitos em maiúsculas ou minúsculas (`strings.ToUpper`).
- Limitação conhecida: como o protocolo usa um `Read` por comando (padrão do
  echo server visto em aula), dois comandos enviados em rajada poderiam chegar
  colados no mesmo pacote. No uso interativo isso não acontece, porque o cliente
  só envia o próximo comando depois de receber a resposta.
