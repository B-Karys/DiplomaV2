import { useState, useEffect } from 'react';
import './profile-page.css';

interface User {
    id: number;
    created_at: string;
    name: string;
    surname: string;
    username: string;
    telegram: string;
    discord: string;
    email: string;
    skills: string[] | null;
    activated: boolean;
    profileImage: string; // Change to the GCS URL
}

export function ProfilePage() {
    const [user, setUser] = useState<User | null>(null);
    const [loading, setLoading] = useState<boolean>(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        fetchUserData();
    }, []);

    const fetchUserData = async () => {
        try {
            const response = await fetch('http://localhost:4000/v2/users/my', {
                method: 'GET',
                headers: {
                    'Content-Type': 'application/json',
                },
                credentials: 'include', // Ensure cookies are included in the request
            });

            if (!response.ok) {
                throw new Error('Failed to fetch user data');
            }

            const userData = await response.json();
            setUser(userData);
            setLoading(false);
        } catch (error) {
            setError((error as Error).message);
            setLoading(false);
            console.error('Error fetching user data:', error);
        }
    };

    if (loading) {
        return <p>Loading...</p>;
    }

    if (error) {
        return <p>{error}</p>;
    }

    return (
        <div className="profile-container">
            <div className="profile-content">
                <div className="profile-sidebar">
                    <div className="profile-photo">
                        {user?.profileImage ? (
                            <img src={user.profileImage} alt="Profile" /> // Update to use the GCS URL
                        ) : (
                            'Profile Photo Of User'
                        )}
                    </div>
                    <div className="profile-contacts">
                        <p>@{user?.telegram}</p>
                        <p>@{user?.discord}</p>
                    </div>
                </div>
                <div className="profile-main">
                    <h1>{user?.name} ({user?.username}) {user?.surname}</h1>
                    <p>{user?.activated ? 'Active' : 'Inactive'} user interested in various technologies and projects.</p>
                    <div className="profile-skills">
                        {user?.skills?.map((skill, index) => (
                            <p key={index}>{skill}</p>
                        ))}
                    </div>
                    <div className="profile-team">
                        <h2>Team Name</h2>
                        <p>Team Description</p>
                    </div>
                </div>
            </div>
        </div>
    );
}

export default ProfilePage;
