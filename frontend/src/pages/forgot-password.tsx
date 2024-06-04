import { TextInput, Paper, Title, Container, Button, Text } from '@mantine/core';
import { SetStateAction, useState} from 'react';
import axios from "axios";

export const ForgotPassword = () => {
    const [email, setEmail] = useState('');
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState(null);
    const [success, setSuccess] = useState(null);

    const handleEmailChange = (event: { target: { value: SetStateAction<string>; }; }) => {
        setEmail(event.target.value);
    };

    const handleSubmit = async (e: { preventDefault: () => void; }) => {
        e.preventDefault();
        setLoading(true);
        setError(null);
        setSuccess(null);

        try {
            await axios.post(
                'http://localhost:4000/v2/users/forgot-password',
                { email }
            );
// @ts-ignore
            setSuccess('Password reset email sent successfully');
        } catch (error) {
            // @ts-ignore
            setError(error.response?.data?.message || 'An error occurred');
        } finally {
            setLoading(false);
        }
    };

    return (
        <Container size={420} my={40}>
            <Title ta="center">
                TeamFinder Forgot Password
            </Title>
            <Paper withBorder shadow="md" p={30} mt={30} radius="md">
                <TextInput label="Email" placeholder="example@mail.ru" value={email} onChange={handleEmailChange} required />
                <Button fullWidth mt="xl" onClick={handleSubmit} disabled={loading}>
                    {loading ? 'Sending...' : 'Send request'}
                </Button>
                {error && <Text c="red" size="sm">{error}</Text>}
                {success && <Text c="green" size="sm">{success}</Text>}
            </Paper>
        </Container>
    );
}

export default ForgotPassword;
