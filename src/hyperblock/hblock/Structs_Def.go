package hblock

import (
	"time"
)

type InitParams struct {
	name     string
	size     int64
	format   string
	output   string
	checkout bool
}

type CheckoutParams struct {
	layer    string
	volume   string
	output   string
	template string
	branch   string
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
	} `json:"snapshots"`
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
	configPath  string
	checkoutFlg bool
	layerUUID   string
	branch      string
	hardLink    bool
	protocol    int
}

type RemoteParams struct {
	verbose bool
	add     struct {
		name string
		url  string
	}
	remove string
	rename struct {
		oldName string
		newName string
	}
	setUrl struct {
		name string
		url  string
	}
	backingFile string
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

type BranchParams struct {
	show_all   bool
	list       bool
	volumePath string
}

type YamlBranch struct {
	Name  string `yaml:"name"`
	Local int    `yaml:"local"`
	Head  string `yaml:"head"`
}

type YamlRemote struct {
	Name string `yaml:"name"`
	Url  string `yaml:"url"`
}

type YamlBackingFileConfig struct {
	Name        string       `yaml:"name"`
	Format      string       `yaml:"format"`
	VirtualSize int64        `yaml:"virtual size"`
	DefaultHead string       `yaml:"head"`
	Remote      []YamlRemote `yaml:"remote"`
	Branch      []YamlBranch `yaml:"branch"`
}

type YamlVolumeConfig struct {
	Branch    string `yaml:"branch.name"`
	NewBranch bool   `yaml:"branch.create"`
}

type PushParams struct {
	remote string
	volume string
	branch string
	url    string
}

type PullParams struct {
	pullList       []string
	branch         string
	protocol       int
	all            bool
	remoteRepoPath string
	configPath     string
	localRepoPath  string
}

type DumpParams struct {
	backngFile string
	layerUUID  string
	output     string
}

type RebaseParams struct {
	backingfile string
	volumePath  string
	parentLayer string
}
