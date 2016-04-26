package misc

import (
	"log"
	"os"
)

func CheckError(err error, msg string, out bool) {
	if err != nil {
		log.Println("Fatal error: ", err.Error(), "From", msg)
		if out {
			os.Exit(1)
		}
	}
}
