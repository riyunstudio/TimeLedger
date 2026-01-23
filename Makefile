# TimeLedger Makefile

.PHONY: test test-coverage test-watch test-report clean help

# 運行所有測試
test:
	@echo "🚀 運行 TimeLedger 測試..."
	go test ./testing/test/... -v -count=1

# 運行測試並生成覆蓋率報告
test-coverage:
	@echo "📊 運行測試並生成覆蓋率報告..."
	go test ./testing/test/... -v -count=1 -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "✅ 覆蓋率報告已生成: coverage.html"
	@open coverage.html 2>/dev/null || echo "請手動打開 coverage.html 查看報告"

# 運行單一測試檔案
test-file:
	@if [ -z "$(file)" ]; then \
		echo "❌ 請指定測試檔案: make test-file file=auth_test"; \
	else \
		echo "🔍 運行測試檔案: $(file)"; \
		go test ./testing/test/$(file).go -v; \
	fi

# 運行特定名稱的測試
test-name:
	@if [ -z "$(name)" ]; then \
		echo "❌ 請指定測試名稱: make test-name name=TestAuth"; \
	else \
		echo "🔍 運行測試: $(name)"; \
		go test ./testing/test/... -v -run "$(name)"; \
	fi

# 監控模式 - 檔案變動時自動重新測試
test-watch:
	@echo "👀 監控模式 - 檔案變動時自動重新測試 (按 Ctrl+C 結束)"
	@while true; do \
		inotifywait -e modify -e create -e delete -r ./testing/test/ ./app/ 2>/dev/null || true; \
		echo "📝 偵測到變動，重新運行測試..."; \
		go test ./testing/test/... -v -count=1 2>&1 | head -50; \
		echo "---"; \
	done

# 生成測試報告 (JUnit XML格式)
test-report:
	@echo "📋 生成測試報告..."
	go test ./testing/test/... -v -count=1 -xml=test-results.xml 2>&1 | tee test-output.txt
	@echo "✅ 測試報告已生成: test-results.xml"
	@if command -v powershell &> /dev/null; then \
		Start-Process "test-results.xml"; \
	fi

# 清理測試產生的檔案
clean:
	@echo "🧹 清理測試檔案..."
	rm -f coverage.out coverage.html test-results.xml test-output.txt
	@echo "✅ 清理完成"

# 快速測試 (安靜模式)
test-quick:
	@echo "⚡ 快速測試..."
	go test ./testing/test/... -count=1 -short

# 顯示幫助
help:
	@echo "📚 TimeLedger 測試命令"
	@echo ""
	@echo "可用命令:"
	@echo "  make test           - 運行所有測試"
	@echo "  make test-coverage  - 運行測試並生成覆蓋率報告"
	@echo "  make test-file      - 運行單一測試檔案 (需指定 file=xxx)"
	@echo "  make test-name      - 運行特定名稱的測試 (需指定 name=xxx)"
	@echo "  make test-report    - 生成測試報告 (JUnit XML)"
	@echo "  make test-quick     - 快速測試 (安靜模式)"
	@echo "  make clean          - 清理測試檔案"
	@echo "  make help           - 顯示此幫助"
	@echo ""
	@echo "範例:"
	@echo "  make test file=auth_test"
	@echo "  make test name=TestAuthController"
