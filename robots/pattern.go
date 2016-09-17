package robots

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"
)

type Movement struct {
	Direction string // forward/backward/left/right
	Duration  time.Duration
}

type Data_new struct {
	fileData  []byte
	filePath  string
	patternOn bool
	movements []Movement
}

func (ia *AI) stopPattern(map[string]interface{}) interface{} {

	//find a way to stop pattern
	ia.pattern <- true
	return fmt.Sprintf("Pattern stopped at %v", time.Now())
}

func (ia *AI) startPattern(map[string]interface{}) interface{} {

	var data = new(Data_new)
	var err error

	data.fileData = make([]byte, 1024)
	data.movements = []Movement{}
	data.patternOn = true
	data.filePath = "./patternDirectory/default.json"

	if _, err := os.Stat(data.filePath); os.IsNotExist(err) {
		fmt.Println("Error default.json does not exist")
	}
	data.fileData, err = ioutil.ReadFile(data.filePath)
	if err != nil {
		log.Println(err)
		return fmt.Sprintf("Cannot start pattern %v: %v", data.filePath, err.Error())
	}
	err = json.Unmarshal([]byte(data.fileData), &data.movements)
	if err != nil {
		log.Println("error : " + err.Error())
		return fmt.Sprintf("Cannot start pattern %v: %v", data.filePath, err.Error())
	}

	var it bool = true

	defer ia.m.stopMoving()
	for it {
		for _, mov := range data.movements {
			fmt.Println("Prise en chage d'une nouvelle ligne json: [", mov.Direction, ",", mov.Duration, "]")
			switch mov.Direction {
			case "forward":
				fmt.Println("Avancer pendant ", time.Millisecond*mov.Duration)
				ia.m.moveForward()
				it = ia.interruptTimer(time.Millisecond * mov.Duration)
				break
			case "backward":
				fmt.Println("Recule pendant ", time.Millisecond*mov.Duration)
				ia.m.moveBackward()
				it = ia.interruptTimer(time.Millisecond * mov.Duration)
				break
			case "right":
				fmt.Println("Tourner droite pendant ", time.Millisecond*mov.Duration)
				ia.m.turnRight()
				it = ia.interruptTimer(time.Millisecond * mov.Duration)
				break
			case "left":
				fmt.Println("Tourner gauche pendant ", time.Millisecond*mov.Duration)
				ia.m.turnLeft()
				it = ia.interruptTimer(time.Millisecond * mov.Duration)
				break
			default:
				fmt.Println("Unknown comand in this pattern : [", mov.Direction, "]")
				return fmt.Sprintf("Unknown comand in %v : [%v]", data.filePath, mov.Direction)
			}
			ia.m.stopMoving()
			if !it {
				return fmt.Sprintf("Pattern interrupted %v", data.filePath)
			}
		}
	}
	return fmt.Sprintf("Pattern finished")
}

func (ia *AI) interruptTimer(t time.Duration) bool {
	select {
	case <-time.After(t):
		fmt.Println(t)
		return true
	case <-ia.pattern:
		fmt.Println("interrupted")
		return false
	}
}