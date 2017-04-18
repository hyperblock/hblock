package hblock

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type CheckoutParams struct {
	layer    string
	volume   string
	output   string
	template string
}

func volume_checkout(obj CheckoutParams, logger *log.Logger) (int, error) {

	checkoutArgs := []string{}
	rollback := false
	tmpOutput := ""
	if obj.volume != "" {
		cmdVolInfo := exec.Command("qcow2-img", "info", obj.volume)
		volInfoBuf, err := cmdVolInfo.Output()
		if err != nil {
			//msg := format_Error(err.Error())
			print_Error(err.Error(), logger)
			return FAIL, err
		}
		print_Log(string(volInfoBuf), logger)
		volInfoList := strings.Split(string(volInfoBuf), "\n")
		strValue := get_InfoValue(volInfoList, "backing file")
		backingFile := get_StringBefore(
			get_StringAfter(strValue, "qcow2://"), "?")
		layer := obj.layer
		backingFile, err = confirm_BackingFilePath(backingFile)
		if err != nil {
			print_Error(err.Error(), logger)
		}
		if backingFile == "" {
			msg := "Can't find backing file."
			print_Log(msg, logger)
			return FAIL, fmt.Errorf(msg)
		}
		tmpOutput = "tmp" + strconv.Itoa(rand.Int())
		rollback = true
		checkoutArgs = []string{"create", "-t", backingFile, "-l", layer, tmpOutput}
	} else if obj.template != "" {
		backingFile, err := confirm_BackingFilePath(obj.template)
		if err != nil {
			if strings.Index(err.Error(), "env") != -1 {
				//fmt.Println(err.Error(), err.Error()[0)
				print_Log(err.Error(), logger)
			} else {
				print_Error(err.Error(), logger)
			}
		}
		if backingFile == "" {
			msg := "Can't find backing file."
			print_Error(msg, logger)
			return FAIL, fmt.Errorf(msg)
		}
		checkoutArgs = []string{"create", "-t", backingFile, "-l", obj.layer, obj.output}
	}

	//	createArgs := []string{"-l", layer, "-t", backingFile}

	print_Log(strings.Join(checkoutArgs, " "), logger)
	cmdCreate := exec.Command("qcow2-img", checkoutArgs[0:]...)
	result, err := cmdCreate.Output()
	if err != nil {
		print_Panic(err.Error(), logger)
		return FAIL, err
	}
	print_Log(string(result), logger)
	if rollback {
		rmErr := os.Remove(obj.volume)
		if rmErr != nil {
			os.Remove(tmpOutput)
			print_Panic(rmErr.Error(), logger)
			return FAIL, rmErr
		}
		mvErr := os.Rename(tmpOutput, obj.volume)
		if mvErr != nil {
			os.Remove(tmpOutput)
			print_Panic(mvErr.Error(), logger)
			return FAIL, mvErr
		}
	}
	//	fmt.Println(backingFile)
	print_Log(format_Success("Finished."), logger)
	return OK, nil
}
