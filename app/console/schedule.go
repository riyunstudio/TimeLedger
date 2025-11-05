package console

import (
	"akali/app"
	"akali/libs/logs"
	"fmt"
	"log"
	"sync"

	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	cron *cron.Cron
	app  *app.App

	wg sync.WaitGroup // For 優雅退出
}

// 初始化 scheduler
func Initialize(app *app.App) *Scheduler {
	return &Scheduler{
		cron: cron.New(cron.WithSeconds()), // 支援秒級
		app:  app,
	}
}

// 註冊任務
func (s *Scheduler) addJob(spec string, job Job) {
	_, err := s.cron.AddFunc(spec, func() {
		s.wg.Add(1)
		defer s.wg.Done()

		// 初始化 TraceLog
		traceLog := logs.TraceLogInit()
		traceLog.SetTopic("Schedule")
		traceLog.SetMethod(job.Name())
		tid, err := s.app.Tools.NewTraceId()
		if err != nil {
			tid = ""
		}
		traceLog.SetTraceID(tid)

		defer func() {
			if r := recover(); r != nil {
				e := s.app.Tools.PanicParser(r)
				traceLog.PrintPanic(e)
			}
		}()

		// 讀取需要的 Repository
		job.Repositories()

		// 執行任務
		if err := job.Handle(tid); err != nil {
			// 發送 Tg警報
			if s.app.Env.TelegramBotToken != "" && s.app.Env.TelegramChatID != "" {
				_ = s.app.Api.PostTelegramMessage(tid, s.app.Tools.TgSchedulerAlert(job.Name(), err, tid))
			}

			// 寫 TraceLog
			traceLog.PrintSchedule(err)
		}
	})
	if err != nil {
		panic(fmt.Errorf("Scheduler AddJob error: %v", err))
	}
}

// 秒 分 時 日 月 星期 * * * * * *
func (s *Scheduler) loadJobs() {

}

// 啟動排程
func (s *Scheduler) Start() {
	s.loadJobs()
	s.cron.Start()
}

// 停止排程
func (s *Scheduler) Stop() {
	log.Println("Scheduler: shutting down gracefully...")

	// 停止接收新任務
	s.cron.Stop()

	// 等待所有任務完成
	s.wg.Wait()

	log.Println("All scheduled jobs completed. Exiting.")
}
