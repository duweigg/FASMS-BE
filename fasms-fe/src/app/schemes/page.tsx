'use client';

import { ApiCallFunction } from "@/api/callFunction";
import { GetSchemes } from "@/api/urls";
import { BackButton } from "@/components/button";
import { AddSchemeComponent, SchemeComponent } from "@/components/scheme";
import { Scheme, SchemeResponse } from "@/models/schemes";
import { useEffect, useState } from "react";

export default function Home() {
    const [reload, setReload] = useState(true)
    // get all applicants data
    const [schemes, setSchemes] = useState<Scheme[]>([]);
    useEffect(() => {
        const fetchSchemes = async () => {
            const data = await ApiCallFunction<SchemeResponse>('GET', GetSchemes);
            if (!data) {
                console.error('Failed to fetch applicants');
                setSchemes([])
                return
            }
            setSchemes(data.schemes);
        };
        if (reload) {
            fetchSchemes()
            setReload(false)
        }
    }, [reload])
    return (
        <>
            <BackButton text="back" />
            <div className="font-bold text-5xl text-center">Scheme Info</div>
            <div className="flex flex-row flex-wrap">
                <AddSchemeComponent setReload={() => { setReload(true) }} />
                {schemes.map((scheme) => {
                    return (
                        <div key={scheme.id} className="p-4">
                            <SchemeComponent scheme={scheme} setReload={() => { setReload(true) }} />
                        </div>
                    )
                })}
            </div>
        </>
    );
}
