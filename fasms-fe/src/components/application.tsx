import { Applicant } from "@/models/applicants"
import { Button, CardActionArea, CardActions, CardContent, MenuItem, Select, SelectChangeEvent, Typography } from "@mui/material";
import * as React from 'react';
import { styled } from '@mui/material/styles';
import Card from '@mui/material/Card';
import IconButton, { IconButtonProps } from '@mui/material/IconButton';
import AddIcon from '@mui/icons-material/Add';
import { useState } from "react";
import { DeleteApplications, SubmitApplications, UpdateApplications } from "@/api/urls";
import { ApiCallFunction } from "@/api/callFunction";
import { Application, ApplicationsPayload } from "@/models/applications";
import { ApplicationStatusMapping } from "@/utils/const";
import { Scheme } from "@/models/schemes";


interface ApplicationProps {
    application: Application;
    setReload: () => void;
}

interface ExpandMoreProps extends IconButtonProps {
    expand: boolean;
}

export const ExpandMore = styled((props: ExpandMoreProps) => {
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    const { expand, ...other } = props;
    return <IconButton {...other} />;
})<ExpandMoreProps>(({ theme, expand }) => ({
    marginLeft: 'auto',
    transform: expand ? 'rotate(180deg)' : 'rotate(0deg)',
    transition: theme.transitions.create('transform', {
        duration: theme.transitions.duration.shortest,
    }),
}));

export const ApplicationComponent = ({ application, setReload }: ApplicationProps) => {
    const handleDeleteClick = async (id: string) => {
        await ApiCallFunction('DELETE', DeleteApplications(id));
        setReload()
    };
    const handleApproveClick = async (id: string) => {
        if (id != "") {
            await ApiCallFunction('PUT', UpdateApplications(id), { application_status: 2 });
            setReload()
        }
    };
    const handleRejectClick = async (id: string) => {
        if (id != "") {
            await ApiCallFunction('PUT', UpdateApplications(id), { application_status: 3 });
            setReload()
        }
    };
    return (
        <Card sx={{ width: 250 }}>
            <CardContent>
                <Typography gutterBottom variant="h5" component="div">
                    Applicant
                </Typography>
                <Typography variant="body2" sx={{ color: 'text.secondary' }}>
                    {application.applicant.name}
                </Typography>
                <Typography variant="body2" sx={{ color: 'text.secondary' }}>
                    {application.applicant.ic}
                </Typography>
                <Typography gutterBottom variant="h5" component="div">
                    Scheme
                </Typography>
                <Typography variant="body2" sx={{ color: 'text.secondary' }}>
                    {application.scheme.name}
                </Typography>
                <Typography gutterBottom variant="h5" component="div">
                    Status
                </Typography>
                <Typography variant="body2" sx={{ color: application.application_status == 1 ? 'orange' : application.application_status == 2 ? "green" : "red" }}>
                    {ApplicationStatusMapping[application.application_status].toUpperCase()}
                </Typography>
            </CardContent>
            <CardActions disableSpacing>
                <Button size="small" sx={{ height: 40 }} onClick={() => { handleApproveClick(application.id ?? "") }}>Approve</Button>
                <Button size="small" sx={{ height: 40 }} onClick={() => { handleRejectClick(application.id ?? "") }}>Reject</Button>
                <Button size="small" sx={{ height: 40, color: "red" }} onClick={() => { handleDeleteClick(application.id ?? "") }}>Delete</Button>
            </CardActions>
        </Card>
    )
}

interface AddApplicationProps {
    applicants: Applicant[]
    schemes: Scheme[]
    setReload: () => void;
    setSelectedApplicant: (id: string) => void;
}
export const AddApplicationComponent = ({ applicants, schemes, setReload, setSelectedApplicant }: AddApplicationProps) => {
    const [addApplication, setAddApplication] = useState(false)
    const [applicantId, setApplicantId] = useState('');
    const [schemeId, setSchemeId] = useState('');
    const handleApplicantChange = (event: SelectChangeEvent) => {
        setApplicantId(event.target.value as string);
        setSelectedApplicant(event.target.value as string)
    };
    const handleSchemeChange = (event: SelectChangeEvent) => {
        setSchemeId(event.target.value as string);
    };

    const handleCLick = () => {
        setAddApplication((previousStatus) => {
            return !previousStatus;
        });
    }
    const handleSubmit = async () => {
        await ApiCallFunction<unknown, ApplicationsPayload>('POST', SubmitApplications, { applicant_id: applicantId, scheme_id: schemeId });
        setAddApplication((previousStatus) => {
            return !previousStatus;
        });
        setReload()
    }


    return (
        <Card sx={{ width: 250, height: 345 }}>
            {addApplication ?
                <div className="p-5">
                    <Typography>Select Applicant</Typography>
                    <Select
                        value={applicantId}
                        onChange={handleApplicantChange}
                        sx={{ width: 200 }}
                    >
                        {applicants.map((applicant) => {
                            return (
                                <MenuItem value={applicant.id} key={applicant.id}>{applicant.name}, {applicant.ic}</MenuItem>
                            )
                        })}
                    </Select>
                    <Typography>Select Scheme</Typography>
                    <Select
                        value={schemeId}
                        onChange={handleSchemeChange}
                        sx={{ width: 200 }}
                    >
                        {schemes.map((scheme) => {
                            return (
                                <MenuItem value={scheme.id} key={scheme.id}>{scheme.name}</MenuItem>
                            )
                        })}
                    </Select>

                    <Button size="small" sx={{ height: 40 }} onClick={handleCLick}>Cancel</Button>
                    <Button size="small" sx={{ height: 40 }} onClick={handleSubmit}>Submit</Button>
                </div> :
                <CardActionArea sx={{ height: 345 }} onClick={handleCLick}>
                    <div className="w-full flex flex-col items-center">
                        <AddIcon sx={{ width: 200, height: 200 }} />
                        Create New Application
                    </div>
                </CardActionArea>
            }
        </Card>
    )
}