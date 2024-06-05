import { useState, useEffect } from 'react';
import { IconBrandTelegram, IconBrandDiscord } from '@tabler/icons-react';
import axios from 'axios';
import '../styles/profile-page.module.css';

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
    profileImage: string;
}

interface Post {
    id: number;
    createdAt: string;
    name: string;
    description: string;
    type: string;
    skills: string[];
}

export function MyProfilePage() {
    const [user, setUser] = useState<User | null>(null);
    const [posts, setPosts] = useState<Post[]>([]);
    const [loading, setLoading] = useState<boolean>(true);
    const [error, setError] = useState<string | null>(null);
    const [postsLoading, setPostsLoading] = useState<boolean>(true);
    const [postsError, setPostsError] = useState<string | null>(null);

    useEffect(() => {
        fetchUserData();
        fetchUserPosts();
    }, []);

    const fetchUserData = async () => {
        try {
            const response = await fetch('http://localhost:4000/v2/users/my', {
                method: 'GET',
                headers: {
                    'Content-Type': 'application/json',
                },
                credentials: 'include',
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

    const fetchUserPosts = async () => {
        setPostsLoading(true);
        setPostsError(null);
        try {
            const response = await axios.get('http://localhost:4000/v2/posts/my', {
                withCredentials: true,
            });
            setPosts(response.data.posts || []); // Ensure posts is an array
            setPostsLoading(false);
        } catch (error) {
            setPostsLoading(false);
            if (axios.isAxiosError(error)) {
                setPostsError(error.message);
            } else {
                setPostsError('An unknown error occurred');
            }
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
                            <img src={user.profileImage} alt="Profile" />
                        ) : (
                            'Profile Photo Of User'
                        )}
                    </div>
                    <div className="profile-contacts">
                        <div className="profile-contact-item">
                            <IconBrandTelegram size={24} />
                            <p>@{user?.telegram}</p>
                        </div>
                        <div className="profile-contact-item">
                            <IconBrandDiscord size={24} />
                            <p>@{user?.discord}</p>
                        </div>
                    </div>
                </div>
                <div className="profile-main">
                    <h2>{user?.name} ({user?.username}) {user?.surname}</h2>
                    <div className="profile-details">
                        <p><strong>Username:</strong> {user?.username}</p>
                        <p><strong>Surname:</strong> {user?.surname}</p>
                        <p><strong>Name:</strong> {user?.name}</p>
                    </div>
                    <div className="profile-skills">
                        <p><strong>Skills:</strong></p>
                        {user?.skills?.map((skill, index) => (
                            <p key={index}>{skill}</p>
                        ))}
                    </div>
                    <div className="profile-posts">
                        <h2>Author's Posts</h2>
                        {postsLoading && <p>Loading posts...</p>}
                        {postsError && <p>Error loading posts: {postsError}</p>}
                        <div className="profile-posts-container">
                            {Array.isArray(posts) && posts.length > 0 ? (
                                posts.map(post => (
                                    <div key={post.id} className="profile-post">
                                        <h3>{post.name}</h3>
                                        <p>{post.description}</p>
                                        <p><strong>Type:</strong> {post.type}</p>
                                        <p><strong>Skills:</strong> {post.skills.join(', ')}</p>
                                        <p><strong>Created At:</strong> {new Date(post.createdAt).toLocaleString()}</p>
                                    </div>
                                ))
                            ) : (
                                <p>No posts available.</p>
                            )}
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
}

export default MyProfilePage;
