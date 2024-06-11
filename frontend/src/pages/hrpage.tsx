import React, { useState, useEffect } from 'react';
import axios from 'axios';
import {
    Table,
    ScrollArea,
    UnstyledButton,
    Group,
    Text,
    Center,
    TextInput,
    rem,
} from '@mantine/core';
import { Link } from 'react-router-dom'
import { IconSelector, IconChevronDown, IconChevronUp, IconSearch } from '@tabler/icons-react';
import '../styles/hrpage.css';

interface User {
    id: number;
    name: string;
    surname: string;
    username: string;
    telegram: string;
    discord: string;
    email: string;
    skills: string[];
}

interface ThProps {
    children: React.ReactNode;
    reversed: boolean;
    sorted: boolean;
    onSort(): void;
}

function Th({ children, reversed, sorted, onSort }: ThProps) {
    const Icon = sorted ? (reversed ? IconChevronUp : IconChevronDown) : IconSelector;
    return (
        <Table.Th className="th">
            <UnstyledButton onClick={onSort} className="control">
                <Group justify="space-between">
                    <Text fw={500} fz="sm">
                        {children}
                    </Text>
                    <Center className="icon">
                        <Icon style={{ width: rem(16), height: rem(16) }} stroke={1.5} />
                    </Center>
                </Group>
            </UnstyledButton>
        </Table.Th>
    );
}

function filterData(data: User[], search: string) {
    const query = search.toLowerCase().trim();
    return data.filter((item) =>
        Object.keys(item).some((key) => {
            if (Array.isArray(item[key])) {
                return item[key].join(',').toLowerCase().includes(query);
            }
            return item[key].toString().toLowerCase().includes(query);
        })
    );
}

function sortData(
    data: User[],
    payload: { sortBy: keyof User | null; reversed: boolean; search: string }
) {
    const { sortBy } = payload;

    if (!sortBy) {
        return filterData(data, payload.search);
    }

    return filterData(
        [...data].sort((a, b) => {
            if (payload.reversed) {
                return b[sortBy].toString().localeCompare(a[sortBy].toString());
            }

            return a[sortBy].toString().localeCompare(b[sortBy].toString());
        }),
        payload.search
    );
}

export function HRPage() {
    const [data, setData] = useState<User[]>([]);
    const [search, setSearch] = useState('');
    const [sortedData, setSortedData] = useState<User[]>([]);
    const [sortBy, setSortBy] = useState<keyof User | null>(null);
    const [reverseSortDirection, setReverseSortDirection] = useState(false);

    useEffect(() => {
        axios
            .get('http://localhost:4000/v2/users/', { withCredentials: true })
            .then((response) => {
                setData(response.data);
                setSortedData(response.data);
            })
            .catch((error) => {
                console.error('Error fetching data:', error);
            });
    }, []);

    const setSorting = (field: keyof User) => {
        const reversed = field === sortBy ? !reverseSortDirection : false;
        setReverseSortDirection(reversed);
        setSortBy(field);
        setSortedData(sortData(data, { sortBy: field, reversed, search }));
    };

    const handleSearchChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        const { value } = event.currentTarget;
        setSearch(value);
        setSortedData(sortData(data, { sortBy, reversed: reverseSortDirection, search: value }));
    };

    const rows = sortedData.map((row) => (
        <Table.Tr key={row.id}>
            <Table.Td>{row.name}</Table.Td>
            <Table.Td>{row.surname}</Table.Td>
            <Table.Td>
                <Link to={`/profile/${row.id}`}>{row.username}</Link>
            </Table.Td>
            <Table.Td>{row.telegram}</Table.Td>
            <Table.Td>{row.discord}</Table.Td>
            <Table.Td>{row.email}</Table.Td>
            <Table.Td>{row.skills.join(', ')}</Table.Td>
        </Table.Tr>
    ));

    return (
        <div className="hr-container">
            <div className="hr-table-wrapper">
                <ScrollArea>
                    <TextInput
                        placeholder="Search by any field"
                        mb="md"
                        leftSection={<IconSearch style={{ width: rem(16), height: rem(16) }} stroke={1.5} />}
                        value={search}
                        onChange={handleSearchChange}
                    />
                    <Table horizontalSpacing="md" verticalSpacing="xs" miw={700} layout="fixed">
                        <Table.Tbody>
                            <Table.Tr>
                                <Th
                                    sorted={sortBy === 'name'}
                                    reversed={reverseSortDirection}
                                    onSort={() => setSorting('name')}
                                >
                                    Name
                                </Th>
                                <Th
                                    sorted={sortBy === 'surname'}
                                    reversed={reverseSortDirection}
                                    onSort={() => setSorting('surname')}
                                >
                                    Surname
                                </Th>
                                <Th
                                    sorted={sortBy === 'username'}
                                    reversed={reverseSortDirection}
                                    onSort={() => setSorting('username')}
                                >
                                    Username
                                </Th>
                                <Th
                                    sorted={sortBy === 'telegram'}
                                    reversed={reverseSortDirection}
                                    onSort={() => setSorting('telegram')}
                                >
                                    Telegram
                                </Th>
                                <Th
                                    sorted={sortBy === 'discord'}
                                    reversed={reverseSortDirection}
                                    onSort={() => setSorting('discord')}
                                >
                                    Discord
                                </Th>
                                <Th
                                    sorted={sortBy === 'email'}
                                    reversed={reverseSortDirection}
                                    onSort={() => setSorting('email')}
                                >
                                    Email
                                </Th>
                                <Th
                                    sorted={sortBy === 'skills'}
                                    reversed={reverseSortDirection}
                                    onSort={() => setSorting('skills')}
                                >
                                    Skills
                                </Th>
                            </Table.Tr>
                        </Table.Tbody>
                        <Table.Tbody>
                            {rows.length > 0 ? (
                                rows
                            ) : (
                                <Table.Tr>
                                    <Table.Td colSpan={data.length > 0 ? Object.keys(data[0]).length : 7}>
                                        <Text fw={500} ta="center">
                                            Nothing found
                                        </Text>
                                    </Table.Td>
                                </Table.Tr>
                            )}
                        </Table.Tbody>
                    </Table>
                </ScrollArea>
            </div>
        </div>
    );
}
