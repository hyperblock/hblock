package hblock

import (
	"fmt"
	"log"
	"os"
	"path"
)

func clone_Repo(obj *CloneParams, logger *log.Logger) (int, error) {

	obj.protocol = return_RepoPath_Type(obj.repoPath)
	var err error
	checkoutObj := CheckoutParams{}

	print_Log("Initializating local hb directory...", logger)
	hbDir, err := hb_Init()
	if err != nil {
		return FAIL, err
	}
	localRepoPath := hbDir + path.Base(obj.repoPath)
	if PathFileExists(localRepoPath) {
		return FAIL, fmt.Errorf("Repo ( %s ) is exists!", localRepoPath)
	}
	if obj.protocol == REPO_PATH_LOCAL {
		checkoutObj, err = clone_Local(obj, logger)
	} else if obj.protocol == REPO_PATH_HTTP {
		if obj.hardLink {
			print_Log("Http protocol get, ignore --hardlink.", logger)
		}
		checkoutObj, err = clone_Http(obj, logger)
	} else if obj.protocol == REPO_PATH_SSH {
		msg := "ssh clone does not support at this time."
		return FAIL, fmt.Errorf(msg)
	}
	if err != nil {
		return FAIL, err
	}
	if obj.checkoutFlg {
		defaultBranch := checkoutObj.branch
		checkoutObj.branch = ""
		//	_, err= volume_checkout(&checkout, logger)
		print_Log("Checkout volume...", logger)
		_, err = volume_checkout(&checkoutObj, logger)
		if err != nil {
			//		print_Error(err.Error(), logger)
			return FAIL, err
		}
		yamlVolConfig := YamlVolumeConfig{}
		print_Log("Write volume config...", logger)
		volConfigPath := return_Volume_ConfigPath(&checkoutObj.output)
		err = LoadConfig(&yamlVolConfig, &volConfigPath)
		if err != nil {
			return FAIL, err
		}
		yamlVolConfig.Branch = defaultBranch
		err = WriteConfig(&yamlVolConfig, &volConfigPath)
		if err != nil {
			return FAIL, err
		}
	}
	print_Log(Format_Success("Clone finished."), logger)
	return OK, nil
}

func clone_Local(obj *CloneParams, logger *log.Logger) (CheckoutParams, error) {

	print_Log(fmt.Sprintf("Clone from local: %s", obj.repoPath), logger)
	obj.repoPath = return_AbsPath(obj.repoPath)
	currentDir, err := return_CurrentDir()
	print_Log("Initializating local hb directory...", logger)
	checkoutRet := CheckoutParams{}
	hbDir, err := hb_Init()
	if err != nil {
		return checkoutRet, err
	}

	jsonVolume, err := return_JsonVolume(obj.repoPath)
	if jsonVolume.BackingFile != "" {
		volumeInfo := convert_to_VolumeInfo(&jsonVolume)
		backingfile := volumeInfo.backingFile
		obj.repoPath = backingfile
		obj.layerUUID = volumeInfo.layer
		if err != nil {
			return checkoutRet, err
		}
		//	volFlag = true
	}
	if dRet := VerifyBackingFile(obj.repoPath); dRet != OK {
		msg := fmt.Sprintf("backing file '%s' info can not be verified. (ErrCode: %d)", obj.repoPath, dRet)
		return checkoutRet, fmt.Errorf(msg)
	}
	if obj.hardLink {
		print_Log("Create hard link...", logger)
		repoLinkPath := hbDir + path.Base(obj.repoPath)
		configLinkpath := hbDir + path.Base(obj.configPath)
		err = os.Link(obj.repoPath, repoLinkPath)
		if err != nil {
			return checkoutRet, err
		}
		err = os.Link(obj.configPath, configLinkpath)
		if err != nil {
			return checkoutRet, err
		}
		obj.repoPath, obj.configPath = repoLinkPath, configLinkpath
	} else {
		targetRepoPath := hbDir + path.Base(obj.repoPath)
		targetConfigPath := return_BackingFileConfig_Path(&targetRepoPath) // targetRepoPath + ".yaml"
		print_Log(fmt.Sprintf("Start clone repo '%s' to '%s'..", obj.repoPath, targetRepoPath), logger)

		_, err = CopyFile(targetConfigPath, obj.configPath)
		if err != nil {
			if PathFileExists(targetConfigPath) {
				os.Remove(targetConfigPath)
			}
			return checkoutRet, err
		}
		print_Log("set remote tag of each branch.", logger)
		if err = setBranchRemoteTag(&targetConfigPath, "origitn"); err != nil {
			return checkoutRet, err
		}

		pullObj := PullParams{
			//branch:   []string{obj.branch},
			all:            false,
			remoteRepoPath: obj.repoPath,
			localRepoPath:  targetRepoPath,
			configPath:     targetConfigPath,
			protocol:       REPO_PATH_LOCAL,
		}
		if obj.branch == "" {
			err = PullDefaultBranch(&pullObj, logger)
		} else {
			pullObj.branch = obj.branch
			err = PullBranch(&pullObj, logger)
		}
		if err != nil {
			return checkoutRet, err
		}
		print_Log("Add remote origin..", logger)
		err = setRemoteOrigin(&targetConfigPath, &obj.repoPath)
		if err != nil {
			//print_Error(err.Error(), logger)
			return checkoutRet, err
		}
		//	setLocalBranchTag(&pullObj.configPath, []string{pullObj.branch})
		if !obj.checkoutFlg {
			return checkoutRet, nil
		}
		obj.branch = pullObj.branch
		obj.layerUUID = pullObj.pullList[0] //will return branch head commit id
		obj.repoPath, obj.configPath = targetRepoPath, targetConfigPath
	}
	checkoutRet = CheckoutParams{branch: obj.branch, template: obj.repoPath, output: currentDir + path.Base(obj.repoPath), layer: obj.layerUUID}
	print_Log(fmt.Sprintf("Ready to checkout from backing file. (volume name: %s)", checkoutRet.output), logger)

	return checkoutRet, err

}

func clone_Http(obj *CloneParams, logger *log.Logger) (CheckoutParams, error) {

	print_Log("Initializating local hb directory...", logger)
	checkoutRet := CheckoutParams{}
	hbDir, err := hb_Init()
	if err != nil {
		return checkoutRet, err
	}
	print_Log(fmt.Sprintf("Downloading repo's config from url: %s", obj.configPath), logger)
	targetConfigPath, err := downloadConfig(&obj.configPath)
	if err != nil {
		return checkoutRet, err
	}

	print_Log("set remote tag of each branch.", logger)
	if err = setBranchRemoteTag(&targetConfigPath, "origin"); err != nil {
		return checkoutRet, err
	}

	branch, err := return_BranchInfo(&targetConfigPath, obj.branch)
	if err != nil {
		return checkoutRet, err
	}
	pullObj := PullParams{
		pullList:       []string{branch.Head},
		branch:         branch.Name,
		protocol:       REPO_PATH_HTTP,
		all:            false,
		remote:         "origin",
		remoteRepoPath: obj.repoPath,
		configPath:     targetConfigPath,
		localRepoPath:  hbDir + path.Base(obj.repoPath),
	}

	err = PullBranch(&pullObj, logger)
	if err != nil {
		return checkoutRet, err
	}
	print_Log("Add remote origin..", logger)
	err = setRemoteOrigin(&targetConfigPath, &obj.repoPath)
	if err != nil {
		//print_Error(err.Error(), logger)
		return checkoutRet, err
	}
	if !obj.checkoutFlg {
		return checkoutRet, nil
	}
	checkoutRet = CheckoutParams{
		template: pullObj.localRepoPath,
		output:   path.Base(pullObj.localRepoPath),
		layer:    branch.Head,
	}
	//checkoutRet.layer, err = return_LayerUUID(targetPath, obj.layerUUID, true)
	if err != nil {
		return checkoutRet, err
	}
	return checkoutRet, nil

}

func downloadConfig(configPath *string) (string, error) {

	currentDir, err := return_CurrentDir()
	if err != nil {
		return "", err
	}

	targetConfigPath := currentDir + DEFALUT_BACKING_FILE_DIR + "/" + path.Base(*configPath)

	if err = downloadFile(configPath, &targetConfigPath); err != nil {
		return "", err
	}
	return targetConfigPath, nil
}
