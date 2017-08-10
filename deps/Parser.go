package hblock

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"strings"

	"strconv"

	flags "github.com/jessevdk/go-flags"
)

// OptSelector decide the appropriate function to execute
type OptSelector struct {
	logger *log.Logger
}

// Create an OptSelector object, and all message will print to _log
func CreateOptSelector(_log *log.Logger) *OptSelector {

	p := &OptSelector{logger: _log}

	return p
}

// SendCommand : Execute args, call this after Create()
func (p OptSelector) SendCommand(args []string) (int, error) {

	if len(args) == 0 {
		//	print_Error("invalid option.", p.logger)
		return FAIL, fmt.Errorf("invalid option.")
	}
	option := args[0]
	//	args = args[1:]
	switch option {
	case "init":
		return p.init(args)
	case "branch":
		return p.branch(args)
	case "checkout":
		return p.checkout(args)
	case "commit":
		return p.commit(args)
	case "clone":
		return p.clone(args)
	case "pull":
		return p.pull(args)
	case "push":
		return p.push(args)
	case "save":
		return p.save(args)
	case "log":
		return p.log(args)
	case "rebase":
		return p.rebase(args)
	case "reset":
		return p.reset(args)
	case "remote":
		return p.remote(args)
	case "tag":
		return p.tag(args)
	case "config":
		return p.config(args)
	case "sh":
		return p.Sh(args)
	case "before_commit_hooks":
		return p.before_commit_hooks(args)
	case "post_checkout_hooks":
		return p.post_checkout_hooks(args)
	case "launch":
		return p.launch(args)
	case "list":
		return p.list(args)
	case "show":
		return p.show(args)
	case "-h":
		fmt.Println(GLOBAL_HELP)
		return OK, nil
	case "--help":
		fmt.Println(GLOBAL_HELP)
		return OK, nil
	default:
		//msg := "invalid option, and there is no --help :)"
		msg := GLOBAL_HELP
		//	print_Error(msg, p.logger)
		return FAIL, fmt.Errorf(msg)
	}
}

func (p OptSelector) init(args []string) (int, error) {

	var options struct {
		Size   string `long:"size" description:"[required] Disk size(M/G) of template.\n\t\t\t    eg.\n\t\t\t\thblock init template0 --size=500M -f qcow2.\n"`
		Output string `short:"o" description:"[optional] Output volume name.\n"`
		Format string `short:"f" long:"format" description:"'qcow2' of 'lvm'.\n"`
	}
	os.Args = custom_Args(args, "<backingfile name>")
	args, err := flags.ParseArgs(&options, args[1:])
	if err != nil {
		//	print_Error(err.Error(), p.logger)
		return FAIL, nil
	}
	templateName := ""
	if len(args) >= 1 {
		templateName = args[0]
		if len(args) > 1 {
			msg := fmt.Sprintf("Get template name '%s' and ignore excrescent arguments.", templateName)
			print_Log(msg, p.logger)
		} else {
			msg := fmt.Sprintf("Get template name '%s'.", templateName)
			print_Log(msg, p.logger)
		}
	} else {
		msg := "Can't get backingfile name."
		p.Usage(&options)
		return FAIL, fmt.Errorf(msg)
	}

	if options.Size == "" {
		msg := "--size is required."
		p.Usage(&options)
		return FAIL, fmt.Errorf(msg)
	}
	sizeI64 := return_Size(options.Size)
	if sizeI64 < 0 {
		msg := "Invalid --size set"
		//	print_Error(msg, p.logger)
		//flags.ParseArgs(&options, []string{"-h"})
		p.Usage(&options)
		return FAIL, fmt.Errorf(msg)
	}
	//directPathFlg := false
	index := strings.LastIndex(templateName, "/")
	if index != -1 {
		// 	directPathFlg = true
		// } else {
		templateName = templateName[index+1:]
	}
	if options.Format == "lvm" && options.Output == "" {
		p.Usage(&options)
		return FAIL, fmt.Errorf("need specify -o if format is 'lvm'")
	}
	if options.Output == "" {
		options.Output = templateName
	}
	if options.Format == "" {
		p.Usage(&options)
		return FAIL, fmt.Errorf("Need specify format.")
	}
	//	fmt.Println(templateName)
	// if !directPathFlg {
	// 	templateDir, _ := return_TemplateDir()
	// 	templateName = templateDir + "/" + templateName
	// }
	obj := InitParams{name: templateName, size: sizeI64, output: options.Output, checkout: true, format: options.Format}
	msg := fmt.Sprintf("Init template named '%s' and new volume '%s'", templateName, obj.output)
	print_Log(msg, p.logger)
	return create_empty_template(obj, p.logger)
}

func (p OptSelector) branch(args []string) (int, error) {

	var options struct {
		List bool ` long:"list" description:"list branch names.\n"`

		All bool `short:"a" long:"all" description:"list both remote-tracking and local branches.\n"`

		Move string `short:"m" long:"move" description:"<exist_branch> <new_branch> move/rename a branch.\n"`

		BackingFile string `short:"t" long:"backingfile"`

		Volume string `short:"v" long:"volume"`
	}
	//os.Args = custom_Args(args, "highlight last commit branch head. ")
	os.Args = custom_Args(args, "")
	args, err := flags.ParseArgs(&options, args[1:])
	if err != nil {
		return FAIL, nil
	}
	if options.Volume == "" && (options.Move != "" && options.BackingFile == "") {
		msg := "Too few arguments. Need specify <volume name>"
		p.Usage(&options)
		return FAIL, fmt.Errorf(msg)
	}
	if !(options.List || options.All) {
		options.List = true
	}
	if options.BackingFile != "" && options.Volume != "" {
		p.Usage(&options)
		//	fmt.Println(options.Volume, options.BackingFile)
		return FAIL, fmt.Errorf("Can not use both '-v' and '-t'.")
	}

	branchParams := BranchParams{
		list: options.List, show_all: options.All, optTag: BRANCH_OPT_SHOW,
	}
	backingfile := options.BackingFile
	if options.Volume != "" {
		branchParams.volumePath = return_AbsPath(options.Volume)
		if !PathFileExists(branchParams.volumePath) {
			msg := fmt.Sprintf("Volume '%s' not found.", branchParams.volumePath)
			//	print_Error(msg, p.logger)
			//flags.ParseArgs(&options, []string{"-h"})
			p.Usage(&options)
			return FAIL, fmt.Errorf(msg)
		}
		backingfile, err = return_Volume_BackingFile(&branchParams.volumePath)
		if err != nil {
			return FAIL, fmt.Errorf("Can't load backingfile info.")
		}
	}
	if options.Move != "" {
		if options.Volume == "" && options.BackingFile == "" {
			msg := "need use '-v' or '-t'."
			p.Usage(&options)
			return FAIL, fmt.Errorf(msg)
		}
		branchParams.move.src = options.Move
		branchParams.move.dst = args[0]
		branchParams.optTag = BRANCH_OPT_MV
		branchParams.backingfile = return_Backingfile_AbsPath(backingfile)
	}
	return volume_Branch(&branchParams, p.logger)

}

func (p OptSelector) checkout(args []string) (int, error) {

	var options struct {
		Volume string `short:"v" long:"vol" description:"<volume_name> <layer | branch> Specify the volume name which needs to be update(restore).\n"`

		Backingfile string `short:"t" long:"backingfile" description:"<backingfile> <layer | branch> Create a new volume from <backingfile>.\n"`

		Output string `short:"o" long:"output" description:"<output_volume_path>.\n"`

		Branch string `short:"b" long:"base" description:"<branch> Create a new branch of base on the specify volume.\n"`

		Force bool `short:"f" long:"force.\n"`
	}
	os.Args = custom_Args(args, "")
	args, err := flags.ParseArgs(&options, args[1:])
	if err != nil {
		return FAIL, nil
	}
	// if len(args) < 1 {
	// 	msg := "Too few arguments"
	// 	p.Usage(&options)
	// 	return FAIL, fmt.Errorf(msg)
	// }
	if options.Backingfile != "" && options.Volume != "" {
		msg := "Can't use both -v and -t."
		p.Usage(&options)
		return FAIL, fmt.Errorf(msg)
	}
	if options.Branch != "" && options.Backingfile != "" {
		msg := "Can't use both -t and -b."
		p.Usage(&options)
		return FAIL, fmt.Errorf(msg)
	}
	// if options.Volume != "" && options.Branch != "" {
	// 	msg := "Can't use both -v and -b."
	// 	p.Usage(&options)
	// 	return FAIL, fmt.Errorf(msg)
	// }

	if options.Branch == "" && (options.Backingfile != "" && options.Output == "") {
		msg := "use '-o' to set output volume file. "
		p.Usage(&options)
		return FAIL, fmt.Errorf(msg)
	}
	if options.Volume != "" {
		options.Volume = return_AbsPath(options.Volume)
		_, err := os.Stat(options.Volume)
		if err != nil {
			msg := fmt.Sprintf("Can't locate volume_name '%s'", options.Volume)
			p.Usage(&options)
			return FAIL, fmt.Errorf(msg)
		}
		if options.Output == "" {
			if !options.Force {
				msg := "Need use '-o' to set output volume or '-f' to reset current volume."
				p.Usage(&options)
				return FAIL, fmt.Errorf(msg)
			}
			options.Output = options.Volume
		}
	}
	// if options.Branch != "" {
	// 	options.Volume = return_AbsPath(args[0])
	// }
	if options.Volume == "" && options.Backingfile == "" {
		msg := "Need specify <volume> or <backingfile>."
		p.Usage(&options)
		return FAIL, fmt.Errorf(msg)
	}
	checkoutObj := CheckoutParams{
		volume:   options.Volume,
		output:   options.Output,
		template: options.Backingfile,
	}
	if options.Branch != "" {
		checkoutObj.branch = options.Branch
	} else {
		checkoutObj.layer = args[0]
	}

	return volume_checkout(&checkoutObj, p.logger)
}

func (p OptSelector) commit(args []string) (int, error) {

	var options struct {
		CommitMsg string `short:"m" description:"commit message"`
		UUID      string `long:"uuid" description:"set uuid by manual instead of auto-generate."`
	}
	os.Args = custom_Args(args, "<volume name>")
	args, err := flags.ParseArgs(&options, args[1:])
	if err != nil {
		return FAIL, nil
	}
	if len(args) < 1 {
		msg := "No volume name specified."
		p.Usage(&options)
		return FAIL, fmt.Errorf(msg)
	}
	if options.CommitMsg == "" {
		msg := "Empty commit message. Use -m to set commit message."
		p.Usage(&options)
		return FAIL, fmt.Errorf(msg)
	}
	commitObj := CommitParams{
		commitMsg: options.CommitMsg, volumeName: args[0], layerUUID: options.UUID,
	}
	commitObj.genUUID = commitObj.layerUUID == ""
	return volume_commit(commitObj, p.logger)

}

func (p OptSelector) clone(args []string) (int, error) {

	var options struct {
		//			Volume string `short:"v" long:"vol" description:"<volume_name>\tSpecify the volume name which needs to be update(restore).\n"`
		Layer string `short:"l" long:"layer" description:"Checkout <layer> instead of the HEAD\n"`

		HardLink bool `long:"hardlink" description:"use local hardlinks.\n"`

		Branch string `short:"b" long:"branch" description:"Clone the specified <branch> instead of default ('master').\n"`

		CheckoutFlag bool `short:"n" long:"no-checkout" description:"No checkout of HEAD is performed after clone is complete.\n"`
	}
	os.Args = custom_Args(args, "<repo path>")
	args, err := flags.ParseArgs(&options, args[1:])

	if err != nil {
		return FAIL, nil
	}
	if len(args) != 1 {
		msg := "Invalid arguments. Use '-h' for help."
		p.Usage(&options)
		return FAIL, fmt.Errorf(msg)
	}
	if options.HardLink && options.Branch != "" {
		msg := "Can't use both --hardlink and -b."
		//	print_Error(msg, p.logger)
		p.Usage(&options)
		return FAIL, fmt.Errorf(msg)
	}
	cloneObj := CloneParams{
		repoPath:    args[0],
		configPath:  return_BackingFileConfig_Path(&args[0]), // args[0] + ".yaml",
		checkoutFlg: !options.CheckoutFlag,
		hardLink:    options.HardLink,
		branch:      options.Branch,
		layerUUID:   options.Layer,
	}
	return clone_Repo(&cloneObj, p.logger)
}

func (p OptSelector) pull(args []string) (int, error) {

	//	p.logger.Println("pull", args)
	//return FAIL, fmt.Errorf("Option unfinished.")
	var options struct {
		Volume string `short:"v" long:"volume"`
	}
	os.Args = custom_Args(args, "<remote> <branch>")
	args, err := flags.ParseArgs(&options, args[1:])
	if err != nil {
		return FAIL, nil
	}
	if len(args) < 2 || options.Volume == "" {
		p.Usage(&options)
		return FAIL, fmt.Errorf("Too few arguments")
	}

	pullObj := PullParams{
		branch: args[1],
		volume: return_AbsPath(options.Volume),
		remote: args[0],
		all:    false,
	}
	return volume_PullBranch(&pullObj, p.logger)
	//return 0, nil
}

func (p OptSelector) push(args []string) (int, error) {

	var options struct {
		Volume string `short:"v" long:"volume" description:"<volume>"`
	}
	os.Args = custom_Args(args, "<repository> <refspec>")
	args, err := flags.ParseArgs(&options, args[1:])
	if err != nil {
		return FAIL, nil
	}
	if len(args) < 2 {
		msg := "Too few arguments."
		//print_Error(msg, p.logger)
		p.Usage(&options)
		return FAIL, fmt.Errorf(msg)
	}
	if options.Volume == "" {
		msg := "Need specify <volume>."
		//print_Error(msg, p.logger)
		p.Usage(&options)
		return FAIL, fmt.Errorf(msg)
	}
	pushObj := PushParams{
		remote: args[0],
		branch: args[1],
		volume: return_AbsPath(options.Volume),
	}
	return push_volume(pushObj, p.logger)
}

func (p OptSelector) save(args []string) (int, error) {

	p.logger.Println("save", args)
	return FAIL, fmt.Errorf("Option unfinished.")
	//return 0, nil
}

func (p OptSelector) log(args []string) (int, error) {

	// p.logger.Println("log", args)
	//	return FAIL, fmt.Errorf("Option unfinished.")
	var options struct {
		// Volume string `short:"v" long:"vol" description:"<volume_name>\tSpecify the volume name which needs to be update(restore).\n"`

		// Layer string `short:"l" long:"layer" description:"<layer>\tSpecify the <layer> that this volume will restore. If the <layer> does\\'not exist, it will create a new layer from current volume."`

		// Output string `short:"o" long:"output" description:"[required if use \\'-t\\'] <output_volume_path>.\n"`

		// Template string `short:"t" long:"template" description:"<template_name>\t Create a new volume from template.\n"`

		// Force bool `short:"f" long:"force.\n"`
		//Last int `short:"l" long:"last" description:"<number> Show the last <number> of commit logs.`
	}
	os.Args = custom_Args(args, "<volume name>")
	if len(args) <= 1 {
		msg := "Too few arguments."
		//print_Error(msg, p.logger)
		//	fmt.Println(OPT_LOG_USAGE)
		return FAIL, fmt.Errorf(msg)
	}
	args, err := flags.ParseArgs(&options, args[1:])
	if err != nil {
		//print_Error(err.Error(), p.logger)
		//	fmt.Println(OPT_LOG_USAGE)
		return FAIL, nil
	}
	//fmt.Println(args[0])
	//volume, err := confirm_BackingFilePath(args[0])
	volume := return_AbsPath(args[0])

	if volume == "" || !PathFileExists(volume) {
		//print_Error(err.Error(), p.logger)
		//fmt.Println(OPT_LOG_USAGE)
		return FAIL, fmt.Errorf("Volume '%s' can not found.", volume)
		//print_Error(msg, p.logger)
		//return FAIL, fmt.Errorf(msg)
	}
	return volume_commit_history(volume, p.logger)
}

func (p OptSelector) rebase(args []string) (int, error) {

	var options struct {
		Backingfile string `short:"b" long:"backingfile" description:"<backingfile>"`
		Layer       string `short:"l" long:"layer" description:"<layer>"`
	}
	os.Args = custom_Args(args, "<volume_name>")
	args, err := flags.ParseArgs(&options, args[1:])
	if err != nil {
		return FAIL, nil
	}
	if options.Backingfile == "" || options.Layer == "" || len(args) == 0 {
		p.Usage(&options)
		return FAIL, fmt.Errorf("Too few arguments.")
	}
	obj := RebaseParams{
		backingfile: return_AbsPath(options.Backingfile),
		parentLayer: options.Layer,
		volumePath:  return_AbsPath(args[0]),
	}
	return volume_Rebase(&obj, p.logger)
}

func (p OptSelector) reset(args []string) (int, error) {

	//	p.logger.Println("reset", args)
	//	return FAIL, fmt.Errorf("Option unfinished.")
	var options struct {
		// Volume string `short:"v" long:"vol" description:"<volume_name>\tSpecify the volume name which needs to be update(restore).\n"`

		// Layer string `short:"l" long:"layer" description:"<layer>\tSpecify the <layer> that this volume will restore. If the <layer> does\\'not exist, it will create a new layer from current volume."`

		// Output string `short:"o" long:"output" description:"[required if use \\'-t\\'] <output_volume_path>.\n"`

		// Template string `short:"t" long:"template" description:"<template_name>\t Create a new volume from template.\n"`

		// Force bool `short:"f" long:"force.\n"`
		//Last int `short:"l" long:"last" description:"<number> Show the last <number> of commit logs.`
	}
	os.Args = custom_Args(args, "")
	if len(args) <= 2 {
		msg := "Too few arguments."
		//print_Error(msg, p.logger)
		fmt.Println(OPT_RESET_USAGE)
		return FAIL, fmt.Errorf(msg)
	}
	args, err := flags.ParseArgs(&options, args[1:])
	if err != nil {
		//print_Error(err.Error(), p.logger)
		fmt.Println(OPT_RESET_USAGE)
		return FAIL, err
	}
	resetObj := ResetParams{time: -1, volume: return_AbsPath(args[0])}
	if !PathFileExists(resetObj.volume) {
		msg := "volume can not find."
		//print_Error(msg, p.logger)
		//fmt.Printf(OPT_RESET_USAGE)
		return FAIL, fmt.Errorf(msg)
	}
	suffix := get_StringAfter(args[1], "HEAD")
	if len(suffix) == 0 {
		resetObj.time = 0
	} else if suffix[0] == '^' {
		count := 0
		for _, ch := range suffix {
			if ch != '^' {
				msg := "invalid options."
				//print_Error(msg, p.logger)
				fmt.Printf(OPT_RESET_USAGE)
				return FAIL, fmt.Errorf(msg)
			}
			count++
		}
		resetObj.time = count
	} else if suffix[0] == '~' {
		num, err := strconv.Atoi(suffix[1:])
		if err != nil {
			msg := "invalid options."
			//	print_Error(msg, p.logger)
			fmt.Printf(OPT_RESET_USAGE)
			return FAIL, fmt.Errorf(msg)
		}
		resetObj.time = num
	} else {
		resetObj.uuid = suffix
	}
	//	fmt.Println(resetObj)
	print_Trace(resetObj)
	return reset_volume(&resetObj, p.logger)

}

func (p OptSelector) remote(args []string) (int, error) {

	var options struct {
		Verbose bool   `short:"a" long:"verbose" description:"show remotes verbose"`
		Volume  string `short:"v" long:"volume" description:"set <volume> whose repo remotes need to be edited."`
		Add     bool   `long:"add" description:"<name> <url>\tAdd a new remote-host to local remote-host list."`
		Remove  string `short:"d" long:"remove" description:"<name>\tDelete a host from local remote-host list."`
		Rename  bool   `long:"rename" description:"<old_name> <new_name>\t Rename an exsiting host name."`
		SetUrl  bool   `long:"set-url" descripion:"<name> <url>\tChange an exists remote-host's url."`
	}
	os.Args = custom_Args(args, "")
	args, err := flags.ParseArgs(&options, args[1:])
	if err != nil {
		//flags.ParseArgs(&options, []string{"-h"})
		return FAIL, nil
	}

	if options.Volume == "" {
		p.Usage(&options)
		return FAIL, fmt.Errorf("Need use -v to set volume")
	}
	volume := return_AbsPath(options.Volume)
	if !PathFileExists(volume) {
		p.Usage(&options)
		return FAIL, fmt.Errorf("The specified volume '%s' can not be found.", volume)
	}
	backingfile, err := return_Volume_BackingFile(&volume)
	if err != nil || VerifyBackingFile(backingfile) != OK {
		p.Usage(&options)
		return FAIL, fmt.Errorf("Can not verify backing file of '%s'", volume)
	}
	remoteObj := RemoteParams{
		verbose:     options.Verbose,
		backingFile: backingfile,
	}
	if !options.Verbose {
		if (options.Add || options.Rename || options.SetUrl) && len(args) < 2 {
			//flags.ParseArgs(&options, []string{"-h"})
			p.Usage(&options)
			return FAIL, fmt.Errorf("Too few arguments..")
		}
		if options.Add {
			remoteObj.add.name = args[0]
			remoteObj.add.url = args[1]
		} else if options.SetUrl {
			remoteObj.setUrl.name = args[0]
			remoteObj.setUrl.url = args[1]
		} else if options.Remove != "" {
			remoteObj.remove = options.Remove
		} else if options.Rename {
			remoteObj.rename.oldName = args[0]
			remoteObj.rename.newName = args[1]
		}
	}
	return Remote(remoteObj, p.logger)
}

func (p OptSelector) tag(args []string) (int, error) {

	p.logger.Println("tag", args)
	return FAIL, fmt.Errorf("Option unfinished.")
	//return 0, nil
}

func (p OptSelector) config(args []string) (int, error) {

	var options struct {
		Global string `long:"global" description:"[user.name|user.email] set global configuration."`
		Get    string `long:"get" description:"<name>\tGet value : <name>"`
	}
	os.Args = custom_Args(args, "")
	args, err := flags.ParseArgs(&options, args[1:])
	if err != nil {
		return FAIL, nil
	}
	if options.Global != "" {
		if len(args) < 1 {
			msg := "Too few arguments."
			//print_Error(msg, p.logger)
			return FAIL, fmt.Errorf(msg)
		}
		configObj := GlobalConfig{}
		configPath := return_hb_ConfigPath()
		err := LoadConfig(&configObj, &configPath)
		//configObj := stConfigObj.(GlobalConfig)
		if err != nil {
			p.Usage(&options)
			return FAIL, fmt.Errorf("Load configuration failed. Please check file '~/.hb/config.yaml' (%s)", err.Error())
			//print_Error(msg, p.logger)

		}
		if options.Global == "user.name" {
			configObj.UserName = args[0]
			print_Log(fmt.Sprintf("Set user.name as '%s'", configObj.UserName), p.logger)
			err = WriteConfig(&configObj, &configPath)
			if err != nil {
				msg := fmt.Sprintf("Write config failed. (%s)", err.Error())
				return FAIL, fmt.Errorf(msg)
			}
		} else if options.Global == "user.email" {
			configObj.UserEmail = args[0]
			print_Log(fmt.Sprintf("Set user.email as '%s'", configObj.UserEmail), p.logger)
			err = WriteConfig(&configObj, &configPath)
			if err != nil {
				msg := fmt.Sprintf("Write config failed. (%s)", err.Error())
				return FAIL, fmt.Errorf(msg)
			}
		} else {
			msg := "unknow option."
			//print_Error(msg, p.logger)
			//flags.ParseArgs(&options, []string{"-h"})
			p.Usage(&options)
			return FAIL, fmt.Errorf(msg)
		}
		msg := Format_Success("Done.")
		print_Log(msg, p.logger)
		return OK, nil
	} else if options.Get != "" {
		configObj := GlobalConfig{}
		configPath := return_hb_ConfigPath()
		err := LoadConfig(&configObj, &configPath)
		if err != nil {
			p.Usage(&options)
			err := fmt.Errorf("Load configuration failed. Please check file '~/.hb/config.yaml' (%s)", err.Error())
			return FAIL, err
		}
		interface_value, err := return_ConfigValue(&configObj, options.Get)
		if err != nil {
			p.Usage(&options)
			return FAIL, fmt.Errorf("Load value of '%s' failed. (%s)", options.Get, err.Error())
			//		print_Error(msg, p.logger)

		}
		value := interface_value.(GlobalConfig)

		msg := Format_Success(fmt.Sprintf("%s: %v", options.Get, value))
		print_Log(msg, p.logger)
		return OK, nil
	} else {
		msg := "Invalid option."
		//	print_Error(msg, p.logger)
		//flags.ParseArgs(&options, []string{"-h"})
		p.Usage(&options)
		return FAIL, fmt.Errorf(msg)
	}
	return FAIL, nil
	//p.logger.Println("config", args)

}

func (p OptSelector) Sh(args []string) (int, error) {

	print_Log(strings.Join(args, " "), p.logger)
	var cmd *exec.Cmd
	if len(args) > 1 {
		cmd = exec.Command("guestfish", args[1:]...)
	} else {
		cmd = exec.Command("guestfish")
	}
	//	cmd := exec.Command("ls", "-l")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		p.logger.Println(err)

	}
	return OK, nil
}

func (p OptSelector) before_commit_hooks(args []string) (int, error) {

	p.logger.Println("before_commit_hooks", args)
	return FAIL, fmt.Errorf("Option unfinished.")

}

func (p OptSelector) post_checkout_hooks(args []string) (int, error) {

	p.logger.Println("post_checkout_hooks", args)
	return FAIL, fmt.Errorf("Option unfinished.")

}

func (p OptSelector) launch(args []string) (int, error) {

	p.logger.Println("launch", args)
	return FAIL, fmt.Errorf("Option unfinished.")
}

func (p OptSelector) list(args []string) (int, error) {

	if len(args) >= 2 && (args[1] == "-h" || args[1] == "--help") {
		fmt.Println(OPT_LIST_USAGE)
		return OK, nil
	}
	targetdir, _ := return_CurrentDir()
	targetdir += "/" + DEFALUT_BACKING_FILE_DIR
	if len(args) > 1 {
		targetdir = return_AbsPath(args[1])
	}
	dir, err := ioutil.ReadDir(targetdir)
	if err != nil {
		return FAIL, err
	}
	for _, file := range dir {
		if file.IsDir() {
			continue
		}
		fmt.Println(file.Name())
	}
	return 0, nil
}

func (p OptSelector) show(args []string) (int, error) {

	p.logger.Println("show", args)
	//	return FAIL, fmt.Errorf("Option unfinished.")
	if len(args) < 2 {
		msg := "Too few arguments.\n"
		fmt.Println(OPT_SHOW_USAGE)
		return FAIL, fmt.Errorf(msg)
	}
	if args[1] == "-h" || args[1] == "--help" {
		fmt.Println(OPT_SHOW_USAGE)
		return OK, nil
	}

	image, err := confirm_BackingFilePath(args[1])
	if image == "" {
		//	print_Error(err.Error(), p.logger)
		usage := OPT_SHOW_USAGE
		fmt.Println(usage)
		return FAIL, err
	}

	return show_template(image, p.logger)
}

func (p OptSelector) Usage(options interface{}) {

	flags.ParseArgs(options, []string{"-h"})
}
