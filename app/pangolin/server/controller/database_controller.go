package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pangolin/app/pangolin/config"
	"pangolin/app/pangolin/model/rest"
)

func (ctl *Controller) GetCurrentConfig(c *gin.Context) {
	c.JSON(http.StatusOK, &rest.BasicResponse{
		Meta: &rest.ResponseMeta{},
		Data: ctl.srv.GetCurrentDBConfig(),
	})

}

func (ctl *Controller) ConnectDB(c *gin.Context) {
	conf := &config.DBConfig{}
	if err := c.ShouldBindJSON(conf); err != nil {
		_ = c.Error(err)
		return
	}
	err := ctl.srv.ConnectDatabase(conf)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, rest.BasicResponse{})
}
