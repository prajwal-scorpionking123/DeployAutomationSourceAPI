package main

import (
	"github.com/team_six/SOURCE_API/controllers/deploycontroller"

	"github.com/gin-gonic/gin"
	"github.com/team_six/SOURCE_API/controllers"
	"github.com/team_six/SOURCE_API/controllers/authcontroller"
	"github.com/team_six/SOURCE_API/controllers/emailcontroller"
)

func main() {
	router := gin.Default()
	router.Static("../SOURCE", "../SOURCE")
	router.POST("/api/postlink", controllers.PostLink)
	router.GET("/api/getSources", controllers.GetSources)
	router.POST("/api/deployFiles", controllers.DeployFiles)
	router.POST("/api/login", authcontroller.Auth)
	router.POST("/api/otpmail", emailcontroller.OtpMail)
	router.PUT("/api/isverified/:id", controllers.ToggleVarified)
	router.PUT("/api/isrequested/:id", controllers.ToggleRequested)
	router.PUT("/api/isapproved/:id", controllers.ToggleApproved)

	//test route
	router.POST("/api/deployMultiple", deploycontroller.DeployFiles)

	router.Run(":3001")
}
