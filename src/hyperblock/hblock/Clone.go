package hblock

import (
	"fmt"
	"log"
	"os"
	"path"
)

func clone_Repo(obj *CloneParams, logger *log.Logger) (int, error) {

	obj.protocol = return_RepoPath_Type(obj.repoPath)
	if obj.protocol == REPO_PATH_LOCAL {
		_, err := clone_Local(obj, logger)
		if err != nil {
			print_Error(err.Error(), logger)
			return FAIL, err
		}
		print_Log(format_Success("Clone finished!"), logger)
	}
	return OK, nil
}

func clone_Local(obj *CloneParams, logger *log.Logger) (int, error) {

	currentDir, err := return_CurrentDir()
	if obj.repoPath[0] != '/' {
		if err == nil {
			obj.repoPath = currentDir + obj.repoPath
		}
	}
	if PathFileExists(obj.repoPath) == false {
		msg := fmt.Sprintf("repo '%s' couldn't be found", obj.repoPath)
		//	print_Error(msg, logger)
		return FAIL, fmt.Errorf(msg)
	}
	print_Log("Initializating local hb directory...", logger)
	_, err = hb_Init()
	if err != nil {
		//	print_Error(err.Error(), logger)
		return FAIL, err
	}
	print_Log("Done.", logger)
	jsonVolume, err := return_JsonVolume(obj.repoPath)
	checkoutObj := CheckoutParams{output: currentDir + path.Base(obj.repoPath)}
	volFlag := true
	if jsonVolume.BackingFile == "" {
		volFlag = false
		targetPath := currentDir + DEFALUT_BACKING_FILE_DIR + "/" + path.Base(obj.repoPath)
		print_Log(fmt.Sprintf("Start clone repo '%s' to '%s'..", obj.repoPath, targetPath), logger)
		_, err := CopyFile(targetPath, obj.repoPath)
		if err != nil {
			//	print_Error(err.Error(), logger)
			if PathFileExists(targetPath) {
				os.Remove(targetPath)
			}
			return FAIL, err
		}
		checkoutObj.template = targetPath

		print_Log("Done.", logger)
	} else {
		print_Log("The clone repo is a volume, will find the backing file.", logger)
		volumeInfo := return_VolumeInfo(&jsonVolume)
		backingfile := volumeInfo.backingFile
		print_Log(fmt.Sprintf("Backing file: %s", backingfile), logger)
		if PathFileExists(backingfile) == false {
			msg := fmt.Sprintf("backing file '%s' not found.", backingfile)
			//	print_Error(msg, logger)
			return FAIL, fmt.Errorf(msg)
		}
		targetPath := currentDir + DEFALUT_BACKING_FILE_DIR + "/" + path.Base(backingfile)
		print_Log(fmt.Sprintf("Start clone repo '%s' to '%s'..", backingfile, targetPath), logger)
		_, err := CopyFile(targetPath, backingfile)
		if err != nil {
			//	print_Error(err.Error(), logger)
			if PathFileExists(targetPath) {
				os.Remove(targetPath)
			}
			return FAIL, err
		}
		checkoutObj.template = targetPath
		checkoutObj.layer = volumeInfo.layer
		print_Log("Done.", logger)
	}
	if !obj.checkoutFlg {
		return OK, nil
	}
	print_Log("Ready to checkout from backing file.", logger)
	if !volFlag || obj.layerUUID != "" {
		jsonBackingFile, err := return_JsonBackingFile(obj.repoPath)
		if err != nil {
			msg := fmt.Sprintf("Can't get backing file info (%s).", err.Error())
			//	print_Error(msg, logger)
			return FAIL, fmt.Errorf(msg)
		}
		layerList := return_Snapshots(&jsonBackingFile)
		if obj.layerUUID == "" {
			if len(layerList) > 0 {
				checkoutObj.layer = layerList[len(layerList)-1].uuid
			} else {
				checkoutObj.layer = ""
			}
		} else {
			layer, err := return_LayerUUID_from_Snapshots(layerList, obj.layerUUID)
			if err != nil {
				//print_Error(err.Error(), logger)
				return FAIL, err
			}
			obj.layerUUID = layer
		}
	}
	return volume_checkout(checkoutObj, logger)
	//volumeInfo := return_VolumeInfo(&jsonVolume)
	//fmt.Println(jsonBackingFile)
}
