package robots

import (
	"github.com/Happykat/R.O.C-CONTROLS"
	"github.com/hybridgroup/gobot"
)

type Video struct {
	*roc.RocRobot
}

const (
	VIDEO_TAG  = 0x10
	CANNY      = VIDEO_TAG | 1
	FACE       = VIDEO_TAG | 2
	ARROW      = VIDEO_TAG | 3
	ZOOM       = VIDEO_TAG | 4
	FOCUS      = VIDEO_TAG | 5
	BRIGHT     = VIDEO_TAG | 6
	GAIN       = VIDEO_TAG | 7
	CONTRAST   = VIDEO_TAG | 8
	SATURATION = VIDEO_TAG | 9

	PARAM_ERR   = "MISSING %s in parameters"
	MIN_ERROR   = "min"
	MAX_ERROR   = "max"
	SCALE_ERROR = "scale"
)

func NewVideo() *Video {

	vid := new(Video)
	vid.RocRobot = roc.NewRocRobot(nil)
	vid.Robot = gobot.NewRobot("video")
	/*	vid.AddFunc(nil, 0, vid.startCannyEdgeAPI, "startCanny")
		vid.AddFunc(nil, 0, vid.stopCannyEdgeAPI, "stopCanny")
		vid.AddFunc(nil, 0, vid.startFaceDetectAPI, "startFaceDetect")
		vid.AddFunc(nil, 0, vid.stopFaceDetectAPI, "stopFaceDetect")*/
	return vid
}

/*
func (vid *Video) startCannyEdge(Min, Max float64) error {

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
	return err
}

func (vid *Video) startCannyEdgeAPI(params map[string]interface{}) interface{} {

	err := vid.CheckAPIParams(params, []types.BasicKind{types.Float64, types.Float64}, "min", "max")
	err = vid.startCannyEdge(params["min"].(float64), params["max"].(float64))
	if err != nil {
		return err.Error()
	}
	return "Image treatment canny activated"
}

func (vid *Video) stopCannyEdge() error {

	b := []byte{roc.DST_L | roc.CMD, CANNY, 0}
	return vid.Send(b)
}

func (vid *Video) stopCannyEdgeAPI(params map[string]interface{}) interface{} {

	err := vid.stopCannyEdge()
	if err != nil {
		return err.Error()
	}
	return "Image treatment canny desactivated"
}

func (vid *Video) startFaceDetect(Scale float64) error {

	scale, err := misc.EncodeBytes(Scale)
	if err != nil {
		log.Println("Error setting minimum to start canny", err.Error())
		return errors.New("Error setting minimum to start canny. " + err.Error())
	}
	b := append([]byte{roc.DST_L | roc.CMD, FACE}, scale...)
	return vid.Send(b)
}

func (vid *Video) startFaceDetectAPI(params map[string]interface{}) interface{} {

	err := vid.CheckAPIParams(params, []types.BasicKind{types.Float64}, "scale")
	if err != nil {
		log.Println(err.Error())
		return err.Error()
	}
	err = vid.startFaceDetect(params["scale"].(float64))
	if err != nil {
		return err.Error()
	}
	return "Face detection activated"
}

func (vid *Video) stopFaceDetect() error {

	b := []byte{roc.DST_L | roc.CMD, FACE, 0}
	return vid.Send(b)
}

func (vid *Video) stopFaceDetectAPI(params map[string]interface{}) interface{} {

	err := vid.stopFaceDetect()
	if err != nil {
		return err.Error()
	}
	return "Face detection stopped"
}
*/
