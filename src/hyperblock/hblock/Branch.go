package hblock

import (
	"fmt"
	"log"
)

func volume_Branch(obj *BranchParams, logger *log.Logger) (int, error) {

	backingfile := obj.backingfile
	layer := ""
	if obj.volumePath != "" {
		volumeInfo, err := return_VolumeInfo(&obj.volumePath)
		if err != nil {
			//	print_Error(err.Error(), logger)
			return FAIL, err
		}
		backingfile = volumeInfo.backingFile
		layer = volumeInfo.layer
	}
	configPath := return_BackingFileConfig_Path(&backingfile) //backingfile + ".yaml"
	err := LoadConfig(&obj.backingFileConfig, &configPath)
	if err != nil {
		return FAIL, fmt.Errorf("Load backing file '%s' config failed. ( %s )", configPath, err.Error())
	}
	if obj.optTag == BRANCH_OPT_SHOW {
		return show_Branch(obj, &layer, logger)
	}
	if obj.optTag == BRANCH_OPT_MV {
		print_Log(fmt.Sprintf("Move branch '%s' to '%s'", obj.move.src, obj.move.dst), logger)
		if err = move_Branch(obj); err != nil {
			return FAIL, nil
		}
		if err = WriteConfig(obj.backingFileConfig, &configPath); err != nil {
			return FAIL, nil
		}
	}
	return OK, nil

}

func move_Branch(obj *BranchParams) error {

	for i := range obj.backingFileConfig.Branch {
		branch := &obj.backingFileConfig.Branch[i]
		if branch.Name == obj.move.src {
			branch.Name = obj.move.dst
			return nil
		}
	}
	return fmt.Errorf("Branch '%s' not found.", obj.move.src)
}

func show_Branch(obj *BranchParams, _layer *string, logger *log.Logger) (int, error) {

	yamlConfig := obj.backingFileConfig
	//branchs := yamlConfig.Branch
	msg := ""
	found := false
	layer := *_layer
	remoteBranch := []string{}
	for _, item := range yamlConfig.Branch {
		info := ""
		if layer == item.Head && item.Local == 1 {
			info += "* "
		} else {
			info += "  "
		}

		if item.Local == 1 {
			if info == "  " {
				msg += fmt.Sprintf("%s%s\n", info, item.Name)
			} else {
				msg += fmt.Sprintf("%s%s\n", info, yellow(item.Name))
				found = true
			}

		} else if obj.show_all {
			remoteBranch = append(remoteBranch, fmt.Sprintf("%sRemotes/%s/%s\n", info, item.Remote, item.Name))
		}
	}
	if !found {
		msg = "* " + yellow(fmt.Sprintf("(Head detached at '%s')\n", layer[0:7])) + msg
	}
	for _, branch := range remoteBranch {
		msg += red(branch)
	}
	print_Log(Format_Success("-------- Result Info -------------:\n")+msg, logger)
	return OK, nil
}

func add_Branch(branch *YamlBranch, configPath *string) error {

	yamlConfig := YamlBackingFileConfig{}
	err := LoadConfig(&yamlConfig, configPath)
	if err != nil {
		return err
	}
	yamlConfig.Branch = append(yamlConfig.Branch, *branch)
	//fmt.Printf("%v\n", yamlConfig.Branch)
	err = WriteConfig(&yamlConfig, configPath)
	if err != nil {
		return err
	}
	return nil
}

func reset_BranchHead(obj CommitParams) error {

	backingFilePath, _ := return_Volume_BackingFile(&obj.volumeName)
	jsonBackingFile, _ := return_JsonBackingFile(&backingFilePath)
	layerList := return_LayerList(&jsonBackingFile)
	parentUUID := ""
	for _, layer := range layerList {
		if layer.uuid == obj.layerUUID {
			parentUUID = layer.parent_uuid
			break
		}
	}
	configPath := return_BackingFileConfig_Path(&backingFilePath) //backingFilePath + ".yaml"
	yamlConfig := YamlBackingFileConfig{}
	err := LoadConfig(&yamlConfig, &configPath)
	if err != nil {
		return err
	}
	found := false
	for i := 0; i < len(yamlConfig.Branch); i++ {
		branch := &yamlConfig.Branch[i]
		if branch.Head == parentUUID {
			branch.Head = obj.layerUUID
			found = true
			break
		}
	}
	if found {
		err = WriteConfig(&yamlConfig, &configPath)
		return err
	}
	return nil
}
