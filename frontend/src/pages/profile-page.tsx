import { useState, useEffect } from 'react';
import { IconBrandTelegram, IconBrandDiscord } from '@tabler/icons-react';
import axios from 'axios';
import { Link } from 'react-router-dom';
import './profile-page.css';

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

export function ProfilePage() {
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
            const response = await axios.get<Post[]>('http://localhost:4000/v2/posts/my', {
                withCredentials: true,
            });
            setPosts(response.data);
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
                        <div className="contact-item">
                            <IconBrandTelegram size={24} />
                            <p>@{user?.telegram}</p>
                        </div>
                        <div className="contact-item">
                            <IconBrandDiscord size={24} />
                            <p>@{user?.discord}</p>
                        </div>
                    </div>
                </div>
                <div className="profile-main">
                    <h1>{user?.name} ({user?.username}) {user?.surname}</h1>
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
                        <h2>My Posts</h2>
                        {postsLoading && <p>Loading posts...</p>}
                        {postsError && <p>Error loading posts: {postsError}</p>}
                        <div className="posts-container">
                            {posts.map(post => (
                                <div key={post.id} className="post">
                                    <h3>{post.name}</h3>
                                    <p>{post.description}</p>
                                    <p><strong>Type:</strong> {post.type}</p>
                                    <p><strong>Skills:</strong> {post.skills.join(', ')}</p>
                                    <p><strong>Created At:</strong> {new Date(post.createdAt).toLocaleString()}</p>
                                </div>
                            ))}
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
}

export default ProfilePage;
