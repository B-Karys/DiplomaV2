import React, { createContext, useState, useEffect, useContext, ReactNode } from 'react';

type AuthContextType = {
    isAuthenticated: boolean;
    setIsAuthenticated: React.Dispatch<React.SetStateAction<boolean>>;
};

const AuthContext = createContext<AuthContextType | null>(null);

export const AuthProvider = ({ children }: { children: ReactNode }) => {
    const [isAuthenticated, setIsAuthenticated] = useState(false);

    useEffect(() => {
        const storedAuthenticated = localStorage.getItem('authenticated') === 'true';
        if (storedAuthenticated) {
            setIsAuthenticated(true);
            return;
        }

        fetch('http://localhost:4000/v2/users/check-auth', {
            method: 'GET',
            credentials: 'include',
        })
            .then(response => response.json())
            .then(data => {
                if (data.authenticated === "true") {
                    setIsAuthenticated(true);
                    localStorage.setItem('authenticated', 'true');
                } else {
                    setIsAuthenticated(false);
                    localStorage.removeItem('authenticated');
                }
            })
            .catch(error => {
                console.error('Error:', error);
                setIsAuthenticated(false);
            });
    }, []);

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
