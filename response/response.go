package response

import "github.com/gin-gonic/gin"

const (
	SUCCESS       = 0
	ERROR         = 1000 // 通用错误
	INVALID_PARAM = 1001 // 参数错误
	UNAUTHORIZED  = 1002 // 未认证
	FORBIDDEN     = 1003 // 无权限
	NOT_FOUND     = 1004 // 资源不存在
	SERVER_ERROR  = 1005 // 服务器错误
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(200, Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

// 失败响应
func Error(c *gin.Context, code int, msg string) {
	c.JSON(200, Response{
		Code:    code,
		Message: msg,
		Data:    nil,
	})
}
