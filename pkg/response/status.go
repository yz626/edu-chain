package response

// 业务错误码常量
const (
	CodeSuccess      = 0   // 成功
	CodeBadRequest   = 400 // 请求错误
	CodeUnauthorized = 401 // 未授权
	CodeForbidden    = 403 // 禁止访问
	CodeNotFound     = 404 // 资源不存在
	CodeServerError  = 500 // 服务器错误
)
