# 控制器
## 数据解析绑定
请求体 --绑定--> 数据类型(JSON, XML, FormData)  
- 使用 `tag` 设置绑定的类型的标签, 如: `json: "fieldname"`, `form: "fieldname"`  
- 使用 `binding: "required"` 标签来设置一个字段是必选的, 如果为空值, 请求会失败并返回错误.  

使用 `context.Bind()` 方法会根据请求头中的 `Content-Type` 自动推断类型. 或者使用 `context.BindWith()` 或者 `context.MustBindWith()` 来指定要绑定的类型.   

```golang
// Bind to JSON
type Login struct {
	User     string `form:"user" json:"user" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

// 绑定为 JSON
router.POST("/loginJSON", func(ctx *gin.Context) {
	var json Login
	if ctx.BindJSON(&json) == nil {
		if json.User == "manu" && json.Password == "123" {
			ctx.JSON(http.StatusOK, gin.H{"status": "you are logged in"})
		} else {
			ctx.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
		}
	}
})

// 绑定为 FORM_DATA
router.POST("/loginForm", func(ctx *gin.Context) {
	var form Login
	if ctx.Bind(&form) == nil {
		if form.User == "manu" && form.Password == "123" {
			ctx.JSON(http.StatusOK, gin.H{"status": "you are logged in"})
		} else {
			ctx.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
		}
	}
})
```

## 文件上传
```golang
router.POST("/upload", func(ctx *gin.Context) {
	file, header, _ := ctx.Request.FormFile("upload")

	filename := header.Filename
	fmt.Println("filename: ", filename)

	// the dir must exsit
	out, err := os.Create("./tmp/" + filename)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

    // write file to dest dir
	_, err = io.Copy(out, file)
	if err != nil {
		log.Fatal(err)
	}
})
```