/**
 * Title: 			Labo3 - Election
 * File:			administrator.go
 * Date:			18.12.12
 * Authors:			Le Guillou Benjamin, Reis de Carvalho Luca
 *
 * Description:		File containing the network interface
 */

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

}

type Echo struct {

}

func ClientWriter(remoteId uint,buf bytes.Buffer)bool {

	var remoteAddress = config.GetAdressById(remoteId)
	buffer := make([]byte, 1024)
	conn, err := net.DialUDP("udp", nil,remoteAddress)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	_,err = conn.Write(buf.Bytes())

	deadline := time.Now().Add(config.GetWaitingAckdelay())
	err = conn.SetReadDeadline(deadline)

	var ack acknowledge

	for {

		n ,_,err := conn.ReadFromUDP(buffer)

		if err != nil {

			return false
		}

		if err := gob.NewDecoder(bytes.NewReader(buffer[:n])).Decode(&ack); err != nil {
			return false
		}
		return true
	}
}

func ClientReader(localId uint, msgChannel chan Message) {

	var address = config.GetAdressById(localId)

	conn, err := net.ListenUDP("udp", address)

	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	decrypt(conn,msgChannel)

}

func decrypt(conn *net.UDPConn ,msgChannel chan Message) {

	buf := make([]byte, 1024)
	for {

		var result Message
		var echo Echo
		n, ip, err := conn.ReadFrom(buf)
		if err != nil {
			log.Fatal(err)
		}
		var buffer bytes.Buffer

		if err := gob.NewEncoder(&buffer).Encode(acknowledge{}); err != nil {
		fmt.Println(err)
		}

		if err := gob.NewDecoder(bytes.NewReader(buf[:n])).Decode(&result); err == nil {
			_,err = conn.WriteTo(buffer.Bytes(),ip)
			msgChannel <- result
		}

		if err := gob.NewDecoder(bytes.NewReader(buf[:n])).Decode(&echo); err == nil {

			_,err = conn.WriteTo(buffer.Bytes(),ip)
		}
	}

}
