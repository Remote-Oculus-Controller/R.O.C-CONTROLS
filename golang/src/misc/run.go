package misc

import (
	"fmt"
)

type Runner interface {
	Start() error
	Stop() error
}

type Run struct {
	Runners []Runner
}

func (r *Run) Start() error {
	for _, runner := range r.Runners {
		err := runner.Start()
		if err != nil {
			r.Stop()
			fmt.Println("Error during runners start")
			return err
		}
	}
	return nil
}

func (r *Run) Stop() error {
	for _, runner := range r.Runners {
		err := runner.Stop()
		if err != nil {
			fmt.Println("Error during runners stop")
			return err
		}
	}
	return nil
}

func PrintRunners(s string, x []Runner) {
	fmt.Printf("%s len=%d cap=%d %v\n",
		s, len(x), cap(x), x)
}
