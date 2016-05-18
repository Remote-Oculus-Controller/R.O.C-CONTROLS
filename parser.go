package roc

import (
	"io/ioutil"
	"fmt"
	"encoding/json"
)

type Cmd struct {
	Name              string
	Code              byte
	Min, Max, Default int
}

type Cmds struct {
	Commands []Cmd
}

func RobotCommand(f_path string) (map[string]Cmd, error) {

	var cmds Cmds

	c := make(map[string]Cmd)
	data, err := ioutil.ReadFile(f_path)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	err = json.Unmarshal(data, &cmds)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	for _, kk := range cmds.Commands {
		c[kk.Name] = kk
	}
	return c, nil
}

func Decode(fp string) (map[string]interface{}, error) {

	var c interface{}

	b, err := ioutil.ReadFile(fp)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, &c)
	if err != nil {
		return nil, err
	}
	return c.(map[string]interface{}), nil
}

