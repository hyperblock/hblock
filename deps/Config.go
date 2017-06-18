package hblock

import (
	"fmt"
	"os"
	"reflect"

	yaml "gopkg.in/yaml.v2"
)

func return_UserInfo() UserInfo {

	userInfo := UserInfo{}
	configPath := return_hb_ConfigPath()
	fileInfo, err := os.Stat(configPath)
	if err != nil {
		//fmt.Println(err.Error())
		return userInfo
	}
	buffer := make([]byte, fileInfo.Size())
	f, err := os.Open(configPath)
	if err != nil {
		//fmt.Println(err.Error())
		return userInfo
	}
	defer f.Close()
	_, err = f.Read(buffer)
	if err != nil {
		//	fmt.Println(err.Error())
		return userInfo
	}
	obj := GlobalConfig{}
	err = yaml.Unmarshal(buffer, &obj)
	if err != nil {
		//		fmt.Println(err.Error())
		return userInfo
	}
	//	fmt.Println(obj)
	userInfo.email = obj.UserEmail
	userInfo.name = obj.UserName
	return userInfo
}

func LoadConfig(ret interface{}, configPath *string) error {

	fileInfo, err := os.Stat(*configPath)
	if err != nil {
		return err
	}
	buffer := make([]byte, fileInfo.Size())
	f, err := os.Open(*configPath)
	if err != nil {
		//msg:="Open file error. o"
		return err
	}
	defer f.Close()
	_, err = f.Read(buffer)
	if err != nil {
		//	fmt.Println(err.Error())
		return err
	}
	// fmt.Println("hahahahah")

	switch ret := ret.(type) {
	case *YamlVolumeConfig:
		err = yaml.Unmarshal(buffer, &ret)
	case *YamlBackingFileConfig:
		err = yaml.Unmarshal(buffer, &ret)
	case *GlobalConfig:
		err = yaml.Unmarshal(buffer, &ret)
	default:
		err = fmt.Errorf("Unassert type: %s", reflect.TypeOf(ret))
	}
	if err != nil {
		msg := fmt.Sprintf("Config unmarshal failed. ( %s )", err.Error())
		return fmt.Errorf(msg)
	}

	//	return ret, nil
	return nil
}

//write global config or backing_file config to yaml file
func WriteConfig(configObj interface{}, configPath *string) error {

	tmpPath := *configPath + ".bak"

	buffer, err := yaml.Marshal(configObj)
	if err != nil {
		return err
	}
	//configPath := return_hb_ConfigPath()
	file, err := os.OpenFile(tmpPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	_, err = file.Write(buffer)
	if err != nil {
		file.Close()
		return err
	}
	file.Close()
	if PathFileExists(*configPath) {
		err = os.Remove(*configPath)
		if err != nil {
			msg := fmt.Sprintf("Remove old config failed ( %s )", err.Error())
			return fmt.Errorf(msg)
		}
	}
	err = os.Rename(tmpPath, *configPath)
	if err != nil {
		msg := fmt.Sprintf("Replace old config failed ( %s )", err.Error())
		return fmt.Errorf(msg)
	}
	return nil
}
