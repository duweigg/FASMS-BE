'use client';

import { ApiCallFunction } from "@/api/callFunction";
import { GetApplicants, GetApplications, GetEligibleSchemes } from "@/api/urls";
import { AddApplicationComponent, ApplicationComponent } from "@/components/application";
import { BackButton } from "@/components/button";
import { Applicant, ApplicantsResponse } from "@/models/applicants";
import { Application, ApplicationsResponse } from "@/models/applications";
import { Scheme, SchemeResponse } from "@/models/schemes";
import { useEffect, useState } from "react";

export default function Home() {
  const [reload, setReload] = useState(true)
  // get all applicants data
  const [applicantions, setApplications] = useState<Application[]>([]);
  useEffect(() => {
    const fetchApplicants = async () => {
      const data = await ApiCallFunction<ApplicationsResponse>('GET', GetApplications);
      if (!data) {
        console.error('Failed to fetch applicants');
        setApplications([])
        return
      }
      setApplications(data.applications);
    };
    if (reload) {
      fetchApplicants()
      setReload(false)
    }
  }, [reload])


  const [applicants, setApplicants] = useState<Applicant[]>([]);
  useEffect(() => {
    const fetchApplicants = async () => {
      const data = await ApiCallFunction<ApplicantsResponse>('GET', GetApplicants);
      if (!data) {
        console.error('Failed to fetch applicants');
        setApplicants([])
        return
      }
      setApplicants(data.applicants);
    };
    fetchApplicants()
  }, [])


  const [selectedApplicant, setSelectedApplicant] = useState<string>("");
  const [schemes, setSchemes] = useState<Scheme[]>([]);
  useEffect(() => {
    const fetchSchemes = async () => {
      const data = await ApiCallFunction<SchemeResponse>('GET', GetEligibleSchemes(selectedApplicant));
      if (!data) {
        console.error('Failed to fetch eligibel scheme');
        setSchemes([])
        return
      }
      setSchemes(data.schemes);
    };
    if (selectedApplicant != "") {
      fetchSchemes()
    }
  }, [selectedApplicant])
  return (
    <>
      <BackButton text="back" />
      <div className="font-bold text-5xl text-center">Applications Info</div>
      <div className="flex flex-row flex-wrap">
        <AddApplicationComponent setReload={() => { setReload(true) }} applicants={applicants} setSelectedApplicant={(id: string) => { setSelectedApplicant(id) }} schemes={schemes} />
        {applicantions.map((application) => {
          return (
            <div key={application.id} className="p-4">
              <ApplicationComponent application={application} setReload={() => { setReload(true) }} />
            </div>
          )
        })}
      </div>
    </>
  );
}
