package i18n

import "net/http"

var HTTPStatusZHCN = StatusDict{
	http.StatusOK:                  "成功",
	http.StatusBadRequest:          "请求参数错误",
	http.StatusUnauthorized:        "请先登录",
	http.StatusForbidden:           "没有访问权限",
	http.StatusNotFound:            "资源未找到",
	http.StatusInternalServerError: "内部错误",
}
