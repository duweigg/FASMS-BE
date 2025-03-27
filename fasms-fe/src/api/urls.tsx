const domain = "http://13.228.252.37:8000"
export const GetApplicants = domain + "/api/applicants/"
export const DeleteApplicants = (id: string) => {
    return domain + "/api/applicants/" + id
}
export const AddApplicants = domain + "/api/applicants/"


export const GetSchemes = domain + "/api/schemes/"
export const GetEligibleSchemes = (id: string) => {
    return domain + "/api/schemes/eligible?applicant=" + id
}
export const DeleteScheme = (id: string) => {
    return domain + "/api/schemes/" + id
}
export const AddSchemes = domain + "/api/schemes/"



export const GetApplications = domain + "/api/applications/"
export const UpdateApplications = (id: string) => {
    return domain + "/api/applications/" + id
}
export const DeleteApplications = (id: string) => {
    return domain + "/api/applications/" + id
}
export const SubmitApplications = domain + "/api/applications/"
