package util

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetParamInt(c *gin.Context, param string) (int, error) {
	idStr := c.Param(param)
	if idStr == "" {
		return -1, fmt.Errorf("empty param")
	}
	id, err := strconv.Atoi(idStr) // convert ascii to int
	if err != nil {
		return -1, fmt.Errorf("conversion of %s to int not possible", param)
	}
	return id, nil
}
