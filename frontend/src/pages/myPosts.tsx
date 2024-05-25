import { useState, useEffect } from 'react';
import axios from 'axios';
import '@mantine/core/styles.css';
import './home.css'; // Import the CSS file
import { Link } from 'react-router-dom'; // Import Link for routing

interface Post {
    id: number;
    createdAt: string;
    name: string;
    description: string;
    type: string;
    skills: string[];
}

export function MyPosts() {
    const [posts, setPosts] = useState<Post[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    const fetchPosts = async () => {
        setLoading(true);
        setError(null);
        try {
            const response = await axios.get<Post[]>('http://localhost:4000/v2/posts/my', {
                withCredentials: true  // Include cookies in the request
            });
            setPosts(response.data);
        } catch (error) {
            if (axios.isAxiosError(error)) {
                setError(error.message);
            } else {
                setError('An unknown error occurred');
            }
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchPosts();
    }, []);

    if (loading) return <div>Loading...</div>;
    if (error) return <div>Error loading posts: {error}</div>;

    return (
        <div className="container">
            {/*<div className="filter-container">*/}
            {/*    <Filter />*/}
            {/*</div>*/}
            <div className="posts-container">
                <h1>This is My Posts page</h1>
                <div>
                    {posts.map(post => (
                        <div key={post.id} className="post">
                            <h2>{post.name}</h2>
                            <p>{post.description}</p>
                            <p>Type: {post.type}</p>
                            <p>Skills: {post.skills.join(', ')}</p>
                            <p>Created At: {new Date(post.createdAt).toLocaleString()}</p>
                            <Link to={`/manage-post/${post.id}`}>
                                <button>Manage Post</button>
                            </Link>
                        </div>
                    ))}
                </div>
            </div>
        </div>
    );
}
