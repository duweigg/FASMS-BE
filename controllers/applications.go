package controllers

import (
	"FASMS/models"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Define a struct to hold the database instance
type ApplicationController struct {
	DB *gorm.DB
}

// Constructor function to create a new BookController
func NewApplicationController(db *gorm.DB) *ApplicationController {
	return &ApplicationController{DB: db}
}

func (ac *ApplicationController) GetApplicationList(c *gin.Context) {
	var applications []models.Applications
	var applicationsRequest models.GetApplicationsRequest

	// Bind JSON and return 422 Unprocessable Entity on failure
	if err := c.ShouldBindQuery(&applicationsRequest); err != nil {
		log.Printf("Invalid request payload: %v\n", err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid request payload"})
		return
	}
	// assign default value to pageSize
	if applicationsRequest.PageSize == 0 {
		applicationsRequest.PageSize = 10
	}

	// Fetch applicants and return 500 Internal Server Error on failure
	var query = ac.DB.Offset(applicationsRequest.Page * applicationsRequest.PageSize).Limit(applicationsRequest.PageSize).Order("id")
	if err := query.Preload("Scheme").
		Preload("Scheme.CriteriaGroups.Criterias").
		Preload("Scheme.Benefits").
		Preload("Applicant").
		Preload("Applicant.Households").
		Find(&applications).
		Error; err != nil {
		log.Printf("Database error fetching applicants list: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch application list"})
		return
	}

	if len(applications) == 0 {
		log.Printf("No applicants found")
		c.JSON(http.StatusNotFound, gin.H{"error": "No applications found"})
		return
	}

	// find total number of applications for pagenation
	var total int64
	if err := ac.DB.Model(&models.Applications{}).Count(&total).Error; err != nil {
		log.Printf("Database error counting total aplications: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch total aplications"})
		return
	}

	var ret []models.ApplicationsResponse
	for _, application := range applications {
		ret = append(ret, application.ConvertToResponse())
	}

	c.JSON(http.StatusOK, gin.H{"applications": ret, "total": total})
}

func (ac *ApplicationController) CreateApplication(c *gin.Context) {
	var scheme models.Schemes
	var applicant models.Applicants
	var applicationsRequest models.CreateApplicationRequest

	// Bind JSON and return 422 Unprocessable Entity on failure
	if err := c.ShouldBindJSON(&applicationsRequest); err != nil {
		log.Printf("Invalid request payload: %v\n", err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid request payload"})
		return
	}

	// Check if this applicant has already applied for this scheme
	var count int64
	if err := ac.DB.
		Model(&models.Applications{}).
		Where("applicant_id = ?", applicationsRequest.ApplicantID).
		Where("scheme_id = ?", applicationsRequest.SchemeID).
		Count(&count).Error; err != nil {
		log.Printf("Database error checking application existence: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if count > 0 {
		log.Printf("Duplicate application: applicant_id=%s, scheme_id=%s\n", applicationsRequest.ApplicantID, applicationsRequest.SchemeID)
		c.JSON(http.StatusConflict, gin.H{"error": "Applicant has already applied for this scheme"})
		return
	}

	// Fetch applicants and return 500 Internal Server Error on failure
	if err := ac.DB.Preload("Households").Where("id = ?", applicationsRequest.ApplicantID).First(&applicant).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("applicant with id: %s did not found, %v\n", applicationsRequest.ApplicantID, err)
			c.JSON(http.StatusNotFound, gin.H{"error": "Applicant not found"})
			return
		}
		log.Printf("Database error fetching applicants: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch applicants list"})
		return
	}
	// Fetch scheme and return 500 Internal Server Error on failure
	if err := ac.DB.Preload("CriteriaGroups.Criterias").Preload("Benefits").Where("id = ?", applicationsRequest.SchemeID).First(&scheme).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("scheme with id: %s did not found, %v\n", applicationsRequest.SchemeID, err)
			c.JSON(http.StatusNotFound, gin.H{"error": "Scheme not found"})
			return
		}
		log.Printf("Database error fetching scheme: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch scheme list"})
		return
	}
	if models.CheckEligiblity(applicant, scheme) {
		newApplication := applicationsRequest.ConvertToModel()
		if err := ac.DB.Create(&newApplication).Error; err != nil {
			log.Printf("create applicants failed: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create applicant"})
			return
		}
		newApplication.Applicant = applicant
		newApplication.Scheme = scheme
		c.JSON(http.StatusCreated, newApplication.ConvertToResponse())
	} else {
		c.JSON(http.StatusForbidden, gin.H{"error": "Applicant is not eligible for the selected scheme"})
	}
}

func (ac *ApplicationController) UpdateApplication(c *gin.Context) {
	applicationID := c.Param("id")
	var applicationsRequest models.UpdateApplicationRequest

	// Bind JSON and return 422 Unprocessable Entity on failure
	if err := c.ShouldBindJSON(&applicationsRequest); err != nil {
		log.Printf("Invalid request payload: %v\n", err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid request payload"})
		return
	}

	// Fetch applicants and return 500 Internal Server Error on failure
	var application models.Applications
	if err := ac.DB.Preload("Scheme").
		Preload("Scheme.CriteriaGroups.Criterias").
		Preload("Scheme.Benefits").
		Preload("Applicant").
		Preload("Applicant.Households").
		Where("id = ?", applicationID).First(&application).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("application with id: %s did not found, %v\n", applicationID, err)
			c.JSON(http.StatusNotFound, gin.H{"error": "Application not found"})
			return
		}
		log.Printf("Database error fetching Application: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch Application"})
		return
	}

	// update applications involved
	result := ac.DB.Model(&models.Applications{}).
		Where("id = ?", applicationID).
		Update("application_status", applicationsRequest.ApplicationStatus)

	// Check if any rows were affected
	if result.Error != nil {
		log.Printf("Error updating application with id: %s, %v\n", applicationID, result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update application"})
		return
	}

	application.ApplicationStatus = applicationsRequest.ApplicationStatus
	c.JSON(http.StatusOK, gin.H{"application": application.ConvertToResponse()})
}

func (ac *ApplicationController) DeleteApplication(c *gin.Context) {
	applicationID := c.Param("id")

	// Fetch applicants and return 500 Internal Server Error on failure
	if err := ac.DB.Where("id = ?", applicationID).First(&models.Applications{}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("application with id: %s did not found, %v\n", applicationID, err)
			c.JSON(http.StatusNotFound, gin.H{"error": "Application not found"})
			return
		}
		log.Printf("Database error fetching Application: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch Application"})
		return
	}
	// delete application
	if err := ac.DB.Where("id = ?", applicationID).Delete(&models.Applications{}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("Database error deleting households belonging to applicant id: %v, %v\n", applicationID, err)
			c.JSON(http.StatusNotFound, gin.H{"error": "Failed to delete applicant"})
			return
		}
		log.Printf("Database error fetching applicants list: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch applicants list"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "OK"})
}
