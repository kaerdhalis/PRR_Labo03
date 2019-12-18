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

	fmt.Println(remoteId)
	var remoteAddress = config.GetAdressById(remoteId)
	buffer := make([]byte, 1024)
	conn, err := net.Dial("udp", remoteAddress.String())
	if err != nil {

		log.Fatal(err)
	}
	defer conn.Close()

	_, err = conn.Write(buf.Bytes())
	if err != nil {

		log.Fatal(err)
	}

	deadline := time.Now().Add(200* time.Millisecond)
	conn.SetReadDeadline(deadline)
	var ack acknowledge
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if e, ok := err.(net.Error); !ok || !e.Timeout() {
				fmt.Println(err)
				fmt.Println(remoteId)
				//	log.Fatal(err)
			}
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
	conn, err := net.ListenPacket("udp",address.String())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Listen on " + address.String())
	defer conn.Close()
	buffer := make([]byte, 1024)
	for {
		n, cliAddr, err := conn.ReadFrom(buffer)
		if err != nil {
			log.Fatal(err)
		}
		var result Message
		if err := gob.NewDecoder(bytes.NewReader(buffer[:n])).Decode(&result); err != nil {
		}



		var buffer bytes.Buffer

		if err := gob.NewEncoder(&buffer).Encode(acknowledge{true}); err != nil {
			fmt.Println(err)
		}
			_,err = conn.WriteTo(buffer.Bytes(),cliAddr)
			fmt.Println(err)

		msgChannel <- result
		}


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
