package util

import (
	"os"
	"os/signal"
	"syscall"
)

func WaitQuit(fns ...func()){
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	for _,fn := range fns{
		fn()
	}
}
