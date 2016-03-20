package linker

import (
	"fmt"
	"net"
)

//Better error gestion panick and defer
func handleConnection(conn net.Conn, in, out chan byte) {

	go sender(conn, out)
	go reader(conn, in)
}

//sender wait for data on channel and send them over the network
//TODO bufferies and optimise
func sender(conn net.Conn, out chan byte) {

	buff := make([]byte, 1)
	for {
		select {
		case b := <-out:
			buff[0] = b
			fmt.Println(buff)
			_, err := conn.Write(buff)
			if err != nil {
				fmt.Print("Error : " + err.Error() + "\n")
			}
		}
	}
}

//reader wait for data to arrive and send them on the channel
//TODO optimise
func reader(conn net.Conn, in chan byte) {

	buff := make([]byte, 8)
	for {
		l, err := conn.Read(buff)
		if l > 0 {
			fmt.Println(buff)
			for i := 0; i < l; i++ {
				in <- buff[i]
			}
		}
		if err != nil {
			fmt.Print("Error : " + err.Error() + "\n")
			conn.Close()
			return
		}
	}
}
