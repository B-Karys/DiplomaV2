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
        setIsAuthenticated(storedAuthenticated);
    }, []);

    const contextValue: AuthContextType = {
        isAuthenticated,
        setIsAuthenticated,
    };

    return (
        <AuthContext.Provider value={contextValue}>
            {children}
        </AuthContext.Provider>
    );
};

export const useAuth = () => useContext(AuthContext) as AuthContextType;
