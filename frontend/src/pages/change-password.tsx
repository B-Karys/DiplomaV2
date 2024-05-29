import { Container, Paper, TextInput, PasswordInput, Title, Text, Button } from '@mantine/core';
import { useState } from 'react';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';
import classes from "../styles/authentication.module.css"
export function ChangePassword() {
    const [currentPassword, setCurrentPassword] = useState('');
    const [newPassword, setNewPassword] = useState('');
    const [repeatNewPass, setRepeatNewPass] = useState('');
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [success, setSuccess] = useState<string | null>(null);
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
        <Container size={420} my={40}>
            <Title ta="center" className={classes.cpTitle}>Change Password</Title>
            <Paper withBorder shadow="md" p={30} mt={30} radius="md">
                <form onSubmit={handleChangePassword}>
                    <TextInput
                        type={"password"}
                        label="Current Password"
                        value={currentPassword}
                        onChange={(e) => setCurrentPassword(e.target.value)}
                        required
                    />
                    <PasswordInput
                        label="New Password"
                        value={newPassword}
                        onChange={(e) => setNewPassword(e.target.value)}
                        required
                    />
                    <PasswordInput
                        label="Repeat New Password"
                        value={repeatNewPass}
                        onChange={(e) => setRepeatNewPass(e.target.value)}
                        required
                    />
                    {error && <Text color="red" mt={2}>{error}</Text>}
                    {success && <Text color="green" mt={2}>{success}</Text>}
                    <Button
                        type="submit"
                        variant="outline"
                        loading={loading}
                        mt={4}
                        fullWidth
                    >
                        {loading ? 'Changing...' : 'Change Password'}
                    </Button>
                </form>
            </Paper>
        </Container>
    );
}

export default ChangePassword;
