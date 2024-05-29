import { SetStateAction, useState} from 'react';
import {TextInput, PasswordInput, Paper, Title, Text, Container, Group, Button} from '@mantine/core';
import {Link} from 'react-router-dom';
import classes from '../styles/authentication.module.css';

export function Registration() {
    const [email, setEmail] = useState('');
    const [name, setName] = useState('');
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const [emailError, setEmailError] = useState('');
    const [passwordError, setPasswordError] = useState('');
    const [generalError, setGeneralError] = useState('');

    const handleEmailChange = (event: { target: { value: SetStateAction<string>; }; }) => {
        setEmail(event.target.value);
        setEmailError('');
        setGeneralError('');
    };

    const handleNameChange = (event: { target: { value: SetStateAction<string>; }; }) => {
        setName(event.target.value);
        setGeneralError('');
    };

    const handleUsernameChange = (event: { target: { value: SetStateAction<string>; }; }) => {
        setUsername(event.target.value);
        setGeneralError('');
    };

    const handlePasswordChange = (event: { target: { value: SetStateAction<string>; }; }) => {
        setPassword(event.target.value);
        setPasswordError('');
        setGeneralError('');
    };

    const validateEmail = (email: string) => {
        const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
        return emailRegex.test(email);
    };

    const handleSubmit = () => {
        let valid = true;

        if (!validateEmail(email)) {
            setEmailError('Please enter a valid email');
            valid = false;
        }

        if (password.length < 8) {
            setPasswordError('Password should be at least 8 characters long');
            valid = false;
        }

        if (!valid) {
            return;
        }

        const userData = {
            email: email,
            name: name,
            username: username,
            password: password
        };

        // Send the user data via POST request
        fetch('http://localhost:4000/v2/users/registration', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(userData)
        })
            .then(async response => {
                if (response.status === 201) {
                    // Successful registration
                    window.location.href = '/login'; // Redirect to login
                    return;
                }
                const data = await response.json();
                if (response.status === 400) {
                    if (data.error && data.error.includes('duplicate key value violates unique constraint "uni_users_email"')) {
                        setEmailError('Email is already registered');
                    } else if (data.error && data.error.includes('duplicate key value violates unique constraint "uni_users_username"')) {
                        setGeneralError('Username is already taken');
                    } else if (data.error && data.error.includes('Validation error')) {
                        setEmailError('Use existing email');
                        setPasswordError('Password should contain more than 8 characters');
                    } else {
                        setGeneralError('Unexpected response');
                    }
                } else {
                    setGeneralError('Unexpected response');
                }
            })
            .catch(error => {
                // Handle errors
                console.error('Error:', error);
                setGeneralError('Unexpected error occurred');
            });
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
