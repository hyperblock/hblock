package hblock

import (
	"fmt"

	"log"

	"libguestfs.org/guestfs"
)

type InitParams struct {
	name string
	size int64
}

func create_empty_template(obj InitParams, logger *log.Logger) (int, error) {

	//output := obj.name
	g, errno := guestfs.Create()
	if errno != nil {
		return FAIL, errno
	}
	//	defer
	//fmt.Println(size)
	if errCreate := g.Disk_create(obj.name, "qcow2", obj.size, nil); errCreate != nil {
		//return FAIL, errCreate
		g.Close()
		print_Panic(errCreate.Errmsg, logger)
	}
	msg := fmt.Sprintf("Create template '%s' finished.", obj.name)
	g.Close()
	print_Log(format_Success(msg), logger)
	//fmt.Println(msg)

	return OK, nil
}
