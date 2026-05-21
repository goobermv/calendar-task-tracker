import React from 'react';
import { Navigate, Outlet } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';

export const PrivateRoute: React.FC = () => {
    const { isAuthenticated, isLoading } = useAuth();

    if (isLoading) return <div>Загрузка...</div>;

    return isAuthenticated ? <Outlet /> : <Navigate to="/auth" replace />;
};

export const AdminRoute: React.FC = () => {
    const { user, isAuthenticated, isLoading } = useAuth();

    if (isLoading) return <div>Загрузка...</div>;

    if (!isAuthenticated) return <Navigate to="/auth" replace />;
    
    return user?.user_type === 'admin' ? <Outlet /> : <Navigate to="/" replace />;
};