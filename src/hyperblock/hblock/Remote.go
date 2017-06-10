package hblock

import (
	"fmt"
	"log"
)

func Remote(obj RemoteParams, logger *log.Logger) (int, error) {

	configPath := return_BackingFileConfig_Path(&obj.backingFile) //obj.backingFile + ".yaml"
	configInfo := YamlBackingFileConfig{}
	err := LoadConfig(&configInfo, &configPath)
	if err != nil {
		return FAIL, fmt.Errorf("Load config '%s' failed. (%s)", configPath, err.Error())
		//		print_Error(msg, logger)
		//	return FAIL, err
	}
	if obj.verbose {
		ret := "\n"
		for _, item := range configInfo.Remote {
			ret += (item.Name + ": " + item.Url + "\n")
		}
		ret = Format_Success("Done.") + ret
		print_Log(ret, logger)
		return OK, nil
	}
	if obj.add.name != "" {
		print_Log(fmt.Sprintf("Add a remote host named: %s, url: %s", obj.add.name, obj.add.url), logger)
		for _, item := range configInfo.Remote {
			if item.Name == obj.add.name {
				return FAIL, fmt.Errorf("Host name '%s' exists.", item.Name)
			}
		}
		remoteHost := YamlRemote{Name: obj.add.name, Url: obj.add.url}
		configInfo.Remote = append(configInfo.Remote, remoteHost)
	}
	if obj.rename.oldName != "" {
		print_Log(fmt.Sprintf("Rename host '%s' to '%s'.", obj.rename.oldName, obj.rename.newName), logger)
		found := false
		for i := 0; i < len(configInfo.Remote); i++ {
			if configInfo.Remote[i].Name == obj.rename.oldName {
				configInfo.Remote[i].Name = obj.rename.newName
				found = true
				break
			}
		}
		if !found {
			return FAIL, fmt.Errorf("Remote host '%s' doesn't exist.", obj.remove)
			//	print_Error(msg, logger)
			//	return FAIL, fmt.Errorf(msg)
		}
	}
	if obj.setUrl.name != "" {
		print_Log(fmt.Sprintf("Set host '%s' url as: %s, ", obj.setUrl.name, obj.setUrl.url), logger)
		found := false
		for i := 0; i < len(configInfo.Remote); i++ {
			if configInfo.Remote[i].Name == obj.setUrl.name {
				configInfo.Remote[i].Url = obj.setUrl.url
				found = true
				break
			}
		}
		if !found {
			return FAIL, fmt.Errorf("Remote host '%s' doesn't exist.", obj.remove)
		}
	}
	if obj.remove != "" {
		print_Log(fmt.Sprintf("Remove a remote host from local list named: %s", obj.remove), logger)
		found := false
		for i := 0; i < len(configInfo.Remote); i++ {
			if configInfo.Remote[i].Name == obj.remove {
				configInfo.Remote = append(configInfo.Remote[:i], configInfo.Remote[i+1:]...)
				found = true
				break
			}
		}
		if !found {
			return FAIL, fmt.Errorf("Remote host '%s' doesn't exist.", obj.remove)
		}
	}
	err = WriteConfig(&configInfo, &configPath)
	if err != nil {
		//	print_Error(err.Error(), logger)
		return FAIL, err
	}
	print_Log(Format_Success("done."), logger)
	return OK, nil
}

func setRemoteOrigin(configPath *string, url *string) error {

	config := YamlBackingFileConfig{}
	err := LoadConfig(&config, configPath)
	if err != nil {
		return err
	}
	origin := YamlRemote{Name: "origin", Url: *url}
	config.Remote = []YamlRemote{origin}
	err = WriteConfig(&config, configPath)
	return err
}
