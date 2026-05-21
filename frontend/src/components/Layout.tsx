import React, { useState } from 'react';
import { Link, useNavigate, Outlet } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { usersApi } from '../api/users';

export const Layout: React.FC = () => {
    const { user, logout } = useAuth();
    const navigate = useNavigate();

    const [isProfileModalOpen, setIsProfileModalOpen] = useState(false);
    const [isPasswordModalOpen, setIsPasswordModalOpen] = useState(false);
    
    const [username, setUsername] = useState(user?.username || '');
    const [email, setEmail] = useState(user?.email || '');
    const [newPassword, setNewPassword] = useState('');

    const modalButtonStyle = (isPrimary: boolean): React.CSSProperties => ({
        padding: '10px 20px',
        borderRadius: '6px',
        border: 'none',
        cursor: 'pointer',
        fontSize: '14px',
        fontWeight: '500',
        transition: 'all 0.2s',
        backgroundColor: isPrimary ? '#4a90e2' : '#e0e0e0',
        color: isPrimary ? 'white' : '#333'
    });

    const handleLogout = () => {
        logout();
        navigate('/auth');
    };

    const handleUpdateInfo = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!user?.id) return;
        try {
            await usersApi.updateInfo(user.id, { username, email });
            alert("Профиль успешно обновлен!");
            setIsProfileModalOpen(false);
            window.location.reload(); 
        } catch (err) { 
            console.error(err);
            alert("Ошибка обновления: возможно, это имя или email уже заняты."); 
        }
    };

    const handleUpdatePassword = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!user?.id || !newPassword) return;
        
        try {
            await usersApi.updatePassword(user.id, { password: newPassword });
            alert("Пароль успешно изменен!");
            setNewPassword('');
            setIsPasswordModalOpen(false);
        } catch (err) { 
            console.error(err);
            alert("Ошибка изменения пароля."); 
        }
    };

    const handleSelfDelete = async () => {
        if (!user?.id || !confirm("Удалить аккаунт навсегда?")) return;
        try {
            await usersApi.deleteAccount(user.id);
            handleLogout();
        } catch (err) { alert("Ошибка удаления"); }
    };

    const menuButtonStyle: React.CSSProperties = {
        background: '#3a3b40',
        padding: '10px 14px',
        border: 'none',
        borderRadius: '6px',
        color: 'white',
        cursor: 'pointer',
        textAlign: 'left',
        fontSize: '14px',
        transition: 'background 0.2s'
    };

    return (
        <div style={{ display: 'flex', minHeight: '100vh', position: 'relative' }}>
            <aside style={{ width: '260px', backgroundColor: 'var(--panel-bg)', padding: '24px', display: 'flex', flexDirection: 'column' }}>
                <h3 style={{ color: 'white' }}>⚡ Tracker App</h3>
                <div style={{ fontSize: '14px', color: 'var(--text-muted)', marginBottom: '30px' }}>
                    Вошел как: <strong>{user?.username}</strong>
                </div>

                <nav style={{ display: 'flex', flexDirection: 'column', gap: '15px', flex: 1 }}>
                    <Link to="/" style={{ color: 'white', textDecoration: 'none', fontSize: '16px' }}>📋 Список задач</Link>
                    <Link to="/calendar" style={{ color: 'white', textDecoration: 'none', fontSize: '16px' }}>📅 Календарь</Link>
                    {user?.user_type === 'admin' && (
                        <Link to="/admin" style={{ color: 'var(--warning)', textDecoration: 'none', fontSize: '16px' }}>🛡️ Панель админа</Link>
                    )}
                    
                    <div style={{ display: 'flex', flexDirection: 'column', gap: '12px', marginTop: '20px' }}>
                        <button onClick={() => { setIsProfileModalOpen(true); setIsPasswordModalOpen(false); }} style={menuButtonStyle}>
                            ✎ Редактировать профиль
                        </button>
                        <button onClick={() => { setIsPasswordModalOpen(true); setIsProfileModalOpen(false); }} style={menuButtonStyle}>
                            🔑 Сменить пароль
                        </button>
                        <button onClick={handleSelfDelete} style={{ ...menuButtonStyle, background: 'transparent', border: '1px solid var(--danger)', color: 'var(--danger)' }}>
                            🗑️ Удалить аккаунт
                        </button>
                    </div>
                </nav>

                <button onClick={handleLogout} style={{ marginTop: 'auto', padding: '10px', cursor: 'pointer' }}>Выйти</button>
            </aside>

            <main style={{ flex: 1, paddingTop: '40px', paddingBottom: '40px', paddingLeft: '40px', paddingRight: '40px', overflowY: 'auto' }}>
                <Outlet />
            </main>

            {/* Контейнер для модалок, чтобы они накладывались поверх контента */}
            {(isProfileModalOpen || isPasswordModalOpen) && (
                <div style={{ position: 'fixed', top: 0, left: 0, right: 0, bottom: 0, backgroundColor: 'rgba(0,0,0,0.7)', display: 'flex', alignItems: 'center', justifyContent: 'center', zIndex: 1000 }}>
                    
                    {isProfileModalOpen && (
                        <form onSubmit={handleUpdateInfo} style={{ backgroundColor: '#1e1e1e', padding: '30px', borderRadius: '12px', display: 'flex', flexDirection: 'column', gap: '15px', color: 'white', width: '400px' }}>
                            <h3>Редактировать профиль</h3>
                            <label>Новое имя</label>
                            <input value={username} onChange={e => setUsername(e.target.value)} style={{ padding: '8px' }} />
                            <label>Новый Email</label>
                            <input type="email" value={email} onChange={e => setEmail(e.target.value)} style={{ padding: '8px' }} />
                            <div style={{ display: 'flex', gap: '10px', marginTop: '10px' }}>
                                <button type="submit" style={modalButtonStyle(true)}>Сохранить</button>
                                <button type="button" onClick={() => setIsProfileModalOpen(false)} style={modalButtonStyle(false)}>Отмена</button>
                            </div>
                        </form>
                    )}

                    {isPasswordModalOpen && (
                        <form onSubmit={handleUpdatePassword} style={{ backgroundColor: '#1e1e1e', padding: '30px', borderRadius: '12px', display: 'flex', flexDirection: 'column', gap: '15px', color: 'white', width: '400px' }}>
                            <h3>Смена пароля</h3>
                            <label>Новый пароль</label>
                            <input type="password" value={newPassword} onChange={e => setNewPassword(e.target.value)} required style={{ padding: '8px' }} />
                            <div style={{ display: 'flex', gap: '10px', marginTop: '10px' }}>
                                <button type="submit" style={modalButtonStyle(true)}>Обновить</button>
                                <button type="button" onClick={() => setIsPasswordModalOpen(false)} style={modalButtonStyle(false)}>Отмена</button>
                            </div>
                        </form>
                    )}
                </div>
            )}
        </div>
    );
};