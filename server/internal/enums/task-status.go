package enums

import (
	"log"
)

type TaskStatus int64

const (
	New TaskStatus = iota
	InProcess
	Done
	Error
)

func (t TaskStatus) Itoa() (str string) {
	switch t {
	case 0:
		return "new"
	case 1:
		return "in_process"
	case 2:
		return "done"
	case 3:
		return "error"
	default:
		log.Printf("incorrect status: %d", t)
		return "wrong_status"
	}
}
