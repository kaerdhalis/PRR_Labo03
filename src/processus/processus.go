/**
 * Title: 			Labo2 - Mutual exclusion
 * File:			processus.go
 * Date:			18.12.12
 * Authors:			Le Guillou Benjamin, Reis de Carvalho Luca
 *
 * Description:		File containing the client side of the process. It can read the inputs of the users and read or
 *                  modify the shared value.
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


func main() {

	//read the id of the process passed in argument
	args := os.Args[1:]

	if len(args)!=1{
		log.Fatal("Number of arguments invalid, you need to pass the id of the Process")
	}
	id,_ := strconv.Atoi(args[0])

	elected := make(chan uint)

	launchElection := make(chan bool)

	//get the global configuration of the application
	config.SetConfiguration()

	go administrator.Run(elected, uint(id),launchElection)

	go getElectedProcess(elected)
	fmt.Println("launch election")
	launchElection <-true


	for{

		time.Sleep(500* config.GetArtificialDelay())
				fmt.Println("launch election")
				launchElection<- true

	}
}

func getElectedProcess(elected chan uint){
	for   {
		fmt.Printf("processus elu = %d\n",<-elected)

	}

}


