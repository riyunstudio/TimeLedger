package app

import (
	"akali/apis"
	"akali/configs"
	"akali/database/mysql"
	"akali/database/redis"
	"akali/libs/tools"
	rpcClient "akali/rpc/clients"
)

type App struct {
	Env       *configs.Env
	Tools     *tools.Tools
	Mysql     *mysql.DB
	Redis     *redis.Redis
	Api       *apis.Api
	RpcClient *rpcClient.RpcClient
}

func Initialize() *App {
	env := configs.LoadEnv()

	tools := tools.Initialize(env.AppTimezone)

	mysql := mysql.Initialize(env)
	mysql.AutoMigrate()
	mysql.Seeds(tools)

	redis := redis.Initialize(env)

	api := apis.Initialize(env, tools)

	rpcClient := rpcClient.Initialize(env, tools)

	return &App{
		Env:       env,
		Tools:     tools,
		Mysql:     mysql,
		Redis:     redis,
		Api:       api,
		RpcClient: rpcClient,
	}
}
