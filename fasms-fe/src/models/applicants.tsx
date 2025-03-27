export type ApplicantsResponse = {
    applicants: Applicant[];
    total: number;
}

export type ApplicantsPayload = {
    applicants: Applicant[];
}

export type Applicant = {
    id: string | undefined;
    name: string;
    marital_status: number;
    ic: string;
    employment_status: number;
    sex: number;
    dob: string;
    households: Household[];
}

export type Household = {
    id: string | undefined;
    name: string;
    marital_status: number;
    ic: string;
    employment_status: number;
    sex: number;
    dob: string;
    relation: number;
}

export const default_applicant: Applicant = {
    id: undefined,
    name: "",
    marital_status: 1,
    ic: "",
    employment_status: 1,
    sex: 1,
    dob: "",
    households: [],
}

export const default_household: Household = {
    id: undefined,
    name: "",
    marital_status: 1,
    ic: "",
    employment_status: 1,
    sex: 1,
    dob: "",
    relation: 1,
}