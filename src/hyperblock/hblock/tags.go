package hblock

const (
	OK   = 0
	FAIL = -1
)

const (
	color_red = uint8(iota + 91)
	color_green
	color_yellow
	color_blue
	color_magenta //洋红

	info = "[INFO]"
	trac = "[TRAC]"
	erro = "[ERRO]"
	warn = "[WARN]"
	succ = "[SUCC]"
)

const (
	REPO_PATH_LOCAL = 0
	REPO_PATH_HTTP  = 1
	REPO_PATH_SSH   = 2
)

const USER_INFO_EMPTY = `user.name or user.email is emtpy.
Use 'hb config --global user.name <user name>' to set user.name
    'hb config --global user.email <email address>' to set user.email
`

const DEFALUT_BACKING_FILE_DIR = ".hb"

const OPT_SHOW_USAGE = `Usage:
	hb show <backing file> 	show backing file details.`

const OPT_LOG_USAGE = `Usage:
	hb log	<backing file>		show commit log of specify backing file.`

const HB_LOG_FORMAT = `
%s
Author: %s <%s>
Date: %s

    %s
	
	`

const SHOW_FORMAT = `-------------Backing File Details----------------

Name: %s
Disk Size: %.2fG (%d bytes)
Actual Size: %.2fMB 

Snapshots info:
	`
const SNAPSHOT_FORMAT = `
	Index: %s
	Create Date: %s
	UUID: %s
	Parent-UUID: %s
	Disk Size: %.2fG (%d bytes)
	Commit Message:
%s
`

const OPT_RESET_USAGE = `Usage:
	hb reset <volume> [<commit_uuid>] | [HEAD point]	reset <volume> and discard changes.
	eg.
		hb reset volume0 3f2ed7		reset 'volume' to specified commit 3f2ed7
		hb reset volume0 HEAD^^		reset 'volume' to the last 2 commits
		hb reset volume0 HEAD~5		reset 'volume' to the last 5 commits
	`

// const DEFALUT_TEMPLATE_DIR = "backing_file"
// const DEFAULT_VL_LOG_DIR = "volume_logs"
