package server

import (
	"github.com/dollarkillerx/SimpleDns/pkg/model"
	"github.com/gin-gonic/gin"
)

func (s *SimpleDns) api() {
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
