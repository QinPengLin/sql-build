package debug

import (
	"log"
)

// 错误调试
var Debug bool

func Printf(format string, v ...interface{}) {
	if !Debug {
		return
	}
	log.Printf(format, v...)
}

func Warning(v ...interface{}) {
	if !Debug {
		return
	}
	log.Println("Warning:", v)
}

func Error(v ...interface{}) {
	if !Debug {
		return
	}
	log.Println("Error:", v)
}

func Println(v ...interface{}) {
	if !Debug {
		return
	}
	log.Println(v...)
}

func Fatal(v ...interface{}) {
	if !Debug {
		return
	}
	log.Fatal(v...)
}
