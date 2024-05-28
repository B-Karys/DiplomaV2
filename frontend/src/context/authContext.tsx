import React, { createContext, useState, useEffect, useContext, ReactNode } from 'react';

type AuthContextType = {
    isAuthenticated: boolean;
    setIsAuthenticated: React.Dispatch<React.SetStateAction<boolean>>;
    logout: () => void; // Add logout function type
};

const AuthContext = createContext<AuthContextType | null>(null);

export const AuthProvider = ({ children }: { children: ReactNode }) => {
    const [isAuthenticated, setIsAuthenticated] = useState(false);

    useEffect(() => {
        const storedAuthenticated = getItemWithExpiration('authenticated');
        setIsAuthenticated(storedAuthenticated === 'true');
    }, []);

    const logout = () => {
        fetch("http://localhost:4000/v2/users/logout", {
            method: "POST",
            credentials: "include"
        })
            .then(response => {
                if (response.ok) {
                    setIsAuthenticated(false);
                    localStorage.removeItem("authenticated");
                    window.location.reload();
                } else {
                    throw new Error("Logout failed");
                }
            })
            .catch(error => console.error("Logout error:", error));
    };

    const contextValue: AuthContextType = {
        isAuthenticated,
        setIsAuthenticated,
        logout,
    };

    return (
        <AuthContext.Provider value={contextValue}>
            {children}
        </AuthContext.Provider>
    );
};

export const useAuth = () => useContext(AuthContext) as AuthContextType;

function getItemWithExpiration(key: string) {
    const itemStr = localStorage.getItem(key);
    if (!itemStr) {
        return null;
    }
    const item = JSON.parse(itemStr);
    const now = new Date();
    if (now.getTime() > item.expiry) {
        localStorage.removeItem(key);
        return null;
    }
    return item.value;
}
