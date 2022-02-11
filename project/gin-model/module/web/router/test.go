package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zqlpaopao/tool/gin-model/common"
	"net/http"
)

type Login struct {

}

func (_ Login) Load(r *gin.Engine){
	l := r.Group("/fire-timer")
	l.POST("/login", func(g *gin.Context) {
		fmt.Println(g.GetString(common.ContextKey))
		g.JSON(http.StatusOK, gin.H{
			"message": "ok",
			"code":    200,
			"data":    map[string]string{
				"aa":"bb",
			},
		})
		return
	})

}