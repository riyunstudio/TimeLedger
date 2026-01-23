package app

import (
	"timeLedger/apis"
	"timeLedger/configs"
	"timeLedger/database/mysql"
	"timeLedger/database/redis"
	"timeLedger/global/errInfos"
	"timeLedger/rpc"

	"gitlab.en.mcbwvx.com/frame/teemo/tools"
)

type App struct {
	Env   *configs.Env
	Err   *errInfos.ErrInfo
	Tools *tools.Tools
	MySQL *mysql.DB
	Redis *redis.Redis
	Api   *apis.Api
	Rpc   *rpc.Rpc
}

func Initialize() *App {
	env := configs.LoadEnv()

	tools := tools.Initialize(env.AppTimezone)

	e := errInfos.Initialize(env.AppID)

	mysql := mysql.Initialize(env)
	mysql.AutoMigrate()
	mysql.Seeds(tools)

	redis := redis.Initialize(env)

	api := apis.Initialize(env, tools)

	rpc := rpc.Initialize(env, tools)

	return &App{
		Env:   env,
		Err:   e,
		Tools: tools,
		MySQL: mysql,
		Redis: redis,
		Api:   api,
		Rpc:   rpc,
	}
}
