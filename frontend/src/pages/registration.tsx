import { useState } from 'react';
import { TextInput, PasswordInput, Paper, Title, Text, Container, Group, Button } from '@mantine/core';
import { Link } from 'react-router-dom';
import classes from '../styles/authentication.module.css';

export function Registration() {
    const [email, setEmail] = useState('');
    const [name, setName] = useState('');
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const [emailError, setEmailError] = useState('');
    const [passwordError, setPasswordError] = useState('');
    const [generalError, setGeneralError] = useState('');

    const handleEmailChange = (event) => {
        setEmail(event.target.value);
        setEmailError('');
        setGeneralError('');
    };

    const handleNameChange = (event) => {
        setName(event.target.value);
        setGeneralError('');
    };

    const handleUsernameChange = (event) => {
        setUsername(event.target.value);
        setGeneralError('');
    };

    const handlePasswordChange = (event) => {
        setPassword(event.target.value);
        setPasswordError('');
        setGeneralError('');
    };

    const validateEmail = (email) => {
        const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
        return emailRegex.test(email);
    };

    const handleSubmit = () => {
        // Your form submission logic here
    };

    const handleLoginClick = () => {
        // Redirect to '/login'
        window.location.href = '/login';
    };

    return (
        <Container size={420} my={40}>
            <Title ta="center" className={classes.title}>
                TeamFinder Registration
            </Title>
            <Text c="dimmed" size="sm" ta="center" mt={5}>
                Already have an account?{' '}
                <Link to="/login" className={`${classes.createAccount} ${classes.blueText}`}>
                    Login
                </Link>
            </Text>

            <Paper withBorder shadow="md" p={30} mt={30} radius="md">
                <TextInput
                    label="Email"
                    placeholder="example@mail.ru"
                    value={email}
                    onChange={handleEmailChange}
                    required
                    error={emailError}
                />
                <TextInput
                    label="Name"
                    placeholder="Bekarys"
                    value={name}
                    onChange={handleNameChange}
                    required
                />
                <TextInput
                    label="Username"
                    placeholder="b.karys"
                    value={username}
                    onChange={handleUsernameChange}
                    required
                />
                <PasswordInput
                    label="Password"
                    placeholder="Your password"
                    value={password}
                    onChange={handlePasswordChange}
                    required
                    mt="md"
                    error={passwordError}
                />
                {generalError && <Text c="red" size="sm">{generalError}</Text>}
                <Group justify="space-between" mt="lg">
                    <Link to="/forgot-password" className={`${classes.forgotPassword} ${classes.blueText}`}>
                        Forgot password?
                    </Link>
                </Group>
                <Button fullWidth mt="xl" onClick={handleSubmit}>
                    Sign up
                </Button>
            </Paper>
        </Container>
    );
}
