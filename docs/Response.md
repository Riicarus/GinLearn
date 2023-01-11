# 响应

## 字符串响应
```golang
c.String(http.StatusOK, "some message")
```

## JSON/XML/YAML 响应
可以绑定自定义 `struct`, 对 `field` 使用 `tag` 标签来指定序列化的 `key` 名称, 序列化时会自动更改.  
```golang
r.GET("/moreJSON", func(c *gin.Context) {
	// 使用 struct 来定义返回值
	var msg struct {
		Name    string `json:"user" xml:"user"`
		Message string
		Number  int
	}
	msg.Name = "Lena"
	msg.Message = "hey"
	msg.Number = 123
	// 注意 msg.Name 变成了 "user" 字段
	// 以下方式都会输出 :   {"user": "Lena", "Message": "hey", "Number": 123}
	c.JSON(http.StatusOK, gin.H{"user": "Lena", "Message": "hey", "Number": 123})
	c.XML(http.StatusOK, gin.H{"user": "Lena", "Message": "hey", "Number": 123})
	c.YAML(http.StatusOK, gin.H{"user": "Lena", "Message": "hey", "Number": 123})
	c.JSON(http.StatusOK, msg)
	c.XML(http.StatusOK, msg)
	c.YAML(http.StatusOK, msg)
})
```

## 重定向
```golang
r.GET("/redirect", func(c *gin.Context) {
	// 支持内部和外部的重定向
    c.Redirect(http.StatusMovedPermanently, "http://www.baidu.com/")
})
```

## 文件响应
```golang
// 获取当前文件的相对路径
router.Static("/assets", "./assets")
// 功能类似 Statis, 但是可以使用 FileSystem
router.StaticFS("/more_static", http.Dir("my_file_system"))
// 获取相对路径下的文件
router.StaticFile("/favicon.ico", "./resources/favicon.ico")
```

## 异步
可以使用异步进行响应处理, 但是一定注意 **不要直接使用 `context`, 而是使用 `ctxCopy := context.Copy()`(只读上下文) 来进行异步处理**  
```golang
func main() {
	r := gin.Default()
	//1. 异步
	r.GET("/long_async", func(c *gin.Context) {
		// goroutine 中只能使用只读的上下文 c.Copy()
		cCp := c.Copy()
		go func() {
			time.Sleep(5 * time.Second)

			// 注意使用只读上下文
			log.Println("Done! in path " + cCp.Request.URL.Path)
		}()
	})
	//2. 同步
	r.GET("/long_sync", func(c *gin.Context) {
		time.Sleep(5 * time.S econd)

		// 注意可以使用原始上下文
		log.Println("Done! in path " + c.Request.URL.Path)
	})

	// Listen and serve on 0.0.0.0:8080
	r.Run(":8080")
}
```