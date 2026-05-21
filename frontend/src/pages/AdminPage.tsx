import React, { useEffect, useState } from 'react';
import { tasksApi } from '../api/tasks';
import { eventsApi } from '../api/events';
import { usersApi } from '../api/users';
import type { Task, Event, TaskPriority, TaskStatus, EventType } from '../types';

interface AdminUser {
    id: string;
    username: string;
    email: string;
    user_type: string;
}

const EVENT_TYPE_DETAILS: Record<EventType, { label: string; color: string }> = {
    personal: { label: 'Личное', color: '#3b82f6' },
    business: { label: 'Бизнес / Работа', color: '#10b981' },
    formal: { label: 'Официальное', color: '#6366f1' },
    important: { label: 'Важное', color: '#ef4444' },
    celebration: { label: 'Праздник / Торжество', color: '#f59e0b' },
    holiday: { label: 'Выходной', color: '#ec4899' },
    casual: { label: 'Повседневное', color: '#6b7280' }
};

const TASK_PRIORITY_DETAILS: Record<TaskPriority, { label: string; color: string }> = {
    low: { label: 'Низкий', color: '#10b981' },
    medium: { label: 'Средний', color: '#3b82f6' },
    high: { label: 'Высокий', color: '#f59e0b' },
    urgent: { label: 'Критический', color: '#ef4444' }
};

export const AdminPage: React.FC = () => {
    const [tasks, setTasks] = useState<Task[]>([]);
    const [events, setEvents] = useState<Event[]>([]);
    const [users, setUsers] = useState<AdminUser[]>([]);
    
    const [totalTasks, setTotalTasks] = useState(0);
    const [totalEvents, setTotalEvents] = useState(0);
    
    const [usersMap, setUsersMap] = useState<Record<string, string>>({});

    const [editingTask, setEditingTask] = useState<Task | null>(null);
    const [editingEvent, setEditingEvent] = useState<Event | null>(null);

    const loadAdminData = async () => {
        try {
            const usersRes = await usersApi.getAllUsers() || [];
            setUsers(usersRes);

            const mapping: Record<string, string> = {};
            usersRes.forEach((u: AdminUser) => {
                mapping[u.id] = u.username;
            });
            setUsersMap(mapping);

            const tasksRes = await tasksApi.adminGetAllTasks?.(1, 100);
            if (tasksRes) {
                const tasksArray = Array.isArray(tasksRes) ? tasksRes : (tasksRes.tasks || []);
                setTasks(tasksArray);
                setTotalTasks(tasksRes.total_count || tasksArray.length);
            }

            const eventsRes = await eventsApi.adminGetAllEvents?.(1, 100);
            if (eventsRes) {
                const eventsArray = Array.isArray(eventsRes) ? eventsRes : (eventsRes.events || []);
                setEvents(eventsArray);
                setTotalEvents(eventsRes.total_count || eventsArray.length);
            }
        } catch (err) {
            console.error("Ошибка при загрузке данных админ-панели", err);
        }
    };

    useEffect(() => {
        loadAdminData();
    }, []);

    const handlePromote = async (id: string) => {
        try {
            await usersApi.promoteToAdmin(id);
            loadAdminData();
        } catch (err) {
            alert("Не удалось повысить пользователя до админа");
        }
    };

    const handleDemote = async (id: string) => {
        try {
            await usersApi.demoteToUser(id);
            loadAdminData();
        } catch (err) {
            alert("Не удалось разжаловать админа");
        }
    };

    const handleDeleteTask = async (id: string) => {
        if (!window.confirm("Удалить эту задачу глобально?")) return;
        try {
            await tasksApi.adminDeleteTask(id);
            loadAdminData();
        } catch (err) {
            alert("Не удалось удалить задачу");
        }
    };

    const handleUpdateTaskSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!editingTask) return;
        try {
            await tasksApi.adminUpdateTask(editingTask.id, {
                title: editingTask.title,
                priority: editingTask.priority,
                status: editingTask.status
            });
            setEditingTask(null);
            loadAdminData();
        } catch (err) {
            alert("Не удалось обновить задачу");
        }
    };

    const handleDeleteEvent = async (id: string) => {
        if (!window.confirm("Удалить это событие глобально?")) return;
        try {
            await eventsApi.adminDeleteEvent(id);
            loadAdminData();
        } catch (err) {
            alert("Не удалось удалить событие");
        }
    };

    const handleUpdateEventSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!editingEvent) return;
        try {
            await eventsApi.adminUpdateEvent(editingEvent.id, {
                title: editingEvent.title,
                event_type: editingEvent.event_type
            });
            setEditingEvent(null);
            loadAdminData();
        } catch (err) {
            alert("Не удалось обновить событие");
        }
    };

    const thStyle: React.CSSProperties = { padding: '12px', textAlign: 'left', borderBottom: '1px solid #2d2d38', color: '#9ca3af', fontSize: '14px' };
    const tdStyle: React.CSSProperties = { padding: '12px', borderBottom: '1px solid #2d2d38', fontSize: '14px', color: 'white' };
    const inputStyle: React.CSSProperties = { padding: '8px', background: '#202024', border: '1px solid #2d2d38', borderRadius: '4px', color: 'white', width: '100%' };

    return (
        <div style={{ padding: '24px', display: 'flex', flexDirection: 'column', gap: '40px' }}>
            
            {/* ТАБЛИЦА ЗАДАЧ */}
            <div style={{ backgroundColor: '#16161a', padding: '24px', borderRadius: '8px' }}>
                <h2 style={{ display: 'flex', alignItems: 'center', gap: '10px', color: '#f59e0b', fontSize: '22px', margin: 0 }}>
                    🛡️ Панель глобального администрирования
                </h2>
                <p style={{ color: '#9ca3af', marginTop: '10px', marginBottom: '20px' }}>
                    Всего задач во всей системе: <strong style={{ color: 'white' }}>{totalTasks}</strong>
                </p>

                <div style={{ overflowX: 'auto' }}>
                    <table style={{ width: '100%', borderCollapse: 'collapse' }}>
                        <thead>
                            <tr style={{ backgroundColor: '#202024' }}>
                                <th style={thStyle}>Пользователь</th>
                                <th style={thStyle}>Название задачи</th>
                                <th style={thStyle}>Приоритет</th>
                                <th style={thStyle}>Статус</th>
                                <th style={thStyle}>Действия</th>
                            </tr>
                        </thead>
                        <tbody>
                            {tasks.map(task => {
                                const priorityInfo = TASK_PRIORITY_DETAILS[task.priority] || { label: task.priority, color: '#3b82f6' };
                                const usernameDisplay = usersMap[task.user_id] || `ID: ${task.user_id.slice(0, 8)}...`;

                                return (
                                    <tr key={task.id} style={{ backgroundColor: '#1a1a1e' }}>
                                        <td style={{ ...tdStyle, color: '#f59e0b', fontWeight: '500' }}>👤 {usernameDisplay}</td>
                                        <td style={{ ...tdStyle, fontWeight: 'bold' }}>{task.title}</td>
                                        <td style={tdStyle}>
                                            <span style={{ padding: '3px 8px', borderRadius: '12px', fontSize: '11px', background: `${priorityInfo.color}22`, color: priorityInfo.color, border: `1px solid ${priorityInfo.color}44`, fontWeight: 'bold' }}>
                                                {priorityInfo.label}
                                            </span>
                                        </td>
                                        <td style={tdStyle}>{task.status}</td>
                                        <td style={{ ...tdStyle, display: 'flex', gap: '12px' }}>
                                            <button onClick={() => setEditingTask(task)} style={{ background: 'none', border: 'none', color: '#f59e0b', cursor: 'pointer' }}>✏️</button>
                                            <button onClick={() => handleDeleteTask(task.id)} style={{ background: 'none', border: 'none', color: '#ef4444', cursor: 'pointer' }}>🗑️</button>
                                        </td>
                                    </tr>
                                );
                            })}
                        </tbody>
                    </table>
                </div>
            </div>

            {/* ТАБЛИЦА СОБЫТИЙ */}
            <div style={{ backgroundColor: '#16161a', padding: '24px', borderRadius: '8px' }}>
                <h3 style={{ color: '#f59e0b', fontSize: '18px', margin: '0 0 10px 0' }}>📅 Глобальное управление событиями</h3>
                <p style={{ color: '#9ca3af', marginBottom: '20px' }}>
                    Всего событий во всей системе: <strong style={{ color: 'white' }}>{totalEvents}</strong>
                </p>

                <div style={{ overflowX: 'auto' }}>
                    <table style={{ width: '100%', borderCollapse: 'collapse' }}>
                        <thead>
                            <tr style={{ backgroundColor: '#202024' }}>
                                <th style={thStyle}>Пользователь</th>
                                <th style={thStyle}>Название события</th>
                                <th style={thStyle}>Тип события</th>
                                <th style={thStyle}>Временной интервал</th>
                                <th style={thStyle}>Действия</th>
                            </tr>
                        </thead>
                        <tbody>
                            {events.map(event => {
                                const currentType = event.event_type as EventType;
                                const typeDetails = EVENT_TYPE_DETAILS[currentType] || { label: event.event_type, color: '#6366f1' };
                                const usernameDisplay = usersMap[event.user_id] || `ID: ${event.user_id.slice(0, 8)}...`;

                                return (
                                    <tr key={event.id} style={{ backgroundColor: '#1a1a1e' }}>
                                        <td style={{ ...tdStyle, color: '#f59e0b', fontWeight: '500' }}>👤 {usernameDisplay}</td>
                                        <td style={{ ...tdStyle, fontWeight: 'bold' }}>{event.title}</td>
                                        <td style={tdStyle}>
                                            <span style={{ padding: '3px 8px', borderRadius: '12px', fontSize: '11px', background: `${typeDetails.color}22`, color: typeDetails.color, border: `1px solid ${typeDetails.color}44`, fontWeight: 'bold' }}>
                                                {typeDetails.label}
                                            </span>
                                        </td>
                                        <td style={{ ...tdStyle, fontSize: '12px', color: '#9ca3af' }}>
                                            ⏰ {new Date(event.start_time).toLocaleDateString()} — {new Date(event.end_time).toLocaleDateString()}
                                        </td>
                                        <td style={{ ...tdStyle, display: 'flex', gap: '12px' }}>
                                            <button onClick={() => setEditingEvent(event)} style={{ background: 'none', border: 'none', color: '#f59e0b', cursor: 'pointer' }}>✏️</button>
                                            <button onClick={() => handleDeleteEvent(event.id)} style={{ background: 'none', border: 'none', color: '#ef4444', cursor: 'pointer' }}>🗑️</button>
                                        </td>
                                    </tr>
                                );
                            })}
                        </tbody>
                    </table>
                </div>
            </div>

            {/* НОВАЯ СЕКЦИЯ: ГЛОБАЛЬНОЕ УПРАВЛЕНИЕ ПОЛЬЗОВАТЕЛЯМИ */}
            <div style={{ backgroundColor: '#16161a', padding: '24px', borderRadius: '8px' }}>
                <h3 style={{ color: '#f59e0b', fontSize: '18px', margin: '0 0 10px 0' }}>👥 Глобальное управление пользователями</h3>
                <p style={{ color: '#9ca3af', marginBottom: '20px' }}>
                    Всего зарегистрировано пользователей: <strong style={{ color: 'white' }}>{users.length}</strong>
                </p>

                <div style={{ overflowX: 'auto' }}>
                    <table style={{ width: '100%', borderCollapse: 'collapse' }}>
                        <thead>
                            <tr style={{ backgroundColor: '#202024' }}>
                                <th style={thStyle}>ID пользователя</th>
                                <th style={thStyle}>Username</th>
                                <th style={thStyle}>Email</th>
                                <th style={thStyle}>Роль (Тип)</th>
                                <th style={thStyle}>Изменить права</th>
                            </tr>
                        </thead>
                        <tbody>
                            {users.map(u => (
                                <tr key={u.id} style={{ backgroundColor: '#1a1a1e' }}>
                                    <td style={{ ...tdStyle, color: '#6b7280', fontSize: '12px' }}>{u.id}</td>
                                    <td style={{ ...tdStyle, fontWeight: 'bold' }}>{u.username}</td>
                                    <td style={{ ...tdStyle, color: '#9ca3af' }}>{u.email}</td>
                                    <td style={tdStyle}>
                                        <span style={{
                                            padding: '4px 10px',
                                            borderRadius: '4px',
                                            fontSize: '12px',
                                            fontWeight: 'bold',
                                            background: u.user_type === 'admin' ? '#ef444422' : '#10b98122',
                                            color: u.user_type === 'admin' ? '#ef4444' : '#10b981',
                                            border: u.user_type === 'admin' ? '1px solid #ef444444' : '1px solid #10b98144'
                                        }}>
                                            {u.user_type === 'admin' ? 'ADMIN' : 'REGULAR'}
                                        </span>
                                    </td>
                                    <td style={tdStyle}>
                                        {u.user_type === 'admin' ? (
                                            <button 
                                                onClick={() => handleDemote(u.id)}
                                                style={{ padding: '4px 12px', background: '#ef4444', color: 'white', border: 'none', borderRadius: '4px', cursor: 'pointer', fontWeight: 'bold', fontSize: '12px' }}
                                            >
                                                Demote to User
                                            </button>
                                        ) : (
                                            <button 
                                                onClick={() => handlePromote(u.id)}
                                                style={{ padding: '4px 12px', background: '#3b82f6', color: 'white', border: 'none', borderRadius: '4px', cursor: 'pointer', fontWeight: 'bold', fontSize: '12px' }}
                                            >
                                                Promote to Admin
                                            </button>
                                        )}
                                    </td>
                                </tr>
                            ))}
                        </tbody>
                    </table>
                </div>
            </div>

            {/* МОДАЛКИ (ОСТАЮТСЯ БЕЗ ИЗМЕНЕНИЙ) */}
            {editingTask && (
                <div style={{ position: 'fixed', top: 0, left: 0, right: 0, bottom: 0, background: 'rgba(0,0,0,0.7)', display: 'flex', justifyContent: 'center', alignItems: 'center', zIndex: 1000 }}>
                    <form onSubmit={handleUpdateTaskSubmit} style={{ background: '#16161a', padding: '24px', borderRadius: '8px', width: '400px', display: 'flex', flexDirection: 'column', gap: '15px', border: '1px solid #2d2d38' }}>
                        <h4 style={{ color: 'white', margin: 0 }}>Редактировать задачу (Админ)</h4>
                        <input type="text" style={inputStyle} value={editingTask.title} onChange={e => setEditingTask({...editingTask, title: e.target.value})} required />
                        <select style={inputStyle} value={editingTask.priority} onChange={e => setEditingTask({...editingTask, priority: e.target.value as TaskPriority})}>
                            {Object.entries(TASK_PRIORITY_DETAILS).map(([k, v]) => <option key={k} value={k}>{v.label}</option>)}
                        </select>
                        <select style={inputStyle} value={editingTask.status} onChange={e => setEditingTask({...editingTask, status: e.target.value as TaskStatus})}>
                            <option value="pending">pending</option>
                            <option value="in_progress">in_progress</option>
                            <option value="completed">completed</option>
                            <option value="cancelled">cancelled</option>
                        </select>
                        <div style={{ display: 'flex', gap: '10px', justifyContent: 'flex-end', marginTop: '10px' }}>
                            <button type="button" onClick={() => setEditingTask(null)} style={{ background: 'none', border: 'none', color: 'white', cursor: 'pointer' }}>Отмена</button>
                            <button type="submit" style={{ padding: '6px 12px', background: '#f59e0b', border: 'none', borderRadius: '4px', color: 'black', fontWeight: 'bold', cursor: 'pointer' }}>Сохранить</button>
                        </div>
                    </form>
                </div>
            )}

            {editingEvent && (
                <div style={{ position: 'fixed', top: 0, left: 0, right: 0, bottom: 0, background: 'rgba(0,0,0,0.7)', display: 'flex', justifyContent: 'center', alignItems: 'center', zIndex: 1000 }}>
                    <form onSubmit={handleUpdateEventSubmit} style={{ background: '#16161a', padding: '24px', borderRadius: '8px', width: '400px', display: 'flex', flexDirection: 'column', gap: '15px', border: '1px solid #2d2d38' }}>
                        <h4 style={{ color: 'white', margin: 0 }}>Редактировать событие (Админ)</h4>
                        <input type="text" style={inputStyle} value={editingEvent.title} onChange={e => setEditingEvent({...editingEvent, title: e.target.value})} required />
                        <select style={inputStyle} value={editingEvent.event_type} onChange={e => setEditingEvent({...editingEvent, event_type: e.target.value as EventType})}>
                            {Object.entries(EVENT_TYPE_DETAILS).map(([k, v]) => <option key={k} value={k}>{v.label}</option>)}
                        </select>
                        <div style={{ display: 'flex', gap: '10px', justifyContent: 'flex-end', marginTop: '10px' }}>
                            <button type="button" onClick={() => setEditingEvent(null)} style={{ background: 'none', border: 'none', color: 'white', cursor: 'pointer' }}>Отмена</button>
                            <button type="submit" style={{ padding: '6px 12px', background: '#f59e0b', border: 'none', borderRadius: '4px', color: 'black', fontWeight: 'bold', cursor: 'pointer' }}>Сохранить</button>
                        </div>
                    </form>
                </div>
            )}

        </div>
    );
};