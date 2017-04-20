package main

import (
	"hyperblock/hblock"
	"log"

	"os"
)

func main() {

	//home
	_, err := os.Stat("/var/hyperblock/log/hblock.log")
	if err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir("/var/hyperblock/log", 0777)
		}
	}
	file, err := os.OpenFile("/var/hyperblock/log/hblock.log", os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		log.Fatalln("fail to create log file!")
	}
	optSelector := hblock.Create(log.New(file, "", log.LstdFlags))
	//	args := []string{"init", "--name", "hehe", "--size", "2G"}
	args := os.Args
	//	fmt.Println(args)
	//status, err :=
	optSelector.SendCommand(args[1:])
	//fmt.Println(status, err)

}
