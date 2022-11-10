package handlers

// Содержит отбработчики для gin

import (
	"durationCount/internal/core"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

/*
REST API:GET Функция возвращает кратчайшее время для работ

~/duration/<string>
*/
func GetBrowserOptDuration(c *gin.Context) {
	Order_name := c.Param("Order")
	type returnstruct struct {
		Duration float64
		Path     []string
	}
	//lint:ignore SA4006 (выражение используется далее)
	path := []string{}
	i, path := core.GetOptDuration(Order_name, 10, 100000)
	ret := returnstruct{Duration: i, Path: path}
	if i == -1 {
		c.IndentedJSON(http.StatusBadRequest, returnstruct{i, path})
	} else {
		c.IndentedJSON(http.StatusOK, ret)
	}
}
