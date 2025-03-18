package controllers

import (
	"FASMS/models"
	"FASMS/utils"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Define a struct to hold the database instance
type ApplicantController struct {
	DB *gorm.DB
}

// Constructor function to create a new BookController
func NewApplicantController(db *gorm.DB) *ApplicantController {
	return &ApplicantController{DB: db}
}

func (ac *ApplicantController) GetApplicantsList(c *gin.Context) {
	var applicants []models.Applicants
	var applicantsRequest models.GetApplicantsRequest

	// Bind JSON and return 422 Unprocessable Entity on failure
	if err := c.ShouldBindQuery(&applicantsRequest); err != nil {
		log.Printf("Invalid request payload: %v\n", err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid request payload"})
		return
	}
	// assign default value to pageSize
	if applicantsRequest.PageSize == 0 {
		applicantsRequest.PageSize = 10
	}

	// Fetch applicants and return 500 Internal Server Error on failure
	var query = ac.DB.Offset(applicantsRequest.Page * applicantsRequest.PageSize).Limit(applicantsRequest.PageSize)
	if err := query.Preload("Households").Find(&applicants).Error; err != nil {
		log.Printf("Database error fetching applicants list: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch applicants list"})
		return
	}

	if len(applicants) == 0 {
		log.Printf("No applicants found")
		c.JSON(http.StatusNotFound, gin.H{"error": "No applicants found"})
		return
	}

	// find total number of applicants for pagenation
	var total int64
	if err := ac.DB.Model(&models.Applicants{}).Count(&total).Error; err != nil {
		log.Printf("Database error counting total applicants: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch total applicants"})
		return
	}

	var ret []models.ApplicantsResponse
	for _, applicant := range applicants {
		ret = append(ret, applicant.ConvertToResponse())
	}

	c.JSON(http.StatusOK, gin.H{"applicants": ret, "total": total})
}

func (ac *ApplicantController) CreateApplicants(c *gin.Context) {
	var req models.CreateApplicantsRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var applicants []models.Applicants

	for _, appReq := range req.Applicants {
		applicant := models.Applicants{
			ID:               utils.GenerateUUID(),
			Name:             appReq.Name,
			EmploymentStatus: appReq.EmploymentStatus,
			Sex:              appReq.Sex,
			DOB:              appReq.DOB.ToTime(),
		}
		var households []models.Households

		for _, household := range appReq.Households {
			households = append(households, models.Households{
				ID:               utils.GenerateUUID(),
				Name:             household.Name,
				EmploymentStatus: household.EmploymentStatus,
				Sex:              household.Sex,
				DOB:              household.DOB.ToTime(),
				Relation:         household.Relation,
				ApplicantID:      applicant.ID,
			})
		}
		applicant.Households = households
		applicants = append(applicants, applicant)
	}

	tx := ac.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback() // Ensure rollback in case of panic
		}
	}()
	if err := tx.Create(&applicants).Error; err != nil {
		tx.Rollback()
		log.Printf("create applicants failed: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create applicant"})
		return
	}
	if err := tx.Commit().Error; err != nil {
		log.Printf("Transaction commit failed: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction commit failed"})
		return
	}
	c.JSON(http.StatusCreated, applicants)
}

func (ac *ApplicantController) UpdateApplicant(c *gin.Context) {
	applicantID := c.Param("id")
	var applicant models.Applicants
	var updatedApplicant models.CreateApplicants

	if err := c.ShouldBindJSON(&updatedApplicant); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx := ac.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback() // Ensure rollback in case of panic
		}
	}()
	// update applicant
	if err := tx.Preload("Households").Where("id = ?", applicantID).First(&applicant).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			log.Printf("applicant with id: %s did not found, %v\n", applicantID, err)
			c.JSON(http.StatusNotFound, gin.H{"error": "Applicant not found"})
			return
		}
		tx.Rollback()
		log.Printf("Database error fetching applicants list: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch applicants list"})
		return
	}
	applicant.Name = updatedApplicant.Name
	applicant.EmploymentStatus = updatedApplicant.EmploymentStatus
	applicant.Sex = updatedApplicant.Sex
	applicant.DOB = updatedApplicant.DOB.ToTime()

	if err := tx.Save(applicant).Error; err != nil {
		tx.Rollback()
		log.Printf("update applicant error: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete applicant"})
		return
	}

	existingHouseholds := make(map[string]models.Households)
	for _, household := range applicant.Households {
		existingHouseholds[household.ID] = household
	}
	newHouseholds, err := updateHouseHolds(tx, updatedApplicant.Households, existingHouseholds, applicantID)
	if err != nil {
		tx.Rollback()
		log.Printf("update household failed: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "update household failed"})
		return
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		log.Printf("Transaction commit failed: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction commit failed"})
		return
	}

	applicant.Households = newHouseholds

	c.JSON(http.StatusOK, applicant.ConvertToResponse())
}

func (ac *ApplicantController) DeleteApplicant(c *gin.Context) {
	applicantID := c.Param("id")

	// check if applicant exist
	var applicant models.Applicants
	if err := ac.DB.Where("id = ?", applicantID).First(&applicant).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("applicant with id: %s did not found, %v\n", applicantID, err)
			c.JSON(http.StatusNotFound, gin.H{"error": "Applicant not found"})
			return
		}
		log.Printf("Database error fetching applicants: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch applicants list"})
		return
	}

	tx := ac.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback() // Ensure rollback in case of panic
		}
	}()

	// delete household
	if err := tx.Where("applicant_id = ?", applicantID).Delete(&models.Households{}).Error; err != nil {
		tx.Rollback()
		log.Printf("Database error deleting households belonging to applicant: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete households"})
		return
	}

	// delete applicant
	if err := tx.Where("id = ?", applicantID).Delete(&models.Applicants{}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			log.Printf("Database error deleting applicant id: %v, %v\n", applicantID, err)
			c.JSON(http.StatusNotFound, gin.H{"error": "Failed to delete applicant"})
			return
		}
		tx.Rollback()
		log.Printf("Database error fetching applicants list: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch applicants list"})
		return
	}

	// update applications involved
	result := tx.Model(&models.Applications{}).
		Where("applicant_id = ?", applicantID).
		Update("application_status", 3)

	// Check if any rows were affected
	if result.Error != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "while deleting applicant, Failed to delete applications"})
		return
	}

	if result.RowsAffected == 0 {
		log.Printf("while deleting applicant: %v, No applications found for the applicant\n", applicantID)
		// c.JSON(http.StatusNotFound, gin.H{"error": "No applications found for the applicant"})
		// return
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		log.Printf("Transaction commit failed: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction commit failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OK"})
}

func updateHouseHolds(tx *gorm.DB, updatedHousehold []models.CreateHouseholds, existingHouseholdsMapping map[string]models.Households, applicantID string) (newHouseholds []models.Households, err error) {

	var createHouseholds []models.Households
	for _, newHousehold := range updatedHousehold {
		householdID := newHousehold.ID
		if householdID == "" {
			householdID = utils.GenerateUUID()
		}

		// Check if benefit exists
		existingHousehold, exists := existingHouseholdsMapping[householdID]

		if exists {
			// Update benefit
			updateData := map[string]interface{}{
				"name":              newHousehold.Name,
				"employment_status": newHousehold.EmploymentStatus,
				"sex":               newHousehold.Sex,
				"dob":               newHousehold.DOB.ToTime(),
				"relation":          newHousehold.Relation,
			}
			if err := tx.Model(&existingHousehold).Updates(updateData).Error; err != nil {
				tx.Rollback()
				log.Println("Error updating benefit:", err)
				return nil, err
			}

			// keep track of the updated benefit for later return in response
			existingHousehold.Name = newHousehold.Name
			existingHousehold.EmploymentStatus = newHousehold.EmploymentStatus
			existingHousehold.Sex = newHousehold.Sex
			existingHousehold.DOB = newHousehold.DOB.ToTime()
			newHouseholds = append(newHouseholds, existingHousehold)

			delete(existingHouseholdsMapping, householdID) // Remove from map to track deletions
		} else {
			// New benefit
			createHouseholds = append(createHouseholds, models.Households{
				ID:               householdID,
				Name:             newHousehold.Name,
				EmploymentStatus: newHousehold.EmploymentStatus,
				Sex:              newHousehold.Sex,
				DOB:              newHousehold.DOB.ToTime(),
				ApplicantID:      applicantID,
				Relation:         newHousehold.Relation,
			})
		}
		if len(createHouseholds) > 0 {
			if err := tx.Save(&createHouseholds).Error; err != nil {
				tx.Rollback()
				log.Printf("Save new benefits failed: %v\n", err)
				return nil, err
			}
		}
	}

	// Delete removed benefits
	for _, household := range existingHouseholdsMapping {
		tx.Delete(&household)
	}
	return append(newHouseholds, createHouseholds...), nil
}
