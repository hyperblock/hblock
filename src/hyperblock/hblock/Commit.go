package hblock

import (
	"log"

	"fmt"

	"os/exec"

	"strings"

	"github.com/satori/go.uuid"
)

type CommitParams struct {
	commitMsg  string
	volumeName string
	snapshot   string
}

func volume_commit(obj CommitParams, logger *log.Logger) (int, error) {

	obj.snapshot = fmt.Sprintf("%s", uuid.NewV4())
	print_Log("Generate uuid: "+obj.snapshot, logger)
	commitArgs := []string{"commit", "-m", obj.commitMsg, "-s", obj.snapshot, obj.volumeName}
	commitCmd := exec.Command("qcow2-img", commitArgs[0:]...)
	print_Log("qcow2-img "+strings.Join(commitArgs, " "), logger)
	result, err := commitCmd.Output()
	if err != nil {
		print_Error(err.Error(), logger)
		return FAIL, err
	}
	print_Log(format_Success(string(result)), logger)
	return OK, nil
}
