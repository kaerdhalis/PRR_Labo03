/**
 * Title: 			Labo2 - Mutual exclusion
 * File:			client.go
 * Date:			20.11.12
 * Authors:			Le Guillou Benjamin, Reis de Carvalho Luca
 *
 * Description:		File containing the configurationof the application. It contains the adress of every process and other
 *                  values like the waiting time in the critical section
 */
package config

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"
)

// contains all the values of configuration
var config configuration

//struct representing the adress from the json file
type Process struct {

	Ip string
	Port uint
	Aptitude uint

}

// Configuration struct represents all configurations from the json file
type configuration struct {
	NumberOfProcesses uint
	Process []Process
	ArtificialDelay uint
	TransmitDelay uint
}


func GetAdressById(id uint) *net.UDPAddr{

	id = id%config.NumberOfProcesses
	var localAdrr = new(net.UDPAddr)

	localAdrr.IP = net.ParseIP(config.Process[id].Ip)
	localAdrr.Port =int(config.Process[id].Port)

	return localAdrr
}

func GetAptById(id uint) uint{

	id = id%config.NumberOfProcesses


	return config.Process[id].Aptitude
}

//read the json file and stock the values
func SetConfiguration()  {

	file, err:= os.Open("src/config/config.json")

	if err != nil{
		fmt.Println(err)
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)

	if err != nil{

		fmt.Println(err)
	}
}

func GetNumberOfProc() uint{

	return config.NumberOfProcesses
}

func GetArtificialDelay() time.Duration{

	return time.Duration(config.ArtificialDelay)*time.Second
}

func GetTransmitDelay()  time.Duration{

	return time.Duration(config.TransmitDelay)*time.Second
}

func GetWaitingAckdelay() time.Duration{

	return time.Duration(2*config.TransmitDelay)*time.Second
}

