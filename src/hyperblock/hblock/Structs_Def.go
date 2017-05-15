package hblock

import (
	"time"
)

type InitParams struct {
	name   string
	size   int64
	output string
}

type CheckoutParams struct {
	layer    string
	volume   string
	output   string
	template string
}

type JsonBackingFile struct {
	Filename    string
	ActualSize  int64 `json:"actual-size"`
	VirtualSize int64 `json:"virtual-size"`
	Layers      []struct {
		Name     string
		Id       string
		DiskSize int64 `json:"disk-size"`
		DateSec  int64 `json:"date-sec"`
		DateNSec int64 `json:"date-nsec"`
	}
}

type Layer struct {
	id          string
	uuid        string
	diskSize    int64
	createDate  time.Time
	parent_uuid string
	commit_msg  string
}

type JsonVolume struct {
	Filename    string
	VirutalSize int64  `json:"virtual-size"`
	ActualSize  int64  `json:"actual-size"`
	BackingFile string `json:"full-backing-filename"`
}

type VolumeInfo struct {
	fileName    string
	virtualSize int64
	actualSize  int64
	backingFile string
	layer       string
}

type CommitParams struct {
	commitMsg  string
	volumeName string
	layerUUID  string
	genUUID    bool
}

type LogInfo struct {
	UUID      string
	CommitMsg string
	Date      time.Time
}

type ResetParams struct {
	time   int
	uuid   string
	volume string
}

type CloneParams struct {
	repoPath    string
	checkoutFlg bool
	layerUUID   string
	protocol    int
}

type GlobalConfig struct {
	UserName  string `yaml:"user.name"`
	UserEmail string `yaml:"user.email"`
}

type UserInfo struct {
	name  string
	email string
}

type YamlCommitMsg struct {
	Message string `yaml:"msg"`
	Name    string `yaml:"name"`
	Email   string `yaml:"email"`
	Tag     string `yaml:"tag"`
}
