// src/pages/pix/ListaPix.tsx
import React, { useEffect, useState } from 'react';
import {
    Box,
    Typography,
    Container,
    Paper,
    Table,
    TableBody,
    TableCell,
    TableContainer,
    TableHead,
    TableRow,
    CircularProgress,
    Button,
} from '@mui/material';
import { useLocation, useNavigate } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';
import api from '../../services/api';
import Header from '../../components/Menu/Header';
import Sidebar from '../../components/Menu/Sidebar';

interface ChavePix {
    chave: string;
    tipoChave: string;
    status: string;
    cpfCnpj: string;
    nomeProprietario: string;
    nomeBanco: string;
    numeroBanco: string;
    agencia: string;
    numeroConta: string;
}

const ListaPix: React.FC = () => {
    const [loading, setLoading] = useState(false);
    const [requisicoes, setRequisicoes] = useState<any[]>([]);
    const { user } = useAuth();
    const location = useLocation();
    const navigate = useNavigate();

    // Recuperando dados da navegação, se existirem
    const resultados = location.state?.resultados;
    const tipoBusca = location.state?.tipoBusca;
    const chaveBusca = location.state?.chaveBusca;

    useEffect(() => {
        const fetchRequisicoes = async () => {
            setLoading(true);
            try {
                const response = await api.get('/api/bacen/pix/requisicoespix', {
                    params: {
                        cpfCnpj: user?.cpf,
                    },
                });
                setRequisicoes(response.data);
            } catch (error) {
                console.error('Erro ao buscar requisições:', error);
            } finally {
                setLoading(false);
            }
        };

        fetchRequisicoes();
    }, [user]);

    const handleNovaPesquisa = () => {
        navigate('/pix/novo');
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
                    {resultados && (
                        <Paper sx={{ p: 2, mb: 2 }}>
                            <Typography variant="h6" gutterBottom>
                                Resultado da Consulta {tipoBusca === 'chave' ? `(Chave: ${chaveBusca})` : `(CPF/CNPJ: ${chaveBusca})`}
                            </Typography>

                            {Array.isArray(resultados) && resultados.length > 0 ? (
                                <TableContainer>
                                    <Table>
                                        <TableHead>
                                            <TableRow>
                                                <TableCell>Chave</TableCell>
                                                <TableCell>Tipo</TableCell>
                                                <TableCell>Status</TableCell>
                                                <TableCell>CPF/CNPJ</TableCell>
                                                <TableCell>Nome</TableCell>
                                                <TableCell>Banco</TableCell>
                                                <TableCell>Agência</TableCell>
                                                <TableCell>Conta</TableCell>
                                            </TableRow>
                                        </TableHead>
                                        <TableBody>
                                            {resultados.map((item: ChavePix, index: number) => (
                                                <TableRow key={index}>
                                                    <TableCell>{item.chave}</TableCell>
                                                    <TableCell>{item.tipoChave}</TableCell>
                                                    <TableCell>{item.status}</TableCell>
                                                    <TableCell>{item.cpfCnpj}</TableCell>
                                                    <TableCell>{item.nomeProprietario}</TableCell>
                                                    <TableCell>{item.numeroBanco} - {item.nomeBanco}</TableCell>
                                                    <TableCell>{item.agencia}</TableCell>
                                                    <TableCell>{item.numeroConta}</TableCell>
                                                </TableRow>
                                            ))}
                                        </TableBody>
                                    </Table>
                                </TableContainer>
                            ) : (
                                <Typography>Nenhum resultado encontrado</Typography>
                            )}

                            <Button
                                variant="contained"
                                onClick={handleNovaPesquisa}
                                sx={{ mt: 2 }}
                            >
                                Nova Consulta
                            </Button>
                        </Paper>
                    )}

                    <Paper sx={{ p: 2 }}>
                        <Typography variant="h6" gutterBottom>
                            Histórico de Consultas PIX
                        </Typography>

                        {loading ? (
                            <Box sx={{ display: 'flex', justifyContent: 'center', p: 2 }}>
                                <CircularProgress />
                            </Box>
                        ) : (
                            <TableContainer>
                                <Table>
                                    <TableHead>
                                        <TableRow>
                                            <TableCell>Data</TableCell>
                                            <TableCell>Tipo</TableCell>
                                            <TableCell>Chave/CPF/CNPJ</TableCell>
                                            <TableCell>Motivo</TableCell>
                                            <TableCell>Resultado</TableCell>
                                        </TableRow>
                                    </TableHead>
                                    <TableBody>
                                        {requisicoes.map((req) => (
                                            <TableRow key={req.id}>
                                                <TableCell>{new Date(req.data).toLocaleString()}</TableCell>
                                                <TableCell>{req.tipoBusca}</TableCell>
                                                <TableCell>{req.chaveBusca}</TableCell>
                                                <TableCell>{req.motivoBusca}</TableCell>
                                                <TableCell>{req.resultado}</TableCell>
                                            </TableRow>
                                        ))}
                                        {requisicoes.length === 0 && (
                                            <TableRow>
                                                <TableCell colSpan={5} align="center">
                                                    Nenhuma consulta realizada
                                                </TableCell>
                                            </TableRow>
                                        )}
                                    </TableBody>
                                </Table>
                            </TableContainer>
                        )}

                        <Button
                            variant="contained"
                            onClick={handleNovaPesquisa}
                            sx={{ mt: 2 }}
                        >
                            Nova Consulta
                        </Button>
                    </Paper>
                </Container>
            </Box>
        </Box>
    );
};

export default ListaPix;