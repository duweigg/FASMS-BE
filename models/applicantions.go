package models

import (
	"FASMS/utils"
	"log"
)

type Applications struct {
	ID                string     `json:"id" gorm:"primaryKey"`
	ApplicantID       string     `json:"applicant_id" gorm:"index;not null"`
	Applicant         Applicants `json:"-" gorm:"foreignKey:ApplicantID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	SchemeID          string     `json:"scheme_id" gorm:"index;not null"`
	Scheme            Schemes    `json:"-" gorm:"foreignKey:SchemeID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ApplicationStatus uint       `json:"application_status" gorm:"comment:'0: submitted, 1: approved, 2: rejected, 3: need review (due to applicant/scheme updates)'"`
	CommonTime
}

type GetApplicationsRequest struct {
	PaginationQuery
}

type CreateApplicationRequest struct {
	ApplicantID string `json:"applicant_id" binding:"required"`
	SchemeID    string `json:"scheme_id" binding:"required"`
}
type UpdateApplicationRequest struct {
	ApplicationStatus uint `json:"application_status" binding:"oneof=0 1 2 3"`
}
type ApplicationsResponse struct {
	ID                string             `json:"id"`
	Applicant         ApplicantsResponse `json:"applicant"`
	Scheme            SchemesResponse    `json:"scheme"`
	ApplicationStatus uint               `json:"application_status"`
}

func (ar *Applications) ConvertToResponse() ApplicationsResponse {
	return ApplicationsResponse{
		ID:                ar.ID,
		Applicant:         ar.Applicant.ConvertToResponse(),
		Scheme:            ar.Scheme.ConvertToResponse(),
		ApplicationStatus: ar.ApplicationStatus,
	}
}

func (car *CreateApplicationRequest) ConvertToModel() Applications {
	var newApplication = Applications{
		ID:                utils.GenerateUUID(),
		ApplicantID:       car.ApplicantID,
		SchemeID:          car.SchemeID,
		ApplicationStatus: 0,
	}
	return newApplication
}

// the biz logic is that applicant must satisify all the criteria groups.
// a criteria group is considered as satisified if any of the criteria in the groupo is satisified
func CheckEligiblity(applicant Applicants, scheme Schemes) bool {
	if len(scheme.CriteriaGroups) == 0 {
		return true
	}

	var householdCriteriaGroups []CriteriaGroup
	var applicantCriteriaGroup CriteriaGroup
	// check each group
	for _, criteriaGroup := range scheme.CriteriaGroups {
		var isHouseholdCriteriaGroup = true
		// within the group, check if criteria is satisified
		for _, criteria := range criteriaGroup.Criterias {
			if !criteria.IsHouseHold {
				isHouseholdCriteriaGroup = false
			}
		}
		// log.Println(criteriaGroup.ID)
		// log.Println(isHouseholdCriteriaGroup)
		if isHouseholdCriteriaGroup {
			householdCriteriaGroups = append(householdCriteriaGroups, criteriaGroup)
		} else {
			applicantCriteriaGroup = criteriaGroup
		}
	}
	// log.Println(applicantCriteriaGroup)
	// log.Println(householdCriteriaGroups)
	var applicantEligible = IsApplicantEligible(applicant, applicantCriteriaGroup)
	var usedHouseholds = make(map[string]bool)
	var houseHoldEligible = IsHouseholdEligible(householdCriteriaGroups, applicant.Households, usedHouseholds, 0)
	// log.Println(applicantEligible)
	// log.Println(houseHoldEligible)

	return applicantEligible && houseHoldEligible
}

func IsApplicantEligible(applicant Applicants, criteriaGroup CriteriaGroup) bool {
	// check age
	age := applicant.GetAge()
	for _, criteria := range criteriaGroup.Criterias {
		// log.Printf("%v > %v > %v", criteria.AgeLowerLimit, age, criteria.AgeUpperLimit)
		// log.Printf("eployment: %v == %v", criteria.EmploymentStatus, applicant.EmploymentStatus)
		// log.Printf("sex: %v == %v", criteria.Sex, applicant.Sex)
		if (age >= criteria.AgeLowerLimit && age <= criteria.AgeUpperLimit) &&
			(criteria.EmploymentStatus == 99 || applicant.EmploymentStatus == criteria.EmploymentStatus) &&
			(criteria.Sex == 99 || applicant.Sex == criteria.Sex) {
			return true
		}
	}
	return false
}

func IsHouseholdEligible(criteriaGroups []CriteriaGroup, households []Households, usedHouseholds map[string]bool, index int) bool {
	// log.Println(len(criteriaGroups))
	// log.Println(index)
	if index >= len(criteriaGroups) {
		return true
	}
	for _, household := range households {
		if usedHouseholds[household.ID] {
			continue // Skip already used households
		}

		age := household.GetAge()
		log.Println(age)
		for _, criteria := range criteriaGroups[index].Criterias {

			if age >= criteria.AgeLowerLimit && age <= criteria.AgeUpperLimit &&
				(criteria.EmploymentStatus == 99 || household.EmploymentStatus == criteria.EmploymentStatus) &&
				(criteria.Sex == 99 || household.Sex == criteria.Sex) &&
				(criteria.Relation == 99 || household.Relation == criteria.Relation) {

				// Mark household as used
				usedHouseholds[household.ID] = true

				// Recur to check next criteriaGroup
				if IsHouseholdEligible(criteriaGroups, households, usedHouseholds, index+1) {
					return true
				}

				// Backtrack: unmark this household
				delete(usedHouseholds, household.ID)
			}
		}
	}

	return false // No valid assignment found
}
