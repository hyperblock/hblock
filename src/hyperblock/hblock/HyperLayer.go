package hblock

import "fmt"
import "os/exec"

const (
	FMT_UNKNOWN = 0
	FMT_QCOW2   = 1
	FMT_LVM     = 2
)

type HyperLayer struct {
	format int
}

func CreateHyperLayer(fmt_tag int, backingfilePath_Or_Config *string) (HyperLayer, error) {

	if fmt_tag == FMT_UNKNOWN {
		if *backingfilePath_Or_Config == "qcow2" {
			fmt_tag = FMT_QCOW2
		} else if *backingfilePath_Or_Config == "lvm" {
			fmt_tag = FMT_LVM
		} else {
			config := YamlBackingFileConfig{}
			err := LoadConfig(&config, backingfilePath_Or_Config)
			if err != nil {
				return HyperLayer{}, nil
			}
			if config.Format == "qcow2" {
				fmt_tag = FMT_QCOW2
			} else if config.Format == "lvm" {
				fmt_tag = FMT_LVM
			} else {
				return HyperLayer{}, fmt.Errorf("Can't confirm backingfile's format.")
			}
		}
	}
	ret := HyperLayer{format: fmt_tag}
	return ret, nil
}

func (h HyperLayer) return_Command() (string, error) {

	if h.format == FMT_UNKNOWN {
		return "", fmt.Errorf("Format unknow.")
	}
	if h.format == FMT_QCOW2 {
		cmd, err := exec.LookPath("qcow2-img")
		if err != nil {
			return "", fmt.Errorf("Command 'qcow2-img' not found. ( %s )", err.Error())
		}
		return cmd, nil
	}
	if h.format == FMT_LVM {
	}
	return "", nil
}

func (h HyperLayer) Commit(obj *CommitParams) error {

	print_Trace("HyperLayer.Commit.")
	//	exec.LookPath("vgcfgbackup")
	cmd, err := h.return_Command()
	if err != nil {
		return err
	}
	commitArgs := []string{"commit", "-m", obj.commitMsg, "-s", obj.layerUUID, obj.volumeName}
	commitCmd := exec.Command(cmd, commitArgs[0:]...)
	result, err := commitCmd.Output()
	if err != nil {
		return err
	}
	print_Trace(string(result))
	return nil
}

func (h HyperLayer) DumpLayer(obj *DumpParams) error {

	print_Trace("HyperLayer.DumpLayer.")
	return nil
}
