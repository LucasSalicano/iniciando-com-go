package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const tempoMonitoramento = 5

func main() {
	for {
		introducao()
		menu()

		comando := lerComando()

		switch comando {
		case 1:
			monitoramento()
		case 2:
			imprimeLog()
		case 0:
			fmt.Println("Saindo...")
			os.Exit(0)
		default:
			fmt.Println("Comando desconhecido")
			os.Exit(-1)
		}
	}
}

func imprimeLog() {
	arquivo, err := ioutil.ReadFile("log.txt")

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(arquivo))
}

func monitoramento() {
	fmt.Println("Monitorando...")
	sites := sitesParaMonitoramente()

	for i := 0; i < tempoMonitoramento; i++ {
		testStatusSite(sites)
		fmt.Println("")
		time.Sleep(5 * time.Second)
	}
}

func testStatusSite(sites []string) {
	for _, site := range sites {
		resp, err := http.Get(site) // _ operador de identificado em branco ex: _ , err := http.Get(site)

		if err != nil {
			fmt.Println("Erro ao realizar a requisição: ", err)
		}

		if resp.StatusCode == 200 {
			fmt.Println("OK : ", site)
			registraLog(site, true, 200)
		} else {
			fmt.Println("Warning: ", resp.StatusCode, site)
			registraLog(site, false, resp.StatusCode)
		}
	}
}

func lerComando() int {
	var comando int
	fmt.Print("Digite um comando: ")
	fmt.Scan(&comando)
	return comando
}

func menu() {
	fmt.Println("[1] - Iniciar Monitoramente")
	fmt.Println("[2] - Exibir Histórico")
	fmt.Println("[0] - Sair")
}

func introducao() {
	nome := "Lucas"
	versao := 1.1
	fmt.Println("Olá", nome, " versão: ", versao)
}

func sitesParaMonitoramente() []string {
	var sites []string

	arquivo, err := os.Open("sites.txt")
	//arquivo, err := ioutil.ReadFile("sites.txt")

	if err != nil {
		fmt.Println("Erro ao abrir arquivo: ", err)
	}

	leitor := bufio.NewReader(arquivo)

	for {
		linha, err := leitor.ReadString('\n')
		linha = strings.TrimSpace(linha)

		sites = append(sites, linha)

		if err == io.EOF {
			break
		}
	}

	arquivo.Close()
	return sites
}

func registraLog(site string, status bool, statusCode int) {
	arquivo, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		fmt.Println(err)
	}

	arquivo.WriteString(time.Now().Format("02/01/2006 15:04:05:.000") + " | Online: " + strconv.FormatBool(status) + " | StatusCode: " + strconv.Itoa(statusCode) + " | " + site + "\n")

	arquivo.Close()
}
