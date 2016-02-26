package parser

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Cmd struct {
	Name  string
	Code  string
	Other []interface{}
}

type Cmds struct {
	Commands []Cmd
}

func RobotCommand(f_path string) (map[string]Cmd, error) {

	var cmds Cmds

	c := make(map[string]Cmd)
	data, err := ioutil.ReadFile("command.json")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	err = json.Unmarshal(data, &cmds)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Println(cmds)
	for _, kk := range cmds.Commands {
		c[kk.Name] = kk
	}
	fmt.Println(c)
	return c, nil
}
