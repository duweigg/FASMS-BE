package models

import (
	"FASMS/utils"
	"time"
)

type Applicants struct {
	ID               string       `json:"id" gorm:"primaryKey"`
	Name             string       `json:"name"`
	EmploymentStatus uint         `json:"employment_status" gorm:"comment:'0: unemployed, 1: employed, 3: in school'"`
	Sex              uint         `json:"sex" gorm:"comment:'0: male, 1: female"`
	DOB              time.Time    `gorm:"type:date" json:"dob"`
	Households       []Households `json:"households" gorm:"foreignKey:ApplicantID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CommonTime
}

type Households struct {
	ID               string     `json:"id" gorm:"primaryKey"`
	Name             string     `json:"name"`
	EmploymentStatus uint       `json:"employment_status" gorm:"comment:'0: unemployed, 1: employed, 3: in school'"`
	Sex              uint       `json:"sex" gorm:"comment:'0: male, 1: female"`
	DOB              time.Time  `gorm:"type:date" json:"dob"`
	Relation         uint       `json:"relation" gorm:"comment:'0: children, 1: spouse, 2: parents'"`
	ApplicantID      string     `json:"applicant_id" gorm:"index;not null"`
	Applicant        Applicants `json:"-" gorm:"foreignKey:ApplicantID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CommonTime
}

type GetApplicantsRequest struct {
	PaginationQuery
}
type CreateApplicants struct {
	Name             string             `json:"name"`
	EmploymentStatus uint               `json:"employment_status"`
	Sex              uint               `json:"sex"`
	DOB              utils.Date         `json:"dob"`
	Households       []CreateHouseholds `json:"households"`
}
type CreateHouseholds struct {
	ID               string     `json:"id" gorm:"primaryKey"`
	Name             string     `json:"name"`
	EmploymentStatus uint       `json:"employment_status"`
	Sex              uint       `json:"sex"`
	DOB              utils.Date `json:"dob"`
	Relation         uint       `json:"relation"`
}
type CreateApplicantsRequest struct {
	Applicants []CreateApplicants `json:"applicants"`
}

type ApplicantsResponse struct {
	ID               string               `json:"id" gorm:"primaryKey"`
	Name             string               `json:"name"`
	EmploymentStatus uint                 `json:"employment_status" gorm:"comment:'0: unemployed, 1: employed, 3: in school'"`
	Sex              uint                 `json:"sex"`
	DOB              utils.Date           `json:"dob"`
	Households       []HouseholdsResponse `json:"households" gorm:"foreignKey:ApplicantID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
type HouseholdsResponse struct {
	ID               string     `json:"id" gorm:"primaryKey"`
	Name             string     `json:"name"`
	EmploymentStatus uint       `json:"employment_status" gorm:"comment:'0: unemployed, 1: employed, 3: in school'"`
	Sex              uint       `json:"sex" gorm:"comment:'0: male, 1: female"`
	DOB              utils.Date `gorm:"type:date" json:"dob"`
	Relation         uint       `json:"relation" gorm:"comment:'0: children, 1: spouse, 2: parents'"`
}

func (a *Applicants) ConvertToResponse() ApplicantsResponse {
	applicants := ApplicantsResponse{
		ID:               a.ID,
		Name:             a.Name,
		EmploymentStatus: a.EmploymentStatus,
		Sex:              a.Sex,
		DOB:              utils.Date(a.DOB),
	}
	for _, household := range a.Households {
		applicants.Households = append(applicants.Households, HouseholdsResponse{
			ID:               household.ID,
			Name:             household.Name,
			EmploymentStatus: household.EmploymentStatus,
			Sex:              household.Sex,
			DOB:              utils.Date(household.DOB),
			Relation:         household.Relation,
		})
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
