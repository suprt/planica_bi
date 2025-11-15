import React from 'react';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { ThemeProvider } from './contexts/ThemeContext';
import Login from './pages/Login/Login';
import Dashboard from './components/Dashboard/Dashboard';
import Statistics from './pages/Statistics/Statistics';
import Placeholder from './components/Placeholder/Placeholder';

const App: React.FC = () => {
    return (
        <ThemeProvider>
            <BrowserRouter>
                <Routes>
                    <Route path="/login" element={<Login />} />
                    <Route path="/dashboard" element={<Dashboard />}>
                        <Route path="statistics" element={<Statistics />} />
                        <Route index element={<Placeholder />} />
                    </Route>
                    <Route path="/" element={<Navigate to="/login" replace />} />
                    <Route path="*" element={<Navigate to="/login" replace />} />
                </Routes>
            </BrowserRouter>
        </ThemeProvider>
    );
};

export default App;
