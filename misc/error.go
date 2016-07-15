package misc

import (
	"log"
)

func CheckError(err error, msg string, out bool) error {
	if err != nil {
		if out {
			log.Fatalln("Fatal error: ", msg, err.Error())
			panic(err)
		} else {
			log.Println(err.Error(), "From", msg)
		}
	}
	return err
}
