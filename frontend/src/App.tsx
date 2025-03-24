// src/App.tsx
import React from 'react';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { ThemeProvider } from '@mui/material/styles';
import { CssBaseline } from '@mui/material';
import theme from './theme';
import { AuthProvider, useAuth } from './context/AuthContext';
import Login from './pages/auth/Login';
import Dashboard from './pages/dashboard/Dashboard';
import NovoPix from './pages/pix/NovoPix';
import ListaPix from './pages/pix/ListaPix';

// Componente de rota protegida
const PrivateRoute = ({ children }: { children: React.ReactNode }) => {
  const { signed } = useAuth();

  if (!signed) {
    return <Navigate to="/" replace />;
  }

  return <>{children}</>;
};

// Componente principal
const AppContent = () => {
  return (
    <Routes>
      <Route path="/" element={<Login />} />
      <Route
        path="/dashboard"
        element={
          <PrivateRoute>
            <Dashboard />
          </PrivateRoute>
        }
      />
      <Route
        path="/pix/novo"
        element={
          <PrivateRoute>
            <NovoPix />
          </PrivateRoute>
        }
      />
      <Route
        path="/pix"
        element={
          <PrivateRoute>
            <ListaPix />
          </PrivateRoute>
        }
      />
    </Routes>
  );
};

const App: React.FC = () => {
  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <BrowserRouter>
        <AuthProvider>
          <AppContent />
        </AuthProvider>
      </BrowserRouter>
    </ThemeProvider>
  );
};

export default App;