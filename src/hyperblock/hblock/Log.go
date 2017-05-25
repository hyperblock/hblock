package hblock

import (
	"fmt"
	"log"

	yaml "gopkg.in/yaml.v2"
)

func volume_commit_history(path string, logger *log.Logger) (int, error) {

	print_Log("Fetching volume information...", logger)
	// jsonVolume, err := return_JsonVolume(path)

	// if err != nil {
	// 	print_Error(err.Error(), logger)
	// 	return FAIL, err
	// }
	// print_Log("Done.", logger)
	// volumeInfo := return_VolumeInfo(&jsonVolume)
	volumeInfo, err := return_VolumeInfo(&path)
	if err != nil {
		print_Error(err.Error(), logger)
		return FAIL, err
	}
	print_Log("Done.", logger)
	layer := volumeInfo.layer
	backingFile := volumeInfo.backingFile
	print_Log("Locate backing file path done.", logger)
	print_Log("Fetching backing file information...", logger)
	jsonBackingFile, err := return_JsonBackingFile(backingFile)
	if err != nil {
		print_Error(err.Error(), logger)
		return FAIL, err
	}
	print_Log("Done.", logger)

	p := layer
	print_Log(fmt.Sprintf("Analysing related commits...(start layer: %s)", p), logger)
	related_commit := return_commit_history(&jsonBackingFile, layer)
	print_Log(fmt.Sprintf("Done. Related commits count: %d", len(related_commit)), logger)
	result := ""
	for _, commitInfo := range related_commit {
		tmp := LogInfo{UUID: commitInfo.uuid, CommitMsg: commitInfo.commit_msg, Date: commitInfo.createDate}
		//	result = append(result, tmp)
		commitInfo := YamlCommitMsg{}
		err = yaml.Unmarshal([]byte(tmp.CommitMsg), &commitInfo)
		if err != nil {
			result += fmt.Sprintf("%s\nparse error.\n\n", yellow("commit: "+tmp.UUID))
		} else {
			result += fmt.Sprintf(HB_LOG_FORMAT, yellow("commit: "+tmp.UUID), commitInfo.Name, commitInfo.Email, tmp.Date.Format("2006-01-02 15:04:05"), commitInfo.Message)
		}
	}
	fmt.Println(result)
	print_Log(Format_Success("Done."), logger)
	return OK, nil
}
