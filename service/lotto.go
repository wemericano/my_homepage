package service

import (
	common "my-homepage/common"
	dbcall "my-homepage/dbcall"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetLottoList(c *gin.Context) {

	res, err := dbcall.GetLottoList()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": common.API_DB_FAIL, "message": common.API_DB_FAIL_MESSAGE})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": common.API_SUCCESS, "message": common.API_SUCCESS_MESSAGE, "data": res})
}
