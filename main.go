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

func main() {
	const inlineWord = "Go"
	in := bufio.NewScanner(os.Stdin)
	total, ratelimit := 0, 5
	quotaCh := make(chan struct{}, ratelimit)
	wg := &sync.WaitGroup{}
	i := 0

	for in.Scan() {
		txt := in.Text()

		quotaCh <- struct{}{}
		datach := make(chan string)
		wg.Add(1)
		i++
		go func(i int) {
			fmt.Printf("start %d task\n", i)
			body, err := get(txt)
			if err != nil {
				log.Fatalf("cann't get content from %s", txt)
			}
			datach <- body
		}(i)

		go func(i int) {
			fmt.Printf("start count for %d task\n", i)
			cnt := strings.Count(<-datach, inlineWord)
			total += cnt
			fmt.Printf("Count for %s: %d\n", txt, cnt)
			<-quotaCh
			wg.Done()
		}(i)
	}
	wg.Wait()
	fmt.Printf("Total: %d\n", total)
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
