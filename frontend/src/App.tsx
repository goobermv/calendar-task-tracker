import { Routes, Route, Navigate } from 'react-router-dom';
import { AuthPage } from './pages/AuthPage';
import { TasksPage } from './pages/TasksPage';
import { CalendarPage } from './pages/CalendarPage';
import { AdminPage } from './pages/AdminPage';
import { Layout } from './components/Layout';
import { PrivateRoute, AdminRoute } from './components/ProtectedRoute';

function App() {
    return (
        <Routes>
            {/* Вход и Регистрация */}
            <Route path="/auth" element={<AuthPage />} />

            {/* Авторизованная зона приложения */}
            <Route element={<PrivateRoute />}>
                <Route element={<Layout />}>
                    <Route path="/" element={<TasksPage />} />
                    <Route path="/calendar" element={<CalendarPage />} />
                </Route>
            </Route>

            {/* Закрытая зона администрирования */}
            <Route element={<AdminRoute />}>
                <Route element={<Layout />}>
                    <Route path="/admin" element={<AdminPage />} />
                </Route>
            </Route>

            {/* На любой неопознанный URL — кидаем на главную */}
            <Route path="*" element={<Navigate to="/" replace />} />
        </Routes>
    );
}

export default App;