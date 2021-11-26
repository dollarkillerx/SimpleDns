package server

import (
	"github.com/dollarkillerx/SimpleDns/pkg/model"
	"github.com/gin-gonic/gin"

	"log"
	"net/http"
)

func (s *SimpleDns) api() {
	s.app.Use(Cors())

	s.app.GET("/auth", func(ctx *gin.Context) {
		token := ctx.GetHeader("token")
		if token != s.conf.Token {
			ctx.String(401, "401 Token Error")
		}
		ctx.String(200, "sucess")
	})

	s.app.GET("/ip", func(ctx *gin.Context) {
		ctx.JSON(200, ctx.ClientIP())
	})

	api := s.app.Group("/api", s.authMiddleware)
	{
		api.GET("/list", s.listDns)
		api.POST("/del", s.delDns)
		api.POST("/update", s.updateDns)
		api.POST("/add", s.addDns)
	}
}

func (s *SimpleDns) authMiddleware(ctx *gin.Context) {
	token := ctx.GetHeader("token")
	if token != s.conf.Token {
		ctx.String(401, "401 Token Error")
		ctx.Abort()
	}
}

func (s *SimpleDns) listDns(ctx *gin.Context) {
	dns, err := s.storage.APIListDns()
	if err != nil {
		ctx.String(500, err.Error())
		return
	}

	ctx.JSON(200, dns)
}

type delDnsRequest struct {
	ID string `json:"id"`
}

func (s *SimpleDns) delDns(ctx *gin.Context) {
	var req delDnsRequest
	err := ctx.BindJSON(&req)
	if err != nil {
		ctx.String(400, err.Error())
		return
	}

	if req.ID == "" {
		ctx.String(400, "req.ID == \"\"")
		return
	}

	err = s.storage.APIDeleteDns(req.ID)
	if err != nil {
		ctx.String(500, err.Error())
		return
	}

	ctx.String(200, "success")
}

func (s *SimpleDns) updateDns(ctx *gin.Context) {
	var req model.DnsModel
	err := ctx.BindJSON(&req)
	if err != nil {
		ctx.String(400, err.Error())
		return
	}

	if req.ID == "" {
		ctx.String(400, "req id is null")
		return
	}

	err = s.storage.APIUpdateDns(req.ID, &req)
	if err != nil {
		ctx.String(500, err.Error())
		return
	}

	ctx.JSON(200, "success")
}

func (s *SimpleDns) addDns(ctx *gin.Context) {
	var req model.DnsModel
	err := ctx.BindJSON(&req)
	if err != nil {
		ctx.String(400, err.Error())
		return
	}

	err = s.storage.APIStorageDns(req.Domain, &req)
	if err != nil {
		ctx.String(500, err.Error())
		return
	}

	ctx.JSON(200, "success")
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin") //请求头部
		if origin != "" {
			//接收客户端发送的origin （重要！）
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			//服务器支持的所有跨域请求的方法
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
			//允许跨域设置可以返回其他子段，可以自定义字段
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session")
			// 允许浏览器（客户端）可以解析的头部 （重要）
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers")
			//设置缓存时间
			c.Header("Access-Control-Max-Age", "172800")
			//允许客户端传递校验信息比如 cookie (重要)
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		//允许类型校验
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "ok!")
		}

		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic info is: %v", err)
			}
		}()

		c.Next()
	}
}
