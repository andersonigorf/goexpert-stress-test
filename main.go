package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
	"sync"
	"time"
)

var lock sync.Mutex

func main() {
	data := parseFlags()
	flag.Parse()
	if !validateFlags(data) {
		return
	}

	mapOfStatusCodes, statusCodeOk, requestsRealizados := doRequests(data)

	printResults(mapOfStatusCodes, statusCodeOk, requestsRealizados, data.url)
}

type requestData struct {
	url         *string
	requests    *int
	concurrency *int
}

func parseFlags() requestData {
	return requestData{
		url:         flag.String("url", "http://google.com", "URL do serviço a ser testado"),
		requests:    flag.Int("requests", 50, "Número total de requests"),
		concurrency: flag.Int("concurrency", 10, "Número de chamadas simultâneas"),
	}
}

func validateFlags(data requestData) bool {
	if *data.url == "" {
		fmt.Println("--url não pode ser vazia")
		return false
	}
	if *data.requests <= 0 {
		fmt.Println("--requests não pode ser vazia e deve ser maior que 0")
		return false
	}
	if *data.concurrency <= 0 {
		fmt.Println("--concurrency não pode ser vazio e deve ser maior que 0")
		return false
	}
	return true
}

func doRequests(data requestData) (map[int]int, int, int) {
	var requestsRealizados = 0
	var statusCodeOk = 0
	mapOfStatusCodes := make(map[int]int)
	wg := sync.WaitGroup{}
	wg.Add(*data.requests)
	for i := 0; i < *data.requests; i++ {
		go doSingleRequest(i+1, data.url, &wg, &lock, &statusCodeOk, &requestsRealizados, mapOfStatusCodes)
	}
	wg.Wait()

	return mapOfStatusCodes, statusCodeOk, requestsRealizados
}

func doSingleRequest(i int, url *string, wg *sync.WaitGroup, lock *sync.Mutex, statusCodeOk *int, requestsRealizados *int, mapOfStatusCodes map[int]int) {
	lock.Lock()
	defer lock.Unlock()
	defer wg.Done()
	code := makeRequest(i, url)
	println("Request finalizada com status code ", code)
	mapOfStatusCodes[code]++
	if code == http.StatusOK {
		(*statusCodeOk)++
	}
	(*requestsRealizados)++
}

func makeRequest(i int, url *string) int {
	var client = &http.Client{
		Timeout:       3 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error { return http.ErrUseLastResponse },
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	fmt.Printf("#%d Request iniciada \n", i)
	req, err := client.Get(*url)
	if err != nil {
		fmt.Printf("Request falhou com erro %v\n", err)
		return 0
	}
	req.Close = true
	defer req.Body.Close()

	return req.StatusCode
}

func printResults(mapOfStatusCodes map[int]int, statusCodeOk int, requestsRealizados int, url *string) {
	fmt.Printf("RELATÓRIO FINAL:\n")
	fmt.Printf("URL: %s\n", *url)
	fmt.Printf("Quantidade total de requests: %d\n", requestsRealizados)
	fmt.Printf("Quantidade de requests com status HTTP 200: %d\n", statusCodeOk)
	fmt.Printf("Distribuição de códigos de status HTTP:\n")
	for status, count := range mapOfStatusCodes {
		if status != http.StatusOK {
			fmt.Printf("Status %d: %d requests\n", status, count)
		}
	}
}
