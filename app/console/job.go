package console

// 所有排程任務必須實現這個接口
type Job interface {
	Name() string        // 名稱
	Description()        // 說明
	Repositories()       // Repository
	Handle(string) error // 主程式
}
