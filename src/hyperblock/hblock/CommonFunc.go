package hblock

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func print_Log(msg string, logger *log.Logger) {

	if len(msg) == 0 {
		return
	}
	if msg[0] != 27 {
		msg = format_Info(msg)
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
	msg = format_Error(msg)
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
	msg = format_Error(msg)
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
	msg = format_Error(msg)
	fmt.Println(msg)
	if logger != nil {
		msg = msg[15:]
		logger.Fatalln(msg)
	}
}

func print_Trace(a ...interface{}) {

	// msg := format_Trace(fmt.Sprint(a))
	// fmt.Println(msg)
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

	// ret := os.Getenv("HBLOCK_TEMPLATE_DIR")
	// //	fmt.Println("env", fmt.Sprintf("[%s]", ret))
	// if ret == "" {
	// 	msg := format_Warning("env HBLOCK_TEMPLATE_DIR not set. use default dir '/var/hyperblock'")
	// 	return "/var/hyperblock", fmt.Errorf(msg)
	// }
	// return ret, nil
	path, err := exec.Command("pwd").Output()
	if err != nil {
		return "", err
	}
	ret := string(path[:len(path)-1]) + "/" + DEFALUT_BACKING_FILE_DIR
	print_Trace("template dir: " + ret)
	return ret, nil
}

func confirm_BackingFilePath(imgPath string) (string, error) {

	if (len(imgPath) > 0) && (imgPath[0] == '/') {
		_, err := os.Stat(imgPath)
		if err != nil {
			return "", fmt.Errorf("Invalid template path '%s'", imgPath)
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

func return_Snapshots(jsonBackingFile *JsonBackingFile) []SnapShot {

	list := jsonBackingFile.Snapshots
	ret := []SnapShot{}
	for _, item := range list {
		snapShot := SnapShot{
			id: item.Id, diskSize: item.DiskSize, createDate: time.Unix(item.DateSec, item.DateNSec),
		}
		nameInfo := strings.Split(item.Name, ",")
		snapShot.uuid = nameInfo[0]
		snapShot.parent_uuid = nameInfo[1]
		snapShot.commit_msg = nameInfo[2]
		print_Trace(snapShot)
		ret = append(ret, snapShot)
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

func return_VolumeInfo(jsonVolume *JsonVolume) VolumeInfo {

	volInfo := VolumeInfo{
		fileName: jsonVolume.Filename, actualSize: jsonVolume.ActualSize, virtualSize: jsonVolume.VirutalSize,
	}
	args := strings.Split(jsonVolume.BackingFile, "?")
	volInfo.backingFile = get_StringAfter(args[0], "qcow2://")
	volInfo.layer = get_StringAfter(args[1], "layer=")
	return volInfo
}

// func return_LayerUUID_from_Snapshots(snapshots []SnapShot, layerPrefix string) (string, error) {

// 	cnt := 0
// 	ret := ""
// 	for _, item := range snapshots {
// 		uuid := item.uuid
// 		if strings.HasPrefix(uuid, layerPrefix) {
// 			cnt++
// 			ret = uuid
// 		}
// 	}
// 	if cnt == 0 {
// 		return "", fmt.Errorf(
// 			fmt.Sprintf("Can't get coresponding layer from prefix '%s'", layerPrefix))
// 	}
// 	if cnt > 1 {
// 		return "", fmt.Errorf(
// 			fmt.Sprintf("There are more than one layers have prefix '%s'", layerPrefix))
// 	}
// 	return ret, nil
// }

func return_LayerUUID(backingFilePath string, layerPrefix string, returnLast bool) (string, error) {

	if layerPrefix == "" && !returnLast {
		//msg := fmt.Sprintf("layerPrefix is null and not return last layer.")
		return "", nil
	}
	jsonBackingFile, err := return_JsonBackingFile(backingFilePath)
	if err != nil {
		return "", fmt.Errorf("Invalid layer_uuid or backing file path.")
	}
	snapshots := return_Snapshots(&jsonBackingFile)
	if len(snapshots) == 0 {
		return "", fmt.Errorf("There're no any layer in backing file.")
	}
	if returnLast {
		if layerPrefix == "" {
			return snapshots[len(snapshots)-1].uuid, nil
		} else {
			return "", fmt.Errorf("'returnLast' flag is true but LayerPrefix is not null.")
		}
	}
	cnt := 0
	ret := ""
	for _, item := range snapshots {
		uuid := item.uuid
		if strings.HasPrefix(uuid, layerPrefix) {
			cnt++
			ret = uuid
		}
	}
	if cnt == 0 {
		return "", fmt.Errorf(
			fmt.Sprintf("Can't get coresponding layer from prefix '%s'", layerPrefix))
	}
	if cnt > 1 {
		return "", fmt.Errorf(
			fmt.Sprintf("There are more than one layers have prefix '%s'", layerPrefix))
	}
	return ret, nil
	//	return return_LayerUUID_from_Snapshots(snapshots, layerPrefix)
}

func return_commit_history(jsonBackingFile *JsonBackingFile, p string) []SnapShot {

	commitList := return_Snapshots(jsonBackingFile)
	ret := []SnapShot{}

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

// see complete color rules in document in https://en.wikipedia.org/wiki/ANSI_escape_code#cite_note-ecma48-13
func format_Trace(format string, a ...interface{}) string {
	prefix := yellow(trac)
	return fmt.Sprint(formatLog(prefix), fmt.Sprintf(format, a...))
}

func format_Info(format string, a ...interface{}) string {
	prefix := blue(info)
	return fmt.Sprint(formatLog(prefix), fmt.Sprintf(format, a...))
}

func format_Success(format string, a ...interface{}) string {
	prefix := green(succ)
	return fmt.Sprint(formatLog(prefix), fmt.Sprintf(format, a...))
}

func format_Warning(format string, a ...interface{}) string {
	prefix := magenta(warn)
	return fmt.Sprint(formatLog(prefix), fmt.Sprintf(format, a...))
}

func format_Error(format string, a ...interface{}) string {
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
