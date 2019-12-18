package network

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"../config"
	"time"
)


type Aptitude struct {
	 Id uint
	 Apt uint
}

type Acknowledge struct {

}

type Echo struct{

}

type Message struct{
	MsgType bool
	List [] Aptitude
	Elected int
	SenderId uint
}

func ClientWriter(localId uint,remoteId uint,buf bytes.Buffer) {

	var localAddress = config.GetAdressById(localId)
	var remoteAddress = config.GetAdressById(remoteId)

	conn, err := net.DialUDP("udp",localAddress, remoteAddress)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	_, err  = buf.WriteTo(conn)
}

func ClientReader(localId uint, msgChannel chan Message,ackChannel chan Acknowledge,echo chan Echo) {
	// error testing suppressed to compact listing on slides

	var address = config.GetAdressById(localId)
	conn, err := net.ListenUDP("udp", address)

	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	decrypt(conn,msgChannel,ackChannel,echo)

}

func decrypt(conn *net.UDPConn ,msgChannel chan Message, ackChannel chan Acknowledge,echo chan Echo){


	buf := make([]byte, 1024)
	for {

		var ack Acknowledge
		var msg Message
		var result Message
		var echo Echo
		n, _, err := conn.ReadFromUDP(buf) // n,addr, err := p.ReadFrom(buf)
		if err != nil {
			log.Fatal(err)
		}

		if err := gob.NewDecoder(bytes.NewReader(buf[:n])).Decode(&ack); err == nil {

			ackChannel<-ack

		}else if err := gob.NewDecoder(bytes.NewReader(buf[:n])).Decode(&msg); err != nil {

			msgChannel <- result

		}else if err := gob.NewDecoder(bytes.NewReader(buf[:n])).Decode(&echo); err != nil {

		}
	}
}


// Create a connection with a process to check if its ready
func PingAdress(address *net.UDPAddr,id uint) {

	timeout := 1 * time.Second
	for {

		conn, err := net.DialTimeout("tcp", address.String(), timeout)
		if err != nil {

		} else {

			fmt.Printf("Processus %d is Up and Ready\n",id)
			conn.Close()
			break
		}
	}
}

