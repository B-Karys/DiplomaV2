import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { useParams, useNavigate } from 'react-router-dom';
import './manage-posts.css'; // Import the CSS file for styling

interface Post {
    id: number;
    createdAt: string;
    name: string;
    description: string;
    type: string;
    skills: string[];
}

export const ManagePost: React.FC = () => {
    const { id } = useParams<{ id: string }>();
    const navigate = useNavigate(); // Initialize the navigate function
    const [post, setPost] = useState<Post>({
        id: 0,
        createdAt: '',
        name: '',
        description: '',
        type: '',
        skills: [],
    });
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const fetchPost = async () => {
            setLoading(true);
            setError(null);
            try {
                const response = await axios.get<Post>(`http://localhost:4000/v2/posts/${id}`, {
                    withCredentials: true,
                });
                setPost(response.data);
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
        fetchPost();
    }, [id]);

    const handleInputChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement | HTMLTextAreaElement>) => {
        const { name, value } = e.target;
        setPost(prevState => ({
            ...prevState,
            [name]: value,
        }));
    };

    const handleCheckboxChange = (skill: string) => {
        const updatedSkills = post.skills.includes(skill)
            ? post.skills.filter(item => item !== skill)
            : [...post.skills, skill];
        setPost(prevState => ({
            ...prevState,
            skills: updatedSkills,
        }));
    };

    const handleSave = async () => {
        try {
            await axios.patch(`http://localhost:4000/v2/posts/${id}`, post, {
                withCredentials: true,
            });
            // Redirect to the posts page after successful save
            navigate('/posts');
        } catch (error) {
            console.error('Error updating post:', error);
            // Handle error, show error message, etc.
        }
    };

    if (loading) return <div>Loading...</div>;
    if (error) return <div>Error: {error}</div>;

    return (
        <div className="container">
            <h1>Edit Post</h1>
            <form>
                <div className="form-group">
                    <label htmlFor="name">Name:</label>
                    <input type="text" id="name" name="name" value={post.name} onChange={handleInputChange} />
                </div>
                <div className="form-group">
                    <label htmlFor="description">Description:</label>
                    <textarea id="description" name="description" value={post.description} onChange={handleInputChange} />
                </div>
                <div className="form-group">
                    <label htmlFor="type">Type:</label>
                    <select id="type" name="type" value={post.type} onChange={handleInputChange}>
                        <option value="team finding">Team Finding</option>
                        <option value="user finding">User Finding</option>
                    </select>
                </div>
                <div className="form-group">
                    <label>Skills:</label>
                    {['golang', 'python', 'java', 'javascript'].map(skill => (
                        <label key={skill} className="checkbox-label">
                            <input type="checkbox" name="skills" value={skill} checked={post.skills.includes(skill)} onChange={() => handleCheckboxChange(skill)} />
                            {skill}
                        </label>
                    ))}
                </div>
                <div className="form-group">
                    <button type="button" onClick={handleSave}>Save</button>
                </div>
            </form>
        </div>
    );
};

export default ManagePost;
