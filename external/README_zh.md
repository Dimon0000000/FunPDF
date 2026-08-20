<div align="center">
  <img src="./FunPDF-Banner.png" width="420" alt="FunPDF 标志">
</div>

<p align="center">
  <a href="../README.md"><img alt="English README" src="https://img.shields.io/badge/English-DFE0E5"></a>
  <a href="./README_zh.md"><img alt="简体中文 README" src="https://img.shields.io/badge/简体中文-DBEDFA"></a>
</p>

<details open>
<summary>📕 <strong>目录</strong></summary>

- [💡 FunPDF 是什么？](#-funpdf-是什么)
- [🎮 快速开始](#-快速开始)
  - [环境要求](#环境要求)
  - [启动应用](#启动应用)
- [🌟 主要功能](#-主要功能)
- [🎬 使用 Docker 运行](#-使用-docker-运行)
- [🔧 配置](#-配置)
- [🔨 构建 Docker 镜像](#-构建-docker-镜像)
- [🛠️ 从源码启动开发环境](#️-从源码启动开发环境)
- [📖 文档](#-文档)
- [🙌 参与贡献](#-参与贡献)
- [开源许可](#开源许可)

</details>

<a id="what-is-funpdf"></a>
## 💡 FunPDF 是什么？

FunPDF 是一个轻量、可自托管的 PDF 阅读与标注原型。项目使用 Vue 3 构建网页界面，使用 Go 提供 API，并通过 MySQL 保存文档元数据。上传的 PDF 及其编辑状态保存在本地 `Cache` 目录中。

> [!NOTE]
> FunPDF 仍在持续开发中。翻译、服务商管理、参考文献和部分 AI 相关入口目前尚未完成。

<a id="get-started"></a>
## 🎮 快速开始

目前最快的完整开发环境启动方式，是使用 Docker 运行 MySQL 和 Go API，再在本地运行 Vite 前端。

### 环境要求

- Docker 与 Docker Compose
- Node.js 18+ 与 npm

### 启动应用

```bash
docker compose up -d --build
cd web
npm ci
npm run dev
```

打开 <http://localhost:5173>。前端会把 `/api` 请求代理到 `http://localhost:9384` 的 Go 服务。

<a id="key-features"></a>
## 🌟 主要功能

- 使用 PDF.js 打开并渲染本地 PDF
- 翻页、缩放、旋转、搜索、打印和导出 PDF
- 添加画笔、高亮、下划线、删除线和便签标注
- 将 PDF 和编辑状态保存到本地文件库
- 使用合集整理已保存的文件
- 使用 Docker Compose 自托管 API 与 MySQL

<a id="run-with-docker"></a>
## 🎬 使用 Docker 运行

启动 Go API 和 MySQL：

```bash
docker compose up -d --build
```

查看容器状态和应用日志：

```bash
docker compose ps
docker compose logs -f app
```

API 位于 <http://localhost:9384>，例如 `GET /api/files` 可以列出已保存的文件。

目前 Compose 配置只构建**后端**。如需使用浏览器界面，请进入 `web` 目录运行 `npm run dev`。MySQL 数据和上传的 PDF 分别持久化到 `mysql_data`、`app_cache` 命名卷。

停止服务但保留数据：

```bash
docker compose down
```

如需同时删除数据库和 PDF 缓存卷，可运行 `docker compose down -v`。该命令会永久删除本地持久化的 FunPDF 数据。

<a id="configurations"></a>
## 🔧 配置

| 环境变量 | 默认值 | 说明 |
| --- | --- | --- |
| `FUNPDF_ADDR` | `:9384` | Go HTTP 服务监听地址 |
| `FUNPDF_MYSQL_DSN` | 本地 root 连接 | 后端使用的 MySQL DSN |

Docker Compose 中的 `root` / `password` 仅适合本地开发。将服务暴露到不可信网络前，请修改凭据，并确保 `MYSQL_ROOT_PASSWORD` 与 `FUNPDF_MYSQL_DSN` 同步更新。

<a id="build-a-docker-image"></a>
## 🔨 构建 Docker 镜像

在仓库根目录单独构建 FunPDF 后端镜像：

```bash
docker build -f docker/Dockerfile -t funpdf:local .
```

镜像监听 `9384` 端口，将文件保存到 `/app/Cache`，并需要通过 `FUNPDF_MYSQL_DSN` 连接 MySQL。对于大多数本地部署，推荐使用 `docker compose up -d --build`，它会自动创建网络、数据库和持久化卷。

<a id="launch-from-source-for-development"></a>
## 🛠️ 从源码启动开发环境

环境要求：Go 1.25+、Node.js 18+、npm 和 MySQL 8.x。

1. 启动 MySQL，也可以直接复用 Compose 中的服务：

   ```bash
   docker compose up -d mysql
   ```

2. 设置后端配置，并在仓库根目录启动 Go。

   Bash：

   ```bash
   export FUNPDF_MYSQL_DSN='root:password@tcp(127.0.0.1:3306)/funpdf?charset=utf8mb4&parseTime=True&loc=Local'
   export FUNPDF_ADDR=':9384'
   go run .
   ```

   PowerShell：

   ```powershell
   $env:FUNPDF_MYSQL_DSN = 'root:password@tcp(127.0.0.1:3306)/funpdf?charset=utf8mb4&parseTime=True&loc=Local'
   $env:FUNPDF_ADDR = ':9384'
   go run .
   ```

3. 在另一个终端启动前端：

   ```bash
   cd web
   npm ci
   npm run dev
   ```

4. 打开 <http://localhost:5173>。运行时 PDF 数据会写入后端工作目录下的 `./Cache`。

常用检查命令：

```bash
go test ./...
cd web && npm run build
```

<a id="documentation"></a>
## 📖 文档

- [English README](../README.md)
- [当前文件 API 说明](../internal/development.md)

<a id="contributing"></a>
## 🙌 参与贡献

欢迎提交 Issue 和 Pull Request。提交前请尽量保持改动聚焦，并运行 Go 测试与前端构建。

## 开源许可

FunPDF 使用 [GNU General Public License v3.0](../LICENSE) 许可证。
