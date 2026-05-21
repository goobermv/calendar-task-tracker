import React, { useState } from 'react';
import { useAuth } from '../context/AuthContext';
import { useNavigate } from 'react-router-dom';

export const AuthPage: React.FC = () => {
    const { login, register } = useAuth();
    const navigate = useNavigate();
    
    const [isLoginMode, setIsLoginMode] = useState(true);
    const [email, setEmail] = useState('');
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const [error, setError] = useState('');

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setError('');

        try {
            if (isLoginMode) {
                await login({ email, password });
            } else {
                await register({ email, username, password });
            }
            navigate('/');
        } catch (err: any) {
            setError(err.response?.data?.message || 'Произошла ошибка. Проверьте данные.');
        }
    };

    return (
        <div style={{ maxWidth: '400px', margin: '100px auto', padding: '20px', border: '1px solid #ccc', borderRadius: '8px' }}>
            <h2>{isLoginMode ? 'Вход в систему' : 'Регистрация'}</h2>
            
            {error && <div style={{ color: 'red', marginBottom: '10px' }}>{error}</div>}

            <form onSubmit={handleSubmit} style={{ display: 'flex', flexDirection: 'column', gap: '15px' }}>
                <input 
                    type="email" 
                    placeholder="Email" 
                    value={email} 
                    onChange={(e) => setEmail(e.target.value)} 
                    required 
                />
                
                {!isLoginMode && (
                    <input 
                        type="text" 
                        placeholder="Имя пользователя (username)" 
                        value={username} 
                        onChange={(e) => setUsername(e.target.value)} 
                        required 
                    />
                )}

                <input 
                    type="password" 
                    placeholder="Пароль" 
                    value={password} 
                    onChange={(e) => setPassword(e.target.value)} 
                    required 
                />

                <button type="submit" style={{ padding: '10px', cursor: 'pointer' }}>
                    {isLoginMode ? 'Войти' : 'Создать аккаунт'}
                </button>
            </form>

            <button 
                onClick={() => setIsLoginMode(!isLoginMode)} 
                style={{ background: 'none', border: 'none', color: 'blue', marginTop: '15px', cursor: 'pointer', textDecoration: 'underline' }}
            >
                {isLoginMode ? 'Еще нет аккаунта? Зарегистрироваться' : 'Уже есть аккаунт? Войти'}
            </button>
        </div>
    );
};