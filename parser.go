package roc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type cmd struct {
	Name              string
	Code              byte
	Min, Max, Default int
}

type cmds struct {
	Commands []cmd
}

func parseCommands(fPath string) (map[string]cmd, error) {

	var cmds cmds

	c := make(map[string]cmd)
	data, err := ioutil.ReadFile(fPath)
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

func decodeJSONFile(fp string) (map[string]interface{}, error) {

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
