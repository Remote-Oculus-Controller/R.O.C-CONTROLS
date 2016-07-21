package robots

import (
	"github.com/Happykat/R.O.C-CONTROLS"
	"github.com/Happykat/R.O.C-CONTROLS/gpsd"
	"github.com/Happykat/R.O.C-CONTROLS/rocproto"
	"github.com/hybridgroup/gobot"
	"github.com/larsth/go-gpsdjson"
	"log"
	"math"
)

const (
	GPS_TAG   = 0x10
	GET_COORD = GPS_TAG | 1
	TOOGLE    = GPS_TAG | 2
	H_DCV     = uint32(rocproto.Packet_DATA)<<uint32(rocproto.Packet_SHIFT) | uint32(rocproto.Packet_VIDEO_CLIENT)
)

// Simulating Gps, change to real gps
type Gps struct {
	*roc.RocRobot
	*gpsd.GpsdDriver
	coord rocproto.Coord

	xoff, yoff, orioff, dir float64
}

func NewGPS() *Gps {

	gps := new(Gps)
	gps.RocRobot = roc.NewRocRobot(nil)
	gpsdA := gpsd.NewGpsdAdaptor("gpsd", "")
	gps.GpsdDriver = gpsd.NewGpsdDriver(gpsdA, "gpsd_driver")
	gps.Robot = gobot.NewRobot("gps",
		[]gobot.Connection{gpsdA},
		[]gobot.Device{gps.GpsdDriver},
	)

	gobot.On(gps.Event("TPV"), func(data interface{}) {
		var err error

		tpv, ok := data.(gpsdjson.TPV)
		if !ok {
			log.Println("Event TPV, didn't reveice a TPV message gps.go")
		}
		m := &rocproto.Coord{
			Lat:  tpv.Lat + gps.xoff,
			Long: tpv.Lon + gps.yoff,
			Ori:  gps.dir * 180 / math.Pi,
		}
		p := &rocproto.Packet{
			ID:     GPS_TAG,
			Header: H_DCV,
		}
		p.Payload, err = rocproto.PackAny(m)
		if err != nil {
			log.Println("Couldn't pack Gps coor into packet: ", err.Error())
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

	var err error

	s := uint32(r.Header) & (uint32(rocproto.Packet_MASK_DEST) << uint32(rocproto.SHIFT_SEND))
	r.Header = (uint32(rocproto.Packet_DATA) << uint32(rocproto.Packet_SHIFT)) | s>>uint32(rocproto.SHIFT_SEND)
	r.Payload, err = rocproto.PackAny(&gps.coord)
	if err != nil {
		return err
	}
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

	n := params["mv"].(rocproto.Mouv)
	if n.Angle < math.Pi-0.0001 || n.Angle > math.Pi+0.0001 {
		gps.dir -= n.Angle / 180
	} else {
		n.Speed = -n.Speed
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

	coord := &rocproto.Coord{}
	coord.Lat = gps.xoff + math.Cos(a)
	coord.Long = gps.yoff + math.Sin(a)
	return coord
}
