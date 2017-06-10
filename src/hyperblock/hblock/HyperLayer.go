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
	cmd    string
	args   []string
}

func CreateHyperLayer(fmt_tag int, backingfilePath_Or_format *string) (HyperLayer, error) {

	if fmt_tag == FMT_UNKNOWN {
		if *backingfilePath_Or_format == "qcow2" {
			fmt_tag = FMT_QCOW2

		} else if *backingfilePath_Or_format == "lvm" {
			fmt_tag = FMT_LVM
		} else {
			config := YamlBackingFileConfig{}
			err := LoadConfig(&config, backingfilePath_Or_format)
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
	return ret, ret.check_Command()
}

func (h *HyperLayer) SetArgs(_args []string) {

	h.args = _args
}

func (h *HyperLayer) check_Command() (err error) {

	cmd := ""
	if h.format == FMT_UNKNOWN {
		return fmt.Errorf("Format unknow.")
	}

	if h.format == FMT_QCOW2 {
		cmd, err = exec.LookPath("qcow2-img")
		if err != nil {
			return fmt.Errorf("Command 'qcow2-img' not found. ( %s )", err.Error())
		}
		//h.cmd = cmd
		//return nil
	}
	if h.format == FMT_LVM {
		return fmt.Errorf("LVM command unfinished.")
	}
	h.cmd = cmd

	return nil
}

func (h *HyperLayer) runCommand(args []string) (err error) {

	print_Trace(fmt.Sprintf("%s %s", h.cmd, args))
	commitCmd := exec.Command(h.cmd, args[0:]...)
	result, err := commitCmd.Output()
	if err != nil {
		return err
	}
	print_Trace(string(result))
	return nil
}

func (h *HyperLayer) Commit(obj *CommitParams) error {

	print_Trace("HyperLayer.Commit.")
	commitArgs := []string{"commit", "-m", obj.commitMsg, "-s", obj.layerUUID, obj.volumeName}
	return h.runCommand(commitArgs)
}

func (h HyperLayer) DumpLayer(obj *DumpParams) error {

	print_Trace("HyperLayer.DumpLayer.")
	dumpArgs := []string{"layerdump", "-t", obj.backngFile, "-l", obj.layerUUID, obj.output}
	return h.runCommand(dumpArgs)
}

func (h HyperLayer) Rebase(obj *RebaseParams) error {

	print_Trace("HyperLayer.Rebase.")
	if h.format == FMT_QCOW2 {
		h.cmd = "qemu-img"
		backingfileInfo := fmt.Sprintf("qcow2://%s?layer=%s", obj.backingfile, obj.parentLayer)
		cmdArgs := []string{"rebase", "-u", "-b", backingfileInfo, obj.volumePath}
		return h.runCommand(cmdArgs)
	} else if h.format == FMT_LVM {
		return fmt.Errorf("LVM command unfinished.")
	}
	return nil
}

func (h HyperLayer) Checkout(obj *CheckoutParams) error {

	print_Trace("HyperLayer.Checkout.")
	checkoutArgs := []string{"create", "-t", obj.template, "-l", obj.layer, obj.output}
	return h.runCommand(checkoutArgs)
}
