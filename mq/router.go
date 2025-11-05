package mq

import (
	"encoding/json"
	"fmt"
	"sync"
)

type Router struct {
	handlers map[string]func([]byte) error
	mutex    sync.RWMutex
}

func NewRouter() *Router {
	return &Router{handlers: make(map[string]func([]byte) error)}
}

func (r *Router) RegisterBulk(regs []MsgRegistration) {
	for _, reg := range regs {
		r.handlers[reg.Type] = func(msg []byte) error {
			m := reg.NewMsg()
			if err := json.Unmarshal(msg, m); err != nil {
				return err
			}
			return reg.Handler(m)
		}
	}
}

func (r *Router) Route(msg []byte) error {
	var base struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(msg, &base); err != nil {
		return err
	}

	r.mutex.RLock()
	h, ok := r.handlers[base.Type]
	r.mutex.RUnlock()
	if !ok {
		return fmt.Errorf("no handler for type: %s", base.Type)
	}
	return h(msg)
}
