/**
 * Title: 			Labo3 - Election
 * File:			processus.go
 * Date:			18.12.12
 * Authors:			Le Guillou Benjamin, Reis de Carvalho Luca
 *
 * Description:		File containing the applicative task asking periodicaly at the administrator the id of the elected process
 */
package main

import (
	"../administrator"
	"../config"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)
var electedId uint

func main() {

	//read the id of the process passed in argument
	args := os.Args[1:]

	if len(args)!=1{
		log.Fatal("Number of arguments invalid, you need to pass the id of the Process")
	}
	id,_ := strconv.Atoi(args[0])

	//channel used to fetch the id of the elected process
	elected := make(chan uint)

	//channel used to start a new election
	launchElection := make(chan bool)

	//get the global configuration of the application
	config.SetConfiguration()

	//launch the administrator and the network
	go administrator.Run(elected, uint(id),launchElection)

	go getElectedProcess(elected)

	fmt.Println("launch election")
	launchElection <-true

	// main loop, ask periodically for the elected process and launch a new election if the elected is down
	for {

		time.Sleep(config.GetArtificialDelay())
		if !administrator.CheckOnelected(electedId) {

			fmt.Println("launch new election")
			launchElection <- true
		}
	}
}

func getElectedProcess(elected chan uint){


	for   {

		electedId=<-elected
	}

}

