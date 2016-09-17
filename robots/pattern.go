package robots

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

type Movement struct {
	Direction string // forward/backward/left/right
	Duration  int
}

type Data struct {
	fileData  string
	filePath  string
	patternOn bool
	movements []Movement
}

func (ia *AI) stopPattern() {
	//find a way to stop pattern
}

func (ia *AI) startPattern() {

	var data = new(Data)
	data.movements = new([]Movement)
	data.patternOn = true
	data.filePath = "./patternDirectory/default.json"

	if _, err := os.Stat(data.filePath); os.IsNotExist(err) {
		fmt.Println("Error default.json does not exist")
	}
	data.fileData = ioutil.ReadFile(data.filePath)
	err := json.Unmarshal([]byte(data.fileData), &data.movements)
	if err != nil {
		fmt.Println("error : " + err)
	}
	for {
		if data.patternOn == false {
			return
		}
		for _, mov := range data.movements {
			fmt.Println("Prise en chage d'une nouvelle ligne json: [" + mov.Direction + "," + mov.Duration + "]")
			switch mov.Direction {
			case "forward":
				fmt.Println("Avancer pendant " + mov.Duration + " ms")
				ia.m.moveForward()
				<-time.After(time.Millisecond * mov.Duration)
				ia.m.stopMoving()
				break
			case "backward":
				fmt.Println("Avancer pendant " + mov.Duration + " ms")
				ia.m.moveBackward()
				<-time.After(time.Millisecond * mov.Duration)
				ia.m.stopMoving()
				break
			case "right":
				fmt.Println("Avancer pendant " + mov.Duration + " ms")
				ia.m.turnRight()
				<-time.After(time.Millisecond * mov.Duration)
				ia.m.stopMoving()
				break
			case "left":
				fmt.Println("Avancer pendant " + mov.Duration + " ms")
				ia.m.turnLeft()
				<-time.After(time.Millisecond * mov.Duration)
				ia.m.stopMoving()
				break
			default:
				fmt.Println("Unknown comand in this pattern : [" + mov.Direction + "]")
				return
				break
			}
		}
	}
}
