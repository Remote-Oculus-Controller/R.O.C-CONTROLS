package robots

import (
	"github.com/Happykat/R.O.C-CONTROLS"
	"github.com/Happykat/R.O.C-CONTROLS/misc"
	"fmt"
	"log"
	"go/types"
	"github.com/hybridgroup/gobot"
	"debug/dwarf"
	"errors"
)

type Video struct {
	*roc.RocRobot
}

const (
	VIDEO_TAG = 0x10
	CANNY = VIDEO_TAG | 1
	FACE = VIDEO_TAG | 2
	ARROW =	VIDEO_TAG | 3
	ZOOM = VIDEO_TAG | 4
	FOCUS		= VIDEO_TAG | 5
	BRIGHT		= VIDEO_TAG | 6
	GAIN		= VIDEO_TAG | 7
	CONTRAST	= VIDEO_TAG | 8
	SATURATION	= VIDEO_TAG | 9

	PARAM_ERR	= "MISSING %s in parameters"
	MIN_ERROR	= "min"
	MAX_ERROR	= "max"
	SCALE_ERROR	= "scale"
)

func NewVideo() *Video {

	vid := new(Video)
	vid.RocRobot = roc.NewRocRobot(nil)
	vid.Robot = gobot.NewRobot("video")
	vid.AddFunc(nil, 0, vid.startCannyEdgeAPI, "startCanny")
	vid.AddFunc(nil, 0, vid.stopCannyEdgeAPI, "stopCanny")
	vid.AddFunc(nil, 0, vid.startFaceDetectAPI, "startFaceDetect")
	vid.AddFunc(nil, 0, vid.stopFaceDetectAPI, "stopFaceDetect")
	return vid
}

func (vid *Video)startCannyEdge(Min, Max float64) error {

	min, err := misc.EncodeBytes(Min)
	if err != nil {
		log.Println("Error setting minimum to start canny", err.Error())
		return errors.New(fmt.Sprintln("Error setting minimum to start canny", err.Error()))
	}
	max, err := misc.EncodeBytes(Max)
	if err != nil {
		log.Println("Error setting maximum to start canny", err.Error())
		return errors.New(fmt.Sprintln("Error setting maximum to start canny", err.Error()))
	}
	buff := []byte{roc.DST_L | roc.CMD, CANNY}
	b := append(min, max...)
	buff = append(buff, b...)
	err = vid.Send(buff)
	if err != nil {
		return err
	}
	return nil
}

func (vid *Video)startCannyEdgeAPI(params map[string]interface{}) interface{} {
	min, ok := params["min"]
	_, assert := min.(float64)
	if (!ok || !assert) {
		log.Println(MIN_ERROR)
		return MIN_ERROR
	}
	max, ok := params["max"]
	_, assert = max.(float64)
	if (!ok || !assert) {
		log.Println(MAX_ERROR)
		return MAX_ERROR
	}
	err := vid.startCannyEdge(min.(float64), max.(float64))
	if err != nil {
		return nil
	}
	return "OK"
}

func (vid *Video)stopCannyEdge() {
	b := []byte{roc.DST_L | roc.CMD, CANNY, 0}
	vid.Send(b)
}
func (vid *Video)stopCannyEdgeAPI(params map[string]interface{}) interface{} {
	vid.stopCannyEdge()
	return 200
}

func (vid *Video)startFaceDetect(Scale float64) {

	scale, err := misc.EncodeBytes(Scale)
	if err != nil {
		log.Println("Error setting minimum to start canny", err.Error())
		return
	}
	b := append([]byte{roc.DST_L | roc.CMD, FACE}, scale...)
	vid.Send(b)
}

func (vid *Video)startFaceDetectAPI(params map[string]interface{}) interface{} {
	err := vid.CheckAPIParams(params, []types.BasicKind{types.Float64}, "scale")
	if (err != nil) {
		log.Println(err.Error())
		return err.Error()
	}
	vid.startFaceDetect(params["scale"].(float64))
	return 200
}

func (vid *Video)stopFaceDetect() {
	b := []byte{roc.DST_L | roc.CMD, FACE, 0}
	vid.Send(b)
}
func (vid *Video)stopFaceDetectAPI(params map[string]interface{}) interface{} {
	vid.stopFaceDetect()
	return 200
}
