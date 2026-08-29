#### 后端需要补 3 个 API。
1. 创建 PDF 会话
**`POST /api/providers/:provider_id/sessions`**
请求：
```json
{
  "model_id": "模型数据库ID",
  "model_name": "deepseek-chat",
  "system_prompt": "系统提示词",
  "document": {
    "name": "paper.pdf",
    "content": "完整论文文本"
  }
}
```

响应：
```json
{
  "code": 200,
  "data": {
    "id": "session-id"
  }
}
```
这个接口需要保存系统提示词、论文上下文、Provider、模型和消息历史。

2. 发送消息
**`POST /api/providers/:provider_id/sessions/:session_id/messages`**

请求：
```json
{
  "content": "用户的问题",
  "quote": "可选的论文引用",
  "stream": true
}
```
SSE 格式：

```http request
event: message
data: {"content": "增量回答", "reasoning_content": ""}

event: message
data: {"content": "", "reasoning_content": "增量思考"}

event: done
data: [DONE]
```

后端负责从 session 取出历史消息和论文上下文、调用模型、保存本轮 user/assistant 消息。这样后续请求不需要重复提交论文。
3. 删除会话
**`DELETE /api/providers/:provider_id/sessions/:session_id`**
用于 PDF 关闭后释放会话资源。

> [!TIP]
>
> session 暂时存数据库，不要只放进程内存，否则程序重启、多实例部署后会丢失。论文很长时，后端还应负责 token 限制、分块或检索；不要把整篇论文无条件塞进每一次模型请求。
> 验证结果：npm run build 已通过。只有项目原有的 Vite 大 chunk 提示，没有新增编译错误。浏览器 UI 自动检查因当前环境没有可用浏览器实例未执行。

