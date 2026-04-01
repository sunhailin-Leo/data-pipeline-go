# Define os name
OS_NAME=$(OS)

# Define go build
OUTPUT_BIN=data-pipeline-go
BUILD_DIR=cmd
CGO=1

ifeq ($(OS_NAME), Windows_NT)
	CHCP_CMD := chcp 65001 > nul
	TARGET_OS = windows
	TARGET_ARCH = amd64
	CMD_PREFIX = powershell -Command
	CMD_SUFFIX = .exe
	CMD_AND = ;
	BUILD_ARGS := set CGO_ENABLED=${CGO}; set GOOS=${TARGET_OS}; set GOARCH=${TARGET_ARCH}
	CLEAN_CMD := Remove-Item -Path ${BUILD_DIR}/${OUTPUT_BIN}${CMD_SUFFIX}; Remove-Item -Path coverage.out -ErrorAction SilentlyContinue; Remove-Item -Path coverage.html -ErrorAction SilentlyContinue
else
	OS_NAME=$(shell uname)
	CHCP_CMD :=
	TARGET_OS = linux
	TARGET_ARCH = amd64
	CMD_PREFIX =
	CMD_SUFFIX =
	CMD_AND = &&
	BUILD_ARGS := export CGO_ENABLED=${CGO} GOOS=${TARGET_OS} GOARCH=${TARGET_ARCH}
	CLEAN_CMD := rm -f ${BUILD_DIR}/${OUTPUT_BIN}${CMD_SUFFIX} coverage.out coverage.html
endif

# Define golangci-lint
GOLANGCI_LINT=golangci-lint
GOLANGCI_LINT_VERBOSE=--verbose

# Define nilaway
NILAWAY=nilaway

# 默认目标, 显示帮助信息
.PHONY: help
help:
	@$(CHCP_CMD)
	@echo "Available targets:"
	@echo    - help              [显示帮助信息]
	@echo    - lint              [代码静态检查]
	@echo    - nilaway           [代码 nil 检查]
	@echo    - test              [单元测试]
	@echo    - integration-up    [启动集成测试服务]
	@echo    - integration-down  [停止集成测试服务]
	@echo    - integration-test  [运行集成测试]
	@echo    - build             [构建二进制包]
	@echo    - clean             [清理构建文件]

.PHONY: all
all: lint-all build

.PHONY: lint-all
lint-all: lint nilaway

.IGNORE: lint
.PHONY: lint
lint:
	@$(CHCP_CMD)
	-${GOLANGCI_LINT} run ${GOLANGCI_LINT_VERBOSE}

.IGNORE: nilaway
.PHONY: nilaway
nilaway:
	@${CHCP_CMD}
	${NILAWAY} ./...

.PHONY: test
test:
	@${CHCP_CMD}
	go test -v -cover -race ./...

.PHONY: integration-up
integration-up:
	@$(CHCP_CMD)
	docker compose -f docker-compose.integration.yml up -d --wait

.PHONY: integration-down
integration-down:
	@$(CHCP_CMD)
	docker compose -f docker-compose.integration.yml down -v

.PHONY: integration-test
integration-test: integration-up
	@$(CHCP_CMD)
	INTEGRATION_TEST=true \
	INTEGRATION_REDIS_ADDR=localhost:6379 \
	INTEGRATION_KAFKA_ADDR=localhost:9092 \
	INTEGRATION_MYSQL_ADDR=localhost:3306 \
	INTEGRATION_MYSQL_USER=root \
	INTEGRATION_MYSQL_PASS=testpass \
	INTEGRATION_MYSQL_DB=integration_test \
	INTEGRATION_POSTGRES_ADDR=localhost:5432 \
	INTEGRATION_POSTGRES_USER=testuser \
	INTEGRATION_POSTGRES_PASS=testpass \
	INTEGRATION_POSTGRES_DB=integration_test \
	INTEGRATION_CLICKHOUSE_ADDR=localhost:9000 \
	INTEGRATION_CLICKHOUSE_USER=default \
	INTEGRATION_CLICKHOUSE_PASS=testpass \
	INTEGRATION_CLICKHOUSE_DB=integration_test \
	INTEGRATION_RABBITMQ_ADDR=localhost:5672 \
	INTEGRATION_RABBITMQ_USER=testuser \
	INTEGRATION_RABBITMQ_PASS=testpass \
	INTEGRATION_ROCKETMQ_ADDR=127.0.0.1:9876 \
	INTEGRATION_PULSAR_ADDR=localhost:6650 \
	INTEGRATION_ES_ADDR=http://localhost:9200 \
	INTEGRATION_ES_USER=elastic \
	INTEGRATION_ES_PASS=testpass \
	go test -v -race -timeout 180s ./pkg/sink/... ./pkg/source/...

.PHONY: pre-ci-build
pre-ci-build:
	@$(CHCP_CMD)
	export GOFLAGS="-buildvcs=false"

.PHONY: build
build:
	@$(CHCP_CMD)
	${CMD_PREFIX} cd ${BUILD_DIR} ${CMD_AND} ${BUILD_ARGS} ${CMD_AND} go build -o ${OUTPUT_BIN}${CMD_SUFFIX} .

.PHONY: clean
clean:
	@$(CHCP_CMD)
	${CMD_PREFIX} ${CLEAN_CMD}

.PHONY: fmt
fmt:
	@$(CHCP_CMD)
	gofmt -s -w .

.PHONY: docker-build
docker-build:
	@$(CHCP_CMD)
	docker build -t data-pipeline-go .

.PHONY: docker-run
docker-run:
	@$(CHCP_CMD)
	docker run --rm data-pipeline-go

.PHONY: coverage
coverage:
	@$(CHCP_CMD)
	go test -v -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out -o coverage.html

.PHONY: benchmark
benchmark:
	@$(CHCP_CMD)
	go test -bench=. -benchmem -run=^$$ ./pkg/utils/... ./pkg/transform/convert/... ./pkg/middlewares/...
