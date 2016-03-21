package misc

import (
	"fmt"
	"log"
)

type Runner interface {
	Start() error
	Stop() error
}

type Run struct {
	Runners []Runner
}

func (r Run) Start() error {
	for i, runner := range r.Runners {
		err := runner.Start()
		if err != nil {
			log.Println("Runners : Error during runners start")
			r.Stop(i)
			return err
		}
	}
	return nil
}

func (r Run) Stop(j int) error {
	for i, runner := range r.Runners {
		if i > j {
			break
		}
		err := runner.Stop()
		if err != nil {
			log.Println("Runners : Error during runners stop")
			return err
		}
	}
	return nil
}

func PrintRunners(s string, x []Runner) {
	fmt.Printf("%s len=%d cap=%d %v\n",
		s, len(x), cap(x), x)
}
