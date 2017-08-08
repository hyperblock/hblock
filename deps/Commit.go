package hblock

import (
	"log"

	yaml "gopkg.in/yaml.v2"

	"fmt"

	"github.com/satori/go.uuid"
)

func volume_commit(obj CommitParams, logger *log.Logger) (int, error) {

	userInfo := return_UserInfo()
	if userInfo.name == "" || userInfo.email == "" {
		msg := USER_INFO_EMPTY
		//	print_Error(msg, logger)
		return FAIL, fmt.Errorf(msg)
	}
	commitInfo := YamlCommitMsg{
		Name:    userInfo.name,
		Email:   userInfo.email,
		Message: obj.commitMsg,
	}
	byteCommitInfo, err := yaml.Marshal(&commitInfo)

	if err != nil {
		msg := "Marshal commit info failed."
		//	print_Error(msg, logger)
		return FAIL, fmt.Errorf(msg)
	}
	obj.commitMsg = string(byteCommitInfo)
	if obj.genUUID {
		obj.layerUUID = fmt.Sprintf("%s", uuid.NewV4())
		print_Log("Generate uuid: "+obj.layerUUID, logger)
	} else {
		print_Log("UUID set by manual: "+obj.layerUUID, logger)
	}

	print_Log("Confirm backingfile's format...", logger)
	//	backingFilePath, err := return_Volume_BackingFile(&obj.volumeName)
	if err != nil {
		return FAIL, err
	}
	//	backingFileConfig := return_BackingFileConfig_Path(&backingFilePath)
	//h, err := CreateHBM(FMT_UNKNOWN, &backingFileConfig)
	h, err := CreateHBM_fromExistVol(obj.volumeName)
	if err != nil {
		return FAIL, err
	}
	print_Log("Commit volume...", logger)
	err = h.Commit(&obj)
	if err != nil {
		return FAIL, err
	}
	print_Log("Update branch info.", logger)
	volumeConfigPath := return_Volume_ConfigPath(&obj.volumeName)
	if !PathFileExists(volumeConfigPath) {
		msg := fmt.Sprintf("Volume config file '%s' can not found.", volumeConfigPath)
		//print_Error(msg, logger)
		return FAIL, fmt.Errorf(msg)
	}
	yamlVolumeConfig := YamlVolumeConfig{}
	err = LoadConfig(&yamlVolumeConfig, &volumeConfigPath)
	if err != nil {
		//	print_Error(err.Error(), logger)
		return FAIL, err
	}
	fmt.Println(volumeConfigPath, yamlVolumeConfig.Branch)
	if yamlVolumeConfig.NewBranch {
		branch := YamlBranch{
			Name:  yamlVolumeConfig.Branch,
			Head:  obj.layerUUID,
			Local: 1,
		}
		backingFilePath, _ := return_Volume_BackingFile(&obj.volumeName)
		backingFileConfigPath := return_BackingFileConfig_Path(&backingFilePath) // backingFilePath + ".yaml"
		print_Log(fmt.Sprintf("Set branch '%s' head at '%s'", branch.Name, branch.Head), logger)
		err = add_Branch(&branch, &backingFileConfigPath)
		yamlVolumeConfig.NewBranch = false
	} else {
		err = reset_BranchHead(obj)
	}
	if err != nil {
		return FAIL, err
	}
	WriteConfig(&yamlVolumeConfig, &volumeConfigPath)
	print_Log(
		Format_Success("Done."), logger)

	if WAIT_CHANGE_LAYER == 1 {
		checkoutObj := CheckoutParams{
			volume: obj.volumeName,
			layer:  obj.layerUUID,
			output: obj.volumeName,
		}
		volume_checkout(&checkoutObj, logger)
	}
	return OK, nil
}
