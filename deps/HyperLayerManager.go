package hblock

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/hyperblock/lvdiff/lvbackup/lvmutil"

	"github.com/hyperblock/hblock/deps/guestfs"
)

const (
	FMT_UNKNOWN = 0
	FMT_QCOW2   = 1
	FMT_LVM     = 2
)

type HBM struct {
	format        int
	cmd           string
	args          []string
	location      string
	volConfigPath string
	volInfo       YamlVolumeConfig
}

func CreateHBM_fromExistVol(volPath string) (HBM, error) {

	ret := HBM{}
	configPath := return_Volume_ConfigPath(&volPath)
	config := YamlVolumeConfig{}
	err := LoadConfig(&config, &volPath)
	ret.volConfigPath = volPath
	if err != nil {
		return HBM{}, nil
	}
	if config.Format == "qcow2" {
		ret.format = FMT_QCOW2
	} else if config.Format == "lvm" {
		ret.format = FMT_LVM
	} else {
		return HBM{}, fmt.Errorf("Can't confirm Volume's format.")
	}
	return ret, nil
}

func CreateHBM(fmt_tag int, fileFullName string) (HBM, error) {

	ret := HBM{}
	// if fmt_tag == FMT_UNKNOWN {
	// 	if *backingfilePath_Or_format == "qcow2" {
	// 		fmt_tag = FMT_QCOW2

	// 	} else if *backingfilePath_Or_format == "lvm" {
	// 		fmt_tag = FMT_LVM
	// 	} else {
	// 		config := YamlBackingFileConfig{}
	// 		err := LoadConfig(&config, backingfilePath_Or_format)
	// 		ret.backingfileConfig = *backingfilePath_Or_format
	// 		if err != nil {
	// 			return HBM{}, nil
	// 		}
	// 		if config.Format == "qcow2" {
	// 			fmt_tag = FMT_QCOW2
	// 		} else if config.Format == "lvm" {
	// 			fmt_tag = FMT_LVM
	// 		} else {
	// 			return HBM{}, fmt.Errorf("Can't confirm backingfile's format.")
	// 		}
	// 	}
	// }
	ret.format = fmt_tag
	ret.location = return_AbsPath(fileFullName)

	return ret, ret.check_Command()
}

func (h *HBM) SetArgs(_args []string) {

	h.args = _args
}

func (h *HBM) check_Command() (err error) {

	cmd := ""
	if h.format == FMT_UNKNOWN {
		return fmt.Errorf("Format unknow.")
	}

	if h.format == FMT_QCOW2 {
		cmd, err = exec.LookPath("qcow2-img")
		if err != nil {
			return fmt.Errorf("Command 'qcow2-img' not found. ( %s )", err.Error())
		}
		//h.cmd = cmds
		//return nil
	}
	// if h.format == FMT_LVM {
	// 	return fmt.Errorf("LVM command unfinished.")
	// }
	h.cmd = cmd

	return nil
}

func (h *HBM) runCommand(args []string) (err error) {

	print_Trace(fmt.Sprintf("%s %s", h.cmd, args))

	cmd := exec.Command(h.cmd, args[0:]...)
	result, err := cmd.CombinedOutput()
	if err != nil {
		if result != nil {
			return fmt.Errorf("%s (%s)", string(result), err.Error())
		}
		return err
	}
	print_Trace(string(result))
	return nil
}

func (h *HBM) Commit(obj *CommitParams) error {

	print_Trace("HyperLayer.Commit.")

	commitArgs := []string{"commit", "-m", obj.commitMsg, "-s", obj.layerUUID, obj.volumeName}
	return h.runCommand(commitArgs)
}

func (h HBM) DumpLayer(obj *DumpParams) error {

	print_Trace("HyperLayer.DumpLayer.")
	dumpArgs := []string{"layerdump", "-t", obj.backngFile, "-l", obj.layerUUID, obj.output}
	return h.runCommand(dumpArgs)
}

func (h HBM) Rebase(obj *RebaseParams) error {

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

func (h HBM) Checkout(obj *CheckoutParams) error {

	print_Trace("HyperLayer.Checkout.")

	if h.format == FMT_QCOW2 {
		args := []string{"create", "-t", obj.template, "-l", obj.layer, obj.output}
		return h.runCommand(args)
	} else if h.format == FMT_LVM {
		//vgPath:=obj.template
		if err := lvmutil.CreateSnapshotLv(obj.template, obj.layer, obj.output); err != nil {
			return err
		}
		return lvmutil.ActivateLv(obj.template, obj.output)
	}

	return errors.New("format unknow.")
}

func (h HBM) CreateDisk(obj *InitParams) error {

	if h.format == FMT_UNKNOWN {
		return errors.New("format unknow.")
	}
	if h.format == FMT_QCOW2 {
		g, errno := guestfs.Create()
		if errno != nil {
			return errno
		}
		defer g.Close()
		if errCreate := g.Disk_create(obj.name, "qcow2", obj.size, nil); errCreate != nil {
			return fmt.Errorf(errCreate.Errmsg)
		}
	}
	if h.format == FMT_LVM {
		volPath := obj.output
		if strings.Index(obj.output, "/dev") != 0 {
			volPath = "/dev/" + volPath
		}
		token := strings.Split(volPath, "/")
		//token = token[2:]
		if len(token) < 5 {
			return fmt.Errorf("invalid volume path. need 'VolumeGroupName[Path]/ThinPoolLogicalVolume/LogicalVolume' ")
		}
		vg := token[2]
		pool := token[3]
		vol := token[4]
		err := lvmutil.CreateThinLv(vg, pool, vol, obj.size)
		return err
	}
	return nil
}

func (h *HBM) return_BackingFileConfig_Path(_path *string) (string, error) {

	if h.format == FMT_UNKNOWN {
		return "", fmt.Errorf("format unknow.")
	}

	return h.volInfo.BackingFile, nil
	// if h.format == FMT_QCOW2 {
	// 	if h.backingfileConfig == "" {
	// 		h.backingfileConfig = *_path + ".yaml"
	// 	}
	// }
	// if h.format == FMT_LVM {
	// 	if h.backingfileConfig == "" {
	// 		volname := path.Base(*_path)
	// 		usr, err := user.Current()
	// 		hb_Dir := usr.HomeDir + "/.hb"
	// 		info := hb_Dir + "/lvinfo"
	// 		if PathFileExists(info) == false {
	// 			err = os.Mkdir(info, 0755)
	// 			if err != nil {
	// 				return "", fmt.Errorf("create dir '%s' faild. %s", info, err.Error())
	// 			}
	// 		}
	// 		h.backingfileConfig = info + "/" + volname
	// 	}
	// }
	//return h.backingfileConfig, nil
}
