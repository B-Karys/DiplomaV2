import { useState, useEffect } from 'react';
import axios from 'axios';
import '@mantine/core/styles.css';
import { Link } from 'react-router-dom'; // Import Link for routing
import "../styles/my-posts.module.css"

interface Post {
    id: number;
    createdAt: string;
    name: string;
    description: string;
    type: string;
    skills: string[];
}

interface Metadata {
    current_page: number;
    page_size: number;
    total_records: number;
}

export function MyPosts() {
    const [posts, setPosts] = useState<Post[]>([]);
    const [metadata, setMetadata] = useState<Metadata>({ current_page: 1, page_size: 10, total_records: 0 });
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    const fetchPosts = async (page: number = 1) => { // Default page value to 1
        setLoading(true);
        setError(null);
        try {
            const response = await axios.get<{ posts: Post[]; metadata: Metadata }>(`http://localhost:4000/v2/posts/my?page=${page}`, {
                withCredentials: true  // Include cookies in the request
            });
            setPosts(response.data.posts);
            setMetadata(response.data.metadata);
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
        fetchPosts(metadata.current_page);
    }, [metadata.current_page]);

    const totalPages = Math.ceil(metadata.total_records / metadata.page_size);

    const onPageChange = (page: number) => {
        if (page >= 1 && page <= totalPages) { // Ensure the page number is valid
            setMetadata(prevMetadata => ({ ...prevMetadata, current_page: page }));
        }
    };

    if (loading) return <div>Loading...</div>;
    if (error) return <div>Error loading posts: {error}</div>;

    return (
        <div className="my-myContainer">
            <div className="my-posts-container">
                {posts.map(post => (
                    <div key={post.id} className="my-post">
                        <h2>{post.name}</h2>
                        <p>{post.description}</p>
                        <p>Type: {post.type}</p>
                        <p>Skills: {post.skills.join(', ')}</p>
                        <p>Created At: {new Date(post.createdAt).toLocaleString()}</p>
                        <Link to={`/manage-post/${post.id}`}>
                            <button className="blue-button">Manage Post</button>
                        </Link>
                    </div>
                ))}
            </div>
            {totalPages > 1 && (
                <div className="pagination">
                    {[...Array(totalPages)].map((_, index) => (
                        <button key={index} onClick={() => onPageChange(index + 1)}>{index + 1}</button>
                    ))}
                </div>
            )}
        </div>
    );
}
