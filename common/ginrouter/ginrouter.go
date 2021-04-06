package ginrouter

import (
	"github.com/daqnext/meson-common/common/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

type GinAutoRouter struct {
	GinInstance  *gin.Engine
	ApiRouterMap map[string]map[string]*gin.RouterGroup
}

var GinInstanceMap = map[string]*GinAutoRouter{}

func New(name string) *GinAutoRouter {
	newGinAutoRouter := &GinAutoRouter{}
	newGinAutoRouter.GinInstance = gin.Default()
	newGinAutoRouter.ApiRouterMap = make(map[string]map[string]*gin.RouterGroup)
	GinInstanceMap[name] = newGinAutoRouter
	return newGinAutoRouter
}

func GetGinInstance(name string) *GinAutoRouter {
	instance, exist := GinInstanceMap[name]
	if exist {
		return instance
	}

	//if not exist, new a gin.Engine
	newGinAutoRouter := New(name)
	return newGinAutoRouter
}

func (g *GinAutoRouter) AutoConfigRouter() {
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		return
	}
	dir := path.Dir(filename)
	filename = filepath.Base(utils.FileNameWithoutExtension(filename))

	if g.ApiRouterMap[filepath.Base(dir)] == nil {
		g.ApiRouterMap[filepath.Base(dir)] = make(map[string]*gin.RouterGroup)
	}

	index := strings.Index(dir, "routerpath")
	subPath := dir[index+10:]

	g.ApiRouterMap[filepath.Base(dir)][filename] = g.GinInstance.Group(subPath + "/").Group(filename)
}

func (g *GinAutoRouter) GetMyRouter() *gin.RouterGroup {
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		return nil
	}

	dir := path.Dir(filename)
	filename = filepath.Base(utils.FileNameWithoutExtension(filename))
	return g.ApiRouterMap[filepath.Base(dir)][filename]
}

func (g *GinAutoRouter) EnableDefaultCors() {
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowCredentials = true
	corsConfig.AddAllowHeaders("Authorization")
	g.GinInstance.Use(cors.New(corsConfig))
}
