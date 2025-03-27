package main

import (
	"FASMS/controllers"
	"FASMS/initializers"

	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.GetEnvs()
	initializers.ConnectDB()
}

func main() {
	router := gin.Default()

	// Allow CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://13.228.252.37"}, // Change to frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	ApplicantController := controllers.NewApplicantController(initializers.DB)
	ApplicationController := controllers.NewApplicationController(initializers.DB)
	SchemeController := controllers.NewSchemeController(initializers.DB)
	apiRouter := router.Group("/api")
	{
		applicantRouter := apiRouter.Group("/applicants")
		{
			applicantRouter.GET("/", ApplicantController.GetApplicantsList)
			applicantRouter.POST("/", ApplicantController.CreateApplicants)

			applicantRouter.PUT("/:id", ApplicantController.UpdateApplicant)
			applicantRouter.DELETE("/:id", ApplicantController.DeleteApplicant)
		}

		schemesRouter := apiRouter.Group("/schemes")
		{
			schemesRouter.GET("/", SchemeController.GetSchemesList)
			schemesRouter.GET("/eligible", SchemeController.GetEligibleSchemesList) // ?applicant={id}

			schemesRouter.POST("/", SchemeController.AddSchemes)
			schemesRouter.PUT("/:id", SchemeController.UpdateScheme)
			schemesRouter.DELETE("/:id", SchemeController.DeleteScheme)
		}

		applicationRouter := apiRouter.Group("/applications")

		{
			applicationRouter.GET("/", ApplicationController.GetApplicationList)
			applicationRouter.POST("/", ApplicationController.CreateApplication)

			applicationRouter.PUT("/:id", ApplicationController.UpdateApplication)
			applicationRouter.DELETE("/:id", ApplicationController.DeleteApplication)
		}
	}
	router.Run()
}
