package robots

import (
	"encoding/json"
	"fmt"
	"github.com/Happykat/R.O.C-CONTROLS"
	"github.com/Happykat/R.O.C-CONTROLS/protoext"
	"github.com/hybridgroup/gobot"
	"go/types"
	"log"
	"time"
)

const (
	GPS_TAG   = 0xB0
	SET_COORD = GPS_TAG | 1
	SET_DEST  = GPS_TAG | 3
	GET_COORD = GPS_TAG | 2
)

// Simulating Gps, change to real gps
type Gps struct {
	*roc.RocRobot
	coord Coord
}

func NewGPS() *Gps {

	gps := new(Gps)
	gps.RocRobot = roc.NewRocRobot(nil)
	work := func() {
		gobot.Every(time.Second, func() {
			fmt.Printf("GpsCoorinates %v\n", gps.coord)
		})
	}
	gps.Robot = gobot.NewRobot("gps",
		[]gobot.Connection{},
		[]gobot.Device{},
		work)

	gps.AddFunc(gps.setCoordByte, SET_COORD, gps.setCoordApi, "setCoord")
	gps.AddFunc(gps.getCoordByte, GET_COORD, gps.getCoordApi, "getCoord")
	return gps
}

func (gps *Gps) setCoord(lat, long float32) {
	gps.coord.Lat = lat
	gps.coord.Long = long
}

func (gps *Gps) setCoordByte(data *roc.Packet) error {

	err := protoext.UnpackAny(data.GetPayload(), &gps.coord)
	if err != nil {
		log.Println("Impossible conversion Message is not a Coordinate")
		return err
	}
	return nil
}

func (gps *Gps) setCoordApi(params map[string]interface{}) interface{} {

	err := gps.CheckAPIParams(params, []types.BasicKind{types.Float64, types.Float64}, "lat", "lon")
	if err != nil {
		log.Println(err.Error())
		return err.Error()
	}
	gps.setCoord(params["lat"].(float32), params["long"].(float32))
	return "Gps coord changed"
}

/*
func (gps *Gps) setDest(lat, long float32) error {

	pos := roc.Position{}
	pos.Lat = lat
	pos.Long = long
	p := roc.Packet{}
	p.ID = GPS_TAG
	p.Header = uint32(roc.Packet_COMMAND) << uint32(roc.Packet_SHIFT) | uint32(roc.Packet_VIDEO_CLIENT)
	p.GetPayload() =
	return gps.Send(roc.Position{})
}

func (gps *Gps) setDestByte(b []byte) error {

	lat, err := misc.DecodeFloat32(b[:3])
	if err != nil {
		log.Printf(err.Error())
		return err
	}
	long, err := misc.DecodeFloat32(b[4:7])
	if err != nil {
		log.Printf(err.Error())
		return err
	}
	return gps.setDest(lat, long)
}

func (gps *Gps) setDestApi(params map[string]interface{}) interface{} {

	err := gps.CheckAPIParams(params, []types.BasicKind{types.Float64, types.Float64}, "lat", "lon")
	if err != nil {
		log.Println(err.Error())
		return err.Error()
	}
	err = gps.setDest(params["lat"].(float32), params["long"].(float32))
	if err != nil {
		return err.Error()
	}
	return "New Destination !"
}
*/

func (gps *Gps) getCoord() (float32, float32) {
	return gps.coord.Lat, gps.coord.Long
}

func (gps *Gps) getCoordByte(r *roc.Packet) error {

	var err error

	fmt.Printf("getCoordinates")
	s := uint32(r.Header) & (uint32(roc.Packet_MASK_DEST) << uint32(roc.Packet_SHIFT_SENT))
	r.Header = (uint32(roc.Packet_DATA) << uint32(roc.Packet_SHIFT)) | s>>uint32(roc.Packet_SHIFT_SENT)
	fmt.Printf("new Header %b", s)
	r.Payload, err = protoext.PackAny(&gps.coord)
	if err != nil {
		return err
	}
	return gps.Send(r)
}

func (gps *Gps) getCoordApi(params map[string]interface{}) interface{} {
	b, err := json.Marshal(gps.coord)
	if err != nil {
		log.Println(err.Error())
		return err.Error()
	}
	return b
}
