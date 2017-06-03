package hblock

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/user"
	"reflect"
	"strconv"
	"strings"
	"time"

	"path"

	yaml "gopkg.in/yaml.v2"
)

func print_Log(msg string, logger *log.Logger) {

	if len(msg) == 0 {
		return
	}
	if msg[0] != 27 {
		msg = Format_Info(msg)
	}
	fmt.Println(msg)
	if logger != nil {
		msg = msg[15:]
		logger.Println(msg)
	}
}

func print_Error(msg string, logger *log.Logger) {

	if len(msg) == 0 {
		return
	}
	msg = Format_Error(msg)
	fmt.Println(msg)
	if logger != nil {
		msg = msg[15:]
		logger.Println(msg)
	}
}

func print_Panic(msg string, logger *log.Logger) {

	if len(msg) == 0 {
		return
	}
	msg = Format_Error(msg)
	fmt.Println(msg)
	if logger != nil {
		msg = msg[15:]
		logger.Panicln(msg)
	}
}

func print_Fatal(msg string, logger *log.Logger) {

	if len(msg) == 0 {
		return
	}
	msg = Format_Error(msg)
	fmt.Println(msg)
	if logger != nil {
		msg = msg[15:]
		logger.Fatalln(msg)
	}
}

func print_Trace(a ...interface{}) {

	// msg := Format_Trace(fmt.Sprint(a))
	// fmt.Println(msg)
}

func print_ProcessBar(current, total int64) string {

	bar := "["
	base := int((float32(current) / float32(total)) * 100)
	delta := int(float32(base)/float32(5) + 0.5)
	for i := 0; i < delta; i++ {
		bar += "="
	}
	delta = 20 - delta
	for i := 0; i < delta; i++ {
		bar += " "
	}
	bar += "]"
	A, B := current>>20, total>>20
	if A == 0 {
		A = 1
	}
	if B == 0 {
		B = 1
	}
	ret := fmt.Sprintf("%s %d%% (%d/%d)", bar, base, A, B)
	return ret
}

func get_StringAfter(content string, prefix string) string {

	//	print_Trace(fmt.Sprintf("get_StringAfter( %s, %s )", content, prefix))
	q := strings.Index(content, prefix)
	if q == -1 {
		return content
	}
	return content[q+len(prefix):]
}

func get_StringBefore(content string, suffix string) string {

	q := strings.Index(content, suffix)
	if q == -1 {
		return content
	}
	return content[:q]
}

func get_InfoValue(list []string, keyword string) string {

	for _, content := range list {
		exist := strings.HasPrefix(content, keyword)
		if exist {
			return content[len(keyword)+2:]
		}
	}
	return ""
}

func return_CurrentDir() (string, error) {

	pwd := exec.Command("pwd")
	ret, err := pwd.Output()
	if err != nil {
		return "", err
	}
	dir := string(ret)
	return dir[:len(dir)-1] + "/", nil
}

func return_TemplateDir() (string, error) {

	path, err := exec.Command("pwd").Output()
	if err != nil {
		return "", err
	}
	ret := string(path[:len(path)-1]) + "/" + DEFALUT_BACKING_FILE_DIR
	print_Trace("template dir: " + ret)
	return ret, nil
}

func return_hb_ConfigPath() string {

	usr, err := user.Current()
	if err != nil {
		return ""
	}
	hb_configPath := usr.HomeDir + "/.hb/config.yaml"
	if PathFileExists(hb_configPath) {
		return hb_configPath
	}
	return ""
}

func return_Volume_ConfigPath(volumeName *string) string {

	if PathFileExists(*volumeName) == false {
		dir, _ := return_CurrentDir()
		if dir == "" {
			return ""
		}
		ret := fmt.Sprintf("%s.v_%s.yaml", dir, *volumeName)
		return ret
	} else {
		ret := fmt.Sprintf("%s/.v_%s.yaml", path.Dir(*volumeName), path.Base(*volumeName))
		return ret
	}

}

func confirm_BackingFilePath(imgPath string) (string, error) {

	if (len(imgPath) > 0) && (imgPath[0] == '/') {
		_, err := os.Stat(imgPath)
		if err != nil {
			return "", fmt.Errorf("Invalid backing file path '%s'", imgPath)
		}
		return imgPath, nil
	}
	path, errTemp := return_TemplateDir()
	if !strings.HasPrefix(imgPath, path) {
		path += "/" + imgPath
	}
	_, err := os.Stat(path)
	if err != nil {
		return "", err
	} else {
		return path, errTemp
	}
}

func PathFileExists(filePath string) bool {

	_, err := os.Stat(filePath)
	if err != nil {
		return false
	}
	return true
}

func VerifyBackingFile(backingfilePath string) int {

	ret := OK
	configPath := backingfilePath + ".yaml"
	if PathFileExists(configPath) == false {
		ret |= BACKINGFILE_CONFIG_NO_FIND
	}
	if PathFileExists(backingfilePath) == false {
		ret |= BACKINGFILE_NO_FIND
	}
	return ret
}

func custom_Args(args []string, addition string) []string {

	procName := os.Args[0]
	ret := args
	ret[0] = procName + " " + ret[0]
	if addition != "" {
		ret[0] += " " + addition
	}
	return ret
}

func return_AbsPath(path string) string {

	if path[0] == '/' {
		return path
	}
	absPath, err := return_CurrentDir()
	if err != nil {
		return path
	}
	absPath += path
	return absPath
}

func return_Size(strSize string) int64 {

	unit := strSize[len(strSize)-1:]
	_size, err := strconv.Atoi(strSize[0 : len(strSize)-1])

	if err != nil {
		return -1
	}
	var sizeI64 int64
	if unit == "M" {
		sizeI64 = int64(_size * 1024 * 1024)
	} else if unit == "G" {
		sizeI64 = int64(_size*1024*1024) * 1024
	}
	return sizeI64
}

func return_JsonBackingFile(backingFilePath string) (JsonBackingFile, error) {

	args := []string{"info", backingFilePath, "--output", "json"}
	cmd := exec.Command("qcow2-img", args[0:]...)
	print_Trace(fmt.Sprintf("qcow2-img info %s --output json", backingFilePath))
	retBytes, err := cmd.Output()
	if err != nil {
		return JsonBackingFile{}, err
	}
	jsonBackingFile := JsonBackingFile{}
	err = json.Unmarshal(retBytes, &jsonBackingFile)
	print_Trace("Json deserialized.")
	if err != nil {
		return JsonBackingFile{}, err
	}
	return jsonBackingFile, nil
}

func return_LayerList(jsonBackingFile *JsonBackingFile) []Layer {

	list := jsonBackingFile.Layers
	ret := []Layer{}
	for _, item := range list {
		layer := Layer{
			id: item.Id, diskSize: item.DiskSize, createDate: time.Unix(item.DateSec, item.DateNSec),
		}
		nameInfo := strings.Split(item.Name, ",")
		layer.uuid = nameInfo[0]
		layer.parent_uuid = nameInfo[1]
		layer.commit_msg = nameInfo[2]
		print_Trace(layer)
		ret = append(ret, layer)
	}
	return ret
}

func return_JsonVolume(volumePath string) (JsonVolume, error) {

	args := []string{"info", volumePath, "--output", "json"}
	cmd := exec.Command("qcow2-img", args[0:]...)
	retBytes, err := cmd.Output()
	if err != nil {
		return JsonVolume{}, err
	}
	jsonVolume := JsonVolume{}
	err = json.Unmarshal(retBytes, &jsonVolume)
	if err != nil {
		return JsonVolume{}, err
	}
	return jsonVolume, nil
}

func convert_to_VolumeInfo(jsonVolume *JsonVolume) VolumeInfo {

	volInfo := VolumeInfo{
		fileName: jsonVolume.Filename, actualSize: jsonVolume.ActualSize, virtualSize: jsonVolume.VirutalSize,
	}
	args := strings.Split(jsonVolume.BackingFile, "?")
	volInfo.backingFile = get_StringAfter(args[0], "qcow2://")
	volInfo.layer = get_StringAfter(args[1], "layer=")
	return volInfo
}

func return_VolumeInfo(volumePath *string) (VolumeInfo, error) {

	jsonVolume, err := return_JsonVolume(*volumePath)
	if err != nil {
		return VolumeInfo{}, fmt.Errorf(
			fmt.Sprintf("Get volumeInfo failed ( %s )", err.Error()))
	}
	volInfo := VolumeInfo{
		fileName: jsonVolume.Filename, actualSize: jsonVolume.ActualSize, virtualSize: jsonVolume.VirutalSize,
	}
	args := strings.Split(jsonVolume.BackingFile, "?")
	volInfo.backingFile = get_StringAfter(args[0], "qcow2://")
	volInfo.layer = get_StringAfter(args[1], "layer=")
	return volInfo, nil
}

func return_Volume_BackingFile(volumePath *string) (string, error) {

	volumeInfo, err := return_VolumeInfo(volumePath)
	if err != nil {
		return "", err
	}
	backingFilePath, err := confirm_BackingFilePath(volumeInfo.backingFile)
	if err != nil {
		return "", err
	}
	if backingFilePath == "" {
		return "", fmt.Errorf(
			fmt.Sprintf("Can't confirm backing file full path."))
	}
	return backingFilePath, nil
}

func return_LayerUUID(backingFilePath string, layerPrefix string, returnLast bool) (string, error) {

	if layerPrefix == "" && !returnLast {
		//msg := fmt.Sprintf("layerPrefix is null and not return last layer.")
		return "", nil
	}
	jsonBackingFile, err := return_JsonBackingFile(backingFilePath)
	if err != nil {
		msg := "Invalid layer_uuid or backing file path."
		return "", fmt.Errorf(msg)
	}
	layers := return_LayerList(&jsonBackingFile)
	if len(layers) == 0 {
		msg := "There're no any layer in backing file."
		return "", fmt.Errorf(msg)
	}
	if returnLast {
		if layerPrefix == "" {
			return layers[len(layers)-1].uuid, nil
		} else {
			msg := "'returnLast' flag is true but LayerPrefix is not null."
			return "", fmt.Errorf(msg)
		}
	}
	cnt := 0
	ret := ""
	for _, item := range layers {
		uuid := item.uuid
		if strings.HasPrefix(uuid, layerPrefix) {
			cnt++
			ret = uuid
		}
	}

	if cnt > 1 {
		return "", fmt.Errorf(
			fmt.Sprintf("There are more than one layers have prefix '%s'", layerPrefix))
	}
	if cnt == 0 {
		backingFileConfig := backingFilePath + ".yaml"
		yamlConfig := YamlBackingFileConfig{}
		err = LoadConfig(&yamlConfig, &backingFileConfig)
		if err != nil {
			msg := fmt.Sprintf("Load backing file config file error. ( %s )", err.Error())
			return "", fmt.Errorf(msg)
		}
		branchList := yamlConfig.Branch
		for _, item := range branchList {
			if item.Name == layerPrefix {
				return item.Head, nil
			}
		}
		return "", fmt.Errorf(
			fmt.Sprintf("Can't get coresponding layer from prefix '%s'", layerPrefix))
	}
	return ret, nil

}

func return_commit_history(jsonBackingFile *JsonBackingFile, p string) []Layer {

	commitList := return_LayerList(jsonBackingFile)
	ret := []Layer{}

	for i := len(commitList) - 1; i >= 0; i-- {
		if commitList[i].uuid != p {
			continue
		}
		//	fmt.Println(commitList[i].uuid)

		ret = append(ret, commitList[i])
		p = commitList[i].parent_uuid
	}
	return ret
}

func return_RepoPath_Type(path string) int {

	path = strings.ToLower(path)
	if strings.HasPrefix(path, "https://") || strings.HasPrefix(path, "http://") {
		return REPO_PATH_HTTP
	}
	if strings.HasPrefix(path, "ssh://") {
		return REPO_PATH_SSH
	}
	return REPO_PATH_LOCAL
}

func CopyFile(dstPath, srcPath string) (int64, error) {

	src, err := os.Open(srcPath)
	if err != nil {
		return 0, err
	}
	defer src.Close()
	dst, err := os.OpenFile(dstPath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return 0, err
	}
	defer dst.Close()
	return io.Copy(dst, src)
}

func hb_Init() (int, error) {

	currentDir, err := return_CurrentDir()
	if err != nil {
		msg := fmt.Sprintf("Initialization failed. (%s)", err.Error)
		return FAIL, fmt.Errorf(msg)
	}
	hbDir := currentDir + "/" + DEFALUT_BACKING_FILE_DIR
	if PathFileExists(hbDir) {
		return OK, nil
	}
	err = os.Mkdir(hbDir, 0744)
	if err != nil {
		msg := fmt.Sprintf("Create init dir '%s' failed. (%s)", hbDir, err.Error())
		return FAIL, fmt.Errorf(msg)
	}
	return OK, nil
}

func return_ConfigValue(configObj interface{}, tag string) (interface{}, error) {

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

func LoadConfig(ret interface{}, configPath *string) error {

	//	print_Log("Loading configuration...", logger)
	//configPath := return_hb_ConfigPath()
	//ret := GlobalConfig{}
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

	// buffer := []byte{}
	// var err error
	// switch configObj := configObj.(type) {
	// case *YamlVolumeConfig:
	// 	buffer, err = yaml.Marshal(configObj)
	// case *YamlBackingFileConfig:
	// 	buffer, err = yaml.Marshal(configObj)
	// case *GlobalConfig:
	// 	buffer, err = yaml.Marshal(configObj)
	// default:
	// 	err = fmt.Errorf("Unassert type: %s", reflect.TypeOf(configObj))
	// 	return err
	// }

	//	print_Log("Updating configuration...", logger)
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

// see complete color rules in document in https://en.wikipedia.org/wiki/ANSI_escape_code#cite_note-ecma48-13
func Format_Trace(format string, a ...interface{}) string {
	prefix := yellow(trac)
	return fmt.Sprint(formatLog(prefix), fmt.Sprintf(format, a...))
}

func Format_Info(format string, a ...interface{}) string {
	prefix := blue(info)
	return fmt.Sprint(formatLog(prefix), fmt.Sprintf(format, a...))
}

func Format_Success(format string, a ...interface{}) string {
	prefix := green(succ)
	return fmt.Sprint(formatLog(prefix), fmt.Sprintf(format, a...))
}

func Format_Warning(format string, a ...interface{}) string {
	prefix := magenta(warn)
	return fmt.Sprint(formatLog(prefix), fmt.Sprintf(format, a...))
}

func Format_Error(format string, a ...interface{}) string {
	prefix := red(erro)
	return fmt.Sprint(formatLog(prefix), fmt.Sprintf(format, a...))
}

func red(s string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", color_red, s)
}

func green(s string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", color_green, s)
}

func yellow(s string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", color_yellow, s)
}

func blue(s string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", color_blue, s)
}

func magenta(s string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", color_magenta, s)
}

func formatLog(prefix string) string {

	return prefix + " "
}
