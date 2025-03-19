package controllers

import (
	"FASMS/models"
	"FASMS/utils"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Define a struct to hold the database instance
type SchemeController struct {
	DB *gorm.DB
}

// Constructor function to create a new BookController
func NewSchemeController(db *gorm.DB) *SchemeController {
	return &SchemeController{DB: db}
}

func (sc *SchemeController) GetSchemesList(c *gin.Context) {
	var schemes []models.Schemes
	var schemesRequest models.GetSchemesRequest

	// Bind JSON and return 422 Unprocessable Entity on failure
	if err := c.ShouldBindQuery(&schemesRequest); err != nil {
		log.Printf("Invalid request payload: %v\n", err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid request payload"})
		return
	}

	// assign default value to pageSize
	if schemesRequest.PageSize == 0 {
		schemesRequest.PageSize = 10
	}

	// Fetch applicants and return 500 Internal Server Error on failure
	var query = sc.DB.Offset(schemesRequest.Page * schemesRequest.PageSize).Limit(schemesRequest.PageSize)
	if err := query.Preload("CriteriaGroups.Criterias").Preload("Benefits").Find(&schemes).Error; err != nil {
		log.Printf("Database error fetching scheme list: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch scheme list"})
		return
	}

	if len(schemes) == 0 {
		log.Printf("No scheme found")
		c.JSON(http.StatusNotFound, gin.H{"error": "No scheme found"})
		return
	}
	// find total number of applications for pagenation
	var total int64
	if err := sc.DB.Model(&models.Schemes{}).Count(&total).Error; err != nil {
		log.Printf("Database error counting total scheme: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch total total"})
		return
	}
	var ret []models.SchemesResponse
	for _, scheme := range schemes {
		ret = append(ret, scheme.ConvertToResponse())
	}

	c.JSON(http.StatusOK, gin.H{"schemes": ret, "total": total})
}

func (sc *SchemeController) GetEligibleSchemesList(c *gin.Context) {
	var schemes []models.Schemes
	var applicant models.Applicants
	var schemesRequest models.GetEligibleSchemesRequest

	// Bind JSON and return 422 Unprocessable Entity on failure
	if err := c.ShouldBindQuery(&schemesRequest); err != nil {
		log.Printf("Invalid request payload: %v\n", err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid request payload"})
		return
	}

	// assign default value to pageSize
	if schemesRequest.PageSize == 0 {
		schemesRequest.PageSize = 10
	}

	// Fetch applicants and return 500 Internal Server Error on failure
	if err := sc.DB.Preload("Households").Where("id = ?", schemesRequest.ApplicantID).First(&applicant).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("applicant with id: %s did not found, %v\n", schemesRequest.ApplicantID, err)
			c.JSON(http.StatusNotFound, gin.H{"error": "Applicant not found"})
			return
		}
		log.Printf("Database error fetching applicants: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch applicants list"})
		return
	}

	// Fetch schemes and return 500 Internal Server Error on failure
	if err := sc.DB.Preload("CriteriaGroups.Criterias").Preload("Benefits").Find(&schemes).Error; err != nil {
		log.Printf("Database error fetching scheem list: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch scheem list"})
		return
	}

	if len(schemes) == 0 {
		log.Printf("No scheme found")
		c.JSON(http.StatusNotFound, gin.H{"error": "No scheme found"})
		return
	}

	var ret []models.SchemesResponse
	for _, scheme := range schemes {
		if models.CheckEligiblity(applicant, scheme) {
			ret = append(ret, scheme.ConvertToResponse())
		}
	}
	// pagination
	start := schemesRequest.Page * schemesRequest.PageSize
	end := start + schemesRequest.PageSize

	// Ensure bounds are valid
	if start >= len(ret) {
		start = len(ret) // Avoid out-of-bounds error
	}
	if end > len(ret) {
		end = len(ret)
	}
	c.JSON(http.StatusOK, gin.H{"schemes": ret[start:end], "total": len(ret)})
}

func (sc *SchemeController) AddSchemes(c *gin.Context) {
	var addSchemesRequest models.CreateSchemesListRequest

	// Bind JSON and return 422 Unprocessable Entity on failure
	if err := c.ShouldBindJSON(&addSchemesRequest); err != nil {
		log.Printf("Invalid request payload: %v\n", err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid request payload"})
		return
	}

	var newSchemeName []string
	for _, newscheme := range addSchemesRequest.Schemes {
		// validate each new scheme added
		isvalidScheem, err := newscheme.IsValidScheme()
		if !isvalidScheem {
			log.Printf("New scheme is not valid: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("New scheme is not valid, %v", err.Error())})
			return
		}
		newSchemeName = append(newSchemeName, newscheme.Name)
	}

	// check if the scheme with same name exists
	var existingScheme models.Schemes
	if err := sc.DB.Where("name in (?)", newSchemeName).First(&existingScheme).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": fmt.Sprintf("scheme with the same name: %v already exists", existingScheme.Name)})
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	var schemes []models.Schemes

	for _, scheme := range addSchemesRequest.Schemes {
		schemes = append(schemes, scheme.ConvertToModel())
	}

	tx := sc.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback() // Ensure rollback in case of panic
		}
	}()
	if err := tx.Create(&schemes).Error; err != nil {
		tx.Rollback()
		log.Printf("create schemes failed: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create schemes"})
		return
	}
	if err := tx.Commit().Error; err != nil {
		log.Printf("Transaction commit failed: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction commit failed"})
		return
	}

	var schemesResponse []models.SchemesResponse
	for _, scheme := range schemes {
		schemesResponse = append(schemesResponse, scheme.ConvertToResponse())
	}
	c.JSON(http.StatusCreated, schemesResponse)
}

func (sc *SchemeController) DeleteScheme(c *gin.Context) {
	schemeID := c.Param("id")

	// check if applicant exist
	var scheme models.Schemes
	if err := sc.DB.Where("id = ?", schemeID).First(&scheme).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("scheme with id: %s did not found, %v\n", schemeID, err)
			c.JSON(http.StatusNotFound, gin.H{"error": "Scheme not found"})
			return
		}
		log.Printf("Database error fetching scheme: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch scheme list"})
		return
	}

	tx := sc.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback() // Ensure rollback in case of panic
		}
	}()

	var criteriaGroupIds []string
	if err := tx.Model(&models.CriteriaGroup{}).Select("id").Where("scheme_id = ?", schemeID).Scan(&criteriaGroupIds).Error; err != nil {
		tx.Rollback()
		log.Printf("Database error deleting criterias belonging to scheme: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete criterias"})
		return
	}
	// delete criterias
	if err := tx.Where("criteria_group_id in (?)", criteriaGroupIds).Delete(&models.Criterias{}).Error; err != nil {
		tx.Rollback()
		log.Printf("Database error deleting criterias belonging to scheme: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete criterias"})
		return
	}

	// delete criterias group
	if err := tx.Where("scheme_id = ?", schemeID).Delete(&models.CriteriaGroup{}).Error; err != nil {
		tx.Rollback()
		log.Printf("Database error deleting criterias group belonging to scheme: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete criteria group"})
		return
	}

	// delete benefits
	if err := tx.Where("scheme_id = ?", schemeID).Delete(&models.Benefits{}).Error; err != nil {
		tx.Rollback()
		log.Printf("Database error deleting benefits belonging to scheme: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete benefits"})
		return
	}

	// delete scheme
	if err := tx.Where("id = ?", schemeID).Delete(&models.Schemes{}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			log.Printf("Database error deleting scheme id: %v, %v\n", schemeID, err)
			c.JSON(http.StatusNotFound, gin.H{"error": "Failed to delete applicant"})
			return
		}
		tx.Rollback()
		log.Printf("Database error fetching scheme list: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch scheme list"})
		return
	}

	// update applications involved
	result := tx.Model(&models.Applications{}).
		Where("scheme_id = ?", schemeID).
		Update("application_status", 3)

	// Check if any rows were affected
	if result.Error != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "while deleting scheme, Failed to delete applications"})
		return
	}

	if result.RowsAffected == 0 {
		log.Printf("while deleting sceme: %v, No applications found for the scheme\n", schemeID)
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

func (sc *SchemeController) UpdateScheme(c *gin.Context) {
	schemeID := c.Param("id")
	var existingScheme models.Schemes
	var updatedScheme models.CreateSchemesRequest

	if err := c.ShouldBindJSON(&updatedScheme); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	isvalidScheem, err := updatedScheme.IsValidScheme()
	if !isvalidScheem {
		log.Printf("New scheme is not valid: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("New scheme is not valid, %v", err.Error())})
		return
	}

	tx := sc.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback() // Ensure rollback in case of panic
		}
	}()
	// update scheme
	if err := tx.Preload("CriteriaGroups.Criterias").Preload("Benefits").Where("id = ?", schemeID).First(&existingScheme).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			log.Printf("scheme with id: %s did not found, %v\n", schemeID, err)
			c.JSON(http.StatusNotFound, gin.H{"error": "Scheme not found"})
			return
		}
		tx.Rollback()
		log.Printf("Database error fetching scheme list: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch scheme"})
		return
	}

	existingScheme.Name = updatedScheme.Name
	if err := tx.Save(existingScheme).Error; err != nil {
		tx.Rollback()
		log.Printf("update scheme error: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update scheme"})
		return
	}

	existingGroups := make(map[string]models.CriteriaGroup)
	for _, group := range existingScheme.CriteriaGroups {
		existingGroups[group.ID] = group
	}

	var newGroups []models.CriteriaGroup
	var createGroups []models.CriteriaGroup

	for _, newGroup := range updatedScheme.CriteriaGroups {
		groupID := newGroup.ID

		// If it's a new group (not in DB), generate ID
		if groupID == "" {
			groupID = utils.GenerateUUID()
		}

		// Check if group exists in DB
		existingGroup, exists := existingGroups[groupID]
		if exists {
			// Update the existing group
			updatedCriterias, err := UpdateCriterias(tx, existingGroup.Criterias, newGroup.Criterias, groupID)
			existingGroup.Criterias = updatedCriterias
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating criterias"})
				return
			}
			newGroups = append(newGroups, existingGroup)
			delete(existingGroups, groupID) // Remove from map to track deletions
		} else {
			// New group to be inserted
			createGroups = append(createGroups, models.CriteriaGroup{
				ID:        groupID,
				SchemeID:  schemeID,
				Criterias: models.ConvertCriterias(newGroup.Criterias, groupID),
			})
		}
		if len(createGroups) > 0 {
			if err := tx.Save(&createGroups).Error; err != nil {
				tx.Rollback()
				log.Printf("Save new benefits failed: %v\n", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Save new benefits failed"})
				return
			}
		}
	}
	for _, group := range existingGroups {
		tx.Delete(&group)
	}

	existingBenefits := make(map[string]models.Benefits)
	for _, benefit := range existingScheme.Benefits {
		existingBenefits[benefit.ID] = benefit
	}

	var newBenefits []models.Benefits
	var createBenefits []models.Benefits
	for _, newBenefit := range updatedScheme.Benefits {
		benefitID := newBenefit.ID
		if benefitID == "" {
			benefitID = utils.GenerateUUID()
		}

		// Check if benefit exists
		existingBenefit, exists := existingBenefits[benefitID]

		if exists {
			// Update benefit
			updateData := map[string]interface{}{
				"name":   newBenefit.Name,
				"amount": newBenefit.Amount,
			}
			if err := tx.Model(&existingBenefit).Updates(updateData).Error; err != nil {
				tx.Rollback()
				log.Println("Error updating benefit:", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating benefit"})
				return
			}

			// keep track of the updated benefit for later return in response
			existingBenefit.Name = newBenefit.Name
			existingBenefit.Amount = newBenefit.Amount
			newBenefits = append(newBenefits, existingBenefit)

			delete(existingBenefits, benefitID) // Remove from map to track deletions
		} else {
			// New benefit
			createBenefits = append(createBenefits, models.Benefits{
				ID:       benefitID,
				Name:     newBenefit.Name,
				Amount:   newBenefit.Amount,
				SchemeID: schemeID,
			})
		}
		if len(createBenefits) > 0 {
			if err := tx.Save(&createBenefits).Error; err != nil {
				tx.Rollback()
				log.Printf("Save new benefits failed: %v\n", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Save new benefits failed"})
				return
			}
		}
	}

	// Delete removed benefits
	for _, benefit := range existingBenefits {
		tx.Delete(&benefit)
	}
	existingScheme.CriteriaGroups = append(newGroups, createGroups...)
	existingScheme.Benefits = append(newBenefits, createBenefits...)

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		log.Printf("Transaction commit failed: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction commit failed"})
		return
	}
	c.JSON(http.StatusOK, existingScheme.ConvertToResponse())
}

func UpdateCriterias(tx *gorm.DB, existingCriterias []models.Criterias, newCriterias []models.CreateCriteriaRequest, groupID string) ([]models.Criterias, error) {
	existingCriteriasMap := make(map[string]models.Criterias)
	for _, criteria := range existingCriterias {
		existingCriteriasMap[criteria.ID] = criteria
	}

	var updatedCriterias []models.Criterias
	var createCriterias []models.Criterias

	for _, newCriteria := range newCriterias {
		criteriaID := newCriteria.ID
		if criteriaID == "" {
			criteriaID = utils.GenerateUUID()
		}

		existingCriteria, exists := existingCriteriasMap[criteriaID]

		if exists {
			// Update existing criteria

			updateData := map[string]interface{}{
				"employment_status": newCriteria.EmploymentStatus,
				"sex":               newCriteria.Sex,
				"age_upper_limit":   newCriteria.AgeUpperLimit,
				"age_lower_limit":   newCriteria.AgeLowerLimit,
				"relation":          newCriteria.Relation,
				"is_house_hold":     newCriteria.IsHouseHold,
			}

			if err := tx.Model(&existingCriteria).Updates(updateData).Error; err != nil {
				tx.Rollback()
				log.Println("Error updating criteria:", err)
				return nil, err
			}

			existingCriteria.EmploymentStatus = newCriteria.EmploymentStatus
			existingCriteria.Sex = newCriteria.Sex
			existingCriteria.AgeUpperLimit = newCriteria.AgeUpperLimit
			existingCriteria.AgeLowerLimit = newCriteria.AgeLowerLimit
			existingCriteria.Relation = newCriteria.Relation
			existingCriteria.IsHouseHold = *newCriteria.IsHouseHold
			updatedCriterias = append(updatedCriterias, existingCriteria)
			delete(existingCriteriasMap, criteriaID) // Mark as processed
		} else {
			// New criteria
			createCriterias = append(createCriterias, models.Criterias{
				ID:               criteriaID,
				EmploymentStatus: newCriteria.EmploymentStatus,
				Sex:              newCriteria.Sex,
				AgeUpperLimit:    newCriteria.AgeUpperLimit,
				AgeLowerLimit:    newCriteria.AgeLowerLimit,
				Relation:         newCriteria.Relation,
				IsHouseHold:      *newCriteria.IsHouseHold,
				CriteriaGroupID:  groupID,
			})
		}
		if len(createCriterias) > 0 {
			if err := tx.Save(&createCriterias).Error; err != nil {
				tx.Rollback()
				log.Printf("Save new benefits failed: %v\n", err)
				return nil, err
			}
		}
	}

	// Delete criteria that are not in the new request
	for _, criteria := range existingCriteriasMap {
		tx.Delete(&criteria)
	}

	return append(updatedCriterias, createCriterias...), nil
}
