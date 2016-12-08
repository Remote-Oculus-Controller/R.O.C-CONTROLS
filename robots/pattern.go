package robots

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

type Movement struct {
	Direction string // forward/backward/left/right
	Duration  time.Duration
}

type Data_new struct {
	fileData   []byte
	filePath   string
	patternOn  bool
	repeat     bool
	customFile bool
	movements  []Movement
}

func (ia *AI) stopPattern(map[string]interface{}) interface{} {

	//find a way to stop pattern
	ia.pattern <- true
	return fmt.Sprintf("Pattern stopped at %v", time.Now())
}

func (ia *AI) startPattern(m map[string]interface{}) interface{} {

	var data = new(Data_new)
	var err error

	data.fileData = make([]byte, 1024)
	data.movements = []Movement{}
	data.patternOn = true
	data.repeat = false
	data.customFile = false

	log.Println(m)
	if m["path"] != nil {
		data.filePath = m["path"].(string)
	} else {
		data.filePath = "./robots/patternDirectory/default.json"
	}

	if m["repeat"] != nil {
		data.repeat = m["repeat"].(bool)
	}

	if m["custom"] != nil {
		data.customFile = true
		data.filePath = m["custom"].(string)
	}
	if _, err := os.Stat(data.filePath); os.IsNotExist(err) {
		log.Println(err.Error())
		return fmt.Sprintf(err.Error())
	}
	data.fileData, err = ioutil.ReadFile(data.filePath)
	if err != nil {
		log.Println(err.Error())
		return fmt.Sprintf("Cannot read pattern %v: %v", data.filePath, err.Error())
	}
	log.Println(data.fileData)
	err = json.Unmarshal([]byte(data.fileData), &data.movements)
	if err != nil {
		s := fmt.Sprintf("Cannot unmarshallJSON pattern %v: %v", data.filePath, err.Error())
		log.Println(s)
		return s
	} else if data.customFile == true {
		files, err := ioutil.ReadDir("./patternDirectory")
		if err != nil {
			log.Fatal(err)
			return fmt.Sprintf("ERROR DURING READFILE")
		}
		var count = 0
		for _, file := range files {
			if strings.ContainsAny(file.Name(), "custom") == true {
				count += 1
			}
			fmt.Println(fmt.Sprintf("Found %d custom files", count))
		}
		error := ioutil.WriteFile(fmt.Sprintf("./patternDirectory/customFile%d.json", count+1), data.fileData, 0644)
		if error != nil {
			fmt.Println("Cannot create customFile")
		}
	}

	var it bool = false

	defer ia.m.stopMoving()
	for ok := true; ok; ok = data.repeat {
		for _, mov := range data.movements {
			fmt.Println("Prise en chage d'une nouvelle ligne json: [", mov.Direction, ",", mov.Duration, "]")
			switch mov.Direction {
			case "forward":
				fmt.Println("Avancer pendant ", time.Millisecond*mov.Duration)
				ia.m.moveForward()
				break
			case "backward":
				fmt.Println("Recule pendant ", time.Millisecond*mov.Duration)
				ia.m.moveBackward()
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
			it = ia.interruptTimer(time.Millisecond * mov.Duration)
			if it {
				return fmt.Sprintf("Pattern interrupted %v", data.filePath)
			}
		}

	}
	return fmt.Sprint("Pattern finished")
}

func (ia *AI) interruptTimer(t time.Duration) bool {
	select {
	case <-time.After(t):
		fmt.Println(t)
		return false
	case <-ia.pattern:
		fmt.Println("interrupted")
		return true
	}
}
