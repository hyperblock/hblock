package hblock

import (
	"fmt"
	"log"
)

func show_Branch_Info(branchParams *BranchParams, logger *log.Logger) (int, error) {

	volumeInfo, err := return_VolumeInfo(&branchParams.volumePath)
	if err != nil {
		//	print_Error(err.Error(), logger)
		return FAIL, err
	}
	backingfile := volumeInfo.backingFile
	layer := volumeInfo.layer
	configPath := return_BackingFileConfig_Path(&backingfile) //backingfile + ".yaml"
	yamlConfig := YamlBackingFileConfig{}
	err = LoadConfig(&yamlConfig, &configPath)
	if err != nil {
		return FAIL, fmt.Errorf("Load backing file '%s' config failed. ( %s )", configPath, err.Error())
	}
	branchs := yamlConfig.Branch
	msg := ""
	found := false
	for _, item := range branchs {
		info := ""
		if layer == item.Head {
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

		} else if branchParams.show_all {
			msg += fmt.Sprintf("%sRemote/%s\n", info, item.Name)
		}
	}
	if !found {
		msg = "* " + yellow(fmt.Sprintf("(Head detached at '%s')\n", layer[0:7])) + msg
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
