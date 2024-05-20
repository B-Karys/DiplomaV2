import React, { createContext, useState, useEffect, useContext, ReactNode } from 'react';

// Define type for context value
type AuthContextType = {
    isAuthenticated: boolean;
    setIsAuthenticated: React.Dispatch<React.SetStateAction<boolean>>;
};

// Create AuthContext
const AuthContext = createContext<AuthContextType | null>(null);

export const AuthProvider = ({ children }: { children: ReactNode }) => {
    const [isAuthenticated, setIsAuthenticated] = useState(false);

    useEffect(() => {
        // Check if the user is already authenticated from local storage
        const storedAuthenticated = localStorage.getItem('authenticated') === 'true';
        if (storedAuthenticated) {
            setIsAuthenticated(true);
            return; // Exit early if already authenticated
        }

        // If not authenticated, fetch authentication status
        fetch('http://localhost:4000/v2/users/check-auth', {
            method: 'GET',
            credentials: 'include', // Include cookies in the request
        })
            .then(response => response.json())
            .then(data => {
                if (data.authenticated === "true") {
                    setIsAuthenticated(true);
                    // Store authenticated status in local storage
                    localStorage.setItem('authenticated', 'true');
                } else {
                    setIsAuthenticated(false);
                    // Remove authenticated status from local storage if not authenticated
                    localStorage.removeItem('authenticated');
                }
            })
            .catch(error => {
                console.error('Error:', error);
                setIsAuthenticated(false);
            });
    }, []);

    // Define the context value
    const contextValue: AuthContextType = {
        isAuthenticated,
        setIsAuthenticated
    };

    return (
        <AuthContext.Provider value={contextValue}>
            {children}
        </AuthContext.Provider>
    );
};

export const useAuth = () => useContext(AuthContext) as AuthContextType;
