# 路由

Gin 的路由使用的是 httprouter.  
  
**创建路由**:  
```golang
// 创建带有默认中间件的路由:
// 日志与恢复中间件
router := gin.Default()

//创建不带中间件的路由：
r := gin.New()
```
**设置路由**:  
```golang
router.GET("/someGet", getting)
router.POST("/somePost", posting)
router.PUT("/somePut", putting)
router.DELETE("/someDelete", deleting)
router.PATCH("/somePatch", patching)
router.HEAD("/someHead", head)
router.OPTIONS("/someOptions", options)

// 或者使用匿名函数
router.GET("/string/:name", func(c *gin.Context) {
    name := c.Param("name")
    fmt.Println("Hello %s", name)
})
```
**路由参数获取**:  
通过 `context` 的 `Param()` 来获取 restful api 中的参数:  
```golang
router.GET("/string/:name", func(c *gin.Context) {
    name := c.Param("name")
})
```
通过 `DefaultQuery()` 或 `Query()` 获取 URL 中的参数:  
```golang
// context.Query() == context.Request.URL.Query().Get()
router.GET("/welcome", func(ctx *gin.Context) {
	// set default value if there's no value find
	name := ctx.DefaultQuery("name", "Guest")
  	lastname := ctx.Query("lastname")
  	fmt.Println("name: ", name, ", lastname: ", lastname)
})
```
通过 `PostForm()` 获取表单中的参数:  
```golang
router.POST("/form", func(c *gin.Context) {
	// set default value if there's no value find
	type := c.DefaultPostForm("type", "alert")
	msg := c.PostForm("msg")
	title := c.PostForm("title")
	fmt.Println("type is %s, msg is %s, title is %s", type, msg, title)
})
```
路由群组:  
```golang
someGroup := router.Group("/someGroup")
{
	// 这里的 url: "/someGroup/someGet"
    someGroup.GET("/someGet", getting)
	someGroup.POST("/somePost", posting)
}
```