import { useState, useEffect } from 'react';
import axios from 'axios';
import { Link } from 'react-router-dom'; // Import Link for routing
import Filter from '../components/filter.tsx';
import '@mantine/core/styles.css';
import '../styles/home.css'; // Import the CSS file

interface Post {
    id: number;
    createdAt: string;
    name: string;
    description: string;
    authorId: number;
    type: string;
    skills: string[];
}

export function Home() {
    const [posts, setPosts] = useState<Post[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [type, setType] = useState<string>('');
    const [skills, setSkills] = useState<string[]>([]);
    const [sort, setSort] = useState<string>('created_at');

    const fetchPosts = async (type: string, skills: string[], sort: string) => {
        setLoading(true);
        setError(null);
        try {
            const params = new URLSearchParams();
            if (type) params.append('type', type);
            if (skills.length) params.append('skills', skills.join(','));
            if (sort) params.append('sort', sort);

            const response = await axios.get<Post[]>(`http://localhost:4000/v2/posts/?${params.toString()}`);
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
        fetchPosts(type, skills, sort);
    }, [type, skills, sort]);

    const handleFilterChange = (selectedType: string, selectedSkills: string[], selectedSort: string) => {
        setType(selectedType);
        setSkills(selectedSkills);
        setSort(selectedSort);
    };

    if (loading) return <div>Loading...</div>;
    if (error) return <div>Error loading posts: {error}</div>;

    return (
        <div className="home-myContainer">
            <div className="filter-container">
                <Filter
                    initialType={type}
                    initialSkills={skills}
                    initialSort={sort}
                    onFilterChange={handleFilterChange}
                />
            </div>
            <div className="home-posts-container">
                {posts.map(post => (
                    <div key={post.id} className="home-post">
                        <h2>{post.name}</h2>
                        <p>{post.description}</p>
                        <p>
                            <Link to={`/profile/${post.authorId}`}>Author</Link>
                        </p>
                        <p>Type: {post.type}</p>
                        <p>Skills: {post.skills.join(', ')}</p>
                        <p>Created At: {new Date(post.createdAt).toLocaleDateString()}</p>
                    </div>
                ))}
            </div>
        </div>
    );
}
