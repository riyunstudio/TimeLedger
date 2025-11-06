# ===Linter工具來源===
FROM golangci/golangci-lint:v2.4.0 AS lintbin

# ===建置階段===
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN apk add --no-cache git
RUN go mod download

#安裝工具
RUN go install github.com/swaggo/swag/cmd/swag@v1.16.4
COPY --from=lintbin /usr/bin/golangci-lint /usr/local/bin/golangci-lint
ENV PATH=/go/bin:$PATH

#複製原始碼
COPY . .

# === 可選步驟 ===
# 若在 CI/CD 要檢查 Swagger / Lint，可打開以下行
RUN swag init && golangci-lint run --timeout 10m

# === 編譯 ===
RUN go build -mod=vendor -o main .

# ===運行階段 === golang:1.25-alpine（內建 3.21）
FROM alpine:3.21

WORKDIR /app

# 加入時區資料
RUN apk add --no-cache tzdata

# 只帶必要檔案
COPY --from=builder /app/main .
COPY --from=builder /app/docs ./docs

# 設定時區環境變數（可選）
ENV TZ=Asia/Taipei

EXPOSE 8080
CMD ["./main"]