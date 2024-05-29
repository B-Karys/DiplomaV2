import { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { Container, Paper, Title, TextInput, Button, Text, PasswordInput } from '@mantine/core';
import axios from 'axios';

export const ResetPassword = () => {
    const [email, setEmail] = useState('');
    const [token, setToken] = useState('');
    const [password, setPassword] = useState('');
    const [confirmPassword, setConfirmPassword] = useState('');
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState(null);
    const [success, setSuccess] = useState(null);
    const [passwordError, setPasswordError] = useState('');

    // Get the token from the URL parameters
    const { token: urlToken } = useParams();

    // Set the token state when the component mounts
    useEffect(() => {
        setToken(urlToken);
    }, [urlToken]);

    const handleSubmitPassword = async (e: { preventDefault: () => void; }) => {
        e.preventDefault();
        setLoading(true);
        setError(null);
        setSuccess(null);
        setPasswordError('');

        // Validate password length
        if (password.length < 8) {
            setPasswordError('Password should be at least 8 characters long');
            setLoading(false);
            return;
        }

        // Validate password and confirmPassword match
        if (password !== confirmPassword) {
            setPasswordError('Passwords do not match');
            setLoading(false);
            return;
        }

        try {
            await axios.post(
                'http://localhost:4000/v2/users/reset-password',
                { email, token, password, confirmPassword }
            );
            setSuccess('Password reset successfully');
        } catch (error) {
            setError(error.response?.data?.message || 'An error occurred');
        } finally {
            setLoading(false);
        }
    };

    return (
        <Container size={420} my={40}>
            <Title ta="center">Reset Password</Title>
            {success ? (
                <Paper withBorder shadow="md" p={30} mt={30} radius="md">
                    <Text size="md">{success}</Text>
                </Paper>
            ) : (
                <>
                    {token && (
                        <Paper withBorder shadow="md" p={30} mt={30} radius="md">
                            <form onSubmit={handleSubmitPassword}>
                                <TextInput
                                    label="Token"
                                    placeholder="Enter the token from your email"
                                    value={token}
                                    onChange={(e) => setToken(e.target.value)}
                                    required
                                    disabled
                                />
                                <PasswordInput
                                    label="New Password"
                                    placeholder="Enter your new password"
                                    value={password}
                                    onChange={(e) => setPassword(e.target.value)}
                                    required
                                    error={passwordError}
                                />
                                <PasswordInput
                                    label="Confirm Password"
                                    placeholder="Confirm your new password"
                                    value={confirmPassword}
                                    onChange={(e) => setConfirmPassword(e.target.value)}
                                    required
                                    error={passwordError}
                                />
                                <Button fullWidth mt="xl" type="submit" disabled={loading}>
                                    {loading ? 'Resetting...' : 'Reset Password'}
                                </Button>
                                {error && <Text c="red" size="sm">{error}</Text>}
                            </form>
                        </Paper>
                    )}
                </>
            )}
        </Container>
    );
}

export default ResetPassword;
