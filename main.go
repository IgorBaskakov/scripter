package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

const inlineWord = "Go"
const workerlimit = 5

func main() {
	in := bufio.NewScanner(os.Stdin)

	total := 0
	datach := make(chan int)
	dataresult := make(chan int)
	go func() {
		for data := range datach {
			total += data
		}
		dataresult <- total
	}()

	// канал, реуглирующий ограничение на количество одновременно выполняемых запросов
	quotaCh := make(chan struct{}, workerlimit)
	// для синхронизации горутин, после окончания данных в потоке
	wg := &sync.WaitGroup{}

	for in.Scan() {
		txt := in.Text()
		wg.Add(1)
		quotaCh <- struct{}{}
		go startWorker(wg, txt, quotaCh, datach)
	}

	wg.Wait()
	close(datach)
	fmt.Printf("Total: %d\n", <-dataresult)
}

func startWorker(wg *sync.WaitGroup, uri string, quotaCh chan struct{}, datach chan int) {
	defer func() {
		wg.Done()
		<-quotaCh
	}()
	body, err := get(uri)
	if err != nil {
		log.Fatalf("can't get content from %s", uri)
	}

	cnt := strings.Count(body, inlineWord)
	fmt.Printf("Count for %s: %d\n", uri, cnt)
	datach <- cnt
}

func get(uri string) (string, error) {
	resp, err := http.Get(uri)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(result), nil
}
