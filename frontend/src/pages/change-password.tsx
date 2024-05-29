import { useState } from 'react';
import axios from 'axios';
import '../styles/profile-page.css';
import { useNavigate } from 'react-router-dom';

export function ChangePassword() {
    const [currentPassword, setCurrentPassword] = useState('');
    const [newPassword, setNewPassword] = useState('');
    const [repeatNewPass, setRepeatNewPass] = useState('');
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [success, setSuccess] = useState<string | null>(null);
    const [showCurrentPassword, setShowCurrentPassword] = useState(false);
    const [showNewPassword, setShowNewPassword] = useState(false);
    const [showRepeatNewPass, setShowRepeatNewPass] = useState(false);
    const navigate = useNavigate();

    const handleChangePassword = async (e: React.FormEvent) => {
        e.preventDefault();
        setLoading(true);
        setError(null);
        setSuccess(null);

        if (newPassword !== repeatNewPass) {
            setError('New passwords do not match');
            setLoading(false);
            return;
        }

        try {
            const response = await axios.patch(
                'http://localhost:4000/v2/users/password',
                {
                    currentPassword,
                    newPassword,
                    repeatNewPass
                },
                {
                    withCredentials: true, // Include cookies in the request
                }
            );
            setSuccess('Password changed successfully');
            navigate("/");
        } catch (error) {
            if (axios.isAxiosError(error)) {
                setError(error.response?.data?.message || 'An error occurred');
            } else {
                setError('An unknown error occurred');
            }
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="profile-container">
            <div className="profile-content">
                <form onSubmit={handleChangePassword} className="change-password-form">
                    <h2>Change Password</h2>
                    {error && <p className="error">{error}</p>}
                    {success && <p className="success">{success}</p>}
                    <div className="form-group">
                        <label htmlFor="currentPassword">Current Password</label>
                        <div className="password-input-container">
                            <input
                                type={showCurrentPassword ? "text" : "password"}
                                id="currentPassword"
                                value={currentPassword}
                                onChange={(e) => setCurrentPassword(e.target.value)}
                                required
                            />
                            <label>
                                <input
                                    type="checkbox"
                                    checked={showCurrentPassword}
                                    onChange={() => setShowCurrentPassword(!showCurrentPassword)}
                                /> Show Password
                            </label>
                        </div>
                    </div>
                    <div className="form-group">
                        <label htmlFor="newPassword">New Password</label>
                        <div className="password-input-container">
                            <input
                                type={showNewPassword ? "text" : "password"}
                                id="newPassword"
                                value={newPassword}
                                onChange={(e) => setNewPassword(e.target.value)}
                                required
                            />
                            <label>
                                <input
                                    type="checkbox"
                                    checked={showNewPassword}
                                    onChange={() => setShowNewPassword(!showNewPassword)}
                                /> Show Password
                            </label>
                        </div>
                    </div>
                    <div className="form-group">
                        <label htmlFor="repeatNewPass">Repeat New Password</label>
                        <div className="password-input-container">
                            <input
                                type={showRepeatNewPass ? "text" : "password"}
                                id="repeatNewPass"
                                value={repeatNewPass}
                                onChange={(e) => setRepeatNewPass(e.target.value)}
                                required
                            />
                            <label>
                                <input
                                    type="checkbox"
                                    checked={showRepeatNewPass}
                                    onChange={() => setShowRepeatNewPass(!showRepeatNewPass)}
                                /> Show Password
                            </label>
                        </div>
                    </div>
                    <button type="submit" className="blue-button" disabled={loading}>
                        {loading ? 'Changing...' : 'Change Password'}
                    </button>
                </form>
            </div>
        </div>
    );
}

export default ChangePassword;
