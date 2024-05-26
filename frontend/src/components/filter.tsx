import React, { useState, useEffect } from 'react';

interface FilterProps {
    onFilterChange: (type: string, skills: string[], sort: string, search: string) => void;
    initialType: string;
    initialSkills: string[];
    initialSort: string;
    initialSearch: string;
}

const Filter: React.FC<FilterProps> = ({ onFilterChange, initialType, initialSkills, initialSort, initialSearch }) => {
    const [type, setType] = useState<string>(initialType);
    const [skills, setSkills] = useState<string[]>(initialSkills);
    const [sort, setSort] = useState<string>(initialSort);
    const [search, setSearch] = useState<string>(initialSearch);

    useEffect(() => {
        setType(initialType);
    }, [initialType]);

    useEffect(() => {
        setSkills(initialSkills);
    }, [initialSkills]);

    useEffect(() => {
        setSort(initialSort);
    }, [initialSort]);

    useEffect(() => {
        setSearch(initialSearch);
    }, [initialSearch]);

    const handleTypeChange = (event: React.ChangeEvent<HTMLSelectElement>) => {
        const selectedType = event.target.value;
        setType(selectedType);
        onFilterChange(selectedType, skills, sort, search);
    };

    const handleSkillChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        const skill = event.target.name;
        const isChecked = event.target.checked;
        const updatedSkills = isChecked
            ? [...skills, skill]
            : skills.filter(s => s !== skill);
        setSkills(updatedSkills);
        onFilterChange(type, updatedSkills, sort, search);
    };

    const handleSortChange = (event: React.ChangeEvent<HTMLSelectElement>) => {
        const selectedSort = event.target.value;
        setSort(selectedSort);
        onFilterChange(type, skills, selectedSort, search);
    };

    const handleSearchChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        const searchQuery = event.target.value;
        setSearch(searchQuery);
        onFilterChange(type, skills, sort, searchQuery);
    };

    return (
        <div style={{ padding: '10px', border: '1px solid #ccc' }}>
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
            <div>
                <label>
                    Sort by:
                    <select value={sort} onChange={handleSortChange}>
                        <option value="asc">Ascending</option>
                        <option value="desc">Descending</option>
                    </select>
                </label>
            </div>
            <div>
                <label>
                    Search:
                    <input
                        type="text"
                        value={search}
                        onChange={handleSearchChange}
                    />
                </label>
            </div>
        </div>
    );
};

export default Filter;
