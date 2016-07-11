package gpsd

import (
	"bufio"
	"fmt"
	"github.com/hybridgroup/gobot"
	"log"
	"net"
	"time"
)

var _ gobot.Adaptor = (*GpsdAdaptor)(nil)
var _ GpsdReader = (*GpsdAdaptor)(nil)
var _ GpsdWriter = (*GpsdAdaptor)(nil)

// Be sure that gpsd daemon is started
// service --status-all

const (
	START     = "?WATCH={\"enable\":true,\"json\":true}"
	STOP      = "?WATCH={\"enable\":false}"
	GPSD_PORT = ":2947"
	TPV       = "TPV"
	ERROR     = "Error"
)

type GpsdAdaptor struct {
	name   string
	ip     string
	conn   net.Conn
	reader *bufio.Reader
}

type GpsdReader interface {
	gobot.Adaptor
	GpsdRead() (string, error)
}

type GpsdWriter interface {
	gobot.Adaptor
	GpsdWrite(string) error
}

func NewGpsdAdaptor(name string, ip string) *GpsdAdaptor {

	gpsd := &GpsdAdaptor{
		name:   name,
		ip:     ip,
		reader: nil,
	}

	return gpsd
}

func (gpsd *GpsdAdaptor) Name() string {
	return gpsd.name
}

func (gpsd *GpsdAdaptor) Connect() (errs []error) {

	var err error

	gpsd.conn, err = net.Dial("tcp", GPSD_PORT)
	if err != nil {
		log.Println("Cannot connect to gpsd", err.Error())
		return []error{err}
	}
	gpsd.reader = bufio.NewReader(gpsd.conn)
	log.Printf("gpsd connected\n")
	return nil
}

// Disconnect closes the io connection to the board
func (gpsd *GpsdAdaptor) Disconnect() (err error) {
	log.Printf("Terminating gpsd connection\n")
	fmt.Fprintf(gpsd.conn, "?WATCH={\"enable\":false}")
	gpsd.conn.Close()
	return nil
}

// Finalize terminates the gpsd connection
func (gpsd *GpsdAdaptor) Finalize() []error {
	err := gpsd.Disconnect()
	if err != nil {
		return []error{err}
	}
	return nil
}

func (gpsd *GpsdAdaptor) GpsdRead() (string, error) {

	gpsd.conn.SetDeadline(time.Now().Add(time.Second * 4))
	line, err := gpsd.reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return line, nil
}

func (gpsd *GpsdAdaptor) GpsdWrite(s string) error {

	_, err := fmt.Fprintf(gpsd.conn, s)
	return err
}
