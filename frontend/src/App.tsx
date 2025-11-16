import React from 'react';
import { BrowserRouter, Routes, Route, Navigate, useLocation } from 'react-router-dom';
import { ThemeProvider } from './contexts/ThemeContext';
import { AuthProvider } from './contexts/AuthContext';
import ProtectedRoute from './components/ProtectedRoute';
import Login from './pages/Login/Login';
import Dashboard from './components/Dashboard/Dashboard';
import Statistics from './pages/Statistics/Statistics';
import Reports from './pages/Reports/Reports';
import Metrics from './pages/Metrics/Metrics';
import Projects from './pages/Projects/Projects';
import Marketing from './pages/Marketing/Marketing';
import Placeholder from './pages/Placeholder/Placeholder';
import AdminApp from './admin/AdminApp';

// Компонент для обработки catch-all роута с проверкой пути
const CatchAllRoute: React.FC = () => {
    const location = useLocation();
    
    // Не редиректим, если путь начинается с /admin
    // Это включает /admin, /admin/users, /admin/projects и т.д.
    if (location.pathname.startsWith('/admin')) {
        // Возвращаем null, чтобы не рендерить ничего и позволить AdminApp обработать путь
        return null;
    }
    
    // Также не редиректим, если путь - это /users или /projects (react-admin использует эти пути)
    // Это происходит, когда react-admin пытается навигировать на /users или /projects
    // но внешний роутер еще не обработал /admin/*
    if (location.pathname === '/users' || location.pathname === '/projects' || 
        location.pathname.startsWith('/users/') || location.pathname.startsWith('/projects/')) {
        // Редиректим на /admin версию этого пути
        const adminPath = `/admin${location.pathname}`;
        console.log('[CatchAllRoute] Redirecting react-admin path to:', adminPath);
        return <Navigate to={adminPath} replace />;
    }
    
    return <Navigate to="/dashboard" replace />;
};

const App: React.FC = () => {
    return (
        <AuthProvider>
            <ThemeProvider>
                <BrowserRouter>
                    <Routes>
                        <Route path="/login" element={<Login />} />
                        <Route path="/admin/*" element={<AdminApp />} />
                        <Route path="/dashboard" element={
                            <ProtectedRoute>
                                <Dashboard />
                            </ProtectedRoute>
                        }>
                            <Route path="statistics" element={<Statistics />} />
                            <Route path="reports" element={<Reports />} />
                            <Route path="metrics" element={<Metrics />} />
                            <Route path="projects" element={<Projects />} />
                            <Route path="marketing" element={<Marketing />} />
                            <Route path="sources" element={<Placeholder />} />
                            <Route path="purchases" element={<Placeholder />} />
                            <Route path="tasks" element={<Placeholder />} />
                            <Route path="resources" element={<Placeholder />} />
                            <Route path="finance" element={<Placeholder />} />
                            <Route path="logistics" element={<Placeholder />} />
                            <Route path="innovation" element={<Placeholder />} />
                            <Route path="production" element={<Placeholder />} />
                            <Route path="company" element={<Placeholder />} />
                            <Route path="documents" element={<Placeholder />} />
                            <Route path="settings" element={<Placeholder />} />
                            <Route index element={<Projects />} />
                        </Route>
                        <Route path="/" element={<Navigate to="/dashboard" replace />} />
                        {/* Catch-all для всех остальных роутов, кроме /admin */}
                        <Route path="*" element={<CatchAllRoute />} />
                    </Routes>
                </BrowserRouter>
            </ThemeProvider>
        </AuthProvider>
    );
};

export default App;
