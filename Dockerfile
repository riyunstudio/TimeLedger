# 先宣告一個 stage，只用來提供 golangci-lint 二進位
FROM golangci/golangci-lint:v2.4.0 AS lintbin

# 真正的建置/執行基底
FROM golang:1.25-alpine

WORKDIR /app

# 先處理依賴與工具（可被快取）
COPY go.mod go.sum ./
RUN apk add --no-cache git
RUN go install github.com/swaggo/swag/cmd/swag@v1.16.4
COPY --from=lintbin /usr/bin/golangci-lint /usr/local/bin/golangci-lint
ENV PATH=/go/bin:$PATH

# 再拷專案原始碼（變動最大）
COPY . .

# 產生 Swagger、跑 Lint、編譯
RUN swag init
RUN golangci-lint run --timeout 10m
RUN go build -mod=vendor -o main .

EXPOSE 8080
CMD ["./main"]