import '@mantine/core/styles.css';
import { MantineProvider } from '@mantine/core';
import { BrowserRouter, Routes, Route } from "react-router-dom";
import './App.css';
import { Authentication } from "./pages/authentication.tsx";
import { Registration } from "./pages/registration.tsx";
import { Home } from "./pages/home.tsx";
import ProfilePage from "./pages/profilepage.tsx"; // Import ProfilePage component

export default function App() {
    return (
        <MantineProvider defaultColorScheme="dark">
            <BrowserRouter>
                <Routes>
                    <Route path="/" element={<Home />} />
                    <Route path="/login" element={<Authentication />} />
                    <Route path="/register" element={<Registration />} />
                    <Route path="/profile" element={<ProfilePage />} /> {/* Add route for profile page */}
                </Routes>
            </BrowserRouter>
        </MantineProvider>
    );
}
