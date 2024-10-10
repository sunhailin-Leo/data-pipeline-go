<h1 align="center">data-pipeline-go</h1>

---

## 项目介绍
* 基于 Golang 实现一个类似 SeaTunnel 的数据同步工具, 主要是为了简便易用
* **文档在 docs 目录下** 【To be continued!】
* **示例在 example 目录下** 【To be continued!】

## 静态检查

* Windows 下需要安装 make 命令
  * https://gnuwin32.sourceforge.net/packages/make.htm
  * 安装完后加环境变量即可

* golangci-lint
  * 安装: `curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.59.1` 或 `brew install golangci-lint`
  * 检查: `golangci-lint --version`

* nilaway
  * 安装: `go install go.uber.org/nilaway/cmd/nilaway@latest`
  * 检查: `nilaway ./...`

## 实现模块

[ROADMAP](ROADMAP.md)

## 版本日志

[CHANGELOG](CHANGELOG.md)
