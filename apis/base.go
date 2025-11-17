package apis

import (
	"akali/configs"

	"gitlab.en.mcbwvx.com/frame/teemo/tools"
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
