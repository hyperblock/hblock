package hblock

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"
)

type ShowDetail struct {
	name string
	//	createDate  time.Time
	virtualSize int64
	actualSize  int64
	layers      []Layer
}

func format_commitMsg(msg string) string {

	contents := strings.Split(msg, "\n")
	for i := 0; i < len(contents); i++ {
		contents[i] = "\t    " + contents[i]
	}
	return strings.Join(contents, "\n")
}

func show_template(image string, logger *log.Logger) (int, error) {

	print_Log("Fetching backing file info...", logger)
	jsonBackingFile, err := return_JsonBackingFile(image)
	if err != nil {
		print_Error(err.Error(), logger)
		return FAIL, err
	}
	print_Log("Done", logger)
	//	fmt.Println(jsonBackingFile)
	layerList := return_LayerList(&jsonBackingFile)
	detail := ShowDetail{
		name: filepath.Base(jsonBackingFile.Filename),
		//createDate:  time.Unix(jsonBackingFile.Snapshots[0].DateSec, jsonBackingFile.Snapshots[0].DateNSec),
		virtualSize: jsonBackingFile.VirtualSize,
		actualSize:  jsonBackingFile.ActualSize,
		layers:      layerList,
	}
	msg := fmt.Sprintf(
		SHOW_FORMAT, detail.name,
		// detail.createDate.Format("2006-01-02 15:04:05"),
		float64(detail.virtualSize>>20)/1024, detail.virtualSize, float64(detail.actualSize>>10)/1024,
	)
	layerInfo := ""
	//	fmt.Println(len(detail.layers))
	for _, item := range detail.layers {
		info := fmt.Sprintf(
			LAYER_INFO_FORMAT,
			item.id,
			item.createDate.Format("2006-01-02 15:04:05"),
			item.uuid,
			item.parent_uuid,
			float64(item.diskSize>>20)/1024,
			item.diskSize,
			format_commitMsg(item.commit_msg))
		layerInfo += info
	}
	print_Log(msg+layerInfo, logger)
	print_Log(Format_Success("Done."), logger)
	return OK, nil
}
