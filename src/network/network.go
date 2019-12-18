package network

import (
	"../config"
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"net"
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

	//var localAddress = config.GetAdressById(localId)
	var remoteAddress = config.GetAdressById(remoteId)

	conn, err := net.DialUDP("udp",nil, remoteAddress)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	_, err  = buf.WriteTo(conn)
}

func ClientReader(localId uint, msgChannel chan Message,ackChannel chan Acknowledge,echo chan Echo) {
	// error testing suppressed to compact listing on slides


	var address = config.GetAdressById(localId)
	fmt.Println(address.String())
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
		var result Message
		var echo Echo
		n, _, err := conn.ReadFromUDP(buf) // n,addr, err := p.ReadFrom(buf)
		if err != nil {
			log.Fatal(err)
		}

		if err := gob.NewDecoder(bytes.NewReader(buf[:n])).Decode(&ack); err == nil {
			fmt.Print(ack)
			fmt.Println("testack")

			ackChannel<-ack

		}else if err := gob.NewDecoder(bytes.NewReader(buf[:n])).Decode(&result); err == nil {

			fmt.Print(result)
			fmt.Println("testresult")
			msgChannel <- result

		}else if err := gob.NewDecoder(bytes.NewReader(buf[:n])).Decode(&echo); err == nil {

		}
	}
}

