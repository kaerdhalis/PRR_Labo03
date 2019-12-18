/**
 * Title: 			Labo2 - Mutual exclusion
 * File:			client.go
 * Date:			20.11.12
 * Authors:			Le Guillou Benjamin, Reis de Carvalho Luca
 *
 * Description:		File containing the client side of the process. It can read the inputs of the users and read or
 *                  modify the shared value.
 */
package main

import (
	"../config"
	"../administrator"
	"../network"
	"bytes"
	"encoding/gob"
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

	election := make(chan uint)
	echo := make(chan network.Echo)

	var elected uint

	//get the global configuration of the application
	config.SetConfiguration()

	go administrator.Run(election,echo, uint(id))

	for{

		time.Sleep(time.Duration(config.GetArtificialDelay()))

		var buf bytes.Buffer

		if err := gob.NewEncoder(&buf).Encode(network.Echo{}); err != nil {
			fmt.Println(err)
		}

		network.ClientWriter(uint(id),elected,buf)

		select {

			case <-echo:

			case <- time.After(time.Duration(config.GetArtificialDelay())):

				election<- 1
		}



	}
}


