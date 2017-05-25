package hblock

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
)

func clone_Repo(obj *CloneParams, logger *log.Logger) (int, error) {

	obj.protocol = return_RepoPath_Type(obj.repoPath)
	var err error
	if obj.protocol == REPO_PATH_LOCAL {
		_, err = clone_Local(obj, logger)
	} else if obj.protocol == REPO_PATH_HTTP {
		_, err = clone_Http(obj, logger)
	} else if obj.protocol == REPO_PATH_SSH {
		//_, err = clone_
	}
	if err != nil {
		print_Error(err.Error(), logger)
		return FAIL, err
	}
	print_Log(Format_Success("Clone finished."), logger)
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
		volumeInfo := convert_to_VolumeInfo(&jsonVolume)
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
		checkoutObj.layer, err = return_LayerUUID(obj.repoPath, obj.layerUUID, true)
		if err != nil {
			return FAIL, err
		}
	}
	return volume_checkout(&checkoutObj, logger)

}

func clone_Http(obj *CloneParams, logger *log.Logger) (int, error) {

	resp, err := http.Get(obj.repoPath)
	print_Log(fmt.Sprintf("Downloading repo from url: %s", obj.repoPath), logger)
	if err != nil {
		msg := fmt.Sprintf("Fetch: %v", err)
		return FAIL, fmt.Errorf(msg)
	}

	buffer, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		msg := fmt.Sprintf("Fetch: reading %s: %v\n", obj.repoPath, err)
		return FAIL, fmt.Errorf(msg)
	}
	print_Log(fmt.Sprintf("Read buffer done. (length: %d)", len(buffer)), logger)
	print_Log("Initializating local hb directory...", logger)
	_, err = hb_Init()
	if err != nil {
		return FAIL, err
	}

	currentDir, err := return_CurrentDir()
	if err != nil {
		print_Log(Format_Warning("Can't get pwd."), logger)
	}
	targetPath := currentDir + DEFALUT_BACKING_FILE_DIR + "/" + path.Base(obj.repoPath)
	dst, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_CREATE, 0644)
	defer dst.Close()
	if err != nil {
		return 0, err
	}
	_, err = io.Copy(dst, bytes.NewReader(buffer))
	if err != nil {
		msg := fmt.Sprintf("Write buffer to file failed. (%v)", err)
		return FAIL, fmt.Errorf(msg)
	}
	if !obj.checkoutFlg {
		return OK, nil
	}
	checkoutObj := CheckoutParams{
		template: targetPath,
		output:   path.Base(targetPath),
	}
	checkoutObj.layer, err = return_LayerUUID(targetPath, obj.layerUUID, true)
	if err != nil {
		return FAIL, err
	}
	return volume_checkout(&checkoutObj, logger)

}
