import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';
import Resizer from 'react-image-file-resizer';
import '../styles/manage-profile.css';

interface User {
    id: number;
    name: string;
    surname: string;
    username: string;
    telegram: string;
    discord: string;
    email: string;
    skills: string[] | null;
    profileImage: string;
}

export const ManageProfile: React.FC = () => {
    const [user, setUser] = useState<User | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [validationError, setValidationError] = useState<string | null>(null);
    const [profileImage, setProfileImage] = useState<File | null>(null);
    const [setResizedProfileImage] = useState<string | null>(null);
    const navigate = useNavigate();

    useEffect(() => {
        fetchUserData();
    }, []);

    const fetchUserData = async () => {
        try {
            const response = await axios.get<User>('http://localhost:4000/v2/users/my', {
                withCredentials: true,
            });

            const skillsArray = response.data.skills ? response.data.skills[0].split(',').filter(skill => skill) : [];

            setUser({ ...response.data, skills: skillsArray });
            setLoading(false);

            if (response.data.profileImage) {
                fetchImageAsBlob(response.data.profileImage);
            }
        } catch (error) {
            console.error('Error fetching user data:', error);
            setError('Failed to fetch user data');
            setLoading(false);
        }
    };

    const fetchImageAsBlob = async (imageUrl: string) => {
        try {
            const response = await fetch(imageUrl);
            const blob = await response.blob();
            resizeImage(blob);
        } catch (error) {
            console.error('Error fetching image:', error);
            setLoading(false);
        }
    };

    const resizeImage = (imageBlob: Blob) => {
        Resizer.imageFileResizer(
            imageBlob,
            150,
            150,
            'JPEG',
            100,
            0,
            (uri) => {
                setResizedProfileImage(uri as string);
                setLoading(false);
            },
            'base64'
        );
    };

    const handleInputChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
        const { name, value } = e.target;
        setUser(prevUser => prevUser ? { ...prevUser, [name]: value } : null);
    };

    const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        if (e.target.files && e.target.files[0]) {
            setProfileImage(e.target.files[0]);
        }
    };

    const handleCheckboxChange = (skill: string) => {
        if (!user) return;

        const updatedSkills = user.skills ? [...user.skills] : [];

        if (updatedSkills.includes(skill)) {
            updatedSkills.splice(updatedSkills.indexOf(skill), 1);
        } else {
            updatedSkills.push(skill);
        }

        setUser({
            ...user,
            skills: updatedSkills,
        });
    };

    const handleSave = async () => {
        if (!user || !user.name || !user.surname) {
            setValidationError('Name and surname are required.');
            return;
        }

        setValidationError(null);

        const formData = new FormData();
        formData.append('name', user.name);
        formData.append('surname', user.surname);
        formData.append('telegram', user.telegram);
        formData.append('discord', user.discord);
        formData.append('skills', user.skills ? user.skills.join(',') : '');
        if (profileImage) {
            formData.append('profileImage', profileImage);
        }

        try {
            await axios.patch('http://localhost:4000/v2/users/update', formData, {
                withCredentials: true,
                headers: {
                    'Content-Type': 'multipart/form-data',
                },
            });
            navigate('/profile/my');
        } catch (error) {
            console.error('Error updating user:', error);
        }
    };

    if (loading) return <div>Loading...</div>;
    if (error) return <div>Error: {error}</div>;

    return (
        <div className="container">
            <h1>Edit Profile</h1>
            <form>
                <div className="form-group">
                    <label htmlFor="name">Name:</label>
                    <input type="text" id="name" name="name" value={user?.name || ''} onChange={handleInputChange} />
                </div>
                <div className="form-group">
                    <label htmlFor="surname">Surname:</label>
                    <input type="text" id="surname" name="surname" value={user?.surname || ''} onChange={handleInputChange} />
                </div>
                <div className="form-group">
                    <label htmlFor="username">Username:</label>
                    <input type="text" id="username" name="username" value={user?.username || ''} onChange={handleInputChange} />
                </div>
                <div className="form-group">
                    <label htmlFor="telegram">Telegram:</label>
                    <input type="text" id="telegram" name="telegram" value={user?.telegram || ''} onChange={handleInputChange} />
                </div>
                <div className="form-group">
                    <label htmlFor="discord">Discord:</label>
                    <input type="text" id="discord" name="discord" value={user?.discord || ''} onChange={handleInputChange} />
                </div>
                <div className="form-group">
                    <label htmlFor="profileImage">Profile Image:</label>
                    <input type="file" id="profileImage" name="profileImage" onChange={handleFileChange} />
                </div>
                <div className="form-group">
                    <label>Skills:</label>
                    {['golang', 'python', 'java', 'javascript', 'c++', 'c#', 'php', 'rust'].map(skill => (
                        <label key={skill} className="checkbox-label">
                            <input
                                type="checkbox"
                                name="skills"
                                value={skill}
                                checked={user?.skills?.includes(skill) || false}
                                onChange={() => handleCheckboxChange(skill)}
                            />
                            {skill}
                        </label>
                    ))}
                </div>
                <div className="form-group">
                    <button type="button" className="blue-button" onClick={handleSave}>Save</button>
                </div>
            </form>
            {validationError && <div className="error-message">{validationError}</div>}
        </div>
    );
};

export default ManageProfile;
