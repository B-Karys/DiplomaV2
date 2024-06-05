import React, { useState, useEffect } from 'react';
import '../styles/filter.css'; // Import the CSS file for styling

interface FilterProps {
    onFilterChange: (type: string, skills: string[], sort: string) => void;
    initialType: string;
    initialSkills: string[];
    initialSort: string;
}

const Filter: React.FC<FilterProps> = ({ onFilterChange, initialType, initialSkills, initialSort }) => {
    const [type, setType] = useState<string>(initialType);
    const [skills, setSkills] = useState<string[]>(initialSkills);
    const [sortField, setSortField] = useState<string>(initialSort);

    const handleTypeChange = (event: React.ChangeEvent<HTMLSelectElement>) => {
        const selectedType = event.target.value;
        setType(selectedType);
    };

    const handleSkillChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        const skill = event.target.name;
        const isChecked = event.target.checked;
        const updatedSkills = isChecked
            ? [...skills, skill]
            : skills.filter(s => s !== skill);
        setSkills(updatedSkills);
    };

    const handleSortFieldChange = (event: React.ChangeEvent<HTMLSelectElement>) => {
        const selectedSortField = event.target.value;
        setSortField(selectedSortField);
    };

    useEffect(() => {
        onFilterChange(type, skills, sortField);
    }, [type, skills, sortField]);

    return (
        <div>
            <h2>Filter</h2>
            <div className="filter-group">
                <label>
                    Select a type:
                    <select value={type} onChange={handleTypeChange}>
                        <option value="">Select a type</option>
                        <option value="team finding">Team Finding</option>
                        <option value="user finding">User Finding</option>
                    </select>
                </label>
            </div>
            <div className="filter-group">
                Select skills:
                {['golang', 'python', 'java', 'javascript', 'c++','c#','rust','php','kotlin','ruby'].map(skill => (
                    <label className="checkbox-label" key={skill}>
                        <input
                            type="checkbox"
                            name={skill}
                            checked={skills.includes(skill)}
                            onChange={handleSkillChange}
                        />
                        {skill.charAt(0).toUpperCase() + skill.slice(1)}
                    </label>
                ))}
            </div>
            <div className="filter-group">
                <label>
                    Sort by:
                    <select value={sortField} onChange={handleSortFieldChange}>
                        <option value="name">Name</option>
                        <option value="created_at">Created At</option>
                    </select>
                </label>
            </div>
        </div>
    );
};

export default Filter;
