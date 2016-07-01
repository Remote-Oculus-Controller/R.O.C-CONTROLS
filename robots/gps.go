package robots

import (
	"github.com/Happykat/R.O.C-CONTROLS"
	"github.com/Happykat/R.O.C-CONTROLS/gpsd"
	"github.com/Happykat/R.O.C-CONTROLS/protoext"
	"github.com/hybridgroup/gobot"
	"github.com/larsth/go-gpsdjson"
	"log"
)

const (
	GPS_TAG   = 0xB0
	GET_COORD = GPS_TAG | 1
	TOOGLE    = GPS_TAG | 2
	H_DCV     = uint32(roc.Packet_DATA)<<uint32(roc.Packet_SHIFT) | uint32(roc.Packet_VIDEO_CLIENT)
)

// Simulating Gps, change to real gps
type Gps struct {
	*roc.RocRobot
	*gpsd.GpsdDriver
	coord Coord
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
		m := &Coord{
			Lat:  tpv.Lat,
			Long: tpv.Lon,
		}
		p := &roc.Packet{
			ID:     GPS_TAG,
			Header: H_DCV,
		}
		p.Payload, err = protoext.PackAny(m)
		if err != nil {
			log.Println("Couldn't pack Gps coor into packet: ", err.Error())
			return
		}
		gps.Send(p)
	})
	gps.AddFunc(gps.tooglePause, TOOGLE, gps.tooglePauseAPI, "toogle")
	gps.AddFunc(gps.getCoordByte, GET_COORD, gps.getCoordApi, "getCoord")
	return gps
}

func (gps *Gps) getCoord() (float64, float64) {
	return gps.coord.Lat, gps.coord.Long
}

func (gps *Gps) getCoordByte(r *roc.Packet) error {

	var err error

	s := uint32(r.Header) & (uint32(roc.Packet_MASK_DEST) << uint32(roc.Packet_SHIFT_SENT))
	r.Header = (uint32(roc.Packet_DATA) << uint32(roc.Packet_SHIFT)) | s>>uint32(roc.Packet_SHIFT_SENT)
	r.Payload, err = protoext.PackAny(&gps.coord)
	if err != nil {
		return err
	}
	return gps.Send(r)
}

func (gps *Gps) getCoordApi(params map[string]interface{}) interface{} {
	return gps.coord
}

func (gps *Gps) tooglePause(p *roc.Packet) error {
	gps.TooglePause()
	return nil
}

func (gps *Gps) tooglePauseAPI(params map[string]interface{}) interface{} {
	gps.TooglePause()
	return "Gps state toogled"
}
