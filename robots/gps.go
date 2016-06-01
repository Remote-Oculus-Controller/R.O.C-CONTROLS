package robots

import (
	"github.com/hybridgroup/gobot"
	"github.com/Happykat/R.O.C-CONTROLS"
	"github.com/Happykat/R.O.C-CONTROLS/misc"
	"log"
	"fmt"
	"encoding/json"
)

const (
	GPS_TAG = 0xB0
	SET_COORD = GPS_TAG | 1
	SET_DEST = GPS_TAG | 2
	GET_COORD = GPS_TAG
	LAT_ERR = "Missing lat float32 in parameters"
	LONG_ERR = "Missing long float32 in paramters"
)

// Simulating Gps, change to real gps
type Gps struct {
	*roc.RocRobot
	Coord
}

type Coord struct {
	Lat  float32
	Long float32
	Ori  float32
}

func NewGPS() *Gps {

	gps := new(Gps)
	gps.RocRobot = roc.NewRocRobot(nil)
	work := func(){
		/*
		gobot.Every(300*time.Millisecond, func(){
			gps.sendCoord(nil)
		})
		*/
	}
	fmt.Printf("%+v", gps)
	gps.Robot = gobot.NewRobot("gps",
		[]gobot.Connection{},
		[]gobot.Device{},
		work)

	gps.AddFunc(gps.setCoordByte, SET_COORD, gps.setCoordApi, "setCoord")
	gps.AddFunc(gps.setDestByte, SET_DEST, gps.setDestApi, "setDest")
	gps.AddFunc(gps.sendCoord, GET_COORD, gps.getCoordApi, "getCoord")
	return gps
}

func (gps *Gps) setCoord(lat, long float32) {
	gps.Lat = lat
	gps.Long = long
}

func (gps *Gps) setCoordByte(b []byte) error {

	lat, err := misc.DecodeFloat32(b[:3])
	if err != nil {
		log.Printf(err.Error())
		return err
	}
	long, err := misc.DecodeFloat32(b[4:])
	if err != nil {
		log.Printf(err.Error())
		return err
	}
	gps.setCoord(lat, long)
	return nil
}

func (gps *Gps) setCoordApi(params map[string]interface{}) interface{} {

	lat, ok := params["lat"]
	_, assert := lat.(float32)
	if (!ok || assert) {
		log.Println(LAT_ERR)
		return LAT_ERR
	}
	long, ok := params["long"]
	_, assert = long.(float32)
	if (!ok || !assert) {
		log.Println(LONG_ERR)
		return LONG_ERR
	}
	gps.setCoord(lat.(float32), long.(float32))
	return 200
}

func (gps *Gps) setDest(lat, long float32) {
	b, err := misc.EncodeBytes(lat)
	if err != nil {
		fmt.Println("Error setting lattitude for destination", err.Error())
		return
	}
	l, err := misc.EncodeBytes(long)
	if err != nil {
		fmt.Println("Error setting longitude for destination", err.Error())
		return
	}
	b = append([]byte{roc.DST_L | roc.CMD, SET_DEST},b...)
	b = append(b, l...)
	gps.Send(b)
}

func (gps *Gps) setDestByte(b []byte) error {

	lat, err := misc.DecodeFloat32(b[:3])
	if err != nil {
		log.Printf(err.Error())
		return err
	}
	long, err := misc.DecodeFloat32(b[4:])
	if err != nil {
		log.Printf(err.Error())
		return err
	}
	gps.setDest(lat, long)
	return nil
}

func (gps *Gps) setDestApi(params map[string]interface{}) interface{} {

	lat, ok := params["lat"]
	_, assert := lat.(float32)
	if (!ok || !assert) {
		log.Println(LAT_ERR)
		return LAT_ERR
	}
	long, ok := params["long"]
	_, assert = long.(float32)
	if (!ok || !assert) {
		log.Println(LONG_ERR)
		return LONG_ERR
	}
	gps.setDest(lat.(float32), long.(float32))
	return 200
}

func (gps *Gps) getCoord() (float32, float32) {
	return gps.Lat, gps.Long
}

func (gps *Gps) sendCoord([]byte) error {
	b:= gps.getCoordByte()
	gps.Send(b)
	return nil
}

func (gps *Gps) getCoordByte() []byte {
	lat, err := misc.EncodeBytes(gps.Lat)
	misc.CheckError(err, "Encoding latitude", false)
	long, err := misc.EncodeBytes(gps.Lat)
	misc.CheckError(err, "Encoding longitude", false)
	ori, err := misc.EncodeBytes(gps.Ori)
	misc.CheckError(err, "Encoding orientation", false)
	b := append([]byte{roc.DST_RL | roc.DATA, GPS_TAG}, lat...)
	b = append(b, long...)
	b = append(b, ori...)
	return b
}

func (gps *Gps) getCoordApi(params map[string]interface{}) interface{} {
	b, err := json.Marshal(gps.Coord)
	if err != nil {
		log.Println(err.Error())
		return err.Error()
	}
	return b
}