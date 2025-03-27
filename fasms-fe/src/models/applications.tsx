import { Applicant, default_applicant } from "./applicants";
import { default_scheme, Scheme } from "./schemes";

export type ApplicationsResponse = {
    applications: Application[];
    total: number;
}


export type ApplicationsPayload = {
    applicant_id: string;
    scheme_id: string;
}

export type Application = {
    id: string | undefined;
    application_status: number;
    applicant: Applicant;
    scheme: Scheme;
}


export const default_application: Application = {
    id: undefined,
    application_status: 0,
    applicant: default_applicant,
    scheme: default_scheme,
}