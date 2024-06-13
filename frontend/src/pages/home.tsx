import { useState, useEffect } from 'react';
import axios from 'axios';
import { Link } from 'react-router-dom';
import Filter from '../components/filter.tsx';
import '@mantine/core/styles.css';
import '../styles/home.css';

interface Post {
    id: number;
    createdAt: string;
    name: string;
    description: string;
    authorId: number;
    type: string;
    skills: string[];
}

interface Metadata {
    current_page: number;
    page_size: number;
    total_records: number;
}

export function Home() {
    const [posts, setPosts] = useState<Post[]>([]);
    const [metadata, setMetadata] = useState<Metadata>({ current_page: 1, page_size: 10, total_records: 0 });
    const [loading, setLoading] = useState(true);
    const [, setError] = useState<string | null>(null);
    const [type, setType] = useState<string>('');
    const [skills, setSkills] = useState<string[]>([]);
    const [sort, setSort] = useState<string>('created_at');

    const fetchPosts = async (type: string, skills: string[], sort: string, page: number) => {
        setLoading(true);
        setError(null);
        try {
            const params = new URLSearchParams();
            if (type) params.append('type', type);
            if (skills.length) params.append('skills', skills.join(','));
            if (sort) params.append('sort', sort);
            params.append('page', page.toString());

            const response = await axios.get<{ posts: Post[]; metadata: Metadata }>(`http://localhost:4000/v2/posts/?${params.toString()}`);

            const { posts, metadata } = response.data;

            // Handle case where posts array is empty
            if (posts.length === 0) {
                setPosts([]);
                setMetadata({current_page: 1, page_size: 10, total_records: 0});
            } else {
                setPosts(posts);
                setMetadata(metadata);
            }

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
        fetchPosts(type, skills, sort, metadata.current_page);
    }, [type, skills, sort, metadata.current_page]);

    const handleFilterChange = (selectedType: string, selectedSkills: string[], selectedSort: string) => {
        setType(selectedType);
        setSkills(selectedSkills);
        setSort(selectedSort);
    };

    const totalPages = Math.ceil(metadata.total_records / metadata.page_size);

    const onPageChange = (page: number) => {
        setMetadata(prevMetadata => ({ ...prevMetadata, current_page: page }));
    };

    if (loading) return <div>Loading...</div>;

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
                {posts.length === 0 ? (
                    <div>There are no posts matching the current filters.</div>
                ) : (
                    posts.map(post => (
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
                    ))
                )}
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
