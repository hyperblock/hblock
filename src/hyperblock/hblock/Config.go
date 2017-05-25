package hblock

import (
	"os"

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
