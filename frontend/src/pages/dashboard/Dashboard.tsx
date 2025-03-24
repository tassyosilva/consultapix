// src/pages/dashboard/Dashboard.tsx
import React from 'react';
import { Box, CssBaseline, Typography, Container, Grid, Paper } from '@mui/material';
import { useAuth } from '../../context/AuthContext';
import Header from '../../components/Menu/Header';
import Sidebar from '../../components/Menu/Sidebar';

const Dashboard: React.FC = () => {
    const { user } = useAuth();

    return (
        <Box sx={{ display: 'flex' }}>
            <CssBaseline />
            <Sidebar />
            <Box
                component="main"
                sx={{
                    backgroundColor: (theme) =>
                        theme.palette.mode === 'light'
                            ? theme.palette.grey[100]
                            : theme.palette.grey[900],
                    flexGrow: 1,
                    height: '100vh',
                    overflow: 'auto',
                }}
            >
                <Header />
                <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
                    <Grid container spacing={3}>
                        <Grid item xs={12}>
                            <Paper
                                sx={{
                                    p: 2,
                                    display: 'flex',
                                    flexDirection: 'column',
                                    height: 240,
                                }}
                            >
                                <Typography variant="h4" gutterBottom>
                                    Bem-vindo, {user?.nome}
                                </Typography>
                                <Typography variant="body1">
                                    Sistema de Consulta PIX/CCS
                                </Typography>
                            </Paper>
                        </Grid>
                    </Grid>
                </Container>
            </Box>
        </Box>
    );
};

export default Dashboard;