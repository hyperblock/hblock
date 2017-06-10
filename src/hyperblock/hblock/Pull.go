package hblock

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

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
		head := obj.pullList[0]
		for {
			if head == "" {
				break
			}
			print_Log("\rDownload layer (%s)", logger)
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
			os.Remove(layer)
			preLayer = layer
		}
	}
	print_Log("Update branch info...", logger)
	return setLocalBranchTag(&obj.configPath, &obj.branch)
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
