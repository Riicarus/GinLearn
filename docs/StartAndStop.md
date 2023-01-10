# 启动与关闭服务

可以使用如下方法启动一个 web server:  
```golang
router.Run()
```
或者:  
```golang
http.ListenAndServe(":8080", router)
```
完整代码:  
```golang
func main() {
  router := gin.Default()
  router.GET("/ping", func(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
      "message": "pong",
    })
  })

  // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
  // also: http.ListenAndServe(":8080", router)
  router.Run()
}
```