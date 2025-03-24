// src/services/api.ts
import axios from 'axios';

const api = axios.create({
    baseURL: 'http://localhost:8080',
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