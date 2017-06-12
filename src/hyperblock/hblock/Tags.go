package hblock

// type HyperblockCLI struct {
// 	logger     *log.Logger
// 	currentDir string
// 	err        error
// }
const SHOW_TRACE = 0

const WAIT_CHANGE_LAYER = 1

const (
	OK                         = 0
	FAIL                       = -1
	BACKINGFILE_CONFIG_NO_FIND = 0x00001
	BACKINGFILE_NO_FIND        = 0x00002
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
	BRANCH_OPT_SHOW = 0
	BRANCH_OPT_MV   = 1
)

const (
	REPO_PATH_LOCAL = 0
	REPO_PATH_HTTP  = 1
	REPO_PATH_SSH   = 2
)

const (
	CONFIG_G_USER_NAME  = "user.name"
	CONFIG_G_USER_EMAIL = "user.email"
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

Layers info:
	`
const LAYER_INFO_FORMAT = `
	Index: %s
	Create Date: %s
	UUID: %s
	Parent-UUID: %s
	Disk Size: %.2fG (%d bytes)
	Commit Message:
%s
`

const OPT_LIST_USAGE = `Usage:
	hb list
	hb list <dir_path>
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
const GLOBAL_HELP = `hb is hyperblock CLI and SDK library.

Usage:
    hb <command> [options]

	Note. use 'hb <command> -h' to see detail.

    ======== support commands =======

	init            create empty backingfile
	config          get and set global options
	clone           clone a repository from remote or local path
	remote          manage set of tracked repositories
	rebase          reapply volume's backingfile and parent-layer 
	branch          list,create or delete branches
	checkout        switch branches or restore volume
	commit          record volume's changes
	reset           reset current HEAD to the specified state
	pull            fetch from and integrate with another repository of a local branch
	push            update remote repository
	list            list backingfiles in current workspace
	show            show backingfile's detail
	log             show commit logs
`
