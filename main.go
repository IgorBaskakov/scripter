package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	const inlineWord = "Go"
	in := bufio.NewScanner(os.Stdin)

	for in.Scan() {
		txt := in.Text()

		body, err := get(txt)
		if err != nil {
			log.Fatalf("cann't get content from %s", txt)
		}
		cnt := strings.Count(body, inlineWord)
		fmt.Printf("Count for %s: %d\n", txt, cnt)
	}
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
