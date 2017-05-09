package hblock

import (
	"fmt"
	"log"
)

func reset_volume(resetObj *ResetParams, logger *log.Logger) (int, error) {

	//	specifyCommit := false
	print_Log("Confirm volume information...", logger)
	jsonVolume, err := return_JsonVolume(resetObj.volume)
	if err != nil {
		msg := "Can not get volume info."
		print_Error(msg, logger)
		return FAIL, fmt.Errorf(msg)
	}
	volumeInfo := return_VolumeInfo(&jsonVolume)
	print_Log("done", logger)
	if resetObj.time != -1 {
		print_Log(fmt.Sprintf("Reset volume to the last %d commit", resetObj.time+1), logger)
	} else {
		print_Log(fmt.Sprint("Reset volume to commit %s", resetObj.uuid), logger)
		fullUUID, err := return_LayerUUID(volumeInfo.backingFile, resetObj.uuid, false)
		if err != nil {
			print_Error(err.Error(), logger)
			return FAIL, err
		}
		checkoutObj := CheckoutParams{volume: resetObj.volume, layer: fullUUID}
		return volume_checkout(checkoutObj, logger)
		//specifyCommit = true
	}

	jsonBackingFile, err := return_JsonBackingFile(volumeInfo.backingFile)
	if err != nil {
		msg := "Can not get backing file info."
		print_Error(msg, logger)
		return FAIL, fmt.Errorf(msg)
	}

	related_commit := return_commit_history(&jsonBackingFile, volumeInfo.layer)
	if len(related_commit) <= resetObj.time {
		msg := fmt.Sprintf("The reset commit version is invalid.( related commits: %d < %d )", len(related_commit), resetObj.time+1)
		print_Error(msg, logger)
		return FAIL, fmt.Errorf(msg)
	}
	checkoutObj := CheckoutParams{volume: resetObj.volume, layer: related_commit[resetObj.time].uuid}
	return volume_checkout(checkoutObj, logger)
}
