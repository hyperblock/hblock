package hblock

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

func push_volume(obj PushParams, logger *log.Logger) (int, error) {

	var tmpPath string
	volumeInfo, err := return_VolumeInfo(&obj.volume)
	if err != nil {
		//print_Error(err.Error(), logger)
		return FAIL, err
	}
	print_Log("Verify backing file.....", logger)
	if dRet := VerifyBackingFile(volumeInfo.backingFile); dRet != OK {
		return FAIL, fmt.Errorf("Can not verify backing file config.\n ( %s ErrCode: %d )", err.Error(), dRet)
	}
	print_Log("Load config.....", logger)
	configPath := return_BackingFileConfig_Path(&volumeInfo.backingFile) // volumeInfo.backingFile + ".yaml"
	backingFileConfig := YamlBackingFileConfig{}
	err = LoadConfig(&backingFileConfig, &configPath)
	if err != nil {
		return FAIL, err
	}
	obj.url = return_RemoteUrl(&backingFileConfig.Remote, &obj.remote)
	if obj.url == "" {
		return FAIL, fmt.Errorf("Can not found remote '%s', use 'hb remote <volume> --add' to add a new remote host.", obj.remote)
	}
	branchHead := return_BranchHead(&obj.branch, &backingFileConfig.Branch)
	if branchHead == "" {
		return FAIL, fmt.Errorf("Branch '%s' doesn't in volume. use 'hb branch <volume> to see all branches.", obj.branch)
	}
	print_Log("Load backing file info...", logger)
	print_Log(fmt.Sprintf("Trace parent layers of branch '%s'...", obj.branch), logger)
	layerUUIDs, err := trace_Parents(&volumeInfo.backingFile, &branchHead)
	if err != nil {
		//print_Error(err.Error(), logger)
		return FAIL, err
	}
	layerTrace := ""
	for _, uuid := range layerUUIDs {
		layerTrace += uuid + "\n"
	}
	print_Log("layer trace:\n"+layerTrace, logger)
	isConflict, tmpPath, err := branchConflict(obj.url, obj.branch, layerUUIDs)
	if isConflict == BRANCH_CONFLICT || err != nil {
		if err == nil && isConflict == BRANCH_CONFLICT {
			return FAIL, fmt.Errorf("There is a conflict between local branch '%s' and remote. Use 'hb pull' to fetch remote branch or 'hb branch -m' to rename local branch.", obj.branch)
		} else {
			return FAIL, err
		}
	}
	//if isConflict ==
	layerFiles := []string{}
	defer RemoveFiles(layerFiles)
	for _, layer := range layerUUIDs {
		fileName := volumeInfo.backingFile + "." + layer
		print_Log(fmt.Sprintf("\rDump layer ( uuid = %s )......", layer), logger)
		dumpObj := DumpParams{
			backngFile: volumeInfo.backingFile,
			layerUUID:  layer,
			output:     fileName,
		}

		//h, err := CreateHBM(FMT_UNKNOWN, &backingFileConfig.Format)
		h, err := CreateHBM(FMT_UNKNOWN, obj.volume)
		if err != nil {
			return 0, err
		}
		err = h.DumpLayer(&dumpObj)

		if err != nil {
			//	print_Error(fmt.Sprintf("Fail. (%s)", err.Error()), logger)
			return FAIL, err
		}
		layerFiles = append(layerFiles, fileName)
		print_Log(fmt.Sprintf("\rDump layer ( uuid = %s )......OK\n", layer), logger)
	}

	index := 0
	msg := ""
	repoDir := func() string {
		p := strings.LastIndex(obj.url, "/")
		return obj.url[0 : p+1]
	}()
	for _, fileName := range layerFiles {
		index++
		msg = fmt.Sprintf("\rPush layers (%d/%d)......", index, len(layerFiles))
		print_Log(msg, logger)
		url := repoDir + path.Base(fileName)
		putError := httpPut(fileName, url)
		if putError != nil {
			return FAIL, fmt.Errorf("%sFail ( %s )", msg, err.Error())
		}
	}
	print_Log(msg+"OK\n", logger)
	print_Log("Update local config and push to remote...", logger)
	if PathFileExists(tmpPath) == false {
		_, err = CopyFile(tmpPath, configPath)
		if err != nil {
			return FAIL, err
		}
	}
	LoadConfig(&backingFileConfig, &tmpPath)
	defer os.Remove(tmpPath)
	found := false
	for i := range backingFileConfig.Branch {
		backingFileConfig.Branch[i].Local = 0
		backingFileConfig.Branch[i].Remote = ""
		if backingFileConfig.Branch[i].Name == obj.branch {
			backingFileConfig.Branch[i].Head = branchHead
			found = true
		}
	}
	if !found {
		backingFileConfig.Branch = append(backingFileConfig.Branch, YamlBranch{
			Name: obj.branch, Head: branchHead, Local: 0, Remote: "",
		})
	}
	err = WriteConfig(&backingFileConfig, &tmpPath)
	if err != nil {
		return FAIL, err
	}
	remoteConfigUrl := repoDir + path.Base(tmpPath)
	print_Log(fmt.Sprintf("Push config to %s.", remoteConfigUrl), logger)
	if err = httpPut(tmpPath, remoteConfigUrl); err != nil {
		return FAIL, err
	}
	print_Log(Format_Success("Done."), logger)
	return OK, nil
}

func httpPut(filename string, targetUrl string) error {

	fh, err := os.Open(filename)
	if err != nil {
		//	fmt.Println("error opening file")
		return err
	}
	defer fh.Close()
	req, err := http.NewRequest("PUT", targetUrl, fh)
	req.Header.Set("Content-Type", "application/octet-stream")
	resp, err := (&http.Client{}).Do(req)

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	defer resp.Body.Close()

	//	resp_body, err := ioutil.ReadAll(resp.Body)
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	// fmt.Println(resp.Status)
	// fmt.Println(string(resp_body))
	return nil

}
