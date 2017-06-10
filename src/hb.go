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
	defer file.Close()
	logger := log.New(file, "", log.LstdFlags)
	optSelector := hblock.CreateOptSelector(logger)
	//	args := []string{"init", "--name", "hehe", "--size", "2G"}
	args := os.Args
	//	fmt.Println(args)
	//status, err :=
	//	args = strings.Split("hb clone /var/www/html/repo/test", " ")
	//args = strings.Split("hb push origin master -v test/test", " ")
	//	args = strings.Split("hb init t1 --size 10G", " ")
	_, err = optSelector.SendCommand(args[1:])
	if err != nil {
		hblock.Print_Error(err.Error(), logger)
	}
	//fmt.Println(status, err)
}
