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
		print_Error(err.Error(), logger)
		return FAIL, err
	}
	//if PathFileExists(obj.name) {
	if VerifyBackingFile(obj.name) == OK {
		print_Error("Already exist.", logger)
		return FAIL, nil
	}
	print_Log("Create backing file config file.", logger)
	configPath := obj.name + ".yaml"
	yamlConfig := YamlBackingFileConfig{
		Name: path.Base(obj.name),
	}
	err = WriteConfig(&yamlConfig, &configPath)
	if err != nil {
		print_Error(err.Error(), logger)
		return FAIL, err
	}
	g, errno := guestfs.Create()
	if errno != nil {
		os.Remove(configPath)
		return FAIL, errno
	}

	//	defer
	//fmt.Println(size)
	if errCreate := g.Disk_create(obj.name, "qcow2", obj.size, nil); errCreate != nil {
		//return FAIL, errCreate
		g.Close()
		os.Remove(configPath)
		print_Panic(errCreate.Errmsg, logger)
	}
	g.Close()
	msg := fmt.Sprintf("Create template '%s' finished.", obj.name)
	print_Log(Format_Success(msg), logger)
	print_Log("Creating volume named "+obj.output, logger)

	checkoutObj := CheckoutParams{layer: "", output: obj.output, template: obj.name}
	ret, err := volume_checkout(&checkoutObj, logger)
	if ret == OK {
		// yamlConfig := YamlBackingFileConfig{}
		// yamlConfig.Name = path.Base(obj.name)
		// branch := YamlBranch{Name: "master", Local: 1, head_layer: ""}
		// yamlConfig.Branch = append(yamlConfig.Branch, branch)
		// configPath := obj.name + ".yaml"
		// err = WriteConfig(&yamlConfig, &configPath)
		// if err != nil {
		// 	msg := fmt.Sprintf("Write backingfile config failed. ( %s )", err.Error())
		// 	print_Error(msg, logger)
		// 	os.Remove(obj.name)
		// 	return FAIL, fmt.Errorf(msg)
		// }
		checkoutObj.branch = "master"
		checkoutObj.volume = obj.output
		return volume_checkout(&checkoutObj, logger)
	}
	return OK, nil
}
