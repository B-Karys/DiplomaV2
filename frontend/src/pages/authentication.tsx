import { TextInput, PasswordInput, Anchor, Paper, Title, Text, Container, Group, Button } from '@mantine/core';
import '@mantine/core/styles.css';
import classes from './authentication.module.css';
import { SetStateAction, useState} from "react";

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
        fetch('http://localhost:4000/v2/tokens/authentication', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(userData)
        })
            .then(response => {
                if (response.ok) {
                    return response.json();
                } else if (response.status === 401) {
                    setErrorMessage('The provided password or email are incorrect');
                    throw new Error('Authentication failed');
                } else {
                    throw new Error('Unexpected response');
                }
            })
            .then(data => {
                // Extract the token from response
                const { authentication_token: token } = data;

                // Set the token as a cookie
                document.cookie = `jwt=${token}; path=/`;
                

                // Redirect to '/profile'
                window.location.href = '/profile';
            })
            .catch(error => {
                // Handle errors
                console.error('Error:', error);
            });
    };

    const handleCreateAccountClick = () => {
        // Redirect to '/register'
        window.location.href = '/register';
    };

    return (
        <Container size={420} my={40}>
            <Title ta="center" className={classes.title}>
                TeamFinder Login
            </Title>
            <Text c="dimmed" size="sm" ta="center" mt={5}>
                Do not have an account yet?{' '}
                <Anchor size="sm" component="button" onClick={handleCreateAccountClick}>
                    Create account
                </Anchor>
            </Text>

            <Paper withBorder shadow="md" p={30} mt={30} radius="md">
                <TextInput label="Email" placeholder="example@mail.ru" value={email} onChange={handleEmailChange} required />
                {errorMessage && <Text c="red" size="sm">{errorMessage}</Text>}
                <PasswordInput label="Password" placeholder="Your password" value={password} onChange={handlePasswordChange} required mt="md" />
                <Group justify="space-between" mt="lg">
                    <Anchor component="button" size="sm">
                        Forgot password?
                    </Anchor>
                </Group>
                <Button fullWidth mt="xl" onClick={handleSubmit}>
                    Sign in
                </Button>
            </Paper>
        </Container>
    );
}
