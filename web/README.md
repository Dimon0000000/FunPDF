# FunPDF Frontend

FunPDF 的 Vue 3 + TypeScript 前端原型。

## 已包含

- Edge 风格白灰简洁 UI
- 可折叠左侧功能栏
- 顶部 PDF 工具栏
- 铅笔 / 荧光笔 / 橡皮擦 / 下划线 / 删除线 / 备注工具入口
- 页面、目录、批注、翻译、参考文献、AI 功能入口
- PDF.js 本地 PDF 打开与渲染
- 翻页、缩放、旋转
- 搜索 UI
- Pinia 状态管理
- `/api` 到 `localhost:8080` 的开发代理
- 翻译 API 占位封装

> 当前批注工具只有 UI 与工具状态，尚未实现实际写入 PDF 的 annotation 数据层。

## 启动

```bash
npm install
npm run dev
```

默认前端：

```text
http://localhost:5173
```

Go 后端建议：

```text
http://localhost:8080
```

## 后续建议

1. 完成 PDF.js Text Layer
2. 实现 Selection API
3. 增加 Annotation Model
4. SVG/Canvas Annotation Layer
5. Go + SQLite 保存批注
6. DeepL / Google Translation Provider
7. 桌面版迁移到 Wails
