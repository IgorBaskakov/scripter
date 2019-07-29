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

func main() {
	in := bufio.NewScanner(os.Stdin)

	total := 0
	datach := make(chan int)
	go func() {
		for data := range datach {
			total += data
		}
	}()

	ratelimit := 5
	// канал, реуглирующий ограничение на количество одновременно выполняемых запросов
	quotaCh := make(chan struct{}, ratelimit)
	// для синхронизации горутин, после окончания данных в потоке
	wg := &sync.WaitGroup{}

	for in.Scan() {
		txt := in.Text()
		wg.Add(1)
		go startWorker(wg, txt, quotaCh, datach)
	}

	wg.Wait()
	close(datach)
	fmt.Printf("Total: %d\n", total)
}

func startWorker(wg *sync.WaitGroup, uri string, quotaCh chan struct{}, datach chan int) {
	quotaCh <- struct{}{}
	defer wg.Done()
	body, err := get(uri)
	if err != nil {
		log.Fatalf("can't get content from %s", uri)
	}

	cnt := strings.Count(body, inlineWord)
	fmt.Printf("Count for %s: %d\n", uri, cnt)
	datach <- cnt
	<-quotaCh
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
