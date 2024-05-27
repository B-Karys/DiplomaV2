    import { useEffect, useState } from 'react';
    import { Link } from 'react-router-dom';
    import { Button, Group, rem, useMantineTheme, Burger, Drawer, ScrollArea, Center, Box, UnstyledButton, Collapse, Divider, Menu } from '@mantine/core';
    import { useDisclosure } from '@mantine/hooks';
    import {
        IconChevronDown,
        IconHeart, IconLogout,
        IconSettings, IconUser
    } from '@tabler/icons-react';
    import classes from './navbar.module.css';

    export function Navbar() {
        const [drawerOpened, { toggle: toggleDrawer, close: closeDrawer }] = useDisclosure(false);
        const [menuOpened, setMenuOpened] = useState(false);
        const theme = useMantineTheme();
        const [authenticated, setAuthenticated] = useState(false);

        useEffect(() => {
            // Check authentication status from local storage or any other mechanism
            const isAuthenticated = localStorage.getItem('authenticated') === 'true';
            setAuthenticated(isAuthenticated);
        }, []);

        return (
            <Box pb={60}> {/* Adjusted padding-bottom */}
                <header className={classes.header}>
                    <Group justify="space-between" h="100%">
                        <Group h="100%" gap={0} visibleFrom="sm">
                            <Link to="/" className={classes.link}>
                                Home
                            </Link>
                            {authenticated && (
                                <Link to="/posts" className={classes.link}>
                                    My Posts
                                </Link>
                            )}
                        </Group>
                        <Group visibleFrom="sm">
                            {!authenticated ? (
                                <Group>
                                    <Link to="/login">
                                        <Button variant="default">Log in</Button>
                                    </Link>
                                    <Link to="/register">
                                        <Button>Sign up</Button>
                                    </Link>
                                </Group>
                            ) : (
                                <>
                                    <Link to="/create-post">
                                        <Button color="blue">Create post</Button>
                                    </Link>
                                    <Menu
                                        width={260}
                                        position="bottom-end"
                                        transitionProps={{transition: 'pop-top-right'}}
                                        onClose={() => setMenuOpened(false)}
                                        onOpen={() => setMenuOpened(true)}
                                        withinPortal
                                    >
                                        <Menu.Target>
                                            <UnstyledButton className={classes.link}>
                                                Menu <IconChevronDown style={{width: rem(12), height: rem(12)}}
                                                                      stroke={1.5}/>
                                            </UnstyledButton>
                                        </Menu.Target>
                                        <Menu.Dropdown opened={menuOpened} onClose={() => setMenuOpened(false)}
                                                       className={classes.menuDropdown}>
                                            {/* Liked Posts */}
                                            <Link to="/posts" className={classes.dropdownSubLink} color={"red"}>
                                                <Menu.Item
                                                    className={`${classes.dropdownSubLink} ${classes.menuItemRed}`} // Apply the menu-item-red class
                                                    leftSection={
                                                        <IconHeart
                                                            style={{width: rem(16), height: rem(16)}}
                                                            color={theme.colors.red[6]}
                                                            stroke={1.5}
                                                        />
                                                    }
                                                >
                                                    Liked posts
                                                </Menu.Item>
                                            </Link>

                                            <Menu.Divider/>

                                            {/* My Profile */}
                                            <Link to="/profile" className={classes.dropdownSubLink}>
                                                <Menu.Item
                                                    leftSection={
                                                        <IconUser style={{width: rem(16), height: rem(16)}}
                                                                  stroke={1.5}/>
                                                    }
                                                >
                                                    My Profile
                                                </Menu.Item>
                                            </Link>

                                            {/* Manage Account */}
                                            <Link to="/profile/settings">
                                                <Menu.Item
                                                    leftSection={
                                                        <IconSettings style={{width: rem(16), height: rem(16)}}
                                                                      stroke={1.5}/>
                                                    }
                                                >
                                                    Manage Account
                                                </Menu.Item>
                                            </Link>

                                            {/* Logout */}
                                            <Menu.Item
                                                leftSection={
                                                    <IconLogout style={{width: rem(16), height: rem(16)}} stroke={1.5}/>
                                                }
                                                onClick={() => {
                                                    fetch("http://localhost:4000/v2/users/logout", {
                                                        method: "POST",
                                                        credentials: "include" // Ensure cookies are included in the request
                                                    })
                                                        .then(response => {
                                                            if (response.ok) {
                                                                // Clear authenticated state and local storage
                                                                setAuthenticated(false);
                                                                localStorage.removeItem("authenticated");
                                                                // Reload the page
                                                                window.location.reload();
                                                            } else {
                                                                throw new Error("Logout failed");
                                                            }
                                                        })
                                                        .catch(error => console.error("Logout error:", error));
                                                }}
                                            >
                                                Logout
                                            </Menu.Item>
                                        </Menu.Dropdown>
                                    </Menu>
                                </>
                            )}
                        </Group>
                        <Burger opened={drawerOpened} onClick={toggleDrawer} hiddenFrom="sm"/>
                    </Group>
                </header>


                <Drawer
                    opened={drawerOpened}
                    onClose={closeDrawer}
                    size="100%"
                    padding="md"
                    title="Navigation"
                    hiddenFrom="sm"
                    zIndex={1000000}
                    className={classes.drawer} // Apply the drawer class
                >
                    <ScrollArea h={`calc(100vh - ${rem(80)})`} mx="-md">
                        <Divider my="sm"/>
                        <Link to="/" className={classes.link}>
                            Home
                        </Link>
                        <UnstyledButton className={classes.link} onClick={() => setMenuOpened(!menuOpened)}>
                            <Center inline>
                                <Box component="span" mr={5}>
                                    Features
                                </Box>
                                <IconChevronDown style={{width: rem(16), height: rem(16)}}
                                                 color={theme.colors.blue[6]}/>
                            </Center>
                        </UnstyledButton>
                        <Collapse in={menuOpened}>
                            {authenticated && (
                                <>
                                    <Link to="/posts" className={classes.link}>My Posts</Link>
                                    <Link to="/profile" className={classes.link}>My Profile</Link>
                                </>
                            )}
                        </Collapse>
                        <Divider my="sm"/>
                        <Group justify="center" grow pb="xl" px="md">
                            {!authenticated && (
                                <>
                                    <Link to="/login"><Button variant="default">Log in</Button></Link>
                                    <Link to="/register"><Button>Sign up</Button></Link>
                                </>
                            )}
                        </Group>
                    </ScrollArea>
                </Drawer>
            </Box>
        );
    }
