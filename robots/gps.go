package robots

import (
	"log"
	"math"

	"errors"

	"github.com/Remote-Oculus-Controller/R.O.C-CONTROLS"
	"github.com/Remote-Oculus-Controller/R.O.C-CONTROLS/gpsd"
	"github.com/Remote-Oculus-Controller/proto"
	"github.com/hybridgroup/gobot"
	"github.com/larsth/go-gpsdjson"
)

const (
	GPS_TAG   = 0x10
	GET_COORD = GPS_TAG | 1
	TOOGLE    = GPS_TAG | 2
	H_DCV     = uint32(rocproto.Packet_DATA)<<uint32(rocproto.Packet_SHIFT_SEND) | uint32(rocproto.Packet_VIDEO_CLIENT)
)

// Simulating Gps, change to real gps
type Gps struct {
	*roc.Robot
	*gpsd.GpsdDriver
	adaptor *gpsd.GpsdAdaptor
	coord   rocproto.Coord

	xoff, yoff, orioff, dir float64
}

func NewGPS() *Gps {

	gps := new(Gps)
	gps.Robot = roc.NewRocRobot(nil)
	gps.adaptor = gpsd.NewGpsdAdaptor("gpsd", gpsd.GPSD_PORT)
	gps.GpsdDriver = gpsd.NewGpsdDriver(gps.adaptor, "gpsd_driver")

	gps.Robot.Robot = gobot.NewRobot("gps",
		[]gobot.Connection{gps.adaptor},
		[]gobot.Device{gps.GpsdDriver},
	)

	gobot.On(gps.Event("TPV"), func(data interface{}) {
		var err error

		tpv, ok := data.(gpsdjson.TPV)
		if !ok {
			log.Println("Event TPV, didn't reveice a TPV message gps.go")
		}
		p := &rocproto.Packet{
			ID:     GPS_TAG,
			Header: H_DCV,
			Coord: &rocproto.Coord{
				Lat:  tpv.Lat + gps.xoff,
				Long: tpv.Lon + gps.yoff,
				Ori:  gps.dir * 180 / math.Pi,
			},
		}
		gps.coord.Lat = p.Coord.Lat
		gps.coord.Long = p.Coord.Long
		if err != nil {
			err = errors.New("Couldn't pack Gps coor into packet: " + err.Error())
			log.Println(err)
			//gps.Send(err)
			return
		}
		gps.Send(p)
	})

	gps.AddFunc(gps.tooglePause, TOOGLE, gps.tooglePauseAPI, "toogle")
	gps.AddFunc(gps.getCoordByte, GET_COORD, gps.getCoordApi, "getCoord")
	gps.AddFunc(nil, 0, gps.sim, "sim")
	return gps
}

func (gps *Gps) getCoord() (float64, float64) {
	return gps.coord.Lat, gps.coord.Long
}

func (gps *Gps) getCoordByte(r *rocproto.Packet) error {

	s := uint32(r.Header) & (uint32(rocproto.Packet_MASK_DEST) << uint32(rocproto.Packet_SHIFT_SEND))
	r.Header = (uint32(rocproto.Packet_DATA) << uint32(rocproto.Packet_SHIFT_TYPE)) | s>>uint32(rocproto.Packet_SHIFT_SEND)
	r.Coord = &gps.coord
	return gps.Send(r)
}

func (gps *Gps) getCoordApi(params map[string]interface{}) interface{} {
	return gps.coord
}

func (gps *Gps) tooglePause(p *rocproto.Packet) error {
	gps.TooglePause()
	return nil
}

func (gps *Gps) tooglePauseAPI(params map[string]interface{}) interface{} {
	gps.TooglePause()
	return "Gps state toogled"
}

func (gps *Gps) sim(params map[string]interface{}) interface{} {

	n := params["mv"].(rocproto.Mv)
	if n.Angle < math.Pi-0.0001 || n.Angle > math.Pi+0.0001 {
		gps.dir -= n.Angle / 180
	}
	if gps.dir > 2*math.Pi {
		gps.dir = gps.dir - 2*math.Pi
	}
	if gps.dir < 0 {
		gps.dir = 2*math.Pi - gps.dir
	}
	gps.yoff += 0.000001 * math.Sin(gps.dir) * (n.Speed / 100)
	gps.xoff += 0.000001 * math.Cos(gps.dir) * (n.Speed / 100)
	return nil
}

func (gps *Gps) angleDir(a float64) *rocproto.Coord {

	coord := &rocproto.Coord{Lat: gps.xoff + math.Cos(a), Long: gps.yoff + math.Sin(a)}
	return coord
}
