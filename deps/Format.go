package hblock

import (
	"fmt"
	"log"
)

func print_Log(msg string, logger *log.Logger) {

	if len(msg) == 0 {
		return
	}
	if msg[0] != '\r' && msg[len(msg)-1] != '\n' {
		msg += "\n"
	}
	if (msg[0] != 27) || (msg[0] == '\r' && msg[1] != 27) {
		if msg[0] == '\r' {
			msg = "\r" + Format_Info(msg[1:])
		} else {
			msg = Format_Info(msg)
		}
	}

	fmt.Print(msg)
	if logger != nil {
		msg = msg[15:]
		logger.Println(msg)
	}
}

func Print_Error(msg string, logger *log.Logger) {

	if len(msg) == 0 {
		return
	}
	msg = Format_Error(msg)
	fmt.Println(msg)
	if logger != nil {
		msg = msg[15:]
		logger.Println(msg)
	}
}

func print_Panic(msg string, logger *log.Logger) {

	if len(msg) == 0 {
		return
	}
	msg = Format_Error(msg)
	fmt.Println(msg)
	if logger != nil {
		msg = msg[15:]
		logger.Panicln(msg)
	}
}

func print_Fatal(msg string, logger *log.Logger) {

	if len(msg) == 0 {
		return
	}
	msg = Format_Error(msg)
	fmt.Println(msg)
	if logger != nil {
		msg = msg[15:]
		logger.Fatalln(msg)
	}
}

func print_Trace(a ...interface{}) {

	if SHOW_TRACE == 1 {
		msg := Format_Trace(fmt.Sprint(a))
		fmt.Println(msg)
	}
}

func print_ProcessBar(current, total int64) string {

	bar := "["
	base := int((float32(current) / float32(total)) * 100)
	delta := int(float32(base)/float32(5) + 0.5)
	for i := 0; i < delta; i++ {
		bar += "="
	}
	delta = 20 - delta
	for i := 0; i < delta; i++ {
		bar += " "
	}
	bar += "]"
	A, B := current>>20, total>>20
	if A == 0 {
		A = 1
	}
	if B == 0 {
		B = 1
	}
	ret := fmt.Sprintf("%s %d%% (%d/%d)", bar, base, A, B)
	return ret
}

// see complete color rules in document in https://en.wikipedia.org/wiki/ANSI_escape_code#cite_note-ecma48-13
func Format_Trace(format string, a ...interface{}) string {
	prefix := yellow(trac)
	return fmt.Sprint(formatLog(prefix), fmt.Sprintf(format, a...))
}

func Format_Info(format string, a ...interface{}) string {
	prefix := blue(info)
	return fmt.Sprint(formatLog(prefix), fmt.Sprintf(format, a...))
}

func Format_Success(format string, a ...interface{}) string {
	prefix := green(succ)
	return fmt.Sprint(formatLog(prefix), fmt.Sprintf(format, a...))
}

func Format_Warning(format string, a ...interface{}) string {
	prefix := magenta(warn)
	return fmt.Sprint(formatLog(prefix), fmt.Sprintf(format, a...))
}

func Format_Error(format string, a ...interface{}) string {
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
