package main

import (
	sdkadapter "ibfu/robot-docking-gitee/gitee-adapter"
	"net/http"

	"github.com/go-playground/validator/v10"

	"github.com/gin-gonic/gin"
)

//var Controller *gin.RouterGroup

func (bot *robot) registerRoutePath() {
	bot.handlerAddComment()
}

var orgValid validator.Func = func(fl validator.FieldLevel) bool {
	val, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}

	l := len(val)
	if l == 0 || l > 16 {
		return false
	}

	return true
}

type SigReqArgs struct {
	Org  string `form:"org"`
	Repo string `form:"repo"`
}

func (bot *robot) handlerAddComment() {
	bot.ctl.GET("sig/name", func(context *gin.Context) {
		var arg SigReqArgs
		// 如果是 `GET` 请求，只使用 `Form` 绑定引擎（`query`）
		sdkadapter.GetClientInstance([]byte("12313"))
		if err := context.ShouldBind(&arg); err != nil {
			context.String(http.StatusBadRequest, "err, "+err.Error())
		} else {
			context.String(http.StatusOK, "1111")
		}
	})
}
