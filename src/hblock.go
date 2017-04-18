package main

import (
	"hyperblock/hblock"
	"log"

	"os"
)

func main() {

	_, err := os.Stat("log")
	if err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir("log", 0775)
		}
	}
	file, err := os.Create("log/testOption.log")
	if err != nil {
		log.Fatalln("fail to create test.log file!")
	}
	optSelector := hblock.Create(log.New(file, "", log.LstdFlags|log.Llongfile))
	//	args := []string{"init", "--name", "hehe", "--size", "2G"}
	args := os.Args
	//	fmt.Println(args)
	//status, err :=
	optSelector.SendCommand(args[1:])
	//fmt.Println(status, err)

}
