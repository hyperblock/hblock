package hblock

import (
	"fmt"
	"log"
	"path/filepath"
)

type ShowDetail struct {
	name string
	//	createDate  time.Time
	virtualSize int64
	actualSize  int64
	snapshots   []SnapShot
}

func show_template(image string, logger *log.Logger) (int, error) {

	print_Log("Fetching backing file info...", logger)
	jsonBackingFile, err := return_JsonBackingFile(image)
	if err != nil {
		print_Error(err.Error(), logger)
		return FAIL, err
	}
	print_Log("Done", logger)
	snapshotsList := return_Snapshots(&jsonBackingFile)
	detail := ShowDetail{
		name: filepath.Base(jsonBackingFile.Filename),
		//createDate:  time.Unix(jsonBackingFile.Snapshots[0].DateSec, jsonBackingFile.Snapshots[0].DateNSec),
		virtualSize: jsonBackingFile.VirtualSize,
		actualSize:  jsonBackingFile.ActualSize,
		snapshots:   snapshotsList,
	}
	msg := fmt.Sprintf(
		SHOW_FORMAT, detail.name,
		// detail.createDate.Format("2006-01-02 15:04:05"),
		float64(detail.virtualSize>>20)/1024, detail.virtualSize, float64(detail.actualSize>>10)/1024,
	)
	snapshotMsg := ""
	for _, item := range detail.snapshots {
		info := fmt.Sprintf(`
		Index: %s
		Create Date: %s
		UUID: %s
		Parent-UUID: %s
		Disk Size: %.2fG (%d bytes)
		Commit Message: %s
		`, item.id, item.createDate.Format("2006-01-02 15:04:05"), item.uuid, item.parent_uuid, float64(item.diskSize>>20)/1024, item.diskSize, item.commit_msg)
		snapshotMsg += info
	}
	print_Log(msg+snapshotMsg, logger)
	print_Log(format_Success("Done."), logger)
	//	fmt.Println(snapshotList)
	// t := time.Unix(templateInfo.Snapshots[0].DateSec, templateInfo.Snapshots[0].DateNSec)
	// //	t = t.Add(dateSec)
	// fmt.Println(t.Format("2006-01-02 15:04:05"))

	// fmt.Printf("%.2fG (%d bytes)\n", float64(templateInfo.ActualSize>>20)/1024, templateInfo.ActualSize)
	return OK, nil
}
