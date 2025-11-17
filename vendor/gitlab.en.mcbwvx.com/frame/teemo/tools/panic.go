package tools

import (
	"fmt"
	"runtime"
)

type Panic struct {
	Panic      string `json:"panic"`
	StackTrace string `json:"stack_trace"`
}

func (tl *Tools) PanicParser(err any) Panic {
	buf := make([]byte, 4096)
	for {
		n := runtime.Stack(buf, true) // 抓所有 goroutine
		if n < len(buf) {
			var msg string
			switch r := err.(type) {
			case runtime.Error:
				msg = fmt.Sprintf("Runtime Error: %v", r)
			case error:
				msg = fmt.Sprintf("Error: %v", r)
			default:
				msg = fmt.Sprintf("Unknown Panic: %v", r)
			}
			return Panic{
				Panic:      msg,
				StackTrace: string(buf[:n]),
			}
		}
		buf = make([]byte, len(buf)*2) // 自動擴充
	}
}
