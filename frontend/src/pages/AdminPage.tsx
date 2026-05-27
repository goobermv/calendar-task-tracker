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

type TabType = 'tasks' | 'events' | 'users';

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

const TASK_STATUS_DETAILS: Record<TaskStatus, { label: string; color: string }> = {
    pending: { label: 'Ожидает', color: '#6b7280' },
    in_progress: { label: 'В процессе', color: '#3b82f6' },
    completed: { label: 'Выполнено', color: '#10b981' },
    cancelled: { label: 'Отменено', color: '#ef4444' }
};

export const AdminPage: React.FC = () => {
    const [activeTab, setActiveTab] = useState<TabType>('tasks');
    const [tasks, setTasks] = useState<Task[]>([]);
    const [events, setEvents] = useState<Event[]>([]);
    const [users, setUsers] = useState<AdminUser[]>([]);
    
    const [taskPage, setTaskPage] = useState(1);
    const [eventPage, setEventPage] = useState(1);
    const [userPage, setUserPage] = useState(1);
    const limit = 10;
    const [totalTasks, setTotalTasks] = useState(0);
    const [totalEvents, setTotalEvents] = useState(0);
    
    const [searchQuery, setSearchQuery] = useState('');

    const [usersMap, setUsersMap] = useState<Record<string, string>>({});
    const [editingTask, setEditingTask] = useState<Task | null>(null);
    const [editingEvent, setEditingEvent] = useState<Event | null>(null);

    const [deleteTarget, setDeleteTarget] = useState<{ id: string, type: 'task' | 'event' } | null>(null);
    const [deleteConfirmText, setDeleteConfirmText] = useState('');

    const loadUsersData = async () => {
        try {
            const usersRes = await usersApi.getAllUsers() || [];
            setUsers(usersRes);
            const mapping: Record<string, string> = {};
            usersRes.forEach((u: AdminUser) => { mapping[u.id] = u.username; });
            setUsersMap(mapping);
        } catch (err) {
            console.error("Ошибка при загрузке пользователей", err);
        }
    };

    const loadTasksData = async () => {
        try {
            const tasksRes = await tasksApi.adminGetAllTasks?.(taskPage, limit);
            if (tasksRes) {
                const tasksArray = Array.isArray(tasksRes) ? tasksRes : (tasksRes.tasks || []);
                setTasks(tasksArray);
                setTotalTasks(tasksRes.total_count || tasksArray.length);
            }
        } catch (err) {
            console.error("Ошибка при загрузке задач", err);
        }
    };

    const loadEventsData = async () => {
        try {
            const eventsRes = await eventsApi.adminGetAllEvents?.(eventPage, limit);
            if (eventsRes) {
                const eventsArray = Array.isArray(eventsRes) ? eventsRes : (eventsRes.events || []);
                setEvents(eventsArray);
                setTotalEvents(eventsRes.total_count || eventsArray.length);
            }
        } catch (err) {
            console.error("Ошибка при загрузке событий", err);
        }
    };

    useEffect(() => {
        loadUsersData();
    }, []);

    useEffect(() => {
        if (activeTab === 'tasks') loadTasksData();
        if (activeTab === 'events') loadEventsData();
    }, [activeTab, taskPage, eventPage]);

    useEffect(() => {
        setSearchQuery('');
    }, [activeTab]);

    useEffect(() => {
        setUserPage(1);
        setTaskPage(1);
        setEventPage(1);
    }, [searchQuery]);

    const handlePromote = async (id: string) => {
        try {
            await usersApi.promoteToAdmin(id);
            loadUsersData();
        } catch (err) { alert("Ошибка повышения"); }
    };

    const handleDemote = async (id: string) => {
        try {
            await usersApi.demoteToUser(id);
            loadUsersData();
        } catch (err) { alert("Ошибка разжалования"); }
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
            loadTasksData();
        } catch (err) { alert("Ошибка обновления"); }
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
            loadEventsData();
        } catch (err) { alert("Ошибка обновления"); }
    };

    const executeDelete = async () => {
        if (deleteConfirmText !== 'УДАЛИТЬ' || !deleteTarget) return;
        
        try {
            if (deleteTarget.type === 'task') {
                await tasksApi.adminDeleteTask(deleteTarget.id);
                setEditingTask(null);
                loadTasksData();
            } else {
                await eventsApi.adminDeleteEvent(deleteTarget.id);
                setEditingEvent(null);
                loadEventsData();
            }
            setDeleteTarget(null);
            setDeleteConfirmText('');
        } catch (err) {
            alert("Не удалось удалить элемент");
        }
    };

    const thStyle: React.CSSProperties = { padding: '12px', textAlign: 'left', borderBottom: '1px solid #2d2d38', color: '#9ca3af', fontSize: '14px' };
    const tdStyle: React.CSSProperties = { padding: '12px', borderBottom: '1px solid #2d2d38', fontSize: '14px', color: 'white' };
    const inputStyle: React.CSSProperties = { padding: '8px', background: '#202024', border: '1px solid #2d2d38', borderRadius: '4px', color: 'white', width: '100%' };
    const tabStyle = (isActive: boolean): React.CSSProperties => ({
        padding: '10px 20px', cursor: 'pointer', background: isActive ? '#f59e0b' : '#202024',
        color: isActive ? 'black' : 'white', border: 'none', borderRadius: '4px', fontWeight: 'bold'
    });

    const filteredTasks = tasks.filter(t => t.title.toLowerCase().includes(searchQuery.toLowerCase()));
    const filteredEvents = events.filter(e => e.title.toLowerCase().includes(searchQuery.toLowerCase()));
    const filteredUsers = users.filter(u => u.username.toLowerCase().includes(searchQuery.toLowerCase()) || u.email.toLowerCase().includes(searchQuery.toLowerCase()));

    const paginatedUsers = filteredUsers.slice((userPage - 1) * limit, userPage * limit);
    const totalUsersPages = Math.ceil(filteredUsers.length / limit);

    return (
        <div style={{ padding: '24px', display: 'flex', flexDirection: 'column', gap: '20px' }}>
            <h2 style={{ display: 'flex', alignItems: 'center', gap: '10px', color: '#f59e0b', fontSize: '22px', margin: 0 }}>
                🛡️ Панель глобального администрирования
            </h2>

            <div style={{ display: 'flex', gap: '10px', marginBottom: '10px' }}>
                <button style={tabStyle(activeTab === 'tasks')} onClick={() => setActiveTab('tasks')}>Задачи</button>
                <button style={tabStyle(activeTab === 'events')} onClick={() => setActiveTab('events')}>События</button>
                <button style={tabStyle(activeTab === 'users')} onClick={() => setActiveTab('users')}>Пользователи</button>
            </div>

            <div style={{ marginBottom: '10px' }}>
                <input 
                    type="text" 
                    placeholder="Быстрый фильтр по названию / имени..." 
                    value={searchQuery} 
                    onChange={e => setSearchQuery(e.target.value)} 
                    style={{...inputStyle, maxWidth: '400px'}} 
                />
            </div>

            <div style={{ backgroundColor: '#16161a', padding: '24px', borderRadius: '8px' }}>
                
                {/* ТАБЛИЦА ЗАДАЧ */}
                {activeTab === 'tasks' && (
                    <>
                        <p style={{ color: '#9ca3af', marginBottom: '20px' }}>
                            Всего задач в системе: <strong style={{ color: 'white' }}>{totalTasks}</strong>
                        </p>
                        <div style={{ overflowX: 'auto' }}>
                            <table style={{ width: '100%', borderCollapse: 'collapse' }}>
                                <thead>
                                    <tr style={{ backgroundColor: '#202024' }}>
                                        <th style={thStyle}>Пользователь</th>
                                        <th style={thStyle}>Название</th>
                                        <th style={thStyle}>Приоритет</th>
                                        <th style={thStyle}>Статус</th>
                                        <th style={thStyle}>Действия</th>
                                    </tr>
                                </thead>
                                <tbody>
                                    {filteredTasks.map(task => {
                                        const priorityInfo = TASK_PRIORITY_DETAILS[task.priority] || { label: task.priority, color: '#3b82f6' };
                                        const statusInfo = TASK_STATUS_DETAILS[task.status] || { label: task.status, color: '#9ca3af' };
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
                                                <td style={tdStyle}>
                                                    <span style={{ color: statusInfo.color }}>{statusInfo.label}</span>
                                                </td>
                                                <td style={tdStyle}>
                                                    <button onClick={() => setEditingTask(task)} style={{ padding: '6px 12px', background: '#3b82f6', color: 'white', border: 'none', borderRadius: '4px', cursor: 'pointer' }}>
                                                        ⚙️ Управление
                                                    </button>
                                                </td>
                                            </tr>
                                        );
                                    })}
                                </tbody>
                            </table>
                        </div>
                        <div style={{ display: 'flex', justifyContent: 'space-between', marginTop: '20px' }}>
                            <button disabled={taskPage === 1} onClick={() => setTaskPage(p => p - 1)} style={{ padding: '8px 16px', borderRadius: '4px', cursor: taskPage === 1 ? 'not-allowed' : 'pointer', background: '#202024', color: 'white', border: '1px solid #2d2d38' }}>Назад</button>
                            <span style={{ color: 'white' }}>Страница {taskPage}</span>
                            <button disabled={filteredTasks.length < limit} onClick={() => setTaskPage(p => p + 1)} style={{ padding: '8px 16px', borderRadius: '4px', cursor: filteredTasks.length < limit ? 'not-allowed' : 'pointer', background: '#202024', color: 'white', border: '1px solid #2d2d38' }}>Вперед</button>
                        </div>
                    </>
                )}

                {/* ТАБЛИЦА СОБЫТИЙ */}
                {activeTab === 'events' && (
                    <>
                        <p style={{ color: '#9ca3af', marginBottom: '20px' }}>
                            Всего событий в системе: <strong style={{ color: 'white' }}>{totalEvents}</strong>
                        </p>
                        <div style={{ overflowX: 'auto' }}>
                            <table style={{ width: '100%', borderCollapse: 'collapse' }}>
                                <thead>
                                    <tr style={{ backgroundColor: '#202024' }}>
                                        <th style={thStyle}>Пользователь</th>
                                        <th style={thStyle}>Название</th>
                                        <th style={thStyle}>Тип</th>
                                        <th style={thStyle}>Сроки</th>
                                        <th style={thStyle}>Действия</th>
                                    </tr>
                                </thead>
                                <tbody>
                                    {filteredEvents.map(event => {
                                        const typeDetails = EVENT_TYPE_DETAILS[event.event_type as EventType] || { label: event.event_type, color: '#6366f1' };
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
                                                <td style={tdStyle}>
                                                    <button onClick={() => setEditingEvent(event)} style={{ padding: '6px 12px', background: '#3b82f6', color: 'white', border: 'none', borderRadius: '4px', cursor: 'pointer' }}>
                                                        ⚙️ Управление
                                                    </button>
                                                </td>
                                            </tr>
                                        );
                                    })}
                                </tbody>
                            </table>
                        </div>
                        <div style={{ display: 'flex', justifyContent: 'space-between', marginTop: '20px' }}>
                            <button disabled={eventPage === 1} onClick={() => setEventPage(p => p - 1)} style={{ padding: '8px 16px', borderRadius: '4px', cursor: eventPage === 1 ? 'not-allowed' : 'pointer', background: '#202024', color: 'white', border: '1px solid #2d2d38' }}>Назад</button>
                            <span style={{ color: 'white' }}>Страница {eventPage}</span>
                            <button disabled={filteredEvents.length < limit} onClick={() => setEventPage(p => p + 1)} style={{ padding: '8px 16px', borderRadius: '4px', cursor: filteredEvents.length < limit ? 'not-allowed' : 'pointer', background: '#202024', color: 'white', border: '1px solid #2d2d38' }}>Вперед</button>
                        </div>
                    </>
                )}

                {/* ТАБЛИЦА ПОЛЬЗОВАТЕЛЕЙ */}
                {activeTab === 'users' && (
                    <>
                        <p style={{ color: '#9ca3af', marginBottom: '20px' }}>
                            Всего пользователей: <strong style={{ color: 'white' }}>{filteredUsers.length}</strong>
                        </p>
                        <div style={{ overflowX: 'auto' }}>
                            <table style={{ width: '100%', borderCollapse: 'collapse' }}>
                                <thead>
                                    <tr style={{ backgroundColor: '#202024' }}>
                                        <th style={thStyle}>ID</th>
                                        <th style={thStyle}>Имя</th>
                                        <th style={thStyle}>Email</th>
                                        <th style={thStyle}>Роль</th>
                                        <th style={thStyle}>Действия</th>
                                    </tr>
                                </thead>
                                <tbody>
                                    {/* РЕНДЕРИМ ТОЛЬКО ПАГИНИРОВАННЫХ ПОЛЬЗОВАТЕЛЕЙ */}
                                    {paginatedUsers.map(u => (
                                        <tr key={u.id} style={{ backgroundColor: '#1a1a1e' }}>
                                            <td style={{ ...tdStyle, color: '#6b7280', fontSize: '12px' }}>{u.id.slice(0, 8)}...</td>
                                            <td style={{ ...tdStyle, fontWeight: 'bold' }}>{u.username}</td>
                                            <td style={{ ...tdStyle, color: '#9ca3af' }}>{u.email}</td>
                                            <td style={tdStyle}>
                                                <span style={{ padding: '4px 10px', borderRadius: '4px', fontSize: '12px', fontWeight: 'bold', background: u.user_type === 'admin' ? '#ef444422' : '#10b98122', color: u.user_type === 'admin' ? '#ef4444' : '#10b981' }}>
                                                    {u.user_type === 'admin' ? 'АДМИН' : 'ПОЛЬЗОВАТЕЛЬ'}
                                                </span>
                                            </td>
                                            <td style={tdStyle}>
                                                {u.user_type === 'admin' ? (
                                                    <button onClick={() => handleDemote(u.id)} style={{ padding: '4px 12px', background: '#ef4444', color: 'white', border: 'none', borderRadius: '4px', cursor: 'pointer', fontSize: '12px', fontWeight: 'bold' }}>
                                                        Разжаловать
                                                    </button>
                                                ) : (
                                                    <button onClick={() => handlePromote(u.id)} style={{ padding: '4px 12px', background: '#3b82f6', color: 'white', border: 'none', borderRadius: '4px', cursor: 'pointer', fontSize: '12px', fontWeight: 'bold' }}>
                                                        Сделать админом
                                                    </button>
                                                )}
                                            </td>
                                        </tr>
                                    ))}
                                </tbody>
                            </table>
                        </div>
                        {/* ПАГИНАЦИЯ ДЛЯ ПОЛЬЗОВАТЕЛЕЙ */}
                        <div style={{ display: 'flex', justifyContent: 'space-between', marginTop: '20px' }}>
                            <button disabled={userPage === 1} onClick={() => setUserPage(p => p - 1)} style={{ padding: '8px 16px', borderRadius: '4px', cursor: userPage === 1 ? 'not-allowed' : 'pointer', background: '#202024', color: 'white', border: '1px solid #2d2d38' }}>Назад</button>
                            <span style={{ color: 'white' }}>Страница {userPage} из {totalUsersPages || 1}</span>
                            <button disabled={userPage >= totalUsersPages} onClick={() => setUserPage(p => p + 1)} style={{ padding: '8px 16px', borderRadius: '4px', cursor: userPage >= totalUsersPages ? 'not-allowed' : 'pointer', background: '#202024', color: 'white', border: '1px solid #2d2d38' }}>Вперед</button>
                        </div>
                    </>
                )}
            </div>

            {/* МОДАЛКА РЕДАКТИРОВАНИЯ ЗАДАЧИ */}
            {editingTask && (
                <div style={{ position: 'fixed', top: 0, left: 0, right: 0, bottom: 0, background: 'rgba(0,0,0,0.7)', display: 'flex', justifyContent: 'center', alignItems: 'center', zIndex: 1000 }}>
                    <form onSubmit={handleUpdateTaskSubmit} style={{ background: '#16161a', padding: '24px', borderRadius: '8px', width: '400px', display: 'flex', flexDirection: 'column', gap: '15px', border: '1px solid #2d2d38' }}>
                        <h4 style={{ color: 'white', margin: 0 }}>Управление задачей</h4>
                        <input type="text" style={inputStyle} value={editingTask.title} onChange={e => setEditingTask({...editingTask, title: e.target.value})} required />
                        
                        <label style={{color: '#9ca3af', fontSize: '12px'}}>Приоритет:</label>
                        <select style={inputStyle} value={editingTask.priority} onChange={e => setEditingTask({...editingTask, priority: e.target.value as TaskPriority})}>
                            {Object.entries(TASK_PRIORITY_DETAILS).map(([k, v]) => <option key={k} value={k}>{v.label}</option>)}
                        </select>
                        
                        <label style={{color: '#9ca3af', fontSize: '12px'}}>Статус:</label>
                        <select style={inputStyle} value={editingTask.status} onChange={e => setEditingTask({...editingTask, status: e.target.value as TaskStatus})}>
                            {Object.entries(TASK_STATUS_DETAILS).map(([k, v]) => <option key={k} value={k}>{v.label}</option>)}
                        </select>
                        
                        <div style={{ display: 'flex', justifyContent: 'space-between', marginTop: '15px', paddingTop: '15px', borderTop: '1px solid #2d2d38' }}>
                            <button type="button" onClick={() => setDeleteTarget({id: editingTask.id, type: 'task'})} style={{ padding: '6px 12px', background: 'transparent', border: '1px solid #ef4444', color: '#ef4444', borderRadius: '4px', cursor: 'pointer' }}>Удалить</button>
                            <div style={{ display: 'flex', gap: '10px' }}>
                                <button type="button" onClick={() => setEditingTask(null)} style={{ background: 'none', border: 'none', color: 'white', cursor: 'pointer' }}>Отмена</button>
                                <button type="submit" style={{ padding: '6px 12px', background: '#f59e0b', border: 'none', borderRadius: '4px', color: 'black', fontWeight: 'bold', cursor: 'pointer' }}>Сохранить</button>
                            </div>
                        </div>
                    </form>
                </div>
            )}

            {/* МОДАЛКА РЕДАКТИРОВАНИЯ СОБЫТИЯ */}
            {editingEvent && (
                <div style={{ position: 'fixed', top: 0, left: 0, right: 0, bottom: 0, background: 'rgba(0,0,0,0.7)', display: 'flex', justifyContent: 'center', alignItems: 'center', zIndex: 1000 }}>
                    <form onSubmit={handleUpdateEventSubmit} style={{ background: '#16161a', padding: '24px', borderRadius: '8px', width: '400px', display: 'flex', flexDirection: 'column', gap: '15px', border: '1px solid #2d2d38' }}>
                        <h4 style={{ color: 'white', margin: 0 }}>Управление событием</h4>
                        <input type="text" style={inputStyle} value={editingEvent.title} onChange={e => setEditingEvent({...editingEvent, title: e.target.value})} required />
                        
                        <label style={{color: '#9ca3af', fontSize: '12px'}}>Тип:</label>
                        <select style={inputStyle} value={editingEvent.event_type} onChange={e => setEditingEvent({...editingEvent, event_type: e.target.value as EventType})}>
                            {Object.entries(EVENT_TYPE_DETAILS).map(([k, v]) => <option key={k} value={k}>{v.label}</option>)}
                        </select>
                        
                        <div style={{ display: 'flex', justifyContent: 'space-between', marginTop: '15px', paddingTop: '15px', borderTop: '1px solid #2d2d38' }}>
                            <button type="button" onClick={() => setDeleteTarget({id: editingEvent.id, type: 'event'})} style={{ padding: '6px 12px', background: 'transparent', border: '1px solid #ef4444', color: '#ef4444', borderRadius: '4px', cursor: 'pointer' }}>Удалить</button>
                            <div style={{ display: 'flex', gap: '10px' }}>
                                <button type="button" onClick={() => setEditingEvent(null)} style={{ background: 'none', border: 'none', color: 'white', cursor: 'pointer' }}>Отмена</button>
                                <button type="submit" style={{ padding: '6px 12px', background: '#f59e0b', border: 'none', borderRadius: '4px', color: 'black', fontWeight: 'bold', cursor: 'pointer' }}>Сохранить</button>
                            </div>
                        </div>
                    </form>
                </div>
            )}

            {/* ПОДТВЕРЖДЕНИЕ УДАЛЕНИЯ */}
            {deleteTarget && (
                <div style={{ position: 'fixed', top: 0, left: 0, right: 0, bottom: 0, background: 'rgba(0,0,0,0.85)', display: 'flex', justifyContent: 'center', alignItems: 'center', zIndex: 2000 }}>
                    <div style={{ background: '#16161a', padding: '24px', borderRadius: '8px', width: '350px', border: '1px solid #ef4444', textAlign: 'center' }}>
                        <h3 style={{ color: '#ef4444', marginTop: 0 }}>Внимание!</h3>
                        <p style={{ color: '#9ca3af', fontSize: '14px', marginBottom: '20px' }}>
                            Вы собираетесь удалить данные безвозвратно. Для подтверждения введите слово <strong>УДАЛИТЬ</strong> заглавными буквами.
                        </p>
                        <input 
                            type="text" 
                            style={{...inputStyle, textAlign: 'center', marginBottom: '20px', borderColor: deleteConfirmText === 'УДАЛИТЬ' ? '#10b981' : '#2d2d38'}} 
                            placeholder="Введите УДАЛИТЬ" 
                            value={deleteConfirmText} 
                            onChange={e => setDeleteConfirmText(e.target.value)} 
                        />
                        <div style={{ display: 'flex', gap: '10px', justifyContent: 'center' }}>
                            <button onClick={() => { setDeleteTarget(null); setDeleteConfirmText(''); }} style={{ padding: '8px 16px', background: '#202024', border: '1px solid #2d2d38', color: 'white', borderRadius: '4px', cursor: 'pointer' }}>Отмена</button>
                            <button 
                                onClick={executeDelete} 
                                disabled={deleteConfirmText !== 'УДАЛИТЬ'}
                                style={{ padding: '8px 16px', background: deleteConfirmText === 'УДАЛИТЬ' ? '#ef4444' : '#ef444455', border: 'none', color: 'white', borderRadius: '4px', cursor: deleteConfirmText === 'УДАЛИТЬ' ? 'pointer' : 'not-allowed', fontWeight: 'bold' }}
                            >
                                Подтвердить
                            </button>
                        </div>
                    </div>
                </div>
            )}

        </div>
    );
};