package main

import "time"

func main() {
	b := make([]byte, 60*1024*1024)
	for i := range b {
		b[i] = 1
	}
	time.Sleep(time.Hour)
}
