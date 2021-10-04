package ktools

import "github.com/gin-gonic/gin"

type KRoutGroup struct {
	*gin.RouterGroup
}

func (kr *KRoutGroup) RouteName(name string) (absoluteUrl string) {

	absoluteUrl = kr.BasePath()

	return
}
