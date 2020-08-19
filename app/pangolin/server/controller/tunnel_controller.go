package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pangolin/app/pangolin/model/db"
	"pangolin/app/pangolin/model/rest"
	"pangolin/app/pangolin/utils"
)

func (ctl *Controller) ListAllTunnel(c *gin.Context) {
	param := &rest.PageParam{}
	if err := c.ShouldBindJSON(param); err != nil {
		_ = c.Error(err)
		return
	}
	ctx := utils.CreateContextFromGinContext(c)
	tunnelConfigs, totalCount, err := ctl.srv.ListAllConfig(ctx, param.PageNum, param.PageSize)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, &rest.BasicResponse{
		Meta: &rest.ResponseMeta{},
		Data: &rest.ListTunnelConfigResponseData{
			Rows:       tunnelConfigs,
			TotalCount: totalCount,
		},
	})
}

func (ctl *Controller) CreateTunnelConfig(c *gin.Context) {
	config := &db.TunnelConfig{}
	if err := c.ShouldBindJSON(config); err != nil {
		_ = c.Error(err)
		return
	}
	ctx := utils.CreateContextFromGinContext(c)
	err := ctl.srv.CreateTunnel(ctx, config)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, nil)
}
