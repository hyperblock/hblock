package hblock

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

func volume_PullBranch(obj *PullParams, logger *log.Logger) (int, error) {

	print_Log("Verify backingfile info...", logger)
	var err error
	var tmpConfig string
	obj.localRepoPath, err = return_Volume_BackingFile(&obj.volume)
	if err != nil {
		return FAIL, err
	}
	if dRet := VerifyBackingFile(obj.localRepoPath); dRet != OK {
		return FAIL, fmt.Errorf("Fail. ( ErrCode: %d )", dRet)
	}
	obj.configPath = return_BackingFileConfig_Path(&obj.localRepoPath)
	backingfileConfig := YamlBackingFileConfig{}
	print_Log("Load backingfile's config...", logger)
	if err = LoadConfig(&backingfileConfig, &obj.configPath); err != nil {
		return FAIL, err
	}
	branchHead := return_BranchHead(&obj.branch, &backingfileConfig.Branch)
	obj.remoteRepoPath = return_RemoteUrl(&backingfileConfig.Remote, &obj.remote)
	if obj.remoteRepoPath == "" {
		return FAIL, fmt.Errorf("Remote ('%s')'s URL can not be found.", obj.remote)
	}

	if branchHead != "" {
		print_Log("Trace branch commitID", logger)
		layerUUIDs, err := trace_Parents(&obj.localRepoPath, &branchHead)
		if err != nil {
			return FAIL, err
		}
		isConflict, tmp, err := branchConflict(obj.remoteRepoPath, obj.branch, layerUUIDs)

		if err != nil || isConflict == BRANCH_CONFLICT {
			if err == nil && isConflict == BRANCH_CONFLICT {
				return FAIL, fmt.Errorf("There is a conflict between local branch '%s' and remote. Use 'hb branch -m' to rename local branch.", obj.branch)
			} else {
				return FAIL, err
			}
		}
		if isConflict == BRANCH_CONTANS {
			print_Log(fmt.Sprintf("Branch '%s' has been fetched in local.", obj.branch), logger)
			return OK, nil
		}
		tmpConfig = tmp
	} else {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		tmpConfig = fmt.Sprintf("%s/%s.%d", os.TempDir(), path.Base(obj.remoteRepoPath), r.Intn(100000))
	}
	//	return FAIL, fmt.Errorf("debug.")

	if PathFileExists(tmpConfig) == false {

		remoteConfigPath := return_BackingFileConfig_Path(&obj.remoteRepoPath)
		print_Trace(fmt.Sprintf("Download remote config `%s` to local `%s`", remoteConfigPath, tmpConfig))
		if err = downloadFile(&remoteConfigPath, &tmpConfig); err != nil {
			return FAIL, err
		}
	}

	remoteConfig := YamlBackingFileConfig{}
	print_Log("integrate branch info.", logger)
	if err = LoadConfig(&remoteConfig, &tmpConfig); err != nil {
		return FAIL, err
	}
	for i := range remoteConfig.Branch {
		branch := &remoteConfig.Branch[i]
		fmt.Println(branch.Name)
		branch.Remote = obj.remote
		found := false
		for j := range backingfileConfig.Branch {
			if branch.Name == backingfileConfig.Branch[j].Name && branch.Remote == backingfileConfig.Branch[j].Remote {
				found = true
				break
			}
		}
		if found == false {
			backingfileConfig.Branch = append(backingfileConfig.Branch, *branch)
		}
	}

	//	return FAIL, fmt.Errorf("debug.")
	if err = WriteConfig(&backingfileConfig, &obj.configPath); err != nil {
		return FAIL, err
	}
	if err = PullBranch(obj, logger); err != nil {
		return FAIL, err
	}
	print_Log(Format_Success("Done."), logger)
	return OK, nil
}

func PullDefaultBranch(obj *PullParams, logger *log.Logger) error {

	print_Log(fmt.Sprintf("\rCheck default branch (Load config '%s' )......", obj.configPath), logger)
	backingfileConfig := YamlBackingFileConfig{}
	err := LoadConfig(&backingfileConfig, &obj.configPath)
	if err != nil {
		print_Log(fmt.Sprintf("\rCheck default branch (Load config '%s' )......FAIL!\n", obj.configPath), logger)
		return err
	}
	obj.branch = backingfileConfig.DefaultHead
	print_Log(fmt.Sprintf("\rCheck default branch (Load config '%s' )......OK: %s\n", obj.configPath, obj.branch), logger)
	return PullBranch(obj, logger)
}

func PullBranch(obj *PullParams, logger *log.Logger) error {

	print_Log("Pull branch: "+obj.branch, logger)
	backingFileConfig := YamlBackingFileConfig{}
	err := LoadConfig(&backingFileConfig, &obj.configPath)
	if err != nil {
		return err
	}
	h, err := CreateHyperLayer(FMT_UNKNOWN, &backingFileConfig.Format)
	if err != nil {
		return err
	}
	if len(obj.pullList) == 0 {
		print_Log(fmt.Sprintf("Found branch's head.( %s )", obj.branch), logger)
		if obj.all == true {
			for _, item := range backingFileConfig.Branch {
				obj.pullList = append(obj.pullList, item.Head)
			}
		} else {
			head := return_BranchHead(&obj.branch, &backingFileConfig.Branch)
			if head == "" {
				return fmt.Errorf("Branch can not be found.")
			}
			obj.pullList = []string{head}
		}
	}
	existLayers := make(map[string]bool) //check local repo's layer contains
	if PathFileExists(obj.localRepoPath) {
		jsonBackingFile, err := return_JsonBackingFile(&obj.localRepoPath)
		if err != nil {
			return nil
		}
		for _, layer := range return_LayerList(&jsonBackingFile) {
			existLayers[layer.uuid] = true
		}
	} else {
		initObj := InitParams{
			name:     obj.localRepoPath,
			size:     backingFileConfig.VirtualSize,
			checkout: false,
		}
		print_Log("\rCreate disk...", logger)
		err := createDisk(initObj)
		if err != nil {
			print_Log("FAIL\n", logger)
			return err
		}
		print_Log("\rCreate disk...OK\n", logger)
	}
	if obj.protocol == REPO_PATH_LOCAL {
		head := obj.pullList[0]
		print_Log("Trace branch..", logger)
		obj.pullList, err = trace_Parents(&obj.remoteRepoPath, &head)
		preLayer := ""
		for i := len(obj.pullList) - 1; i >= 0; i-- {
			layer := obj.pullList[i]
			_, ok := existLayers[layer]
			if !ok {
				layerFile := obj.localRepoPath + "." + layer
				dumpObj := DumpParams{
					backngFile: obj.remoteRepoPath,
					layerUUID:  layer,
					output:     layerFile,
				}
				print_Log(fmt.Sprintf("\rDumplayer to local.( uuid : %s )......", layer), logger)
				err = h.DumpLayer(&dumpObj)
				if err != nil {
					print_Log(fmt.Sprintf("\rDumplayer to local.( uuid : %s )......FAIL\n", layer), logger)
					return err
				}
				print_Log(fmt.Sprintf("\rDumplayer to local.( uuid : %s )......OK\n", layer), logger)
				rebaseObj := RebaseParams{
					volumePath:  dumpObj.output,
					backingfile: obj.localRepoPath,
					parentLayer: preLayer,
				}
				_, err = volume_Rebase(&rebaseObj, logger)
				if err != nil {
					return err
				}
				commitObj, err := return_CommitInfo(&dumpObj.output)
				if err != nil {
					return err
				}
				h.Commit(&commitObj)
				if err != nil {
					return err
				}
				os.Remove(layer)
				preLayer = layer
			}
		}
	} else if obj.protocol == REPO_PATH_HTTP {
		layers := []string{}
		defer RemoveFiles(layers)
		head := obj.pullList[0]
		for {
			if head == "" {
				break
			}
			print_Log(fmt.Sprintf("\rDownload layer (%s)...", head), logger)
			remoteLayerUrl := return_LayerName(obj.remoteRepoPath, head)
			localLayerPath := return_LayerName(obj.localRepoPath, head)
			err := downloadLayer(&remoteLayerUrl, &localLayerPath)
			if err != nil {
				RemoveFiles(append(layers, obj.configPath))
				return err
			}
			print_Log("OK\n", logger)
			layers = append(layers, localLayerPath)
			layerInfo, err := return_VolumeInfo(&localLayerPath)
			if err != nil {
				RemoveFiles(append(layers, obj.configPath))
				return err
			}
			head = layerInfo.layer
		}
		preLayer := ""
		print_Log("Commit layers...", logger)
		for i := len(layers) - 1; i >= 0; i-- {
			layer := layers[i]
			rebaseObj := RebaseParams{
				volumePath:  layer,
				backingfile: obj.localRepoPath,
				parentLayer: preLayer,
			}
			_, err = volume_Rebase(&rebaseObj, logger)
			if err != nil {
				return err
			}
			commitObj, err := return_CommitInfo(&layer)
			if err != nil {
				return err
			}
			h.Commit(&commitObj)
			if err != nil {
				return err
			}
			//	os.Remove(layer)
			dot := strings.LastIndex(layer, ".")
			preLayer = layer[dot+1:]
		}
	}
	print_Log("Update branch info...", logger)
	return setLocalBranchTag(&obj.configPath, &obj.branch, &obj.remote)
	//	return nil
}

func downloadLayer(srcLayer, dstLayer *string) error {

	resp, err := http.Get(*srcLayer)
	if err != nil {
		msg := fmt.Sprintf("Fetch: %v", err)
		return fmt.Errorf(msg)
	}
	defer resp.Body.Close()

	//	targetPath := currentDir + DEFALUT_BACKING_FILE_DIR + "/" + path.Base(layerName)
	dst, err := os.OpenFile(*dstLayer, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer dst.Close()
	length := resp.ContentLength
	buffer := make([]byte, 4*2<<20)
	cnt := int64(0)
	for {
		if length <= 0 {
			break
		}
		dwRead, err := io.ReadFull(resp.Body, buffer)
		if int64(dwRead) != length && err != nil {
			return err
		}
		length -= int64(dwRead)
		//_, err = io.Copy(dst, bytes.NewReader(buffer))
		_, err = dst.Write(buffer)
		if err != nil {
			return err
		}
		cnt += int64(dwRead)
		bar := print_ProcessBar(cnt, resp.ContentLength)
		fmt.Printf("\rDownload layer ( %s ) %s", *srcLayer, bar)
	}
	fmt.Println()
	if err != nil {
		msg := fmt.Errorf("Write buffer to file failed. (%v)", err)
		return msg
	}
	return nil
}
