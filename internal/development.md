# FunPDF 当前文件 API

前端基地址为 `/api`，开发环境代理到 `http://localhost:9384`。数据库仍使用 MySQL，文件内容保存到 `./Cache/{id}`。

## `GET /api/files`

返回已经保存的文件记录。

## `POST /api/files`

第一次 Ctrl+S 时调用，使用 `multipart/form-data`：

- `file`：原始 PDF 文件。
- `editor_state`：编辑状态 JSON 字符串。

成功返回 `201 Created`，`data` 为创建后的 File 实体。

## `PATCH /api/files/:file_id/state`

后续 Ctrl+S 时调用：

```json
{
  "expected_revision": 1,
  "editor_state": {
    "format": "funpdf-editor-state",
    "version": 1,
    "saved_at": "2026-08-16T00:00:00Z",
    "source": {
      "name": "paper.pdf",
      "mime_type": "application/pdf"
    },
    "editor": {
      "annotations": {},
      "rotation": 0,
      "scale": 1.15,
      "current_page": 1
    }
  }
}
```

版本不一致返回 `409 Conflict`。

## `DELETE /api/files/:file_id`

删除数据库记录和对应的 Cache 目录，成功返回 `204 No Content`。

Album、翻译、Provider 等未完成接口当前不注册，避免 Gin 因空 handler 无法启动。
