package tools

import "log"

func CheckError(err error) {
	if err != nil {
		log.Println("出现错误:", err)
	}
}
