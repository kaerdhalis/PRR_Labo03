/**
 * Title: 			Labo2 - Mutual exclusion
 * File:			administrator.go
 * Date:			18.12.12
 * Authors:			Le Guillou Benjamin, Reis de Carvalho Luca
 *
 * Description:		File containing the administrator implementing the Chang-Roberts algorithm
 */

package administrator

import (
	"../config"
	"../network"
	"bytes"
	"encoding/gob"
	"fmt"
	"time"
)

const(
	ANNOUNCE = 0
	RESULT = 1

	ANN = false
	RES = true
)

var(
	id uint
	aptitude uint
	state int
	electedID  =-1
)


//main loop of the algorithm
func Run(elected chan uint, idProc uint,electionLaunch chan bool){

	//channel used to fetch the messages
	var msgChannel =make(chan network.Message)

	id = idProc
	aptitude = config.GetAptById(id)

	//launch the server side of the application
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
			fmt.Print("elected is: ")
			fmt.Println(electedID)
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

	fmt.Print("received announce: ")
	fmt.Println(msg)
	var aptList = msg.List
	if len(aptList)==0{
		return
	}
	for _,proc:= range aptList{

		if proc.Id ==id && proc.Apt ==aptitude{

			var maxApt uint =0
			for _,proc:= range aptList {

				if proc.Apt >= maxApt{
					maxApt = proc.Apt
					electedID = int(proc.Id)
				}
			}
			var aptList = []network.AptList{{id,aptitude}}
			sendResult(electedID,aptList)
			state = RESULT

			return
		}
	}

	aptList = append(aptList,network.AptList{Id:id,Apt:aptitude})
	sendAnnounce(aptList)
	state = ANNOUNCE
}

func resultHandle(msg  network.Message){

	fmt.Print("received result: ")
	fmt.Println(msg)
	var aptList = msg.List
	for _,proc:= range aptList  {

		if proc.Id ==id && proc.Apt ==aptitude{

			return
		}
	}

	if state == RESULT && electedID !=msg.Elected{
		var aptList = []network.AptList{{id,aptitude}}
		sendAnnounce(aptList)
		state = ANNOUNCE
	} else if state == ANNOUNCE{
		electedID =msg.Elected
		aptList = append(aptList,network.AptList{Id:id,Apt:aptitude})
		sendResult(electedID,aptList)
		state = RESULT
	}
}


func sendMessage(buf bytes.Buffer){

	var remoteId =id+1
	time.Sleep(config.GetTransmitDelay())
	for  {

		remoteId%=config.GetNumberOfProc()
		if network.ClientWriter(remoteId,buf)==false{
			remoteId++
		} else {
			return
		}
	}


}


func sendResult(elected int,aptList [] network.AptList){

	fmt.Println("send result")
	var buf bytes.Buffer

	if err := gob.NewEncoder(&buf).Encode(network.Message{MsgType:true,List: aptList, Elected: elected}); err != nil {
		fmt.Println(err)
	}
	sendMessage(buf)

}

func sendAnnounce(aptList [] network.AptList){

	fmt.Println("send announce")
	var buf bytes.Buffer


	if err := gob.NewEncoder(&buf).Encode(network.Message{MsgType:false,List: aptList, Elected: -1}); err != nil {
		fmt.Println(err)
	}
	sendMessage(buf)
}

func CheckOnelected(electedId uint)bool{

	var buf bytes.Buffer

	if err := gob.NewEncoder(&buf).Encode(network.Echo{}); err != nil {
		fmt.Println(err)
	}

	return network.ClientWriter(electedId,buf)

}


