package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

func main() {
	start := time.Now()
	
	fmt.Println("╔════════════════════════════════════════════════════════════════╗")
	fmt.Println("║              TimeLedger 自動化測試 Runner                 ║")
	fmt.Println("╚════════════════════════════════════════════════════════════════╝")
	fmt.Println()

	// 運行測試
	fmt.Println("🚀 運行測試...")
	fmt.Println()
	
	cmd := exec.Command("go", "test", "./testing/test/...", "-v", "-count=1")
	cmd.Dir = "d:\\project\\TimeLedger"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	err := cmd.Run()
	
	duration := time.Since(start)
	
	fmt.Println()
	fmt.Println("╔════════════════════════════════════════════════════════════════╗")
	fmt.Println("║                    測試執行完成                           ║")
	fmt.Println("╚════════════════════════════════════════════════════════════════╝")
	fmt.Println()
	
	if err != nil {
		fmt.Printf("❌ 測試失敗 (耗時: %s)\n", duration.Round(time.Second))
		fmt.Println("請檢查上方輸出以獲取詳細錯誤信息。")
		os.Exit(1)
	} else {
		fmt.Printf("✅ 所有測試通過 (耗時: %s)\n", duration.Round(time.Second))
		fmt.Println()
		fmt.Println("快速指令:")
		fmt.Println("  make test           - 運行所有測試")
		fmt.Println("  make test-coverage  - 生成覆蓋率報告")
		fmt.Println("  run-test.bat        - Windows 測試腳本")
		fmt.Println("  run-test.bat help   - 查看更多選項")
	}
}

// Helper function to parse test output
func parseTestOutput(output string) (passed, failed int) {
	passedLines := strings.Count(output, "--- PASS:")
	failedLines := strings.Count(output, "--- FAIL:")
	
	// 排除 subtests 的重複計算
	passed = strings.Count(output, "PASS:")
	failed = strings.Count(output, "FAIL:")
	
	if failed > 0 {
		return passed, failed
	}
	return passedLines, failedLines
}
