import { EmploymentStatusMapping, MaritalStatusMapping, RelationMapping, SexMapping } from "@/utils/const";
import { Button, CardActionArea, CardActions, CardContent, TextField, Typography } from "@mui/material";
import Card from '@mui/material/Card';
import { Benefits, CriteriaGroups, Criterias, default_benefits, default_criteria, default_criteria_group, default_scheme, Scheme, SchemePayload } from "@/models/schemes";
import { ApiCallFunction } from "@/api/callFunction";
import { AddSchemes, DeleteScheme } from "@/api/urls";
import { memo, useCallback, useState } from "react";
import AddIcon from '@mui/icons-material/Add';
import Checkbox from '@mui/material/Checkbox';


interface SchemeProps {
    scheme: Scheme;
    setReload: () => void;
}


export const SchemeComponent = ({ scheme, setReload }: SchemeProps) => {

    const handleDeleteClick = async (id: string) => {
        if (id != "") {
            await ApiCallFunction('DELETE', DeleteScheme(id));
            setReload()
        }
    };
    return (
        <Card sx={{ width: 400 }}>
            <CardContent>
                <Typography gutterBottom variant="h5" component="div">
                    {scheme.name}
                </Typography>
            </CardContent>
            <CardContent>
                <Typography gutterBottom variant="h5" component="div">
                    Criterias
                </Typography>
                {scheme.criteria_groups.map((cg) => {
                    return (
                        <div key={cg.id}>
                            Criteria Group
                            {cg.criterias.map((c) => {
                                console.log(c.is_household)
                                return (
                                    <div key={c.id} className="pb-5">
                                        <Typography gutterBottom variant="body2" component="div">
                                            Is Household Criteria: {c.is_household ? "Yes" : "No"}
                                        </Typography>
                                        <Typography gutterBottom variant="body2" component="div">
                                            {SexMapping[c.sex]}
                                        </Typography>
                                        <Typography gutterBottom variant="body2" component="div">
                                            {EmploymentStatusMapping[c.employment_status]}
                                        </Typography>
                                        <Typography gutterBottom variant="body2" component="div">
                                            {MaritalStatusMapping[c.marital_status]}
                                        </Typography>
                                        <Typography gutterBottom variant="body2" component="div">
                                            {RelationMapping[c.relation]}
                                        </Typography>
                                        <Typography gutterBottom variant="body2" component="div">
                                            age:{c.age_lower_limit} - {c.age_upper_limit}
                                        </Typography>
                                    </div>
                                )
                            })}
                        </div>
                    )
                })}
            </CardContent>
            <CardContent>
                <Typography gutterBottom variant="h5" component="div">
                    Benefits
                </Typography>
                {scheme.benefits.map((b) => {
                    return (
                        <div key={b.id}>
                            <Typography gutterBottom variant="body2" component="div">
                                {b.name}
                            </Typography>
                            <Typography gutterBottom variant="body2" component="div">
                                S$:{b.amount}
                            </Typography>
                        </div>
                    )
                })}
            </CardContent>
            <CardActions disableSpacing>
                <Button size="small" sx={{ height: 40, color: "red" }} onClick={() => { handleDeleteClick(scheme.id ?? "") }}>DELETE</Button>
            </CardActions>
        </Card>
    )
}


interface AddSchemeProps {
    setReload: () => void;
}
export const AddSchemeComponent = ({ setReload }: AddSchemeProps) => {
    const [addScheme, setAddScheme] = useState(false)
    const [scheme, setScheme] = useState<Scheme>(default_scheme)
    const [newCriteriaGroup, setNewCriteriaGroup] = useState<CriteriaGroups[]>([])

    const handleAddCriteriaGroupCLick = () => {
        setNewCriteriaGroup([...newCriteriaGroup, default_criteria_group]);
    }
    const [benefit, setBenefit] = useState<Benefits>(default_benefits)
    const handleCLick = () => {
        setAddScheme((previousStatus) => {
            return !previousStatus;
        });
    }
    const handleSubmit = async () => {
        const payload = scheme
        payload.criteria_groups = newCriteriaGroup
        payload.benefits = [benefit]
        await ApiCallFunction<unknown, SchemePayload>('POST', AddSchemes, { schemes: [payload] });
        setAddScheme((previousStatus) => {
            return !previousStatus;
        });
        setNewCriteriaGroup([])
        setReload()
    }

    const handleHouseholdChange = useCallback(
        (field: keyof Criterias, index: number, c_index: number, value: string | number | boolean) => {
            setNewCriteriaGroup((prev) => {
                const updated = [...prev];
                const updatedGroup = { ...updated[index] };
                const updatedCriterias = [...updatedGroup.criterias];
                const updatedCriteria = { ...updatedCriterias[c_index], [field]: value };

                updatedCriterias[c_index] = updatedCriteria;
                updatedGroup.criterias = updatedCriterias;
                updated[index] = updatedGroup;

                return updated;
            });
        },
        []
    );
    return (
        <>
            <Card sx={{ width: 250, height: 300 }}>
                {addScheme ?
                    <div className="p-5">
                        <TextField
                            size="small"
                            label="Scheme Name"
                            onChange={(event: React.ChangeEvent<HTMLInputElement>) => {
                                setScheme({ ...scheme, name: event.target.value });
                            }}
                        />
                        <TextField
                            size="small"
                            label="Benefit Name"
                            onChange={(event: React.ChangeEvent<HTMLInputElement>) => {
                                setBenefit({ ...benefit, name: event.target.value });
                            }}
                        />
                        <TextField
                            size="small"
                            label="Benefit amount"
                            onChange={(event: React.ChangeEvent<HTMLInputElement>) => {
                                setBenefit({ ...benefit, amount: parseFloat(event.target.value) });
                            }}
                        />

                        <Button size="small" sx={{ height: 40 }} onClick={handleSubmit}>Add</Button>
                        <Button size="small" sx={{ height: 40 }} onClick={handleAddCriteriaGroupCLick}>Add Criteria Group</Button>
                        <Button size="small" sx={{ height: 40 }} onClick={handleCLick}>Cancel</Button>
                    </div> :
                    <CardActionArea sx={{ height: 300 }} onClick={handleCLick}>
                        <div className="w-full flex flex-col items-center">
                            <AddIcon sx={{ width: 200, height: 200 }} />
                            Add Scheme
                        </div>
                    </CardActionArea>
                }
            </Card>

            {newCriteriaGroup.map((cg, index) => {
                return (
                    <Card sx={{ width: 250 }} key={index}>
                        <div className="p-2">
                            <Typography>Criteria Group {index + 1}</Typography>
                            {cg.criterias.map((c, c_index) => {
                                return (
                                    <div key={c_index}>
                                        <div className="pl-2">
                                            Is Household
                                            <Checkbox checked={c.is_household} onChange={(e) => {
                                                handleHouseholdChange("is_household", index, c_index, e.target.checked)
                                            }} />
                                        </div>
                                        <CriteriaTextField type="number" value={c.sex} label={"Gender"} field={"sex"} index={index} c_index={c_index} onChange={handleHouseholdChange} />
                                        <CriteriaTextField type="number" value={c.employment_status} label={"employment status"} field={"employment_status"} index={index} c_index={c_index} onChange={handleHouseholdChange} />
                                        <CriteriaTextField type="number" value={c.marital_status} label={"marital status"} field={"marital_status"} index={index} c_index={c_index} onChange={handleHouseholdChange} />
                                        <CriteriaTextField type="number" value={c.relation} label={"Relation"} field={"relation"} index={index} c_index={c_index} onChange={handleHouseholdChange} />
                                        <CriteriaTextField type="number" value={c.age_lower_limit} label={"Age Lower Limit"} field={"age_lower_limit"} index={index} c_index={c_index} onChange={handleHouseholdChange} />
                                        <CriteriaTextField type="number" value={c.age_upper_limit} label={"Age Upper Limit"} field={"age_upper_limit"} index={index} c_index={c_index} onChange={handleHouseholdChange} />

                                        <Button size="small" sx={{ height: 40 }} onClick={() => {
                                            const updated = [...newCriteriaGroup];
                                            const updatedGroup = { ...updated[index] };
                                            const updatedCriterias = [...updatedGroup.criterias];
                                            updatedCriterias.splice(c_index, 1)

                                            updatedGroup.criterias = updatedCriterias;
                                            updated[index] = updatedGroup;
                                            setNewCriteriaGroup(updated);
                                        }}>
                                            Remove
                                        </Button>
                                    </div>
                                )
                            })}
                        </div>
                        <Button size="small" sx={{ height: 40 }} onClick={() => {
                            const updatedHouseholds = [...newCriteriaGroup];
                            updatedHouseholds.splice(index, 1);
                            setNewCriteriaGroup(updatedHouseholds);
                        }}>
                            Cancel
                        </Button>
                        <Button size="small" sx={{ height: 40 }} onClick={() => {
                            const updated = [...newCriteriaGroup];
                            const updatedGroup = { ...updated[index] };
                            const updatedCriterias = [...updatedGroup.criterias];
                            updatedCriterias.push(default_criteria)

                            updatedGroup.criterias = updatedCriterias;
                            updated[index] = updatedGroup;
                            setNewCriteriaGroup(updated);
                        }}>
                            Add Criteria
                        </Button>
                    </Card>
                )
            })}
        </>
    )
}


interface CriteriaTextFieldProps {
    label: string;
    type: string;
    field: keyof Criterias;
    value: string | number | boolean;
    index: number;
    c_index: number;
    onChange: (field: keyof Criterias, index: number, c_index: number, value: string | number | boolean) => void;
}
const CriteriaTextField = memo(function CriteriaTextField({ label, type, field, value, index, c_index, onChange }: CriteriaTextFieldProps) {
    return (
        <div className="p-1">
            <TextField
                size="small"
                label={label}
                value={value}
                onChange={(e) => {
                    if (type == "number") {
                        value = parseInt(e.target.value)
                    } else if (type == "boolean") {
                        value = e.target.value == "true" ? true : false
                    } else {
                        value = e.target.value
                    }
                    onChange(field, index, c_index, value)
                }}
            />
        </div>
    );
});