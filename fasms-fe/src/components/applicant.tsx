import { Applicant, ApplicantsPayload, default_applicant, default_household, Household } from "@/models/applicants"
import { EmploymentStatusMapping, MaritalStatusMapping, RelationMapping, SexMapping } from "@/utils/const";
import { Button, CardActionArea, CardActions, CardContent, CardMedia, Collapse, TextField, Typography } from "@mui/material";
import * as React from 'react';
import { styled } from '@mui/material/styles';
import Card from '@mui/material/Card';
import IconButton, { IconButtonProps } from '@mui/material/IconButton';
import ExpandMoreIcon from '@mui/icons-material/ExpandMore';
import AddIcon from '@mui/icons-material/Add';
import { memo, useCallback, useState } from "react";
import { AddApplicants, DeleteApplicants } from "@/api/urls";
import { ApiCallFunction } from "@/api/callFunction";


interface ApplicantProps {
    applicant: Applicant;
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

export const ApplicantComponent = ({ applicant, setReload }: ApplicantProps) => {
    const [expanded, setExpanded] = React.useState(false);

    const handleExpandClick = () => {
        setExpanded(!expanded);
    };
    const handleDeleteClick = async (id: string) => {
        await ApiCallFunction('DELETE', DeleteApplicants(id));
        setReload()
    };
    return (
        <Card sx={{ width: 250, minHeight: 385 }}>
            <CardMedia
                component="img"
                sx={{ height: 120, width: 120, objectFit: "cover" }}
                image={applicant.sex == 1 ? "/male.jpg" : "/female.jpg"}
                title="green iguana"
            />
            <CardContent>
                <Typography gutterBottom variant="h5" component="div">
                    {applicant.name}
                </Typography>
                <Typography variant="body2" sx={{ color: 'text.secondary' }}>
                    {applicant.ic}
                </Typography>
                <Typography variant="body2" sx={{ color: 'text.secondary' }}>
                    {EmploymentStatusMapping[applicant.employment_status]}
                </Typography>
                <Typography variant="body2" sx={{ color: 'text.secondary' }}>
                    {MaritalStatusMapping[applicant.marital_status]}
                </Typography>
                <Typography variant="body2" sx={{ color: 'text.secondary' }}>
                    {SexMapping[applicant.sex]}
                </Typography>
            </CardContent>
            <div className="h-14"></div>
            <CardActions disableSpacing>
                <Button size="small" sx={{ height: 40, color: "red" }} onClick={() => { handleDeleteClick(applicant.id ?? "") }}>Delete</Button>
                {applicant.households &&
                    <ExpandMore
                        expand={expanded}
                        onClick={handleExpandClick}
                        aria-expanded={expanded}
                        aria-label="show more"
                    >
                        <ExpandMoreIcon />
                    </ExpandMore>
                }
            </CardActions>
            <Collapse in={expanded} timeout="auto" unmountOnExit>
                {applicant.households &&
                    <CardContent>
                        <Typography gutterBottom variant="h5" component="div">
                            Households
                        </Typography>
                        {applicant.households.map((h) => {
                            return (
                                <div key={h.id} className="flex flex-row pb-5 justify-center items-center">
                                    <CardMedia
                                        component="img"
                                        sx={{ height: 100, width: 100, objectFit: "cover" }}
                                        image={h.sex == 1 ? "/male.jpg" : "/female.jpg"}
                                        title="green iguana"
                                    />
                                    <div className="pl-5">
                                        <Typography variant="h6" >
                                            {h.name}
                                        </Typography>
                                        <Typography variant="body2" sx={{ color: 'text.secondary' }}>
                                            {h.ic}
                                        </Typography>
                                        <Typography variant="body2" sx={{ color: 'text.secondary' }}>
                                            {SexMapping[h.sex]}
                                        </Typography>
                                        <Typography variant="body2" sx={{ color: 'text.secondary' }}>
                                            {h.dob}
                                        </Typography>
                                        <Typography variant="body2" sx={{ color: 'text.secondary' }} >
                                            {EmploymentStatusMapping[h.employment_status]}
                                        </Typography>
                                        <Typography variant="body2" sx={{ color: 'text.secondary' }} >
                                            {MaritalStatusMapping[h.marital_status]}
                                        </Typography>
                                        <Typography variant="body2" sx={{ color: 'text.secondary' }}>
                                            {RelationMapping[h.relation]}
                                        </Typography>
                                    </div>
                                </div>
                            )
                        })}
                    </CardContent>
                }
            </Collapse>
        </Card>
    )
}

interface AddApplicantProps {
    setReload: () => void;
}
export const AddApplicantComponent = ({ setReload }: AddApplicantProps) => {
    const [addApplicant, setAddApplicant] = useState(false)
    const [applicant, setApplicant] = useState<Applicant>(default_applicant)
    const [newHouseholds, setNewHouseholds] = useState<Household[]>([])
    const handleCLick = () => {
        setAddApplicant((previousStatus) => {
            return !previousStatus;
        });
    }
    const handleAddHouseholdCLick = () => {
        setNewHouseholds([...newHouseholds, default_household]);
    }
    const handleHouseholdChange = useCallback(
        (field: keyof Household, index: number, value: string | number) => {
            setNewHouseholds((prev) => {
                const updated = [...prev];
                updated[index] = { ...updated[index], [field]: value };
                return updated;
            });
        },
        []
    );
    const handleApplicantChange = useCallback(
        (field: keyof Applicant, value: string | number) => {
            setApplicant({ ...applicant, [field]: value })
        },
        [applicant]
    );


    const handleSubmit = async () => {
        await ApiCallFunction<unknown, ApplicantsPayload>('POST', AddApplicants, { applicants: [{ ...applicant, households: newHouseholds }] });
        setAddApplicant((previousStatus) => {
            return !previousStatus;
        });
        setNewHouseholds([]);
        setReload()
    }
    return (
        <>
            <Card sx={{ width: 250, height: 400 }}>
                {addApplicant ?
                    <div className="p-2">
                        <Typography>Applicant</Typography>
                        <ApplicantTextField type="string" label="Applicant Name" value={applicant.name} field="name" onChange={handleApplicantChange} />
                        <ApplicantTextField type="string" label="Applicant IC" value={applicant.ic} field="ic" onChange={handleApplicantChange} />
                        <ApplicantTextField type="string" label="Applicant DOB" value={applicant.dob} field="dob" onChange={handleApplicantChange} />
                        <ApplicantTextField type="number" label="Applicant employment_status" value={applicant.employment_status} field="employment_status" onChange={handleApplicantChange} />
                        <ApplicantTextField type="number" label="Applicant marital_status" value={applicant.marital_status} field="marital_status" onChange={handleApplicantChange} />
                        <ApplicantTextField type="number" label="Applicant Gender" value={applicant.sex} field="sex" onChange={handleApplicantChange} />
                        <Button size="small" sx={{ height: 30 }} onClick={handleSubmit}>Add</Button>
                        <Button size="small" sx={{ height: 30 }} onClick={handleAddHouseholdCLick}>New Household</Button>
                        <Button size="small" sx={{ height: 30 }} onClick={handleCLick}>Cancel</Button>
                    </div> :
                    <CardActionArea sx={{ height: 400 }} onClick={handleCLick}>
                        <div className="w-full flex flex-col items-center">
                            <AddIcon sx={{ width: 200, height: 200 }} />
                            Add Applicant
                        </div>
                    </CardActionArea>
                }
            </Card>
            {newHouseholds.map((h, index) => {
                return (
                    <Card sx={{ width: 250, height: 400 }} key={index}>
                        <div className="p-2">
                            <Typography>Household {index + 1}</Typography>
                            <HouseholdTextField type="string" value={h.name} label={"Household Name"} field={"name"} index={index} onChange={handleHouseholdChange} />
                            <HouseholdTextField type="string" value={h.ic} label={"Household ic"} field={"ic"} index={index} onChange={handleHouseholdChange} />
                            <HouseholdTextField type="string" value={h.dob} label={"Household DOB"} field={"dob"} index={index} onChange={handleHouseholdChange} />
                            <HouseholdTextField type="number" value={h.employment_status} label={"Household employment_status"} field={"employment_status"} index={index} onChange={handleHouseholdChange} />
                            <HouseholdTextField type="number" value={h.marital_status} label={"Household marital_status"} field={"marital_status"} index={index} onChange={handleHouseholdChange} />
                            <HouseholdTextField type="number" value={h.sex} label={"Household Gender"} field={"sex"} index={index} onChange={handleHouseholdChange} />
                            <HouseholdTextField type="number" value={h.relation} label={"Household Relation"} field={"relation"} index={index} onChange={handleHouseholdChange} />

                            <Button size="small" sx={{ height: 40 }} onClick={() => {
                                const updatedHouseholds = [...newHouseholds];
                                updatedHouseholds.splice(index, 1);
                                setNewHouseholds(updatedHouseholds);
                            }}>
                                Cancel
                            </Button>
                        </div>
                    </Card>
                )
            })}
        </>
    )
}


interface HouseholdTextFieldProps {
    label: string;
    type: string;
    field: keyof Household;
    value: string | number;
    index: number;
    onChange: (field: keyof Household, index: number, value: string | number) => void;
}
const HouseholdTextField = memo(function HouseholdCard({ label, type, field, value, index, onChange }: HouseholdTextFieldProps) {
    return (
        <div className="p-1">
            <TextField
                size="small"
                label={label}
                value={value}
                onChange={(e) => {
                    if (type == "number") {
                        value = parseInt(e.target.value)
                    } else {
                        value = e.target.value
                    }
                    onChange(field, index, value)
                }}
            />
        </div>
    );
});


interface ApplicantTextFieldProps {
    label: string;
    type: string;
    field: keyof Applicant;
    value: string | number;
    onChange: (field: keyof Applicant, value: string | number) => void;
}
const ApplicantTextField = memo(function ApplicantTextField({ label, type, field, value, onChange }: ApplicantTextFieldProps) {
    return (
        <div className="p-1">
            <TextField
                size="small"
                label={label}
                value={value}
                onChange={(e) => {
                    if (type == "number") {
                        value = parseInt(e.target.value)
                    } else {
                        value = e.target.value
                    }
                    onChange(field, value)
                }}
            />
        </div>
    );
});