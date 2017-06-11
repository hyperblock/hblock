package hblock

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"reflect"
	"strconv"
	"strings"
	"time"

	"path"
)

// func CreateCLI(_log *log.Logger) *HyperblockCLI {

// 	obj := &HyperblockCLI{
// 		logger: _log,
// 	}
// 	obj.currentDir, obj.err = return_CurrentDir()
// 	return obj
// }

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

	path, err := return_CurrentDir() //exec.Command("pwd").Output()
	if err != nil {
		return "", err
	}
	ret := path + DEFALUT_BACKING_FILE_DIR
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

	absPath := return_AbsPath(*volumeName)
	baseName := path.Base(absPath)
	dir := path.Dir(absPath)
	ret := fmt.Sprintf("%s/.v_%s.yaml", dir, baseName)
	return ret

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
	absImagePath := return_AbsPath(imgPath)
	if !strings.HasPrefix(absImagePath, path) {
		path += "/" + imgPath
	} else {
		path = absImagePath
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
	configPath := return_BackingFileConfig_Path(&backingfilePath) //backingfilePath + ".yaml"
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

func return_JsonBackingFile(backingFilePath *string) (JsonBackingFile, error) {

	args := []string{"info", *backingFilePath, "--output", "json"}
	cmd := exec.Command("qcow2-img", args[0:]...)
	print_Trace(fmt.Sprintf("qcow2-img info %s --output json", *backingFilePath))
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
	jsonBackingFile, err := return_JsonBackingFile(&backingFilePath)
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
		backingFileConfig := return_BackingFileConfig_Path(&backingFilePath) //backingFilePath + ".yaml"
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

func hb_Init() (string, error) {

	currentDir, err := return_CurrentDir()
	if err != nil {
		msg := fmt.Sprintf("Initialization failed. (%s)", err.Error())
		return "", fmt.Errorf(msg)
	}
	hbDir := currentDir + "/" + DEFALUT_BACKING_FILE_DIR
	if PathFileExists(hbDir) {
		return hbDir + "/", nil
	}
	err = os.Mkdir(hbDir, 0744)
	if err != nil {
		msg := fmt.Sprintf("Create init dir '%s' failed. (%s)", hbDir, err.Error())
		return "", fmt.Errorf(msg)
	}
	return hbDir + "/", nil
}

func return_ConfigValue(configObj interface{}, tag string) (interface{}, error) {

	v := reflect.ValueOf(configObj).Elem()
	for i := 0; i < v.NumField(); i++ {
		fieldInfo := v.Type().Field(i)
		fieldTag := get_StringBefore(get_StringAfter(string(fieldInfo.Tag), "\""), "\"")
		if fieldTag == tag {
			return v.Field(i), nil
		}
	}
	return nil, fmt.Errorf("Tag not found.")
}

func return_Volume_BackingFile_Config(volumePath *string) (YamlBackingFileConfig, error) {

	ret := YamlBackingFileConfig{}
	backingFile, err := return_Volume_BackingFile(volumePath)
	if err != nil {
		return ret, err
	}
	if dRet := VerifyBackingFile(backingFile); dRet != OK {
		return ret, fmt.Errorf("Verify Backingfile failed. ( ErrCode %d )", dRet)
	}
	configPath := return_BackingFileConfig_Path(&backingFile)
	err = LoadConfig(&ret, &configPath)
	if err != nil {
		return ret, err
	}
	return ret, nil
}

func return_RemoteUrl(remotes *[]YamlRemote, name *string) string {

	for _, item := range *remotes {
		if item.Name == *name {
			return item.Url
		}
	}
	return ""
}

func return_BranchHead(branchName *string, branches *[]YamlBranch) string {

	for _, item := range *branches {
		if item.Name == *branchName {
			return item.Head
		}
	}
	return ""
}

func return_BranchInfo(backingfileConfig *string, branchName string) (YamlBranch, error) {

	info := YamlBackingFileConfig{}
	err := LoadConfig(&info, backingfileConfig)
	if err != nil {
		return YamlBranch{}, err
	}
	if branchName == "" {
		branchName = info.DefaultHead
	}
	for _, branch := range info.Branch {
		if branch.Name == branchName {
			return branch, nil
		}
	}
	return YamlBranch{}, fmt.Errorf("Branch '%s' not found.", branchName)
}

func trace_Parents(backingFilePath *string, head *string) ([]string, error) {

	jsonBackingFile, err := return_JsonBackingFile(backingFilePath)
	if err != nil {
		return nil, err
	}
	layerList := return_LayerList(&jsonBackingFile)
	pHead := *head
	ret := []string{pHead}
	for i := len(layerList) - 1; i > 0; i-- {
		if layerList[i].uuid == pHead {
			pHead = layerList[i].parent_uuid
			ret = append(ret, pHead)
		}
	}
	return ret, nil
}

func return_BackingFileConfig_Path(path *string) string {
	return *path + ".yaml"
}

// func DumpLayer(obj *DumpParams) error {

// 	dumpArgs := []string{"layerdump", "-t", obj.backngFile, "-l", obj.layerUUID, obj.output}
// 	cmdDump := exec.Command("qcow2-img", dumpArgs[0:]...)
// 	print_Trace(dumpArgs)
// 	_, err := cmdDump.Output()
// 	return err
// }

func return_CommitInfo(layerPath *string) (CommitParams, error) {

	layerUUID := func() string {
		p := strings.LastIndex(*layerPath, ".")
		return (*layerPath)[p+1:]
	}()
	ret := CommitParams{layerUUID: layerUUID, genUUID: false, commitMsg: "test", volumeName: *layerPath}
	return ret, nil
}

func return_LayerName(repoPath, head string) string {
	return repoPath + "." + head
}

func RemoveFiles(files []string) {

	for _, file := range files {
		if PathFileExists(file) {
			os.Remove(file)
		}
	}
}

func setLocalBranchTag(configPath *string, branch *string) error {

	config := YamlBackingFileConfig{}
	err := LoadConfig(&config, configPath)
	if err != nil {
		return err
	}

	for i := 0; i < len(config.Branch); i++ {
		ok := config.Branch[i].Name == *branch
		if ok {
			config.Branch[i].Local = 1
		}
	}
	if err = WriteConfig(&config, configPath); err != nil {
		return err
	}
	return nil
}

func branchConflict(repoURL, branchName string, localLayers []string) (bool, error) {

	remoteConfigPath := return_BackingFileConfig_Path(&repoURL)
	//r := rand.New(rand.NewSource(time.Now().UnixNano()))
	//tmpConfig := fmt.Sprintf("%s/%s.%d", os.TempDir(), path.Base(remoteConfigPath), r.Intn(100000))
	tmpConfig := fmt.Sprintf("%s/%s", os.TempDir(), path.Base(remoteConfigPath))
	print_Trace(fmt.Sprintf("Download remote config to local.(%s)", tmpConfig))
	if err := downloadFile(&remoteConfigPath, &tmpConfig); err != nil {
		if err.Error() == "404" {
			return false, nil
		}
		return false, err
	}
	remoteConfig := YamlBackingFileConfig{}
	if err := LoadConfig(&remoteConfig, &tmpConfig); err != nil {
		os.Remove(tmpConfig)
		return true, err
	}
	head := return_BranchHead(&branchName, &remoteConfig.Branch)
	if head == "" {
		return false, nil
	}
	for _, layerUUID := range localLayers {
		if head == layerUUID {
			return false, nil
		}
	}
	os.Remove(tmpConfig)
	return true, nil
}

func downloadFile(url, localPath *string) error {

	respConfig, err := http.Get(*url)
	if err != nil {
		msg := fmt.Errorf("Fetch: %v", err)
		return msg
	}
	if respConfig.StatusCode == 404 {
		return fmt.Errorf("404")
	}
	defer respConfig.Body.Close()
	configDst, err := os.OpenFile(*localPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer configDst.Close()
	configBuff, err := ioutil.ReadAll(respConfig.Body)
	if err != nil {
		return err
	}
	_, err = configDst.Write(configBuff)
	if err != nil {
		return err
	}
	return nil
}
