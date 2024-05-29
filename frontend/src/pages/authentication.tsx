import { TextInput, PasswordInput, Paper, Title, Text, Container, Group, Button } from '@mantine/core';
import { Link } from 'react-router-dom';
import classes from '../styles/authentication.module.css';
import { SetStateAction, useState} from 'react';

export function Authentication() {
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [errorMessage, setErrorMessage] = useState('');

    const handleEmailChange = (event: { target: { value: SetStateAction<string>; }; }) => {
        setEmail(event.target.value);
        setErrorMessage('');
    };

    const handlePasswordChange = (event: { target: { value: SetStateAction<string>; }; }) => {
        setPassword(event.target.value);
        setErrorMessage('');
    };

    const handleSubmit = () => {
        const userData = {
            email: email,
            password: password
        };

        // Check if password meets the requirement
        if (password.length < 8) {
            setErrorMessage('Password should be at least 8 characters long');
            return;
        }

        // Send the user data via POST request
        fetch('http://localhost:4000/v2/users/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            credentials: 'include', // Include cookies in the request
            body: JSON.stringify(userData)
        })
            .then(async response => {
                if (response.ok) {
                    // Set the authenticated flag in local storage with expiration of 2 minutes
                    setItemWithExpiration('authenticated', 'true', 23.5 * 60 * 60 * 1000);

                    // Redirect to home page after successful login
                    window.location.href = '/';

                    // Parse the response JSON
                    return response.json();
                } else if (response.status === 400) {
                    setErrorMessage('Invalid email or password');
                    throw new Error('Invalid email or password');
                } else if (response.status === 401) {
                    setErrorMessage('The provided password or email are incorrect');
                    throw new Error('Authentication failed');
                } else if (response.status === 403) {
                    const data = await response.json();
                    if (data.message === 'User is not activated') {
                        setErrorMessage('User is not activated');
                    } else {
                        setErrorMessage('User is not activated, please check you mailbox');
                    }
                    throw new Error('User is not activated');
                } else {
                    throw new Error('Unexpected response');
                }
            })
            .catch(error => {
                // Handle errors
                console.error('Error:', error);
            });
    };

    return (
        <Container size={420} my={40}>
            <Title ta="center" className={classes.title}>
                TeamFinder Login
            </Title>
            <Text c="dimmed" size="sm" ta="center" mt={5}>
                Do not have an account yet?{' '}
                <Link to="/register" className={`${classes.createAccount} ${classes.blueText}`}>
                    Create account
                </Link>
            </Text>

            <Paper withBorder shadow="md" p={30} mt={30} radius="md">
                <TextInput label="Email" placeholder="example@mail.ru" value={email} onChange={handleEmailChange} required />
                {errorMessage && <Text className={classes.error}>{errorMessage}</Text>}
                <PasswordInput
                    label="Password"
                    placeholder="Your password"
                    value={password}
                    onChange={handlePasswordChange}
                    required
                    mt="md"
                />
                <Group justify="space-between" mt="lg">
                    <Link to="/forgot-password" className={`${classes.forgotPassword} ${classes.blueText}`}>
                        Forgot password?
                    </Link>
                </Group>
                <Button fullWidth mt="xl" onClick={handleSubmit}>
                    Sign in
                </Button>
            </Paper>
        </Container>
    );
}


export default Authentication;
// Utility functions to set and get items with expiration in local storage
function setItemWithExpiration(key: string, value: string, ttl: number) {
    const now = new Date();
    const item = {
        value: value,
        expiry: now.getTime() + ttl,
    };
    localStorage.setItem(key, JSON.stringify(item));
}

