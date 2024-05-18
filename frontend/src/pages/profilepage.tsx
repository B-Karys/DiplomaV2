import React, { useState, useEffect } from 'react';

interface User {
    name: string;
    surname: string;
    username: string;
    telegram: string;
    discord: string;
    email: string;
}

function ProfilePage() {
    const [user, setUser] = useState<User | null>(null);
    const [loading, setLoading] = useState<boolean>(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        // Function to parse JWT token from cookie
        const parseJwt = (token: string) => {
            try {
                return JSON.parse(atob(token.split('.')[1]));
            } catch (e) {
                return null;
            }
        };

        // Get JWT token from cookie
        const jwtCookie = document.cookie.split(';').find(cookie => cookie.trim().startsWith('jwt='));

        if (jwtCookie) {
            const jwtToken = jwtCookie.split('=')[1];
            const decodedToken = parseJwt(jwtToken);

            // Set user data from decoded token
            if (decodedToken) {
                const userId = decodedToken.sub; // Assuming 'sub' holds the user ID
                fetchUserData(userId);
            }
        }
    }, []);

    const fetchUserData = async (userId: number) => {
        try {
            const response = await Promise.race([
                fetch(`http://localhost:4000/v2/users/${userId}`, {
                    method: 'GET',
                    headers: {
                        'Content-Type': 'application/json'
                    }
                }),
            ]);
            if (!response.ok) {
                throw new Error('Failed to fetch user data');
            }
            const userData = await response.json();

            // Extracting specific fields from the response
            const { name, surname, username, telegram, discord, email } = userData.user;

            setUser({ name, surname, username, telegram, discord, email });
            setLoading(false); // Set loading to false when data is fetched successfully
        } catch (error) {
            // @ts-ignore
            setError(error.message); // Set error message if there is an error
            setLoading(false); // Set loading to false in case of error
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
                </div>
            ) : (
                <p>No user data available</p>
            )}
        </div>
    );
}

export default ProfilePage;
