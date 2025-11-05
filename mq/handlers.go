package mq

import (
	"fmt"
)

type MsgRegistration struct {
	Type    string
	NewMsg  func() interface{}
	Handler func(msg interface{}) error
}

func GetRegistrations() []MsgRegistration {
	return []MsgRegistration{
		{
			Type:   "normal",
			NewMsg: func() interface{} { return &User{} },
			Handler: func(msg interface{}) error {
				user := msg.(*User)
				fmt.Println("Process normal user:", user.ID, user.Name)
				return nil
			},
		},
		{
			Type:   "delay",
			NewMsg: func() interface{} { return &User{} },
			Handler: func(msg interface{}) error {
				user := msg.(*User)
				fmt.Println("Process delay user:", user.ID, user.Name)
				return nil
			},
		},
	}
}
