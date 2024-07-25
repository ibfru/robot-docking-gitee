package main

import (
	"community-robot-lib/config"
	"community-robot-lib/framework"
	"fmt"
	sdkadapter "ibfu/robot-docking-gitee/gitee-adapter"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type robot struct {
	cli *sdkadapter.ClientTarget
	ctl *gin.RouterGroup
}

func newRobot(cli *sdkadapter.ClientTarget) *robot {
	return &robot{cli: cli}
}

func (bot *robot) NewConfig() config.Config {
	return &configuration{}
}

func (bot *robot) RegisterEventHandler(f framework.HandlerRegister) {
	// custom handle request
}

type header struct {
	Token string `form:"token"`
}

const (
	AuthErrorMessage = "authorized failed"
)

func authRequired(c *gin.Context) {
	var h header

	// error happened, gin framework handler
	_ = c.BindHeader(&h)

	if h.Token == "111" {
		return
	}

	c.AbortWithStatusJSON(http.StatusUnauthorized, AuthErrorMessage)
}

func (bot *robot) setupRouter() *gin.Engine {
	// Disable Console Color
	if gin.Mode() == gin.ReleaseMode {
		gin.DisableConsoleColor()
	}

	fmt.Println("----------------------------------")

	router := gin.New()
	// LoggerWithFormatter 中间件会写入日志到 gin.DefaultWriter
	// 默认 gin.DefaultWriter = os.Stdout
	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// 你的自定义格式
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.DateTime),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))
	router.Use(gin.Recovery())

	// 认证路由组
	bot.ctl = router.Group("/", authRequired)
	bot.registerRoutePath()
	return router
}

func (bot *robot) getFile() []*sdkadapter.ContentInfo {
	a, _ := bot.cli.GetRepoContentsByPath("ibforuorg", "community-fork-to-test", "ci-scripts/check_branch.py")
	return a
}
