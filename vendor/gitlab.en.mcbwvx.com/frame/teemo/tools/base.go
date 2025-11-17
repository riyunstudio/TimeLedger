package tools

import (
	"fmt"
	"time"
)

type Tools struct {
	loc *time.Location
}

func Initialize(serverLocation string) *Tools {
	loc, err := time.LoadLocation(serverLocation)
	if err != nil {
		panic(fmt.Errorf("初始化 自定義工具 錯誤, Err: %v", err))
	}

	return &Tools{
		loc: loc,
	}
}
