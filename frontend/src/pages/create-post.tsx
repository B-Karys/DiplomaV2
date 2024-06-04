import { useState } from 'react';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';
import '../styles/manage-posts.module.css'; // Import the CSS file for styling

export const CreatePost = () => {
    const navigate = useNavigate(); // Initialize the navigate function
    const [post, setPost] = useState({
        name: '',
        description: '',
        type: 'team finding',
        skills: [] as string[], // Initialize skills as an array of strings
    });
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const handleInputChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>) => {
        const { name, value } = e.target;
        setPost((prevState) => ({
            ...prevState,
            [name]: value,
        }));
    };

    const handleCheckboxChange = (skill: string) => {
        const updatedSkills = post.skills.includes(skill)
            ? post.skills.filter((item) => item !== skill)
            : [...post.skills, skill];
        setPost((prevState) => ({
            ...prevState,
            skills: updatedSkills,
        }));
    };

    const handleSave = async () => {
        // Client-side validation
        if (!post.name || !post.description) {
            setError('Name and description are required.');
            return;
        }

        setLoading(true);
        try {
            await axios.post('http://localhost:4000/v2/posts/', post, {
                withCredentials: true,
            });
            // Redirect to the posts page after successful save
            navigate('/posts');
        } catch (error) {
            console.error('Error creating post:', error);
            setError('Error creating post');
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="container">
            <h1>Create New Post</h1>
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
                    {['golang', 'python', 'java', 'javascript', 'c++','c#','rust','php','kotlin','ruby'].map((skill) => (
                        <label key={skill} className="checkbox-label">
                            <input
                                type="checkbox"
                                name="skills"
                                value={skill}
                                checked={post.skills.includes(skill)}
                                onChange={() => handleCheckboxChange(skill)}
                            />
                            {skill}
                        </label>
                    ))}
                </div>
                <div className="form-group">
                    <button type="button" className="blue-button" onClick={handleSave} disabled={loading}>
                        {loading ? 'Saving...' : 'Save'}
                    </button>
                </div>
                {error && <div className="error-message">{error}</div>}
            </form>
        </div>
    );
};

export default CreatePost;
