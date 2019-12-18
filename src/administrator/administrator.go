package administrator

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"../network"
	"../config"
	"time"
)

const(
	ANNOUNCE = 0
	RESULT = 1
	NO =-1

	ANN = false
	RES = true
)

var(

	nbProcess int
	id uint
	aptitude uint
	next int
	state int
	elected int
	aptList [] network.Aptitude
)


func Run(election chan uint, echo chan network.Echo, idProc uint){

	var msgChannel chan network.Message

	var ackChannel chan network.Acknowledge

	id = idProc

	go network.ClientReader(id,msgChannel,ackChannel,echo)

	for   {

		select {

			case <-election:
				electionRequest(ackChannel)

			case message := <- msgChannel:
				if message.MsgType == RES{

					resultHandle(message,ackChannel)
				}else {
					announceHandle(message,ackChannel)


				}


		}
		if state == RESULT {
			election<- uint(elected)
		}
	}
}

func electionRequest(ackChannel chan network.Acknowledge){

	sendAnnounce(ackChannel)
	state = ANNOUNCE

}

func announceHandle(msg  network.Message,ackChannel chan network.Acknowledge){

	sendAck(id,msg.SenderId)

	aptList = msg.List
	for _,apt:= range aptList  {

		if apt.Id ==id && apt.Apt ==aptitude{

			var apti uint =0
			for i,apt:= range aptList {

				if apt.Apt > apti{
					apti = apt.Apt
					elected = i
				}
			}

			sendResult(elected,ackChannel)
			state = RESULT
			return
		}
	}

	sendAnnounce(ackChannel)
	state = ANNOUNCE
}

func resultHandle(msg  network.Message,ackChannel chan network.Acknowledge){

	sendAck(id,msg.SenderId)

	for _,apt:= range aptList  {

		if apt.Id ==id && apt.Apt ==aptitude{

			state = NO
			return
		}
	}

	if state == RESULT && elected !=msg.Elected{

		sendAnnounce(ackChannel)
		state = ANNOUNCE
	} else if state == ANNOUNCE{
		elected =msg.Elected
		sendResult(elected,ackChannel)
		state = RESULT
	}

	sendAnnounce(ackChannel)
	state = ANNOUNCE
}


func sendMessage(localId uint, remoteId uint,buf bytes.Buffer,ackChannel chan network.Acknowledge){

	time.Sleep(time.Duration(config.GetTransmitDelay()) * time.Second)

	network.ClientWriter(localId,remoteId,buf)

	select {

	case <-ackChannel:

	case <- time.After(time.Duration(2*config.GetTransmitDelay())*time.Second):

		sendMessage(localId,remoteId+1,buf,ackChannel)
	}


}


func sendResult(elected int,ackChannel chan network.Acknowledge){


	var buf bytes.Buffer

	aptList = append(aptList,network.Aptitude{id,aptitude})

	if err := gob.NewEncoder(&buf).Encode(network.Message{List: aptList, Elected: elected}); err != nil {
		fmt.Println(err)
	}
	sendMessage(id,id+1,buf,ackChannel)

}

func sendAnnounce(ackChannel chan network.Acknowledge){


	var buf bytes.Buffer

	aptList = append(aptList,network.Aptitude{id,aptitude})

	if err := gob.NewEncoder(&buf).Encode(network.Message{List: aptList, Elected: -1}); err != nil {
		fmt.Println(err)
	}
	sendMessage(id,id+1,buf,ackChannel)
}

func sendAck(localId uint,remoteId uint){



	var buf bytes.Buffer

	if err := gob.NewEncoder(&buf).Encode(network.Acknowledge{}); err != nil {
		fmt.Println(err)
	}

	time.Sleep(time.Duration(config.GetTransmitDelay()) * time.Second)

	network.ClientWriter(localId,remoteId,buf)


}

func pingElected(elected uint){

	var buf bytes.Buffer

	if err := gob.NewEncoder(&buf).Encode(network.Acknowledge{}); err != nil {
		fmt.Println(err)
	}

	network.ClientWriter(id,elected,buf)


}



func checkAllProcessAreReady(){

	for i:=0 ;i<int(config.GetNumberOfProc()) ;i++ {

		if uint(i) != id {
			network.PingAdress(config.GetAdressById(uint(i)), uint(i))
		}
	}

	fmt.Println("All Process are Ready")
}