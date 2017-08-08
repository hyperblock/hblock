package hblock

import "log"

import "fmt"

func volume_Rebase(obj *RebaseParams, logger *log.Logger) (int, error) {

	print_Trace(fmt.Sprintf("rebase volume '%s' ( -> %s )", obj.volumePath, obj.backingfile))
	print_Log("Rebase volume...", logger)

	volumeInfo, err := return_VolumeInfo(&obj.volumePath)
	if err != nil {
		return FAIL, fmt.Errorf("Load volume info failed. ( %s )", err.Error())
	}

	backingFileConfig0 := return_BackingFileConfig_Path(&volumeInfo.backingFile)
	backingFileConfig1 := return_BackingFileConfig_Path(&obj.backingfile)
	config0, config1 := YamlBackingFileConfig{}, YamlBackingFileConfig{}
	print_Log("Check backingfile's format...", logger)
	if err = LoadConfig(&config0, &backingFileConfig0); err != nil {
		return FAIL, err
	}
	if err = LoadConfig(&config1, &backingFileConfig1); err != nil {
		return FAIL, err
	}
	if config0.Format != config1.Format {
		return FAIL, fmt.Errorf("The volume's format (%s) is not same as the rebase backingfile's format (%s).", config0.Format, config1.Format)
	}
	print_Log(fmt.Sprintf("Search full layerUUID (%s)......", obj.parentLayer), logger)
	obj.parentLayer, err = return_LayerUUID(obj.backingfile, obj.parentLayer, false)
	if err != nil {
		//	print_Log(fmt.Sprintf("\rSearch full layerUUID (%s)......FAIL\n", obj.parentLayer), logger)
		return FAIL, err
	}
	print_Log("LayerUUID: "+obj.parentLayer, logger)
	h, err := CreateHBM(FMT_UNKNOWN, config0.Format)
	if err != nil {
		return FAIL, err
	}
	if err = h.Rebase(obj); err != nil {
		return FAIL, err
	}

	return OK, nil
}
