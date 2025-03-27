'use client';

import { ApiCallFunction } from "@/api/callFunction";
import { GetApplicants } from "@/api/urls";
import { AddApplicantComponent, ApplicantComponent } from "@/components/applicant";
import { BackButton } from "@/components/button";
import { ApplicantsResponse, Applicant } from "@/models/applicants";
import { useEffect, useState } from "react";

export default function Home() {
    const [reload, setReload] = useState(true)
    // get all applicants data
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
        if (reload) {
            fetchApplicants()
            setReload(false)
        }
    }, [reload])
    return (
        <>
            <BackButton text="back" />
            <div className="font-bold text-5xl text-center">Applicants Info</div>
            <div className="flex flex-row flex-wrap">
                <AddApplicantComponent setReload={() => { setReload(true) }} />
                {applicants.map((applicant) => {
                    return (
                        <div key={applicant.id} className="p-4">
                            <ApplicantComponent applicant={applicant} setReload={() => { setReload(true) }} />
                        </div>
                    )
                })}
            </div>
        </>
    );
}
