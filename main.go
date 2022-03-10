package main

import (
	"fmt"
	"os"

	"ex_gin_pb/controller"

	"github.com/gin-gonic/gin"
)

var (
	port = os.Getenv("PORT")
)

func main() {
	router := gin.Default()
	router.POST("/get_user", controller.GetUser)
	_ = router.Run(fmt.Sprintf(":%s", port))
}
