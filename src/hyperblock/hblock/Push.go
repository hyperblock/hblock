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
		return FAIL, fmt.Errorf("Can not found remote '%s', use 'hb remote --add' to add a new remote host.", obj.remote)
	}
	branchHead := func() string {
		for _, item := range backingFileConfig.Branch {
			if item.Name == obj.branch {
				return item.Head
			}
		}
		return ""
	}()
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
	layerFiles := []string{}
	for _, layer := range layerUUIDs {
		fileName := volumeInfo.backingFile + "." + layer
		print_Log(fmt.Sprintf("\rDump layer ( uuid = %s )......", layer), logger)
		dumpObj := DumpParams{
			backngFile: volumeInfo.backingFile,
			layerUUID:  layer,
			output:     fileName,
		}
		err = DumpLayer(&dumpObj)
		if err != nil {
			//	print_Error(fmt.Sprintf("Fail. (%s)", err.Error()), logger)
			return FAIL, err
		}
		layerFiles = append(layerFiles, fileName)
		print_Log(fmt.Sprintf("\rDump layer ( uuid = %s )......OK\n", layer), logger)
	}

	index := 0
	msg := ""
	for _, fileName := range layerFiles {
		index++
		msg = fmt.Sprintf("\rPush objects (%d/%d)......", index, len(layerFiles))
		print_Log(msg, logger)
		p := strings.LastIndex(obj.url, "/")
		url := obj.url[0:p+1] + path.Base(fileName)
		putError := httpPut(fileName, url)
		if putError != nil {
			return FAIL, fmt.Errorf("%sFail ( %s )", msg, err.Error())
		}
	}
	print_Log(msg+"OK\n", logger)
	return OK, nil
}

func httpPut(filename string, targetUrl string) error {

	// bodyBuf := &bytes.Buffer{}
	// bodyWriter := multipart.NewWriter(bodyBuf)

	// //关键的一步操作
	// fileWriter, err := bodyWriter.CreateFormFile("uploadfile", filename)
	// if err != nil {
	// 	//	fmt.Println("error writing to buffer")
	// 	return err
	// }

	//打开文件句柄操作
	fh, err := os.Open(filename)
	if err != nil {
		//	fmt.Println("error opening file")
		return err
	}
	defer fh.Close()
	//iocopy
	// _, err = io.Copy(fileWriter, fh)
	// if err != nil {
	// 	return err
	// }

	//	contentType := bodyWriter.FormDataContentType()
	//bodyWriter.Close()

	//resp, err := http.Post(targetUrl, contentType, bodyBuf)
	//resp, err := http.Post(targetUrl, "application/octet-stream", bodyBuf)
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
