package parser

import (
	"encoding/json"
	"io/ioutil"
)

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
