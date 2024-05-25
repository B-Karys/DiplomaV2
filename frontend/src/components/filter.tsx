import React, { useState, useEffect } from 'react';

interface FilterProps {
    onFilterChange: (type: string, skills: string[]) => void;
    initialType: string;
    initialSkills: string[];
}
// TODO: CREATE THE SKILLS AFTER MAKING THE METADATA
const Filter: React.FC<FilterProps> = ({ onFilterChange, initialType, initialSkills }) => {
    const [type, setType] = useState<string>(initialType);
    const [skills, setSkills] = useState<string[]>(initialSkills);

    useEffect(() => {
        setType(initialType);
    }, [initialType]);

    useEffect(() => {
        setSkills(initialSkills);
    }, [initialSkills]);

    const handleTypeChange = (event: React.ChangeEvent<HTMLSelectElement>) => {
        const selectedType = event.target.value;
        setType(selectedType);
        onFilterChange(selectedType, skills);
    };

    const handleSkillChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        const skill = event.target.name;
        const isChecked = event.target.checked;
        const updatedSkills = isChecked
            ? [...skills, skill]
            : skills.filter(s => s !== skill);
        setSkills(updatedSkills);
        onFilterChange(type, updatedSkills);
    };

    return (
        <div style={{padding: '10px', border: '1px solid #ccc' }}>
            <h2>Filter</h2>
            <div>
                <label>
                    Select a type:
                    <select value={type} onChange={handleTypeChange}>
                        <option value="">Select a type</option>
                        <option value="teamFinding">Team Finding</option>
                        <option value="userFinding">User Finding</option>
                    </select>
                </label>
            </div>
            <div>
                <label>
                    <input
                        type="checkbox"
                        name="javascript"
                        checked={skills.includes('javascript')}
                        onChange={handleSkillChange}
                    />
                    JavaScript
                </label>
                <label>
                    <input
                        type="checkbox"
                        name="golang"
                        checked={skills.includes('golang')}
                        onChange={handleSkillChange}
                    />
                    Golang
                </label>
                <label>
                    <input
                        type="checkbox"
                        name="python"
                        checked={skills.includes('python')}
                        onChange={handleSkillChange}
                    />
                    Python
                </label>
            </div>
        </div>
    );
};

export default Filter;
