package network

import (
	"../config"
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"time"
)
type AptList struct {
	Id uint
	Apt uint
}
type Message struct{
	MsgType bool
	List [] AptList
	Elected int
}

type acknowledge struct {
	Ack bool
}

func ClientWriter(remoteId uint,buf bytes.Buffer)bool {

	var remoteAddress = config.GetAdressById(remoteId)
	buffer := make([]byte, 1024)
	conn, err := net.Dial("udp", remoteAddress.String())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	_,err = conn.Write(buf.Bytes())

	deadline := time.Now().Add(config.GetWaitingAckdelay())
	err = conn.SetReadDeadline(deadline)

	var ack acknowledge

	for {

		n ,err := conn.Read(buffer)

		if err != nil {
		fmt.Println(err)
		return false
		}

		if err := gob.NewDecoder(bytes.NewReader(buffer[:n])).Decode(&ack); err != nil {
			fmt.Println("test")
			return false
		}
		fmt.Println(ack)
		return true
	}
}

func ClientReader(localId uint, msgChannel chan Message) {

	var address = config.GetAdressById(localId)

	conn, err := net.ListenPacket("udp", address.String())

	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	decrypt(conn,msgChannel)

}

func decrypt(conn net.PacketConn ,msgChannel chan Message) {

	buf := make([]byte, 1024)
	for {

		var result Message
		n, ip, err := conn.ReadFrom(buf)
		if err != nil {
			log.Fatal(err)
		}

		if err := gob.NewDecoder(bytes.NewReader(buf[:n])).Decode(&result); err == nil {

			msgChannel <- result
			var buffer bytes.Buffer

			if err := gob.NewEncoder(&buffer).Encode(acknowledge{true}); err != nil {
				fmt.Println(err)
			}
			_,err := conn.WriteTo(buffer.Bytes(),ip)
			fmt.Println(err)
		}
	}

}
