import axios from 'axios';

// Determina a URL base com base no ambiente
const getBaseUrl = () => {
    // Verifica se está em ambiente de produção (Docker/Nginx)
    if (window.location.hostname !== 'localhost' || window.location.port === '80') {
        return '/api';
    }
    // Em desenvolvimento, usa a URL direta
    return 'http://localhost:8080';
};

const api = axios.create({
    baseURL: getBaseUrl(),
});

api.interceptors.request.use(
    (config) => {
        const token = localStorage.getItem('@ConsultaPix:token');
        if (token) {
            config.params = {
                ...(config.params || {}),
                token
            };
        }
        return config;
    },
    (error) => {
        return Promise.reject(error);
    }
);

export default api;