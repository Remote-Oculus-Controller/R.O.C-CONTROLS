package misc

import (
	"log"
	"os"
)

func CheckError(err error, msg string, out bool) error {
	if err != nil {
		if out {
			log.Println("Fatal error: ", err.Error(), "From", msg)
			os.Exit(1)
		} else {
			log.Println(err.Error(), "From", msg)
		}
	}
	return err
}
