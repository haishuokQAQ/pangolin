package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"pangolin/app/pangolin/model/db"
	"pangolin/app/pangolin/model/rest"
	"pangolin/app/pangolin/utils"
	"strconv"
)

func (ctl *Controller) ListAllTunnel(c *gin.Context) {
	param := &rest.PageParam{}
	if err := c.ShouldBindJSON(param); err != nil {
		_ = c.Error(err)
		return
	}
	ctx := utils.CreateContextFromGinContext(c)
	tunnelConfigs, totalCount, err := ctl.Srv.ListAllConfig(ctx, param.PageNum, param.PageSize)
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
	config, err := ctl.Srv.CreateConfig(ctx, config)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, rest.BasicResponse{
		Meta: &rest.ResponseMeta{},
		Data: struct {
			Config *db.TunnelConfig `json:"config"`
		}{
			Config: config,
		},
	})
}

func (ctl *Controller) CreateTunnel(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		_ = c.Error(errors.New("Invalid path param 'id'."))
		return
	}
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		_ = c.Error(err)
		return
	}
	ctx := utils.CreateContextFromGinContext(c)
	err = ctl.Srv.CreateTunnelByConfigId(ctx, id)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, rest.BasicResponse{})
}

func (ctl *Controller) GetTunnelStatistic(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		_ = c.Error(errors.New("Invalid path param 'id'."))
		return
	}
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		_ = c.Error(err)
		return
	}
	ctx := utils.CreateContextFromGinContext(c)
	statisticMap, err := ctl.Srv.GetStatistic(ctx, id)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, rest.BasicResponse{
		Meta: &rest.ResponseMeta{},
		Data: struct {
			Statistics map[string]*utils.FlowStatistic `json:"statistics"`
		}{
			Statistics: statisticMap,
		},
	})
}
