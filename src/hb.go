package main

import (
	"fmt"
	"hyperblock/hblock"
	"log"
	"os"
	"os/user"
	"time"
)

func main() {

	//home
	usr, err := user.Current()
	hb_Dir := usr.HomeDir + "/.hb"
	_, err = os.Stat(hb_Dir)
	if err != nil {
		fmt.Println(hblock.Format_Warning("Hyperblock global directory doesn't exists, will create..."))
		if os.IsNotExist(err) {
			err = os.Mkdir(hb_Dir, 0755)
			if err != nil {
				msg := hblock.Format_Error(
					fmt.Sprintf("Create failed. ( %s )", err.Error()))
				fmt.Println(msg)
				return
			}
		}
		fmt.Println(hblock.Format_Info("Done."))
	}
	logDir := hb_Dir + "/log"
	_, err = os.Stat(logDir)
	if err != nil {
		fmt.Println(hblock.Format_Warning("Hyperblock log directory doesn't exists, will create..."))
		if os.IsNotExist(err) {
			err = os.Mkdir(logDir, 0755)
			if err != nil {

				msg := hblock.Format_Error(
					fmt.Sprintf("Create failed. ( %s )", err.Error()))
				fmt.Println(msg)
				return
			}
		}
		fmt.Println(hblock.Format_Info("Done."))
	}

	logFile := fmt.Sprintf("%s/%d_%d_%d.log", logDir, time.Now().Year(), time.Now().Month(), time.Now().Day())
	//fmt.Println(logFile)
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		msg := hblock.Format_Error(
			fmt.Sprintf("Create log failed. ( %s )", err.Error()))
		fmt.Println(msg)
		return
	}
	logger := log.New(file, "", log.LstdFlags)
	optSelector := hblock.Create(logger)
	//	args := []string{"init", "--name", "hehe", "--size", "2G"}
	args := os.Args
	//	fmt.Println(args)
	//status, err :=
	optSelector.SendCommand(args[1:])
	file.Close()
	//fmt.Println(status, err)
}
