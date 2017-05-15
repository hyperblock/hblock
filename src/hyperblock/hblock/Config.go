package hblock

import (
	"os"
	"reflect"

	"fmt"

	"log"

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

func return_ConfigValue(configObj *GlobalConfig, tag string) (interface{}, error) {

	v := reflect.ValueOf(configObj).Elem()
	for i := 0; i < v.NumField(); i++ {
		fieldInfo := v.Type().Field(i)
		fieldTag := get_StringBefore(get_StringAfter(string(fieldInfo.Tag), "\""), "\"")
		if fieldTag == tag {
			return v.Field(i), nil
		}
		//	maps[fieldInfo.Name] = fieldInfo
		//	fmt.Println(d)
		//fmt.Println(tag)
	}
	return nil, fmt.Errorf("Tag not found.")
}

func LoadConfig(logger *log.Logger) (GlobalConfig, error) {

	print_Log("Loading configuration...", logger)
	configPath := return_hb_ConfigPath()
	ret := GlobalConfig{}
	fileInfo, err := os.Stat(configPath)
	if err != nil {
		return ret, err
	}
	buffer := make([]byte, fileInfo.Size())
	f, err := os.Open(configPath)
	if err != nil {
		//msg:="Open file error. o"
		return ret, err
	}
	defer f.Close()
	_, err = f.Read(buffer)
	if err != nil {
		//	fmt.Println(err.Error())
		return ret, err
	}
	err = yaml.Unmarshal(buffer, &ret)
	if err != nil {
		msg := "Config unmarshal failed. %s" + err.Error()
		return ret, fmt.Errorf(msg)
	}
	return ret, nil
}

func WriteConfig(configObj *GlobalConfig, logger *log.Logger) error {

	print_Log("Updating configuration...", logger)

	buffer, err := yaml.Marshal(configObj)
	if err != nil {
		return err
	}
	configPath := return_hb_ConfigPath()
	file, err := os.OpenFile(configPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(buffer)
	if err != nil {
		return err
	}
	return nil
}
