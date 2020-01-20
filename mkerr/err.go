package mkerr

import (
	"fmt"
	"runtime"
	"strings"
)

type Error interface {
	ErrNo() int64
	Error() string
	ErrDetail() string
}

type mkerr struct {
	errNo     int64
	errMsg    string
	errDetail string
}

func (s mkerr) ErrNo() int64 {
	return s.errNo
}

func (s mkerr) Error() string {
	return s.errMsg
}

func (s mkerr) ErrDetail() string {
	return s.errDetail
}

func NewError(errNo int64, errMsg string, errDetail string) Error {
	return mkerr{
		errNo:     errNo,
		errMsg:    errMsg,
		errDetail: fmt.Sprintf("error: %s", errDetail) + getStackInfo(),
	}
}

func getStackInfo() string {
	stackInfo := ""
	pcs := make([]uintptr, 128)
	runtime.Callers(3, pcs)
	for _, pc := range pcs {
		if pc == 0 {
			break
		}
		theFunc := runtime.FuncForPC(pc)
		file, line := theFunc.FileLine(pc)
		if strings.HasPrefix(file, "runtime") {
			continue
		}
		stackInfo += fmt.Sprintf("\n    %s:%d", file, line)
	}
	return stackInfo
}
