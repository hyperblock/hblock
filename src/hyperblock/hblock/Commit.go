package hblock

import (
	"log"

	yaml "gopkg.in/yaml.v2"

	"fmt"

	"os/exec"

	"strings"

	"github.com/satori/go.uuid"
)

func volume_commit(obj CommitParams, logger *log.Logger) (int, error) {

	userInfo := return_UserInfo()
	if userInfo.name == "" || userInfo.email == "" {
		msg := USER_INFO_EMPTY
		print_Error(msg, logger)
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
		print_Error(msg, logger)
		return FAIL, fmt.Errorf(msg)
	}
	if obj.genUUID {
		obj.layerUUID = fmt.Sprintf("%s", uuid.NewV4())
		print_Log("Generate uuid: "+obj.layerUUID, logger)
	} else {
		print_Log("UUID set by manual: "+obj.layerUUID, logger)
	}
	commitArgs := []string{"commit", "-m", string(byteCommitInfo), "-s", obj.layerUUID, obj.volumeName}
	commitCmd := exec.Command("qcow2-img", commitArgs[0:]...)
	print_Log("qcow2-img "+strings.Join(commitArgs, " "), logger)
	result, err := commitCmd.Output()
	if err != nil {
		print_Error(err.Error(), logger)
		return FAIL, err
	}
	print_Log(Format_Success(string(result)), logger)
	// volumeLog := JsonLog{
	// 	Operation:  "commit",
	// 	UUID:       obj.snapshot,
	// 	Info:       obj.commitMsg,
	// 	VolumeName: obj.volumeName,
	// }
	// push_Log(volumeLog, logger)
	return OK, nil
}
