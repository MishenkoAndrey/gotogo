package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

const (
	MaxGoroutineCount = 5
)

type Task struct {
	url   string
	count int
	err   error
}

func main() {
	res := make(chan Task)
	mutex := make(chan struct{}, MaxGoroutineCount)
	total := 0
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		url := scanner.Text()
		if url == "" {
			break
		}
		go Handle(res, &url, mutex)
		tmp := <-res
		if tmp.err != nil {
			fmt.Printf("Ошибка подсчёта в %s. Сообщение об ошибке: %s\n",
				tmp.url,
				tmp.err.Error(),
			)
			continue
		}
		fmt.Printf(
			"Количество слов для %s: %d\n",
			tmp.url,
			tmp.count,
		)
		total += tmp.count
	}
	close(res)
	close(mutex)
	fmt.Printf("Общее количество: %d", total)
}

func Handle(c chan Task, url *string, mutex chan struct{}) {
	// lock
	mutex <- struct{}{}
	resp, err := http.Get(*url)
	if err != nil {
		c <- Task{err: err, url: *url}
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	responseText := string(body)
	count := GetGoCount(responseText)
	c <- Task{url: *url, count: count}
	// release
	<-mutex
}
func GetGoCount(source string) (count int) {
	for i := 0; i < len(source)-1; i++ {
		if source[i] == 'G' && source[i+1] == 'o' {
			count++
			i++
		}
	}
	return
}
