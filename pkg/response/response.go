package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 通用响应结构
type Response struct {
	Code    int         `json:"code"`           // 业务状态码
	Message string      `json:"message"`        // 提示信息
	Data    interface{} `json:"data,omitempty"` // 响应数据
}

// PageData 分页数据结构
type PageData struct {
	List       interface{} `json:"list"`        // 数据列表
	Total      int64       `json:"total"`       // 总数
	Page       int         `json:"page"`        // 当前页码
	PageSize   int         `json:"page_size"`   // 每页大小
	TotalPages int         `json:"total_pages"` // 总页数
}

// ==================== 成功响应 ====================

// Success 成功响应（无数据）
func Success(c *gin.Context) {
	c.JSON(http.StatusOK, Response{Code: CodeSuccess, Message: "操作成功"})
}

// SuccessWithData 成功响应（带数据）
func SuccessWithData(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{Code: CodeSuccess, Message: "操作成功", Data: data})
}

// Created 创建成功响应
func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Response{Code: CodeSuccess, Message: "创建成功", Data: data})
}

// ==================== 分页响应 ====================

// Page 分页响应
func Page(c *gin.Context, list interface{}, total int64, page, pageSize int) {
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: "操作成功",
		Data:    PageData{List: list, Total: total, Page: page, PageSize: pageSize, TotalPages: totalPages},
	})
}

// ==================== 错误响应 ====================

// Error 错误响应（HTTP 状态码正常）
func Error(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, Response{Code: code, Message: message})
}

// ErrorWithData 错误响应（带数据）
func ErrorWithData(c *gin.Context, res *Response) {
	c.JSON(http.StatusOK, res)
}

// ErrorWithStatus 错误响应（HTTP 状态码非正常）
func ErrorWithStatus(c *gin.Context, status int, code int, message string) {
	c.JSON(status, Response{Code: code, Message: message})
}

// ErrorWithStatusAndData 错误响应（HTTP 状态码非正常，带数据）
func ErrorWithStatusAndData(c *gin.Context, status int, res *Response) {
	c.JSON(status, res)
}

// Unauthorized 未授权错误响应
func Unauthorized(c *gin.Context, message string) {
	ErrorWithStatus(c, http.StatusUnauthorized, CodeUnauthorized, message)
}

// Forbidden 禁止访问错误响应
func Forbidden(c *gin.Context, message string) {
	ErrorWithStatus(c, http.StatusForbidden, CodeForbidden, message)
}

// NotFound 资源不存在错误响应
func NotFound(c *gin.Context, message string) {
	ErrorWithStatus(c, http.StatusNotFound, CodeNotFound, message)
}
