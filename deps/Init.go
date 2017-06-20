package hblock

import (
	"fmt"

	"log"
	"os"
	"path"
)

func create_empty_template(obj InitParams, logger *log.Logger) (int, error) {

	//output := obj.name
	print_Log("Init hb directory.", logger)
	_, err := hb_Init()
	if err != nil {
		return FAIL, err
	}
	//if PathFileExists(obj.name) {
	if VerifyBackingFile(obj.name) == OK {
		return FAIL, fmt.Errorf("Already exist.")
		//	return FAIL, nil
	}
	fmtTag := FMT_QCOW2
	if obj.format == "lvm" {
		fmtTag = FMT_LVM
	}
	print_Log("Create hyperlayer object...", logger)
	h, err := CreateHyperLayer(fmtTag, nil)
	if err != nil {
		return FAIL, err
	}
	print_Log("Create backing file config file.", logger)
	configPath, err := h.return_BackingFileConfig_Path(&obj.name) //obj.name + ".yaml"
	if err != nil {
		return FAIL, err
	}
	yamlConfig := YamlBackingFileConfig{
		Name:        path.Base(obj.name),
		VirtualSize: obj.size,
		DefaultHead: "master",
		Format:      obj.format,
	}
	err = WriteConfig(&yamlConfig, &configPath)
	if err != nil {
		return FAIL, err
	}
	err = h.CreateDisk(&obj)
	if err != nil {
		os.Remove(configPath)
		return FAIL, err
	}
	msg := fmt.Sprintf("Create template '%s' finished.", obj.name)
	print_Log(Format_Success(msg), logger)
	if !obj.checkout {
		return OK, nil
	}
	// print_Log("Creating volume named "+obj.output, logger)

	// checkoutObj := CheckoutParams{layer: "", output: obj.output, template: obj.name}
	// ret, err := volume_checkout(&checkoutObj, logger)
	// if err != nil {
	// 	return FAIL, err
	// }
	// if ret == OK {
	// 	checkoutObj.branch = "master"
	// 	checkoutObj.volume = obj.output
	// 	return volume_checkout(&checkoutObj, logger)
	// }
	return OK, nil
}
