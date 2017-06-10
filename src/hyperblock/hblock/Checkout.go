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

func volume_checkout(obj *CheckoutParams, logger *log.Logger) (int, error) {

	checkoutArgs := []string{}
	rollback := false
	tmpOutput := ""
	if obj.branch != "" {
		print_Log(fmt.Sprintf("Create new branch '%s' (cached)", obj.branch), logger)
		volumeConfigPath := return_Volume_ConfigPath(&obj.volume)
		yamlVolumeConfig := YamlVolumeConfig{}
		err := LoadConfig(&yamlVolumeConfig, &volumeConfigPath)
		if err != nil {
			//	print_Error(err.Error(), logger)
			return FAIL, err
		}
		yamlVolumeConfig.Branch = obj.branch
		yamlVolumeConfig.NewBranch = true
		err = WriteConfig(&yamlVolumeConfig, &volumeConfigPath)
		if err != nil {
			//	print_Error(err.Error(), logger)
			return FAIL, err
		}
		print_Log(Format_Success("Finished. (create new branch after commit)\n"), logger)
		return OK, nil
	}
	layer := ""
	print_Trace("volume: " + obj.volume)
	if obj.volume != "" {
		backingFile, err := return_Volume_BackingFile(&obj.volume)
		if err != nil {
			//	print_Error(err.Error(), logger)
			return FAIL, err
		}

		layer, err = return_LayerUUID(backingFile, obj.layer, false)
		if err != nil {
			//print_Error(err.Error(), logger)
			return FAIL, err
		}
		checkoutArgs = []string{"create", "-t", backingFile, "-l", layer}
		if obj.volume == obj.output {
			tmpOutput = obj.output + ".tmp" + strconv.Itoa(rand.Int())[0:4]
			rollback = true
			checkoutArgs = append(checkoutArgs, tmpOutput)
		} else {
			checkoutArgs = append(checkoutArgs, obj.output)
		}
	} else if obj.template != "" {
		backingFile, err := confirm_BackingFilePath(obj.template)
		if err != nil {
			//	print_Error(err.Error(), logger)
			return FAIL, err
		}
		layer, err = return_LayerUUID(backingFile, obj.layer, false)
		if err != nil {
			//	print_Error(err.Error(), logger)
			return FAIL, err
		}
		checkoutArgs = []string{"create", "-t", backingFile, "-l", layer, obj.output}
	}
	print_Log("Create volume's config file...", logger)
	//	createArgs := []string{"-l", layer, "-t", backingFile}
	yamlVolume := YamlVolumeConfig{}
	volumeConfigPath := return_Volume_ConfigPath(&obj.output)
	if volumeConfigPath == "" {
		return FAIL, fmt.Errorf("Can't locate volume config file path.")
	}
	err := WriteConfig(&yamlVolume, &volumeConfigPath)
	if err != nil {
		//print_Error(err.Error(), logger)
		return FAIL, err
	}

	print_Log(strings.Join(checkoutArgs, " "), logger)
	cmdCreate := exec.Command("qcow2-img", checkoutArgs[0:]...)
	result, err := cmdCreate.Output()
	if err != nil {
		//print_Panic(err.Error(), logger)
		return FAIL, err
	}
	print_Log(string(result), logger)
	if rollback {
		rmErr := os.Remove(obj.volume)
		if rmErr != nil {
			os.Remove(tmpOutput)
			//	print_Panic(rmErr.Error(), logger)
			return FAIL, rmErr
		}
		mvErr := os.Rename(tmpOutput, obj.volume)
		if mvErr != nil {
			os.Remove(tmpOutput)
			//	print_Panic(mvErr.Error(), logger)
			return FAIL, mvErr
		}
	}
	//	fmt.Println(backingFile)
	print_Log(Format_Success("Checkout finished."), logger)
	return OK, nil
}
