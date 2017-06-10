package hblock

import "log"
import "os/exec"
import "fmt"

func volume_Rebase(obj *RebaseParams, logger *log.Logger) (int, error) {

	print_Trace(fmt.Sprintf("rebase volume '%s' ( -> %s )", obj.volumePath, obj.backingfile))
	print_Log("\rRebase volume...", logger)

	_, err := return_VolumeInfo(&obj.volumePath)
	if err != nil {
		print_Log("\rRebase volume...FAIL\n", logger)
		return FAIL, fmt.Errorf("Load volume info failed. ( %s )", err.Error())
	}

	backingfileInfo := fmt.Sprintf("qcow2://%s?layer=%s", obj.backingfile, obj.parentLayer)
	cmdArgs := []string{"rebase", "-u", "-b", backingfileInfo, obj.volumePath}
	cmd := exec.Command("qemu-img", cmdArgs[0:]...)
	_, err = cmd.Output()
	if err != nil {
		print_Log("\rRebase volume...FAIL\n", logger)
		return 0, err
	}
	print_Log("\rRebase volume...OK\n", logger)
	return OK, nil
}
