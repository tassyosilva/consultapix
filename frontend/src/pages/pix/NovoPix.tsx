// src/pages/pix/NovoPix.tsx
import React, { useState } from 'react';
import {
    Box,
    Button,
    TextField,
    Typography,
    Container,
    Paper,
    Grid,
    FormControl,
    InputLabel,
    Select,
    MenuItem,
    CircularProgress,
    Alert,
} from '@mui/material';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';
import api from '../../services/api';
import Header from '../../components/Menu/Header';
import Sidebar from '../../components/Menu/Sidebar';

const NovoPix: React.FC = () => {
    const [tipoBusca, setTipoBusca] = useState('chave');
    const [chave, setChave] = useState('');
    const [cpfCnpj, setCpfCnpj] = useState('');
    const [motivo, setMotivo] = useState('');
    const [caso, setCaso] = useState('');
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState('');
    const { user } = useAuth();
    const navigate = useNavigate();

    const handleSubmit = async (event: React.FormEvent) => {
        event.preventDefault();
        setLoading(true);
        setError('');

        try {
            if (tipoBusca === 'chave') {
                const response = await api.get('/api/bacen/pix/chave', {
                    params: {
                        chave,
                        motivo,
                        caso,
                        cpfResponsavel: user?.cpf,
                        lotacao: user?.lotacao,
                    },
                });

                // Redirecionamento para página de resultados com os dados
                navigate('/pix', {
                    state: {
                        resultados: response.data,
                        tipoBusca,
                        chaveBusca: chave
                    }
                });
            } else {
                const response = await api.get('/api/bacen/pix/cpfCnpj', {
                    params: {
                        cpfCnpj,
                        motivo,
                        caso,
                        cpfResponsavel: user?.cpf,
                        lotacao: user?.lotacao,
                    },
                });

                navigate('/pix', {
                    state: {
                        resultados: response.data,
                        tipoBusca,
                        chaveBusca: cpfCnpj
                    }
                });
            }
        } catch (err: any) {
            setError(err.response?.data?.message || 'Erro ao realizar consulta');
        } finally {
            setLoading(false);
        }
    };

    return (
        <Box sx={{ display: 'flex' }}>
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
                    <Paper sx={{ p: 2 }}>
                        <Typography variant="h6" gutterBottom>
                            Nova Consulta PIX
                        </Typography>
                        {error && <Alert severity="error">{error}</Alert>}
                        <Box component="form" noValidate onSubmit={handleSubmit} sx={{ mt: 3 }}>
                            <Grid container spacing={2}>
                                <Grid item xs={12}>
                                    <FormControl fullWidth>
                                        <InputLabel id="tipo-busca-label">Tipo de Busca</InputLabel>
                                        <Select
                                            labelId="tipo-busca-label"
                                            id="tipo-busca"
                                            value={tipoBusca}
                                            label="Tipo de Busca"
                                            onChange={(e) => setTipoBusca(e.target.value)}
                                        >
                                            <MenuItem value="chave">Chave PIX</MenuItem>
                                            <MenuItem value="cpfCnpj">CPF/CNPJ</MenuItem>
                                        </Select>
                                    </FormControl>
                                </Grid>

                                {tipoBusca === 'chave' ? (
                                    <Grid item xs={12}>
                                        <TextField
                                            required
                                            fullWidth
                                            id="chave"
                                            label="Chave PIX"
                                            name="chave"
                                            value={chave}
                                            onChange={(e) => setChave(e.target.value)}
                                        />
                                    </Grid>
                                ) : (
                                    <Grid item xs={12}>
                                        <TextField
                                            required
                                            fullWidth
                                            id="cpfCnpj"
                                            label="CPF/CNPJ"
                                            name="cpfCnpj"
                                            value={cpfCnpj}
                                            onChange={(e) => setCpfCnpj(e.target.value)}
                                        />
                                    </Grid>
                                )}

                                <Grid item xs={12}>
                                    <TextField
                                        required
                                        fullWidth
                                        id="motivo"
                                        label="Motivo da Consulta"
                                        name="motivo"
                                        value={motivo}
                                        onChange={(e) => setMotivo(e.target.value)}
                                    />
                                </Grid>

                                <Grid item xs={12}>
                                    <TextField
                                        fullWidth
                                        id="caso"
                                        label="Caso/Processo"
                                        name="caso"
                                        value={caso}
                                        onChange={(e) => setCaso(e.target.value)}
                                    />
                                </Grid>
                            </Grid>

                            <Button
                                type="submit"
                                fullWidth
                                variant="contained"
                                sx={{ mt: 3, mb: 2 }}
                                disabled={loading}
                            >
                                {loading ? <CircularProgress size={24} /> : 'Consultar'}
                            </Button>
                        </Box>
                    </Paper>
                </Container>
            </Box>
        </Box>
    );
};

export default NovoPix;