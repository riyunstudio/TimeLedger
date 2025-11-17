package rpc

import (
	"akali/configs"
	"fmt"

	"gitlab.en.mcbwvx.com/frame/teemo/tools"
)

type Rpc struct {
	Env   *configs.Env
	Tools *tools.Tools
}

func Initialize(env *configs.Env, tools *tools.Tools) *Rpc {
	defer func() {
		if err := recover(); err != nil {
			panic(fmt.Errorf("初始化 [rpc client] 發生 panic, Err: %v", tools.PanicParser(err)))
		}
	}()

	// for _, cfg := range rpcClientsConfig {
	// 	// 連線到 RPC 節點
	// 	client, err := ethclient.Dial(cfg.EndPoint)
	// 	if err != nil {
	// 		panic(fmt.Errorf("初始化 [rpc] 失敗，無法連線節點: %v", err))
	// 	}
	// }

	return &Rpc{
		Env:   env,
		Tools: tools,
	}
}
