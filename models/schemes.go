package models

import (
	"FASMS/utils"
)

type Schemes struct {
	ID             string          `json:"id" gorm:"primaryKey"`
	Name           string          `json:"name"`
	CriteriaGroups []CriteriaGroup `json:"criteria_groups" gorm:"foreignKey:SchemeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Benefits       []Benefits      `json:"benifits" gorm:"foreignKey:SchemeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CommonTime
}
type CriteriaGroup struct {
	ID        string      `json:"id" gorm:"primaryKey"`
	SchemeID  string      `json:"scheme_id" gorm:"index;not null"`
	Scheme    Schemes     `json:"-" gorm:"foreignKey:SchemeID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Criterias []Criterias `json:"criterias" gorm:"foreignKey:CriteriaGroupID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CommonTime
}
type Criterias struct {
	ID               string        `json:"id" gorm:"primaryKey"`
	EmploymentStatus uint          `json:"employment_status" gorm:"comment:'0: unemployed, 1: employed, 2: in school, 99: no limitation'"`
	Sex              uint          `json:"sex" gorm:"comment:'0: male, 1: female, 99:no limitation"`
	AgeUpperLimit    uint32        `json:"age_upper_limit" gorm:"default:999"`
	AgeLowerLimit    uint32        `json:"age_lower_limit" gorm:"default:0"`
	Relation         uint          `json:"relation" gorm:"comment:'0: children, 1: spouse, 2: parents, 99: no limitation'"`
	IsHouseHold      bool          `json:"is_household"`
	CriteriaGroupID  string        `json:"criteria_group_id" gorm:"index;not null"`
	CriteriaGroup    CriteriaGroup `json:"-" gorm:"foreignKey:CriteriaGroupID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CommonTime
}
type Benefits struct {
	ID       string  `json:"id" gorm:"primaryKey"`
	Name     string  `json:"name"`
	Amount   float32 `json:"amount"`
	SchemeID string  `json:"scheme_id" gorm:"index;not null"`
	Scheme   Schemes `json:"-" gorm:"foreignKey:SchemeID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CommonTime
}

type GetSchemesRequest struct {
	PaginationQuery
}

type GetEligibleSchemesRequest struct {
	ApplicantID string `form:"applicant"  binding:"required"`
	PaginationQuery
}
type CreateSchemesListRequest struct {
	Schemes []CreateSchemesRequest `json:"schemes"`
}
type CreateSchemesRequest struct {
	Name           string                        `json:"name" binding:"required"`
	CriteriaGroups []CreateCriteriaGroupsRequest `json:"criteria_groups"`
	Benefits       []CreateBenefitRequest        `json:"benefits"`
}
type CreateCriteriaGroupsRequest struct {
	ID        string                  `json:"id"`
	Criterias []CreateCriteriaRequest `json:"criterias"`
}
type CreateCriteriaRequest struct {
	ID               string `json:"id"`
	EmploymentStatus uint   `json:"employment_status"`
	Sex              uint   `json:"sex"`
	AgeUpperLimit    uint32 `json:"age_upper_limit"`
	AgeLowerLimit    uint32 `json:"age_lower_limit"`
	Relation         uint   `json:"relation"`
	IsHouseHold      bool   `json:"is_household" binding:"required"`
}
type CreateBenefitRequest struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"  binding:"required"`
	Amount float32 `json:"amount"  binding:"required"`
}

type SchemesResponse struct {
	ID                     string                   `json:"id"`
	Name                   string                   `json:"name"`
	CriteriaGroupsResponse []CriteriaGroupsResponse `json:"criteria_groups"`
	BenefitsResponse       []BenefitsResponse       `json:"benifits"`
}
type CriteriaGroupsResponse struct {
	ID                string              `json:"id"`
	CriteriasResponse []CriteriasResponse `json:"criterias"`
}
type CriteriasResponse struct {
	ID               string `json:"id"`
	EmploymentStatus uint   `json:"employment_status"`
	Sex              uint   `json:"sex"`
	AgeUpperLimit    uint32 `json:"age_upper_limit"`
	AgeLowerLimit    uint32 `json:"age_lower_limit"`
	Relation         uint   `json:"relation"`
	IsHouseHold      bool   `json:"is_household"`
}
type BenefitsResponse struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	Amount float32 `json:"amount"`
}

func (s *Schemes) ConvertToResponse() SchemesResponse {
	SchemesResponse := SchemesResponse{
		ID:   s.ID,
		Name: s.Name,
	}

	// Convert CriteriaGroups and their Criterias
	for _, group := range s.CriteriaGroups {
		groupResponse := CriteriaGroupsResponse{
			ID: group.ID,
		}
		for _, criteria := range group.Criterias {
			groupResponse.CriteriasResponse = append(groupResponse.CriteriasResponse, CriteriasResponse{
				ID:               criteria.ID,
				EmploymentStatus: criteria.EmploymentStatus,
				Sex:              criteria.Sex,
				AgeUpperLimit:    criteria.AgeUpperLimit,
				AgeLowerLimit:    criteria.AgeLowerLimit,
				Relation:         criteria.Relation,
				IsHouseHold:      criteria.IsHouseHold,
			})
		}
		SchemesResponse.CriteriaGroupsResponse = append(SchemesResponse.CriteriaGroupsResponse, groupResponse)
	}

	// Convert Benefits
	for _, benefit := range s.Benefits {
		SchemesResponse.BenefitsResponse = append(SchemesResponse.BenefitsResponse, BenefitsResponse{
			ID:     benefit.ID,
			Name:   benefit.Name,
			Amount: benefit.Amount,
		})
	}

	return SchemesResponse
}

func (s *CreateSchemesRequest) ConvertToModel() Schemes {
	schemeId := utils.GenerateUUID()

	// Convert CriteriaGroups and their Criterias
	criteriaGroups := make([]CriteriaGroup, 0, len(s.CriteriaGroups))
	for _, group := range s.CriteriaGroups {
		groupId := utils.GenerateUUID()

		criterias := make([]Criterias, 0, len(group.Criterias))
		for _, c := range group.Criterias {
			criterias = append(criterias, Criterias{
				ID:               utils.GenerateUUID(),
				EmploymentStatus: c.EmploymentStatus,
				Sex:              c.Sex,
				AgeUpperLimit:    c.AgeUpperLimit,
				AgeLowerLimit:    c.AgeLowerLimit,
				IsHouseHold:      c.IsHouseHold,
				CriteriaGroupID:  groupId,
			})
		}

		criteriaGroups = append(criteriaGroups, CriteriaGroup{
			ID:        groupId,
			SchemeID:  schemeId,
			Criterias: criterias,
		})
	}

	// Convert Benefits
	benefits := make([]Benefits, 0, len(s.Benefits))
	for _, b := range s.Benefits {
		benefits = append(benefits, Benefits{
			ID:       utils.GenerateUUID(),
			Name:     b.Name,
			Amount:   b.Amount,
			SchemeID: schemeId,
		})
	}

	return Schemes{
		ID:             schemeId,
		Name:           s.Name,
		CriteriaGroups: criteriaGroups,
		Benefits:       benefits,
	}
}

func ConvertCriterias(newCriterias []CreateCriteriaRequest, groupID string) []Criterias {
	var convertedCriterias []Criterias

	for _, c := range newCriterias {
		criteriaID := c.ID
		if criteriaID == "" {
			criteriaID = utils.GenerateUUID() // Generate new ID for new criteria
		}

		convertedCriterias = append(convertedCriterias, Criterias{
			ID:               criteriaID,
			EmploymentStatus: c.EmploymentStatus,
			Sex:              c.Sex,
			AgeUpperLimit:    c.AgeUpperLimit,
			AgeLowerLimit:    c.AgeLowerLimit,
			Relation:         c.Relation,
			IsHouseHold:      c.IsHouseHold,
			CriteriaGroupID:  groupID, // Assign to the given group
		})
	}

	return convertedCriterias
}
