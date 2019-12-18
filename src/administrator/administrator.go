package administrator

import (
	"../config"
	"../network"
	"bytes"
	"encoding/gob"
	"fmt"
)



const(
	ANNOUNCE = 0
	RESULT = 1
	NO =-1

	ANN = false
	RES = true
)

var(

	id uint
	aptitude uint
	state int
	electedID  =-1
)


func Run(elected chan uint, idProc uint,electionLaunch chan bool){

	var msgChannel =make(chan network.Message)

	id = idProc
	aptitude = config.GetAptById(id)

	go network.ClientReader(id,msgChannel)

	for   {

		select {

			case <-electionLaunch:
				electionRequest()

			case message := <- msgChannel:
				if message.MsgType == RES{

					resultHandle(message)
				}else {
					announceHandle(message)


				}


		}
		if state == RESULT {
			elected<- uint(electedID)
		}
	}
}

func electionRequest(){

	var aptlist = []network.AptList{{id,aptitude}}
	sendAnnounce(aptlist)
	state = ANNOUNCE

}

func announceHandle(msg  network.Message){
	fmt.Println("received announce")
	fmt.Println(msg)
	var aptList = msg.List
	for _,proc:= range aptList{

		if proc.Id ==id && proc.Apt ==aptitude{

			var maxApt uint =0
			for i,proc:= range aptList {

				if proc.Apt > maxApt{
					maxApt = proc.Apt
					electedID = i
				}
			}
			var aptList = []network.AptList{{id,aptitude}}
			sendResult(electedID,aptList)
			state = RESULT

			return
		}
	}

	sendAnnounce(aptList)
	state = ANNOUNCE
}

func resultHandle(msg  network.Message){

	fmt.Println("received result")
	fmt.Println(msg)
	var aptList = msg.List
	for _,proc:= range aptList  {

		if proc.Id ==id && proc.Apt ==aptitude{

			state = NO
			return
		}
	}

	if state == RESULT && electedID !=msg.Elected{

		sendAnnounce(aptList)
		state = ANNOUNCE
	} else if state == ANNOUNCE{
		electedID =msg.Elected
		sendResult(electedID,aptList)
		state = RESULT
	}

	sendAnnounce(aptList)
	state = ANNOUNCE
}


func sendMessage(buf bytes.Buffer){

	var remoteId =id+1
	//time.Sleep(config.GetTransmitDelay())
	for  {
		remoteId%=config.GetNumberOfProc()
		if network.ClientWriter(remoteId,buf)==false{
			fmt.Println("no connection")
			remoteId++
		} else {
			return
		}
	}


}


func sendResult(elected int,aptList [] network.AptList){

	fmt.Println("send result")
	var buf bytes.Buffer

	if err := gob.NewEncoder(&buf).Encode(network.Message{List: aptList, Elected: elected}); err != nil {
		fmt.Println(err)
	}
	sendMessage(buf)

}

func sendAnnounce(aptList [] network.AptList){

	fmt.Println("send announce")
	var buf bytes.Buffer


	if err := gob.NewEncoder(&buf).Encode(network.Message{List: aptList, Elected: -1}); err != nil {
		fmt.Println(err)
	}
	sendMessage(buf)
}


