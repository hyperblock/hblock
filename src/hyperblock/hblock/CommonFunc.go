package hblock

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func print_Log(msg string, logger *log.Logger) {

	if len(msg) == 0 {
		return
	}
	if msg[0] != 27 {
		msg = format_Info(msg)
	}
	fmt.Println(msg)
	if logger != nil {
		logger.Println(msg)
	}
}

func print_Error(err string, logger *log.Logger) {

	if len(err) == 0 {
		return
	}
	err = format_Error(err)
	fmt.Println(err)
	if logger != nil {
		logger.Println(err)
	}
}

func print_Panic(msg string, logger *log.Logger) {

	if len(msg) == 0 {
		return
	}
	msg = format_Error(msg)
	fmt.Println(msg)
	if logger != nil {
		logger.Panicln(msg)
	}
}

func print_Fatal(msg string, logger *log.Logger) {

	if len(msg) == 0 {
		return
	}
	msg = format_Error(msg)
	fmt.Println(msg)
	if logger != nil {
		logger.Fatalln(msg)
	}
}

func get_StringAfter(content string, prefix string) string {

	q := strings.Index(content, prefix)
	if q == -1 {
		return content
	}
	return content[q+len(prefix):]
}

func get_StringBefore(content string, suffix string) string {

	q := strings.Index(content, suffix)
	if q == -1 {
		return content
	}
	return content[:q]
}

func get_InfoValue(list []string, keyword string) string {

	for _, content := range list {
		exist := strings.HasPrefix(content, keyword)
		if exist {
			return content[len(keyword)+2:]
		}
	}
	return ""
}

func return_TemplateDir() (string, error) {

	ret := os.Getenv("HBLOCK_TEMPLATE_DIR")
	//	fmt.Println("env", fmt.Sprintf("[%s]", ret))
	if ret == "" {
		msg := format_Warning("env HBLOCK_TEMPLATE_DIR not set. use default dir '/var/hyperblock'")
		return "/var/hyperblock", fmt.Errorf(msg)
	}
	return ret, nil
}

func confirm_BackingFilePath(imgPath string) (string, error) {

	_, err := os.Stat(imgPath)
	if err != nil {
		if os.IsNotExist(err) {
			path, errTemp := return_TemplateDir()
			path += "/" + imgPath
			_, err = os.Stat(path)
			if err != nil {
				return "", err
			} else {
				return path, errTemp
			}
		} else {
			return "", err
		}
	}
	return imgPath, nil
}

// see complete color rules in document in https://en.wikipedia.org/wiki/ANSI_escape_code#cite_note-ecma48-13
func format_Trace(format string, a ...interface{}) string {
	prefix := yellow(trac)
	return fmt.Sprint(formatLog(prefix), fmt.Sprintf(format, a...))
}

func format_Info(format string, a ...interface{}) string {
	prefix := blue(info)
	return fmt.Sprint(formatLog(prefix), fmt.Sprintf(format, a...))
}

func format_Success(format string, a ...interface{}) string {
	prefix := green(succ)
	return fmt.Sprint(formatLog(prefix), fmt.Sprintf(format, a...))
}

func format_Warning(format string, a ...interface{}) string {
	prefix := magenta(warn)
	return fmt.Sprint(formatLog(prefix), fmt.Sprintf(format, a...))
}

func format_Error(format string, a ...interface{}) string {
	prefix := red(erro)
	return fmt.Sprint(formatLog(prefix), fmt.Sprintf(format, a...))
}

func red(s string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", color_red, s)
}

func green(s string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", color_green, s)
}

func yellow(s string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", color_yellow, s)
}

func blue(s string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", color_blue, s)
}

func magenta(s string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", color_magenta, s)
}

func formatLog(prefix string) string {

	return prefix + " "
}
