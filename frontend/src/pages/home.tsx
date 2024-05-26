import { useState, useEffect } from 'react';
import axios from 'axios';
import '@mantine/core/styles.css';
import './home.css'; // Import the CSS file
import Filter from '../components/filter.tsx';

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
    const [search, setSearch] = useState<string>('');

    const fetchPosts = async (type: string, skills: string[], sort: string, search: string) => {
        setLoading(true);
        setError(null);
        try {
            let url = 'http://localhost:4000/v2/posts/';
            const params = new URLSearchParams();
            if (type) params.append('type', type);
            if (skills.length) params.append('skills', skills.join(','));
            if (sort) params.append('sort', sort);
            if (search) params.append('search', search);
            url += `?${params.toString()}`;

            const response = await axios.get<Post[]>(url);
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
        fetchPosts(type, skills, sort, search);
    }, [type, skills, sort, search]);

    const handleFilterChange = (selectedType: string, selectedSkills: string[], selectedSort: string, searchQuery: string) => {
        setType(selectedType);
        setSkills(selectedSkills);
        setSort(selectedSort);
        setSearch(searchQuery);
    };

    if (loading) return <div>Loading...</div>;
    if (error) return <div>Error loading posts: {error}</div>;

    return (
        <div className="container">
            <div className="filter-container">
                <Filter
                    initialType={type}
                    initialSkills={skills}
                    initialSort={sort}
                    initialSearch={search}
                    onFilterChange={handleFilterChange}
                />
            </div>
            <div className="posts-container">
                <div>
                    {posts.map(post => (
                        <div key={post.id} className="post">
                            <h2>{post.name}</h2>
                            <p>{post.description}</p>
                            <p>Author ID: {post.authorId}</p>
                            <p>Type: {post.type}</p>
                            <p>Skills: {post.skills.join(', ')}</p>
                            <p>Created At: {new Date(post.createdAt).toLocaleString()}</p>
                        </div>
                    ))}
                </div>
            </div>
        </div>
    );
}