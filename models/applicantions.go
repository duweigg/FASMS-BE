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
	ApplicationStatus uint       `json:"application_status" gorm:"comment:'1: submitted, 2: approved, 3: rejected, 4: need review (due to applicant/scheme updates)'"`
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
	ApplicationStatus uint `json:"application_status" binding:"required,oneof=1 2 3 4"`
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
		ApplicationStatus: 1,
	}
	return newApplication
}

// An applicant must satisfy all criteria groups.
// A group is considered satisfied if any one of its criteria is met.
// Criteria can apply to:
// 	The applicant (e.g., age, employment, sex, marital status)
// 	Household members (with unique household IDs and matching logic)

// biz logic: applicant must satisify all the criteria groups.
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
			(criteria.Sex == 99 || applicant.Sex == criteria.Sex) &&
			(criteria.MaritalStatus == 99 || applicant.MaritalStatus == criteria.MaritalStatus) {
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
				(criteria.Relation == 99 || household.Relation == criteria.Relation) &&
				(criteria.MaritalStatus == 99 || household.MaritalStatus == criteria.MaritalStatus) {

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

// func NewCheckEligiblity(applicant Applicants, scheme Schemes) bool {
// 	if len(scheme.CriteriaGroups) == 0 {
// 		return true
// 	}

// 	for _, group := range scheme.CriteriaGroups {
// 		usedHouseholds := make(map[string]bool)
// 		if isGroupSatisfiedBacktracking(group.Criterias, applicant, usedHouseholds, 0) {
// 			return true // At least one group satisfied
// 		}
// 	}
// 	return false // None of the groups satisfied
// }

// // Backtracking function to satisfy all criteria in a group
// func isGroupSatisfiedBacktracking(criteriaList []Criterias, applicant Applicants, used map[string]bool, index int) bool {
// 	if index >= len(criteriaList) {
// 		return true // All criteria satisfied
// 	}

// 	criteria := criteriaList[index]
// 	if !criteria.IsHouseHold {
// 		if matchApplicantCriteria(applicant, criteria) {
// 			// Recurse to next criteria
// 			if isGroupSatisfiedBacktracking(criteriaList, applicant, used, index+1) {
// 				return true
// 			}
// 		}
// 		return false // Applicant doesn't match this non-household criteria
// 	}
// 	for _, household := range applicant.Households {
// 		if used[household.ID] {
// 			continue
// 		}
// 		if matchHouseholdCriteria(household, criteria) {
// 			// Choose this household
// 			used[household.ID] = true

// 			// Recurse to next criteria
// 			if isGroupSatisfiedBacktracking(criteriaList, applicant, used, index+1) {
// 				return true
// 			}

// 			// Backtrack
// 			delete(used, household.ID)
// 		}
// 	}

// 	return false // No match for this criteria
// }

// // Match logic for a household and a single criteria
// func matchHouseholdCriteria(h Households, c Criterias) bool {
// 	age := h.GetAge()
// 	return (age >= c.AgeLowerLimit && age <= c.AgeUpperLimit) &&
// 		(c.EmploymentStatus == 99 || h.EmploymentStatus == c.EmploymentStatus) &&
// 		(c.Sex == 99 || h.Sex == c.Sex) &&
// 		(c.Relation == 99 || h.Relation == c.Relation) &&
// 		(c.MaritalStatus == 99 || h.MaritalStatus == c.MaritalStatus)
// }

// // Match logic for a household and a single criteria
// func matchApplicantCriteria(a Applicants, c Criterias) bool {
// 	age := a.GetAge()
// 	return (age >= c.AgeLowerLimit && age <= c.AgeUpperLimit) &&
// 		(c.EmploymentStatus == 99 || a.EmploymentStatus == c.EmploymentStatus) &&
// 		(c.Sex == 99 || a.Sex == c.Sex) &&
// 		(c.MaritalStatus == 99 || a.MaritalStatus == c.MaritalStatus)
// }
