package main

import (
	"context"
	"log"
	"time"
)

func main() {
	ctx,cancel := context.WithCancel(context.Background())
	go func() {
		ticker := time.NewTicker(time.Duration(1) * time.Second)
		defer func() {
			ticker.Stop()
			log.Println("defer")
		}()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				log.Println("123")
			default:

			}

		}
	}()
	time.Sleep(5*time.Second)
	cancel()
	time.Sleep(30*time.Second)
}

