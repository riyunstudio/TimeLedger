package apis

import (
	"akali/configs"
	"akali/tools"
)

type Api struct {
	Env   *configs.Env
	Tools *tools.Tools
}

func Initialize(env *configs.Env, tools *tools.Tools) *Api {
	return &Api{
		Env:   env,
		Tools: tools,
	}
}
