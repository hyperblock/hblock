package hblock

import (
	"fmt"

	"libguestfs.org/guestfs"

	"log"
	"os"
	"path"
)

func create_empty_template(obj InitParams, logger *log.Logger) (int, error) {

	//output := obj.name
	print_Log("Init hb directory.", logger)
	_, err := hb_Init()
	if err != nil {
		//	print_Error(err.Error(), logger)
		return FAIL, err
	}
	//if PathFileExists(obj.name) {
	if VerifyBackingFile(obj.name) == OK {
		return FAIL, fmt.Errorf("Already exist.")
		//	return FAIL, nil
	}
	print_Log("Create backing file config file.", logger)
	configPath := return_BackingFileConfig_Path(&obj.name) //obj.name + ".yaml"
	yamlConfig := YamlBackingFileConfig{
		Name:        path.Base(obj.name),
		VirtualSize: obj.size,
		DefaultHead: "master",
		Format:      obj.format,
	}
	err = WriteConfig(&yamlConfig, &configPath)
	if err != nil {
		//print_Error(err.Error(), logger)
		return FAIL, err
	}
	err = createDisk(obj)
	if err != nil {
		os.Remove(configPath)
		//print_Error(err.Error(), logger)
		return FAIL, err
	}
	msg := fmt.Sprintf("Create template '%s' finished.", obj.name)
	print_Log(Format_Success(msg), logger)
	if !obj.checkout {
		return OK, nil
	}
	print_Log("Creating volume named "+obj.output, logger)

	checkoutObj := CheckoutParams{layer: "", output: obj.output, template: obj.name}
	ret, err := volume_checkout(&checkoutObj, logger)
	if err != nil {
		return FAIL, err
	}
	if ret == OK {
		checkoutObj.branch = "master"
		checkoutObj.volume = obj.output
		return volume_checkout(&checkoutObj, logger)
	}
	return OK, nil
}

func createDisk(obj InitParams) error {

	g, errno := guestfs.Create()
	if errno != nil {
		return errno
	}
	defer g.Close()
	if errCreate := g.Disk_create(obj.name, "qcow2", obj.size, nil); errCreate != nil {
		return fmt.Errorf(errCreate.Errmsg)
	}
	return nil
}
