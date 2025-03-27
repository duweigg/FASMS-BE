
export type SchemeResponse = {
    schemes: Scheme[];
    total: number;
}
export type SchemePayload = {
    schemes: Scheme[];
}

export type Scheme = {
    id: string | undefined;
    name: string;
    criteria_groups: CriteriaGroups[];
    benefits: Benefits[];
}
export type CriteriaGroups = {
    id: string | undefined;
    criterias: Criterias[];
}
export type Criterias = {
    id: string | undefined;
    employment_status: number;
    marital_status: number;
    sex: number;
    age_upper_limit: number;
    age_lower_limit: number;
    relation: number;
    is_household: boolean;
}
export type Benefits = {
    id: string | undefined;
    name: string;
    amount: number;
}


export const default_scheme: Scheme = {
    id: undefined,
    name: "",
    criteria_groups: [],
    benefits: [],
}
export const default_criteria: Criterias = {
    id: undefined,
    employment_status: 1,
    marital_status: 1,
    sex: 1,
    age_upper_limit: 999,
    age_lower_limit: 0,
    relation: 1,
    is_household: false,
}

export const default_criteria_group: CriteriaGroups = {
    id: undefined,
    criterias: [default_criteria],
}
export const default_benefits: Benefits = {
    id: undefined,
    name: "",
    amount: 0,
}