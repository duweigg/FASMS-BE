package models

import (
	"FASMS/utils"
	"time"
)

type Applicants struct {
	ID               string       `json:"id" gorm:"primaryKey"`
	Name             string       `json:"name"`
	IC               string       `json:"ic" gorm:"unique,not null"`
	EmploymentStatus uint         `json:"employment_status" gorm:"comment:'1: unemployed, 2: employed, 3: in school'"`
	Sex              uint         `json:"sex" gorm:"comment:'1: male, 2: female"`
	DOB              time.Time    `gorm:"type:date" json:"dob"`
	Households       []Households `json:"households" gorm:"foreignKey:ApplicantID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CommonTime
}

type Households struct {
	ID               string     `json:"id" gorm:"primaryKey"`
	Name             string     `json:"name"`
	IC               string     `json:"ic"`
	EmploymentStatus uint       `json:"employment_status" gorm:"comment:'1: unemployed, 2: employed, 3: in school'"`
	Sex              uint       `json:"sex" gorm:"comment:'1: male, 2: female"`
	DOB              time.Time  `gorm:"type:date" json:"dob"`
	Relation         uint       `json:"relation" gorm:"comment:'1: children, 2: spouse, 3: parents'"`
	ApplicantID      string     `json:"applicant_id" gorm:"index;not null"`
	Applicant        Applicants `json:"-" gorm:"foreignKey:ApplicantID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CommonTime
}

type GetApplicantsRequest struct {
	PaginationQuery
}
type CreateApplicants struct {
	Name             string             `json:"name"  binding:"required"`
	IC               string             `json:"ic"  binding:"required"`
	EmploymentStatus uint               `json:"employment_status" binding:"required,oneof=1 2 3"`
	Sex              uint               `json:"sex" binding:"required,oneof=1 2"`
	DOB              utils.Date         `json:"dob" binding:"required"`
	Households       []CreateHouseholds `json:"households"`
}
type CreateHouseholds struct {
	ID               string     `json:"id" gorm:"primaryKey"`
	Name             string     `json:"name" binding:"required"`
	IC               string     `json:"ic"  binding:"required"`
	EmploymentStatus uint       `json:"employment_status" binding:"required,oneof=1 2 3"`
	Sex              uint       `json:"sex" binding:"required,oneof=1 2"`
	DOB              utils.Date `json:"dob" binding:"required"`
	Relation         uint       `json:"relation" binding:"required,oneof=1 2 3"`
}
type CreateApplicantsRequest struct {
	Applicants []CreateApplicants `json:"applicants" binding:"required,dive"`
}

type ApplicantsResponse struct {
	ID               string               `json:"id"`
	Name             string               `json:"name"`
	IC               string               `json:"ic"`
	EmploymentStatus uint                 `json:"employment_status"`
	Sex              uint                 `json:"sex"`
	DOB              utils.Date           `json:"dob"`
	Households       []HouseholdsResponse `json:"households"`
}
type HouseholdsResponse struct {
	ID               string     `json:"id"`
	Name             string     `json:"name"`
	IC               string     `json:"ic"`
	EmploymentStatus uint       `json:"employment_status"`
	Sex              uint       `json:"sex"`
	DOB              utils.Date `json:"dob"`
	Relation         uint       `json:"relation"`
}

func (a *Applicants) ConvertToResponse() ApplicantsResponse {
	applicants := ApplicantsResponse{
		ID:               a.ID,
		Name:             a.Name,
		IC:               a.IC,
		EmploymentStatus: a.EmploymentStatus,
		Sex:              a.Sex,
		DOB:              utils.Date(a.DOB),
	}
	for _, household := range a.Households {
		applicants.Households = append(applicants.Households, HouseholdsResponse{
			ID:               household.ID,
			Name:             household.Name,
			IC:               household.IC,
			EmploymentStatus: household.EmploymentStatus,
			Sex:              household.Sex,
			DOB:              utils.Date(household.DOB),
			Relation:         household.Relation,
		})
	}
	return applicants
}

func (a *CreateApplicantsRequest) ConvertToModel() []Applicants {

	var applicants []Applicants
	for _, appReq := range a.Applicants {
		applicant := Applicants{
			ID:               utils.GenerateUUID(),
			IC:               appReq.IC,
			Name:             appReq.Name,
			EmploymentStatus: appReq.EmploymentStatus,
			Sex:              appReq.Sex,
			DOB:              appReq.DOB.ToTime(),
		}
		var households []Households

		for _, household := range appReq.Households {
			households = append(households, Households{
				ID:               utils.GenerateUUID(),
				Name:             household.Name,
				IC:               household.IC,
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
	return applicants
}

func (a *Applicants) GetAge() uint32 {
	today := time.Now()
	age := today.Year() - a.DOB.Year()
	// Adjust if birthday hasn't occurred yet this year
	if today.YearDay() < a.DOB.YearDay() {
		age--
	}
	return uint32(age)
}
func (h *Households) GetAge() uint32 {
	today := time.Now()
	age := today.Year() - h.DOB.Year()
	// Adjust if birthday hasn't occurred yet this year
	if today.YearDay() < h.DOB.YearDay() {
		age--
	}
	return uint32(age)
}
