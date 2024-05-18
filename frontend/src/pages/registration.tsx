import {
    TextInput,
    PasswordInput,
    Anchor,
    Paper,
    Title,
    Text,
    Container,
    Group,
    Button,
} from '@mantine/core';
import '@mantine/core/styles.css';
import classes from './authentication.module.css';
import {SetStateAction, useState} from "react";

export function Registration() {
    const [email, setEmail] = useState('');
    const [name, setName] = useState('');
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const [errorMessage, setErrorMessage] = useState('');

    const handleEmailChange = (event: { target: { value: SetStateAction<string>; }; }) => {
        setEmail(event.target.value);
        setErrorMessage('');
    };

    const handleNameChange = (event: { target: { value: SetStateAction<string>; }; }) => {
        setName(event.target.value);
        setErrorMessage('');
    };

    const handleUsernameChange = (event: { target: { value: SetStateAction<string>; }; }) => {
        setUsername(event.target.value);
        setErrorMessage('');
    };

    const handlePasswordChange = (event: { target: { value: SetStateAction<string>; }; }) => {
        setPassword(event.target.value);
        setErrorMessage('');
    };

    const handleSubmit = () => {
        const userData = {
            email: email,
            name: name,
            username: username,
            password: password
        };

        // Send the user data via POST request
        fetch('http://localhost:4000/v1/auth', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(userData)
        })
            .then(response => {
                if (response.status === 202) {
                    // Redirect to '/'
                    window.location.href = '/login';
                } else if (response.status === 422) {
                    setErrorMessage('Email is already registered');
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
                TeamFinder Registration
            </Title>
            <Text c="dimmed" size="sm" ta="center" mt={5}>
                Already have an account?{' '}
                <Anchor size="sm" component="button">
                    Login
                </Anchor>
            </Text>

            <Paper withBorder shadow="md" p={30} mt={30} radius="md">
                <TextInput label="Email" placeholder="example@mail.ru" value={email} onChange={handleEmailChange} required />
                {errorMessage && <Text c="red" size="sm">{errorMessage}</Text>}
                <TextInput label="Name" placeholder="Bekarys" value={name} onChange={handleNameChange} required />
                <TextInput label="Username" placeholder="b.karys" value={username} onChange={handleUsernameChange} required />
                <PasswordInput label="Password" placeholder="Your password" value={password} onChange={handlePasswordChange} required mt="md" />
                <Group justify="space-between" mt="lg">
                    <Anchor component="button" size="sm">
                    </Anchor>
                </Group>
                <Button fullWidth mt="xl" onClick={handleSubmit}>
                    Sign in
                </Button>
            </Paper>
        </Container>
    );
}