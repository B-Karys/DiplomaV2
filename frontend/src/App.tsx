import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { MantineProvider } from '@mantine/core';
import '@mantine/core/styles.css';
import { Authentication } from './pages/authentication.tsx';
import { Registration } from './pages/registration.tsx';
import { Home } from './pages/home.tsx';
import  { ProfilePage } from './pages/profilepage.tsx'; // Import ProfilePage component
import { AuthProvider, useAuth } from './context/authContext.tsx';
import { Navbar } from "./components/navbar.tsx";
import "./components/navbar.module.css"
import {MyPosts} from "./pages/myPosts.tsx";
import ManagePost from "./pages/manage-post.tsx";

export default function App() {
    return (
        <MantineProvider defaultColorScheme="dark">
            <AuthProvider>
                <Router>
                    <Navbar />
                    <Routes>
                        <Route path="/" element={<HomeRoute />} />
                        <Route path="/login" element={<AuthenticationRoute />} />
                        <Route path="/register" element={<Registration />} />
                        <Route path="/profile" element={<ProfileRoute />} />
                        <Route path="/posts" element={<MyPostsRoute />} />
                        <Route path="/manage-post/:id" element={<ManagePostsRoute />} />
                    </Routes>
                </Router>
            </AuthProvider>
        </MantineProvider>
    );
}

// Define routes with authentication checks
const HomeRoute = () => {
    return <Home />;
};

const AuthenticationRoute = () => {
    const { isAuthenticated } = useAuth();
    return isAuthenticated ? <Navigate to="/" /> : <Authentication />;
};

const ProfileRoute = () => {
    const { isAuthenticated } = useAuth();
    return isAuthenticated ? <ProfilePage /> : <Navigate to="/login" />;
};

const MyPostsRoute = () => {
    return <MyPosts />;
}

const ManagePostsRoute = () => {
    return <ManagePost />;
}
