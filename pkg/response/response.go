package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    data,
	})
}

func Error(c *gin.Context, code int, msg string) {
	// 根据业务错误码设置HTTP状态码
	var httpStatus int
	switch code {
	case 400:
		httpStatus = http.StatusBadRequest
	case 401:
		httpStatus = http.StatusUnauthorized
	case 403:
		httpStatus = http.StatusForbidden
	case 404:
		httpStatus = http.StatusNotFound
	case 500:
		httpStatus = http.StatusInternalServerError
	default:
		httpStatus = http.StatusOK // 其他情况保持200
	}
	
	c.JSON(httpStatus, gin.H{
		"code":    code,
		"message": msg,
	})
}

// SuccessWithPagination 带分页的成功响应
func SuccessWithPagination(c *gin.Context, data interface{}, total int64, page interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    data,
		"total":   total,
		"page":    page,
	})
}
