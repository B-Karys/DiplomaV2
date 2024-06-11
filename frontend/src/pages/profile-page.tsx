import { useState, useEffect } from 'react';
import { IconBrandTelegram, IconBrandDiscord } from '@tabler/icons-react';
import axios from 'axios';
import '../styles/profile-page.css';

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

interface ProfilePageProps {
    userId: number; // Prop to accept user ID
}

export function ProfilePage({ userId }: ProfilePageProps) {
    const [user, setUser] = useState<User | null>(null);
    const [posts, setPosts] = useState<Post[]>([]); // Initialize posts state as an empty array
    const [loading, setLoading] = useState<boolean>(true);
    const [error, setError] = useState<string | null>(null);
    const [postsLoading, setPostsLoading] = useState<boolean>(true);
    const [postsError, setPostsError] = useState<string | null>(null);

    useEffect(() => {
        fetchUserData();
        fetchUserPosts();
    }, [userId]); // Fetch data whenever userId changes

    const fetchUserData = async () => {
        setLoading(true);
        setError(null);
        try {
            const response = await fetch(`http://localhost:4000/v2/users/${userId}`, {
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
            setUser(userData); // Set user data
            setLoading(false);
        } catch (error) {
            setLoading(false);
            setError((error as Error).message);
            console.error('Error fetching user data:', error);
        }
    };

    const fetchUserPosts = async () => {
        setPostsLoading(true);
        setPostsError(null);
        try {
            const response = await axios.get<Post[]>(`http://localhost:4000/v2/posts/?author=${userId}`, {
                withCredentials: true,
            });
            // @ts-ignore
            setPosts(response.data.posts);
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
                        {posts && Array.isArray(posts) && posts.length > 0 ? (
                            <div className="profile-posts-container">
                                {posts.map(post => (
                                    <div key={post.id} className="profile-post">
                                        <h3>{post.name}</h3>
                                        <p>{post.description}</p>
                                        <p><strong>Type:</strong> {post.type}</p>
                                        <p><strong>Skills:</strong> {post.skills.join(', ')}</p>
                                        <p><strong>Created At:</strong> {new Date(post.createdAt).toLocaleString()}</p>
                                    </div>
                                ))}
                            </div>
                        ) : (
                            <p>No posts available.</p>
                        )}
                    </div>
                </div>
            </div>
        </div>
    );
}

export default ProfilePage;
