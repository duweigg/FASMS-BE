package utils

var EmploymentStatusMapping = map[int]string{
	0:  "unemployed",
	1:  "employed",
	2:  "in school",
	99: "no employment limitation",
}

var GenderMapping = map[int]string{
	0:  "male",
	1:  "female",
	99: "no gender limitation",
}

var RelationMapping = map[int]string{
	0:  "children",
	1:  "spouse",
	2:  "parents",
	99: "no relation limitation",
}
