package hblock

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
)

func volume_checkout(obj *CheckoutParams, logger *log.Logger) (int, error) {

	checkoutArgs := []string{}
	rollback := ""
	tmpOutput := ""
	if obj.branch != "" {
		print_Log(fmt.Sprintf("Create new branch '%s' (cached)", obj.branch), logger)
		volumeConfigPath := return_Volume_ConfigPath(&obj.volume)
		yamlVolumeConfig := YamlVolumeConfig{Branch: obj.branch, NewBranch: true}
		// err := LoadConfig(&yamlVolumeConfig, &volumeConfigPath)
		// if err != nil {
		// 	//	print_Error(err.Error(), logger)
		// 	return FAIL, err
		// }
		err := WriteConfig(&yamlVolumeConfig, &volumeConfigPath)
		if err != nil {
			//	print_Error(err.Error(), logger)
			return FAIL, err
		}
		print_Log(Format_Success("Finished. (create new branch after commit)\n"), logger)
		return OK, nil
	}
	layer := ""
	backingFileConfig := ""
	print_Trace("volume: " + obj.volume)
	if obj.volume != "" {
		backingFile, err := return_Volume_BackingFile(&obj.volume)
		if err != nil {
			//	print_Error(err.Error(), logger)
			return FAIL, err
		}
		backingFileConfig = return_BackingFileConfig_Path(&backingFile)
		layer, err = return_LayerUUID(backingFile, obj.layer, false)
		if err != nil {
			//print_Error(err.Error(), logger)
			return FAIL, err
		}
		//checkoutArgs = []string{"create", "-t", backingFile, "-l", layer}
		if obj.volume == obj.output {
			tmpOutput = obj.output + ".tmp" + strconv.Itoa(rand.Int())[0:4]
			rollback = obj.output
			obj.output = tmpOutput
			checkoutArgs = append(checkoutArgs, tmpOutput)
		}
		obj.layer = layer
		obj.template = backingFile
	} else if obj.template != "" {
		backingFile, err := confirm_BackingFilePath(obj.template)
		if err != nil {
			//	print_Error(err.Error(), logger)
			return FAIL, err
		}
		backingFileConfig = return_BackingFileConfig_Path(&backingFile)
		layer, err = return_LayerUUID(backingFile, obj.layer, false)
		if err != nil {
			//	print_Error(err.Error(), logger)
			return FAIL, err
		}
		obj.template = backingFile
		obj.layer = layer
		//checkoutArgs = []string{"create", "-t", backingFile, "-l", layer, obj.output}
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
	h, err := CreateHyperLayer(FMT_UNKNOWN, &backingFileConfig)
	if err != nil {
		return FAIL, err
	}
	err = h.Checkout(obj)
	if err != nil {
		return FAIL, err
	}

	if rollback != "" {
		rmErr := os.Remove(rollback)
		if rmErr != nil {
			os.Remove(tmpOutput)
			//	print_Panic(rmErr.Error(), logger)
			return FAIL, rmErr
		}
		mvErr := os.Rename(obj.output, rollback)
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
