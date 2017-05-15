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
func Create(_log *log.Logger) *OptSelector {

	p := &OptSelector{logger: _log}

	return p
}

// SendCommand : Execute args, call this after Create()
func (p OptSelector) SendCommand(args []string) (int, error) {

	//args = strings.Split("config --global user.name yyf", " ")
	//args = []string{"init", "--name", "hehe"}
	if len(args) == 0 {
		print_Error("invalid option.", p.logger)
		return FAIL, fmt.Errorf("invalid option.")
	}
	option := args[0]
	//	args = args[1:]
	switch option {
	case "init":
		return p.init(args)
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
	default:
		msg := "invalid option, and there is no --help :)"
		print_Error(msg, p.logger)
		return FAIL, fmt.Errorf(msg, p.logger)
	}
}

func (p OptSelector) init(args []string) (int, error) {

	var options struct {
		Size   string `long:"size" description:"[required] Disk size(M/G) of template.\n\t\t\t    eg.\n\t\t\t\thblock init template0 --size=500M.\n"`
		Output string `short:"o" description:"[optional] Output volume name.\n"`
	}
	os.Args = custom_Args(args, "<template name>")
	args, err := flags.ParseArgs(&options, args[1:])
	if err != nil {
		//	print_Error(err.Error(), p.logger)
		return FAIL, err
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
		msg := "Can't get template name."
		print_Error(msg, p.logger)
		flags.ParseArgs(&options, []string{"-h"})
		return FAIL, fmt.Errorf(msg)
	}

	if options.Size == "" {
		msg := "--size is required."
		print_Error(msg, p.logger)
		flags.ParseArgs(&options, []string{"-h"})
		return FAIL, fmt.Errorf(msg)
	}
	sizeI64 := return_Size(options.Size)
	if sizeI64 < 0 {
		msg := "Invalid --size set"
		print_Error(msg, p.logger)
		flags.ParseArgs(&options, []string{"-h"})
		return FAIL, fmt.Errorf(msg)
	}
	directPathFlg := false
	index := strings.LastIndex(templateName, "/")
	if index != -1 {
		directPathFlg = true
	}
	if options.Output == "" {
		if directPathFlg {
			options.Output = templateName
		} else {
			options.Output = templateName[index+1:]
		}
	}
	//	fmt.Println(templateName)
	if !directPathFlg {
		templateDir, _ := return_TemplateDir()
		templateName = templateDir + "/" + templateName
	}
	obj := InitParams{name: templateName, size: sizeI64, output: options.Output}
	msg := fmt.Sprintf("Init template named '%s' and new volume '%s'", templateName, obj.output)
	print_Log(msg, p.logger)
	return create_empty_template(obj, p.logger)
}

func (p OptSelector) checkout(args []string) (int, error) {

	var options struct {
		Volume string `short:"v" long:"vol" description:"<volume_name>\tSpecify the volume name which needs to be update(restore).\n"`

		Layer string `short:"l" long:"layer" description:"<layer>\tSpecify the <layer> that this volume will restore. If the <layer> does\\'not exist, it will create a new layer from current volume."`

		Output string `short:"o" long:"output" description:"[required if use \\'-t\\'] <output_volume_path>.\n"`

		Template string `short:"t" long:"template" description:"<template_name>\t Create a new volume from template.\n"`

		Force bool `short:"f" long:"force.\n"`
	}
	os.Args = custom_Args(args, "")
	args, err := flags.ParseArgs(&options, args[1:])

	if err != nil {
		return FAIL, err
	}
	if len(args) > 0 {
		msg := "Invalid options"
		print_Error(msg, p.logger)
		flags.ParseArgs(&options, []string{"-h"})
		return FAIL, fmt.Errorf(msg)
	}
	if options.Template != "" && options.Volume != "" {
		msg := "Can't use both -v and -t."
		print_Error(msg, p.logger)
		flags.ParseArgs(&options, []string{"-h"})
		return FAIL, fmt.Errorf(msg)
	}
	if options.Output != "" && options.Volume != "" {
		msg := "Can't use both -v and -o."
		print_Error(msg, p.logger)
		flags.ParseArgs(&options, []string{"-h"})
		return FAIL, fmt.Errorf(msg)
	}
	if options.Volume != "" {
		options.Volume = return_AbsPath(options.Volume)
		_, err := os.Stat(options.Volume)
		if err != nil {
			msg := fmt.Sprintf("Can't locate volume_name '%s'", options.Volume)
			print_Error(msg, p.logger)
			flags.ParseArgs(&options, []string{"-h"})
			return FAIL, fmt.Errorf(msg)
		}
		if !options.Force {
			msg := "Need commit the exist volume or use -f"
			print_Error(msg, p.logger)
			flags.ParseArgs(&options, []string{"-h"})
			return FAIL, fmt.Errorf(msg)
		}
	}
	if options.Template != "" && options.Output == "" {
		msg := "use -o <output_volume_path> to set output volume file. "
		print_Error(msg, p.logger)
		flags.ParseArgs(&options, []string{"-h"})
		return FAIL, fmt.Errorf(msg)
	}
	checkoutObj := CheckoutParams{
		volume:   options.Volume,
		layer:    options.Layer,
		output:   options.Output,
		template: options.Template,
	}

	return volume_checkout(checkoutObj, p.logger)
}

func (p OptSelector) commit(args []string) (int, error) {

	//p.logger.Println("commit", args)
	//fmt.Println(args)
	//return FAIL, fmt.Errorf("Option unfinished.")
	var options struct {
		CommitMsg string `short:"m" description:"commit message"`
		Uuid      string `long:"uuid" description:"set uuid by manual instead of auto-generate.`
	}
	os.Args = custom_Args(args, "<volume name>")
	//	os.Args += " <volume name>"
	args, err := flags.ParseArgs(&options, args[1:])
	if err != nil {
		//	flags.ParseArgs(&options, []string{"-h"})
		return FAIL, nil
	}
	//	fmt.Println(args)
	if len(args) < 1 {
		msg := "No volume name specified."
		print_Error(msg, p.logger)
		flags.ParseArgs(&options, []string{"-h"})
		return FAIL, fmt.Errorf(msg)
	}
	//fmt.Println(args)
	if options.CommitMsg == "" {
		msg := "Empty commit message. Use -m to set commit message."
		print_Error(msg, p.logger)
		flags.ParseArgs(&options, []string{"-h"})
		return FAIL, fmt.Errorf(msg)
	}
	commitObj := CommitParams{
		commitMsg: options.CommitMsg, volumeName: args[0], layerUUID: options.Uuid,
	}
	commitObj.genUUID = commitObj.layerUUID == ""
	return volume_commit(commitObj, p.logger)

}

func (p OptSelector) clone(args []string) (int, error) {

	var options struct {
		//			Volume string `short:"v" long:"vol" description:"<volume_name>\tSpecify the volume name which needs to be update(restore).\n"`
		Layer string `short:"l" long:"layer" description:"Checkout <layer> instead of the HEAD\n"`

		CheckoutFlag bool `short:"n" long:"no-checkout" description:"No checkout of HEAD is performed after clone is complete.\n"`
	}
	os.Args = custom_Args(args, "<repo path>")
	args, err := flags.ParseArgs(&options, args[1:])

	if err != nil {
		//	flags.ParseArgs(&options, []string{"-h"})
		return FAIL, err
	}
	if len(args) != 1 {
		msg := "Invalid arguments. Use '-h' for help."
		print_Error(msg, p.logger)
		flags.ParseArgs(&options, []string{"-h"})
		return FAIL, fmt.Errorf(msg)
	}
	cloneObj := CloneParams{
		repoPath:    args[0],
		checkoutFlg: !options.CheckoutFlag,
		layerUUID:   options.Layer,
	}
	return clone_Repo(&cloneObj, p.logger)
}

func (p OptSelector) pull(args []string) (int, error) {

	p.logger.Println("pull", args)
	return FAIL, fmt.Errorf("Option unfinished.")
	return 0, nil
}

func (p OptSelector) push(args []string) (int, error) {

	p.logger.Println("push", args)
	return FAIL, fmt.Errorf("Option unfinished.")
	return 0, nil
}

func (p OptSelector) save(args []string) (int, error) {

	p.logger.Println("save", args)
	return FAIL, fmt.Errorf("Option unfinished.")
	return 0, nil
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
		print_Error(msg, p.logger)
		//	fmt.Println(OPT_LOG_USAGE)
		return FAIL, nil
	}
	args, err := flags.ParseArgs(&options, args[1:])
	if err != nil {
		//print_Error(err.Error(), p.logger)
		//fmt.Println(OPT_LOG_USAGE)
		return FAIL, err
	}
	//fmt.Println(args[0])
	//volume, err := confirm_BackingFilePath(args[0])
	volume := return_AbsPath(args[0])

	if volume == "" || !PathFileExists(volume) {
		//print_Error(err.Error(), p.logger)
		//fmt.Println(OPT_LOG_USAGE)
		msg := fmt.Sprintf("Volume '%s' can not found.", volume)
		print_Error(msg, p.logger)
		return FAIL, fmt.Errorf(msg)
	}
	return volume_commit_history(volume, p.logger)
}

func (p OptSelector) rebase(args []string) (int, error) {

	p.logger.Println("rebase", args)
	return FAIL, fmt.Errorf("Option unfinished.")
	return 0, nil
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
		print_Error(msg, p.logger)
		fmt.Println(OPT_RESET_USAGE)
		return FAIL, nil
	}
	args, err := flags.ParseArgs(&options, args[1:])
	if err != nil {
		print_Error(err.Error(), p.logger)
		fmt.Println(OPT_RESET_USAGE)
		return FAIL, err
	}
	resetObj := ResetParams{time: -1, volume: return_AbsPath(args[0])}
	if !PathFileExists(resetObj.volume) {
		msg := "volume can not find."
		print_Error(msg, p.logger)
		fmt.Printf(OPT_RESET_USAGE)
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
				print_Error(msg, p.logger)
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
			print_Error(msg, p.logger)
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

func (p OptSelector) tag(args []string) (int, error) {

	p.logger.Println("tag", args)
	return FAIL, fmt.Errorf("Option unfinished.")
	return 0, nil
}

func (p OptSelector) config(args []string) (int, error) {

	var options struct {
		Global string `long:"global" description:"[user.name|user.email] set global configuration."`
		Get    string `long:"get" description:"<name>\tGet value : <name>"`
	}
	os.Args = custom_Args(args, "")
	args, err := flags.ParseArgs(&options, args[1:])
	if err != nil {
		return FAIL, err
	}
	if options.Global != "" {
		if len(args) < 1 {
			msg := "Too few arguments."
			print_Error(msg, p.logger)
			return FAIL, fmt.Errorf(msg)
		}
		configObj, err := LoadConfig(p.logger)
		if err != nil {
			msg := fmt.Sprintf("Load configuration failed. Please check file '~/.hb/config.yaml' (%s)", err.Error())
			print_Error(msg, p.logger)
			return FAIL, err
		}
		if options.Global == "user.name" {
			configObj.UserName = args[0]
			print_Log(fmt.Sprintf("Set user.name as '%s'", configObj.UserName), p.logger)
			err = WriteConfig(&configObj, p.logger)
			if err != nil {
				msg := fmt.Sprintf("Write config failed. (%s)", err.Error())
				return FAIL, fmt.Errorf(msg)
			}
		} else if options.Global == "user.email" {
			configObj.UserEmail = args[0]
			print_Log(fmt.Sprintf("Set user.email as '%s'", configObj.UserEmail), p.logger)
			err = WriteConfig(&configObj, p.logger)
			if err != nil {
				msg := fmt.Sprintf("Write config failed. (%s)", err.Error())
				return FAIL, fmt.Errorf(msg)
			}
		} else {
			msg := "unknow option."
			print_Error(msg, p.logger)
			flags.ParseArgs(&options, []string{"-h"})
			return FAIL, fmt.Errorf(msg)
		}
		msg := Format_Success("Done.")
		print_Log(msg, p.logger)
		return OK, nil
	} else if options.Get != "" {
		configObj, err := LoadConfig(p.logger)
		if err != nil {
			msg := fmt.Sprintf("Load configuration failed. Please check file '~/.hb/config.yaml' (%s)", err.Error())
			print_Error(msg, p.logger)
			return FAIL, err
		}
		value, err := return_ConfigValue(&configObj, options.Get)
		if err != nil {
			msg := fmt.Sprintf("Load value of '%s' failed. (%s)", options.Get, err.Error())
			print_Error(msg, p.logger)
			return FAIL, err
		}
		msg := Format_Success(fmt.Sprintf("%s: %v", options.Get, value))
		print_Log(msg, p.logger)
		return OK, nil
	} else {
		msg := "Invalid option."
		print_Error(msg, p.logger)
		flags.ParseArgs(&options, []string{"-h"})
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
	return 0, nil
}

func (p OptSelector) post_checkout_hooks(args []string) (int, error) {

	p.logger.Println("post_checkout_hooks", args)
	return FAIL, fmt.Errorf("Option unfinished.")
	return 0, nil
}

func (p OptSelector) launch(args []string) (int, error) {

	p.logger.Println("launch", args)
	return FAIL, fmt.Errorf("Option unfinished.")
	return 0, nil
}

func (p OptSelector) list(args []string) (int, error) {

	print_Trace(args)
	targetdir, _ := return_CurrentDir()
	targetdir += "/" + DEFALUT_BACKING_FILE_DIR
	if len(args) > 1 {
		targetdir = args[1]
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
	if len(args) <= 1 {
		msg := "Too few arguments.\n"
		print_Error(msg, p.logger)
		fmt.Println(OPT_SHOW_USAGE)
		return FAIL, nil
	}
	image, err := confirm_BackingFilePath(args[1])
	if image == "" {
		print_Error(err.Error(), p.logger)
		usage := OPT_SHOW_USAGE
		fmt.Println(usage)
		return FAIL, err
	}

	return show_template(image, p.logger)
}
