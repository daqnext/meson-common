package common

//var ApiRouterMap = make(map[string]map[string]*gin.RouterGroup)
//var GinRouter = func() *gin.Engine {
//	v := gin.Default()
//	config := cors.DefaultConfig()
//	config.AllowAllOrigins = true
//	config.AllowCredentials = true
//	config.AddAllowHeaders("Authorization")
//	v.Use(cors.New(config))
//	return v
//}()
//
//func GetMyRouter() *gin.RouterGroup {
//	_, filename, _, ok := runtime.Caller(1)
//	if !ok {
//		return nil
//	}
//
//	dir := path.Dir(filename)
//	filename = filepath.Base(utils.FileNameWithoutExtension(filename))
//	return ApiRouterMap[filepath.Base(dir)][filename]
//}
//
//func AutoConfigRouter() {
//	_, filename, _, ok := runtime.Caller(1)
//	if !ok {
//		return
//	}
//	dir := path.Dir(filename)
//	filename = filepath.Base(utils.FileNameWithoutExtension(filename))
//
//	if ApiRouterMap[filepath.Base(dir)] == nil {
//		ApiRouterMap[filepath.Base(dir)] = make(map[string]*gin.RouterGroup)
//	}
//
//	index := strings.Index(dir, "routerpath")
//	subPath := dir[index+10:]
//
//	ApiRouterMap[filepath.Base(dir)][filename] = GinRouter.Group(subPath + "/").Group(filename)
//}
