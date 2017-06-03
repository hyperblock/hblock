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
		msg := "ssh clone does not support at this time."
		print_Error(msg, logger)
		return FAIL, fmt.Errorf(msg)
	}
	if err != nil {
		print_Error(err.Error(), logger)
		return FAIL, err
	}
	print_Log(Format_Success("Clone finished."), logger)
	return OK, nil
}

func clone_Local(obj *CloneParams, logger *log.Logger) (int, error) {

	obj.repoPath = return_AbsPath(obj.repoPath)
	currentDir, err := return_CurrentDir()
	// if obj.repoPath[0] != '/' {
	// 	if err == nil {
	// 		obj.repoPath = currentDir + obj.repoPath
	// 		obj.configPath = currentDir + obj.repoPath+".yaml"
	// 	}
	// }
	print_Log("Initializating local hb directory...", logger)
	_, err = hb_Init()
	if err != nil {
		//	print_Error(err.Error(), logger)
		return FAIL, err
	}
	print_Log("Done.", logger)
	jsonVolume, err := return_JsonVolume(obj.repoPath)
	volFlag := false
	if jsonVolume.BackingFile != "" {
		volumeInfo := convert_to_VolumeInfo(&jsonVolume)
		backingfile := volumeInfo.backingFile
		obj.repoPath = backingfile
		obj.layerUUID = volumeInfo.layer
		volFlag = true
	}
	if dRet := VerifyBackingFile(obj.repoPath); dRet != OK {
		msg := fmt.Sprintf("backing file '%s' info can not be verified. (ErrCode: %d)", obj.repoPath, dRet)
		print_Error(msg, logger)
		return FAIL, fmt.Errorf(msg)
	}
	targetRepoPath := currentDir + DEFALUT_BACKING_FILE_DIR + "/" + path.Base(obj.repoPath)
	targetConfigPath := targetRepoPath + ".yaml"
	print_Log(fmt.Sprintf("Start clone repo '%s' to '%s'..", obj.repoPath, targetRepoPath), logger)
	_, err = CopyFile(targetConfigPath, obj.configPath)
	if err != nil {
		//	print_Error(err.Error(), logger)
		if PathFileExists(targetConfigPath) {
			os.Remove(targetConfigPath)
		}
		return FAIL, err
	}
	_, err = CopyFile(targetRepoPath, obj.repoPath)
	if err != nil {
		//	print_Error(err.Error(), logger)
		if PathFileExists(targetRepoPath) {
			os.Remove(targetRepoPath)
		}
		return FAIL, err
	}
	//	checkoutObj.template = targetRepoPath
	print_Log("Done.", logger)
	// } else {
	// 	print_Log("The clone repo is a volume, will find the backing file.", logger)
	// 	volumeInfo := convert_to_VolumeInfo(&jsonVolume)
	// 	backingfile := volumeInfo.backingFile
	// 	configPath := backingfile + ".yaml"
	// 	print_Log(fmt.Sprintf("Backing file: %s", backingfile), logger)
	// 	if PathFileExists(configPath) == false {
	// 		msg := fmt.Sprintf("backing file's config '%s' not found.", configPath)
	// 		return FAIL, fmt.Errorf(msg)
	// 	}
	// 	if PathFileExists(backingfile) == false {
	// 		msg := fmt.Sprintf("backing file '%s' not found.", backingfile)
	// 		return FAIL, fmt.Errorf(msg)
	// 	}
	// 	targetRepoPath := currentDir + DEFALUT_BACKING_FILE_DIR + "/" + path.Base(backingfile)
	// 	targetConfigPath := targetRepoPath + ".yaml"
	// 	print_Log(fmt.Sprintf("Start clone repo '%s' to '%s'..", backingfile, targetRepoPath), logger)
	// 	_, err := CopyFile(targetConfigPath, configPath)
	// 	if err != nil {
	// 		//	print_Error(err.Error(), logger)
	// 		if PathFileExists(targetConfigPath) {
	// 			os.Remove(targetConfigPath)
	// 		}
	// 		return FAIL, err
	// 	}
	// 	_, err = CopyFile(targetRepoPath, backingfile)
	// 	if err != nil {
	// 		//	print_Error(err.Error(), logger)
	// 		if PathFileExists(targetRepoPath) {
	// 			os.Remove(targetRepoPath)
	// 		}
	// 		return FAIL, err
	// 	}
	// 	checkoutObj.template = targetRepoPath
	// 	checkoutObj.layer = volumeInfo.layer
	// 	print_Log("Done.", logger)
	// }
	if !obj.checkoutFlg {
		return OK, nil
	}
	checkoutObj := CheckoutParams{template: obj.repoPath, output: currentDir + path.Base(obj.repoPath), layer: obj.layerUUID}
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

	print_Log("Initializating local hb directory...", logger)
	_, err := hb_Init()
	if err != nil {
		return FAIL, err
	}
	currentDir, err := return_CurrentDir()
	if err != nil {
		print_Log(Format_Warning("Can't get pwd."), logger)
	}

	print_Log(fmt.Sprintf("Downloading repo's config from url: %s", obj.configPath), logger)
	respConfig, err := http.Get(obj.configPath)
	if err != nil {
		msg := fmt.Sprintf("Fetch: %v", err)
		return FAIL, fmt.Errorf(msg)
	}
	defer respConfig.Body.Close()
	targetConfigPath := currentDir + DEFALUT_BACKING_FILE_DIR + "/" + path.Base(obj.configPath)
	configDst, err := os.OpenFile(targetConfigPath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return 0, err
	}
	defer configDst.Close()
	configBuff, err := ioutil.ReadAll(respConfig.Body)
	if err != nil {

		return FAIL, err
	}
	_, err = io.Copy(configDst, bytes.NewBuffer(configBuff))
	if err != nil {
		return FAIL, err
	}
	print_Log("Done.", logger)

	print_Log(fmt.Sprintf("Downloading repo from url: %s", obj.repoPath), logger)
	resp, err := http.Get(obj.repoPath)
	if err != nil {
		msg := fmt.Sprintf("Fetch: %v", err)
		return FAIL, fmt.Errorf(msg)
	}
	print_Log(fmt.Sprintf("Content length: %d", resp.ContentLength), logger)
	defer resp.Body.Close()
	//return FAIL, fmt.Errorf("111")
	if err != nil {
		msg := fmt.Sprintf("Fetch: reading %s: %v\n", obj.repoPath, err)
		return FAIL, fmt.Errorf(msg)
	}
	//	print_Log(fmt.Sprintf("Read buffer done. (length: %d)", len(buffer)), logger)
	targetPath := currentDir + DEFALUT_BACKING_FILE_DIR + "/" + path.Base(obj.repoPath)
	dst, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)

	if err != nil {
		return 0, err
	}
	defer dst.Close()
	//_, err = io.Copy(dst, bytes.NewReader(buffer))
	length := resp.ContentLength
	buffer := make([]byte, 4*2<<20)
	cnt := int64(0)
	for {
		if length <= 0 {
			break
		}
		dwRead, err := io.ReadFull(resp.Body, buffer)
		if int64(dwRead) != length && err != nil {
			return FAIL, err
		}
		length -= int64(dwRead)
		_, err = io.Copy(dst, bytes.NewReader(buffer))
		if err != nil {
			return FAIL, nil
		}
		cnt += int64(dwRead)
		msg := print_ProcessBar(cnt, resp.ContentLength)
		fmt.Printf("\rDownloading %s", msg)
	}
	fmt.Println()
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
