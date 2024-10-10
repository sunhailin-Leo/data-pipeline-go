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
	CLEAN_CMD := Remove-Item -Path ${BUILD_DIR}/${OUTPUT_BIN}${CMD_SUFFIX}
else
	OS_NAME=$(shell uname)
	CHCP_CMD :=
	TARGET_OS = linux
	TARGET_ARCH = amd64
	CMD_PREFIX =
	CMD_SUFFIX =
	CMD_AND = &&
	BUILD_ARGS := export CGO_ENABLED=${CGO} GOOS=${TARGET_OS} GOARCH=${TARGET_ARCH}
	CLEAN_CMD := rm -f ${BUILD_DIR}/${OUTPUT_BIN}${CMD_SUFFIX}
endif

# Define golangci-lint
GOLANGCI_LINT=golangci-lint
GOLANGCI_LINT_VERBOSE=--verbose
GOLANGCI_LINT_ENABLE=--enable=nolintlint,bodyclose,gocritic
GOLANGCI_LINT_EXCLUDE=--exclude=SA1029
GOLANGCI_LINT_OTHER=--max-same-issues 100 --max-issues-per-linter 100

# Define nilaway
NILAWAY=nilaway

# 默认目标, 显示帮助信息
.PHONY: help
help:
	@$(CHCP_CMD)
	@echo "Available targets:"
	@echo    - help     [显示帮助信息]
	@echo    - lint     [代码静态检查]
	@echo    - nilaway  [代码 nil 检查]
	@echo    - test     [单元测试]
	@echo    - build    [构建二进制包]
	@echo    - clean    [清理构建文件]

.PHONY: all
all: lint-all build

.PHONY: lint-all
lint-all: lint nilaway

.IGNORE: lint
.PHONY: lint
lint:
	@$(CHCP_CMD)
	-${GOLANGCI_LINT} run ${GOLANGCI_LINT_VERBOSE} ${GOLANGCI_LINT_ENABLE} ${GOLANGCI_LINT_EXCLUDE} ${GOLANGCI_LINT_OTHER}

.IGNORE: nilaway
.PHONY: nilaway
nilaway:
	@${CHCP_CMD}
	${NILAWAY} ./...

.PHONY: test
test:
	@${CHCP_CMD}
	go test -v -cover -race ./...

.PHONY: pre-ci-build
pre-ci-build:
	@$(CHCP_CMD)
	export GOPROXY=https://maven.jtexpress.com.cn/nexus3/repository/go-proxy GOFLAGS="-buildvcs=false"

.PHONY: build
build:
	@$(CHCP_CMD)
	${CMD_PREFIX} cd ${BUILD_DIR} ${CMD_AND} ${BUILD_ARGS} ${CMD_AND} go build -o ${OUTPUT_BIN}${CMD_SUFFIX} .

.PHONY: clean
clean:
	@$(CHCP_CMD)
	${CMD_PREFIX} ${CLEAN_CMD}
