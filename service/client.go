package service

import (
	common "my-homepage/common"
	dbcall "my-homepage/dbcall"
	model "my-homepage/struct"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 회원가입
func AddSignup(c *gin.Context) {
	var i model.Signup
	if err := c.ShouldBindJSON(&i); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": common.API_FAIL, "message": common.API_FAIL_MESSAGE})
		return
	}

	err := dbcall.InsertSignup(i)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": common.API_DB_FAIL, "message": common.API_DB_FAIL_MESSAGE})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": common.API_SUCCESS, "message": common.API_SUCCESS_MESSAGE})
}

// 로그인
func Login(c *gin.Context) {
	var i model.Login
	if err := c.ShouldBindJSON(&i); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": common.API_FAIL, "message": common.API_FAIL_MESSAGE})
		return
	}

	ok, err := dbcall.Login(i)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": common.API_DB_FAIL, "message": common.API_DB_FAIL_MESSAGE})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": common.API_SUCCESS, "message": common.API_SUCCESS_MESSAGE, "data": ok})
}
