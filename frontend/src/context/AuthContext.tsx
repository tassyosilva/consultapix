import React, { createContext, useState, useContext, ReactNode } from 'react';
import { useNavigate } from 'react-router-dom';
import api from '../services/api';

interface User {
    id: number;
    nome: string;
    cpf: string;
    email: string;
    lotacao: string;
    matricula: string;
    admin: boolean;
}

interface AuthState {
    token: string | null;
    user: User | null;
}

interface AuthContextData {
    signed: boolean;
    user: User | null;
    loading: boolean;
    signIn: (email: string, password: string) => Promise<void>;
    signOut: () => void;
}

const AuthContext = createContext<AuthContextData>({} as AuthContextData);

interface AuthProviderProps {
    children: ReactNode;
}

export const AuthProvider: React.FC<AuthProviderProps> = ({ children }) => {
    const [data, setData] = useState<AuthState>(() => {
        const token = localStorage.getItem('@ConsultaPix:token');
        const user = localStorage.getItem('@ConsultaPix:user');

        if (token && user) {
            api.defaults.headers.common['Authorization'] = `Bearer ${token}`;
            return { token, user: JSON.parse(user) };
        }

        return {} as AuthState;
    });

    const [loading, setLoading] = useState(false);
    const navigate = useNavigate();

    async function signIn(email: string, password: string) {
        try {
            setLoading(true);
            const response = await api.post('/api/user/login', { email, password });

            if (response.data.status === 409) {
                throw new Error(response.data.message);
            }

            const { token, payload } = response.data;

            localStorage.setItem('@ConsultaPix:token', token);
            localStorage.setItem('@ConsultaPix:user', JSON.stringify(payload));

            api.defaults.headers.common['Authorization'] = `Bearer ${token}`;

            setData({ token, user: payload });
            navigate('/dashboard');
        } catch (error) {
            throw error;
        } finally {
            setLoading(false);
        }
    }

    function signOut() {
        localStorage.removeItem('@ConsultaPix:token');
        localStorage.removeItem('@ConsultaPix:user');
        setData({} as AuthState);
        navigate('/');
    }

    return (
        <AuthContext.Provider
            value={{
                signed: !!data.user,
                user: data.user,
                loading,
                signIn,
                signOut,
            }}
        >
            {children}
        </AuthContext.Provider>
    );
};

export function useAuth(): AuthContextData {
    const context = useContext(AuthContext);

    if (!context) {
        throw new Error('useAuth must be used within an AuthProvider');
    }

    return context;
}