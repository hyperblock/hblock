package hblock

import "os/exec"
import "strings"
import "fmt"

func show_template(image string) (int, error) {

	cmd := exec.Command("qcow2-img", "info", image)
	result, err := cmd.Output()
	if err != nil {
		return FAIL, err
	}
	resultInfo := string(result)
	snapshotList := strings.Split(
		get_StringBefore(
			get_StringAfter(resultInfo, "list:\n"), "Format"),
		"\n")

	for _, val := range snapshotList[1:] {
		fmt.Println(val)
		colVal := strings.Join(strings.Split(val, " "), "|")
		fmt.Println(colVal)
	}
	return OK, nil
}
