package weblog

import (
	"bytes"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// extend gin.ResponseWriter to get Response in context used in middlewares
// create a buffer to save a data copy of response's input
type CustomResponseWriter struct {
	gin.ResponseWriter

	// save a response's data copy to buffer
	body *bytes.Buffer
}

func (writer CustomResponseWriter) Write(b []byte) (int, error) {
	writer.body.Write(b)
	return writer.ResponseWriter.Write(b)
}

func (writer CustomResponseWriter) WriteString(s string) (int, error) {
	writer.body.WriteString(s)
	return writer.ResponseWriter.WriteString(s)
}

func WebLogMiddleWare() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// set context's writer to our CustomResponseWriter, to get the response body in context
		writer := &CustomResponseWriter{
			body:           bytes.NewBufferString(""),
			ResponseWriter: ctx.Writer,
		}
		ctx.Writer = writer

		start := time.Now()

		rid := ctx.Request.URL.Path + "-" + strconv.FormatInt(time.Now().UnixMilli(), 10)
		ctx.Set("rid", rid)

		fmt.Println("[" + rid + "]: URL: ", ctx.Request.URL)
		fmt.Println("[" + rid + "]: METHOD:	", ctx.Request.Method)
		fmt.Println("[" + rid + "]: IP:	", ctx.Request.RemoteAddr)

		ctx.Next()

		cost := time.Since(start)
		fmt.Println("[" + rid + "]: STATUS: ", ctx.Writer.Status())
		fmt.Println("[" + rid + "]: VALUE: ", writer.body.String())
		fmt.Println("[" + rid + "]: TIME: ", cost)
	}
}
