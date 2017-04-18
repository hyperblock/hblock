package hblock

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"strconv"

	"strings"

	"github.com/voxelbrain/goptions"
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

	//args = strings.Split("checkout -t /var/hyperblock/centos.img -l sp001 -o test002", " ")
	//args = []string{"init", "--name", "hehe"}
	if len(args) == 0 {
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
		return FAIL, fmt.Errorf("invalid option, and there is no --help :)")
	}
}

func (p OptSelector) init(args []string) (int, error) {

	options := struct {
		TemplateName string        `goptions:"--name, description='[required] Template name.\n'"`
		Size         string        `goptions:"--size, description='[required] Disk size(M/G) of template.\n\t\t\t    eg.\n\t\t\t\thblock init --name template0 --size=500M.\n'"`
		Help         goptions.Help `goptions:"-h,--help, description='Show this help\n'"`
	}{}
	os.Args = args
	goptions.ParseAndFail(&options)
	fmt.Println(options.TemplateName)
	//p.logger.Println("init", args)
	print_Log(strings.Join(args, " "), p.logger)
	if options.TemplateName == "" {
		msg := "--name is required."
		return FAIL, fmt.Errorf(msg)
	}
	if options.Size == "" {
		msg := "--size is required."
		p.logger.Println(msg)
		return FAIL, fmt.Errorf(msg)
	}
	strSize := options.Size
	unit := strSize[len(strSize)-1:]
	_size, err := strconv.Atoi(strSize[0 : len(strSize)-1])

	if err != nil {
		msg := "invalid template size"
		//fmt.Println(options.Help)
		p.logger.Println(msg)
		return FAIL, fmt.Errorf(msg)
	}
	var sizeI64 int64
	if unit == "M" {
		sizeI64 = int64(_size * 1024 * 1024)
	} else if unit == "G" {
		sizeI64 = int64(_size*1024*1024) * 1024
	}
	obj := InitParams{name: options.TemplateName, size: sizeI64}
	return create_empty_template(obj, p.logger)
	//return 0, nil
}

func (p OptSelector) checkout(args []string) (int, error) {

	options := struct {
		Help     goptions.Help `goptions:"-h,--help, description='Show this help.\n'"`
		Volume   string        `goptions:"-v,--vol, description='<volume_name>\tSpecify the volume name which needs to be update(restore).\n'"`
		Layer    string        `goptions:"-l,--layer, description='<layer>\tSpecify the <layer> that this volume will restore. If the <layer> does\\'not exist, it will create a new layer from current volume.'\n"`
		Output   string        `goptions:"-o,--output, description='[required if use \\'-t\\'] <output_volume_path>.\n'"`
		Template string        `goptions:"-t,--template, description='<template_name>\t Create a new volume from template.\n'"`
		Force    bool          `goptions:"-f,--force"`
	}{}
	os.Args = args
	goptions.ParseAndFail(&options)
	if options.Template != "" && options.Volume != "" {
		msg := "Can't use both -l and -t."
		print_Error(msg, p.logger)
		return FAIL, fmt.Errorf(msg)
	}
	if options.Output != "" && options.Volume != "" {
		msg := "Can't use both -v and -o."
		print_Error(msg, p.logger)
		return FAIL, fmt.Errorf(msg)
	}
	if options.Volume != "" {
		_, err := os.Stat(options.Volume)
		if err != nil {
			msg := fmt.Sprintf("Can't locate volume_name '%s'", options.Volume)
			print_Error(msg, p.logger)
			return FAIL, fmt.Errorf(msg)
		}
		if !options.Force {
			msg := "Need commit the exist volume or use -f"
			print_Error(msg, p.logger)
			return FAIL, fmt.Errorf(msg)
		}
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

	p.logger.Println("commit", args)
	return FAIL, fmt.Errorf("Option unfinished.")
	return 0, nil
}

func (p OptSelector) clone(args []string) (int, error) {

	p.logger.Println("clone", args)
	return FAIL, fmt.Errorf("Option unfinished.")
	return 0, nil
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

	p.logger.Println("log", args)
	return FAIL, fmt.Errorf("Option unfinished.")
	return 0, nil
}

func (p OptSelector) rebase(args []string) (int, error) {

	p.logger.Println("rebase", args)
	return FAIL, fmt.Errorf("Option unfinished.")
	return 0, nil
}

func (p OptSelector) reset(args []string) (int, error) {

	p.logger.Println("reset", args)
	return FAIL, fmt.Errorf("Option unfinished.")
	return 0, nil
}

func (p OptSelector) tag(args []string) (int, error) {

	p.logger.Println("tag", args)
	return FAIL, fmt.Errorf("Option unfinished.")
	return 0, nil
}

func (p OptSelector) config(args []string) (int, error) {

	p.logger.Println("config", args)
	return FAIL, fmt.Errorf("Option unfinished.")
	return 0, nil
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

	p.logger.Println("list", args)
	targetdir := "/var/hyperblock"
	if len(args) == 1 {
		targetdir = args[0]
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
	image := args[0]
	cmd := exec.Command("qcow2-img", "info", image)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return FAIL, err
	}
	return OK, err
}
