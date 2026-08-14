package main

import (
	"FunPDF/internal"

	"github.com/gin-gonic/gin"
)

func main() {
	r := internal.NewRouter()
	r.Setup(gin.Default())

	run()
}
