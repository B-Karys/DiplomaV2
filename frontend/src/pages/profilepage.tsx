import { useState, useEffect } from 'react';

interface User {
    id: number;
    created_at: string;
    name: string;
    surname: string;
    username: string;
    telegram: string;
    discord: string;
    email: string;
    skills: string | null;
    activated: boolean;
}

export function ProfilePage() {
    const [user, setUser] = useState<User | null>(null);
    const [loading, setLoading] = useState<boolean>(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const parseJwt = (token: string) => {
            try {
                return JSON.parse(atob(token.split('.')[1]));
            } catch (e) {
                return null;
            }
        };

        const jwtCookie = document.cookie.split(';').find(cookie => cookie.trim().startsWith('jwt='));

        if (jwtCookie) {
            const jwtToken = jwtCookie.split('=')[1];
            const decodedToken = parseJwt(jwtToken);

            if (decodedToken) {
                const userId = decodedToken.sub; // Assuming 'sub' holds the user ID
                fetchUserData(userId, jwtToken);
            }
        }
    }, []);

    const fetchUserData = async (userId: number, token: string) => {
        try {
            const response = await fetch(`http://localhost:4000/v2/users/${userId}`, {
                method: 'GET',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`
                }
            });

            if (!response.ok) {
                throw new Error('Failed to fetch user data');
            }

            const userData = await response.json();
            setUser(userData);
            setLoading(false);
        } catch (error) {
            // @ts-ignore
            setError(error.message);
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
        <div>
            <h1>Profile</h1>
            {user ? (
                <div>
                    <p>Name: {user.name}</p>
                    <p>Surname: {user.surname}</p>
                    <p>Username: {user.username}</p>
                    <p>Telegram: {user.telegram}</p>
                    <p>Discord: {user.discord}</p>
                    <p>Email: {user.email}</p>
                    <p>Skills: {user.skills}</p>
                </div>
            ) : (
                <p>No user data available</p>
            )}
        </div>
    );
}

export default ProfilePage;
