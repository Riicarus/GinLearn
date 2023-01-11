# 中间件
> Gin 的中间件类似于 SpringBoot 中的 过滤器和拦截器， 对整个响应过程进行处理.  
> 但是其中也有一部分类似 AOP 的逻辑.  
> 中间件是链式调用的, 可以在任意位置终止调用链, 或者调用下一个中间件.  
> 执行顺序类似压栈弹栈.  

## 1. 中间件的使用
中间件主要可以分为 全局中间件, 单路由中间件, 路由组中间件.  
```golang
// 全局中间件
router.User(GlobalMiddleWare())

// 单路由中间件
router.GET("/benchmark", SingleRouterMiddleware(), benchEndpoint)

// 路由组中间件
authorized := router.Group("/auth", RouterGroupMiddleware())
// 或者
authorized := router.Group("/auth")
authorized.Use(RouterGroupMiddleware()) {
    authorized.POST("/login", loginEndpoint)
}
```
> 当使用 `gin.Default()` 时, gin 会自动配置两个全局中间件: `gin.Logger()` 和 `gin.Recovery()` 分别用于日志记录和错误处理.  
> 如果使用 `gin.New()`, 则不会配置任何中间件.  


## 2. 自定义中间件
> 中间件其实就是一个返回值为 `gin.HanlderFunc` 的函数.

这里实现一个统计执行时间的中间件:  
```golang
func TimeCountMiddleware() gin.HandlerFunc {
    return func(ctx *gin.Context) {
        start := time.Now()

        // 在这里调用后续中间件来链式处理请求
        c.Next()

        cost := time.Since(start)
        fmt.Println("time cost: ", cost)
    }
}
```
定义好之后, 只需要在路由中注册中间件就行, 如: `router.Use(TimeCountMiddleWare())`

## 3. 中间件调用链
之前说过, 中间件的调用是一个链式的过程, 由 `context.Next()` 来进行链式调用. 每调用一个中间件都将其压入调用栈, 然后按顺序弹栈.  
**同一层次的中间件按照调用/注册顺序执行, 不同层次的中间件按照 全局>路由组>单个路由 的顺序执行.**  

如果我们想在某一处链式调用的中间件处终止剩余中间件的调用, 可以使用 `context.Abort()` 方法, 该方法只会阻止后续中间件和其他函数的处理, 但是不会停止已经压栈的中间件的后续处理. 并且 `Abort()` 方法不会给前端返回任何内容, 如果想要返回一些内容, 可以使用 `AbortWithStatusJSON()` 方法, 给前端返回一个 JSON 串.  
在调用过程中, 我们可以通过 `context.IsAbort()` 来获取一个代表是否被终止调用的布尔值, 来从当前的调用中返回, 但是不会影响整个调用链.  