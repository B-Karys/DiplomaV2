import {BrowserRouter as Router, Routes, Route, Navigate, useParams} from 'react-router-dom';
import { MantineProvider } from '@mantine/core';
import { Authentication } from './pages/authentication';
import { Registration } from './pages/registration';
import { Home } from './pages/home';
import { MyProfilePage } from './pages/my-profile-page.tsx';
import { AuthProvider, useAuth } from './context/authContext';
import { Navbar } from "./components/navbar";
import { MyPosts } from "./pages/my-posts.tsx";
import { ManagePost } from "./pages/manage-post";
import  { CreatePost } from "./pages/create-post";

import '@mantine/core/styles.css';
import "./styles/navbar.module.css";
import  { ManageProfile } from "./pages/manage-profile.tsx";
import  { ChangePassword } from "./pages/change-password.tsx";
import  { ForgotPassword } from "./pages/forgot-password.tsx";
import  { ResetPassword } from "./pages/reset-password.tsx";
import  { ProfilePage } from "./pages/profile-page.tsx";

export default function App() {
    return (
        <MantineProvider defaultColorScheme="dark">
            <AuthProvider>
                <Router>
                    <Navbar />
                    <Routes>
                        <Route path="/" element={<Home />} />
                        <Route path="/login" element={<AuthenticationRoute />} />
                        <Route path="/register" element={<RegistrationRoute />} />
                        <Route path="/profile/my" element={<MyProfileRoute />} />
                        <Route path="/profile/settings" element={<ManageProfileRoute />} />
                        <Route path="/posts" element={<MyPostsRoute />} />
                        <Route path="/create-post" element={<CreatePostsRoute />} />
                        <Route path="/manage-post/:id" element={<ManagePostsRoute />} />
                        <Route path="/profile/change-password" element={<ChangePasswordRoute />} />
                        <Route path="/forgot-password" element={<ForgotPasswordRoute />} />
                        <Route path="/reset-password/:token" element={<ResetPasswordRoute />} />
                        <Route path="/profile/:id" element={<ProfilePageRoute />} /> {/* Update route path */}
                    </Routes>
                </Router>
            </AuthProvider>
        </MantineProvider>
    );
}

// Custom route components for route protection

const AuthenticationRoute = () => {
    const { isAuthenticated } = useAuth();
    return isAuthenticated ? <Navigate to="/" /> : <Authentication />;
};

const RegistrationRoute = () => {
    const { isAuthenticated } = useAuth();
    return isAuthenticated ? <Navigate to="/" /> : <Registration />;
};

const MyPostsRoute = () => {
    const { isAuthenticated } = useAuth();
    return isAuthenticated ? <MyPosts /> : <Navigate to="/login" />;
};

const ProfilePageRoute = () => {
    const { isAuthenticated } = useAuth();
    const { id } = useParams(); // Get the user ID from URL parameters
    // @ts-ignore
    return isAuthenticated ? <ProfilePage userId={id} /> : <Navigate to="/login" />;
}

const CreatePostsRoute = () => {
    const { isAuthenticated } = useAuth();
    return isAuthenticated ? <CreatePost /> : <Navigate to="/login" />;
};

const ManagePostsRoute = () => {
    const { isAuthenticated } = useAuth();
    return isAuthenticated ? <ManagePost /> : <Navigate to="/login" />;
};

const MyProfileRoute = () => {
    const { isAuthenticated } = useAuth();
    return isAuthenticated ? <MyProfilePage /> : <Navigate to="/login" />;
};

const ManageProfileRoute = () => {
    const { isAuthenticated } = useAuth();
    return isAuthenticated ? <ManageProfile /> : <Navigate to="/login" />;
};

const ChangePasswordRoute = () => {
    const { isAuthenticated } = useAuth();
    return isAuthenticated ? <ChangePassword /> : <Navigate to="/login"/>;
};

const ForgotPasswordRoute = () => {
    return <ForgotPassword />
}

const ResetPasswordRoute = () => {
    return <ResetPassword />
}
