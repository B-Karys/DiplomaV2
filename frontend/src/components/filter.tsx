import React, { useState, useEffect } from 'react';
import './filter.module.css'; // Import the CSS file for styling

interface FilterProps {
    onFilterChange: (type: string, skills: string[], sort: string,) => void;
    initialType: string;
    initialSkills: string[];
    initialSort: string;
    initialSearch: string;
}

const Filter: React.FC<FilterProps> = ({ onFilterChange, initialType, initialSkills, initialSort, }) => {
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
        // Submit the filters when type, skills, sortField, sortOrder, or search change
        onFilterChange(type, skills, sortField, );
    }, [type, skills, sortField, ]);

    return (
        <div className="filter-container">
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
                <label className="checkbox-label">
                    <input
                        type="checkbox"
                        name="javascript"
                        checked={skills.includes('javascript')}
                        onChange={handleSkillChange}
                    />
                    JavaScript
                </label>
                <label className="checkbox-label">
                    <input
                        type="checkbox"
                        name="golang"
                        checked={skills.includes('golang')}
                        onChange={handleSkillChange}
                    />
                    Golang
                </label>
                <label className="checkbox-label">
                    <input
                        type="checkbox"
                        name="python"
                        checked={skills.includes('python')}
                        onChange={handleSkillChange}
                    />
                    Python
                </label>
                <label className="checkbox-label">
                    <input
                        type="checkbox"
                        name="java"
                        checked={skills.includes('java')}
                        onChange={handleSkillChange}
                    />
                    Java
                </label>
                <label className="checkbox-label">
                    <input
                        type="checkbox"
                        name="c++"
                        checked={skills.includes('c++')}
                        onChange={handleSkillChange}
                    />
                    C++
                </label>
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
