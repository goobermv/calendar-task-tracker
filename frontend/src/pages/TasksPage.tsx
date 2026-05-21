import React, { useEffect, useState } from 'react';
import { tasksApi } from '../api/tasks';
import type { Task, TaskStatus, TaskPriority } from '../types';

export const TasksPage: React.FC = () => {
    const [tasks, setTasks] = useState<Task[]>([]);
    const [title, setTitle] = useState('');
    const [priority, setPriority] = useState<TaskPriority>('medium');
    const [dueDate, setDueDate] = useState('');

    const [editingTask, setEditingTask] = useState<Task | null>(null);
    const [editTitle, setEditTitle] = useState('');
    const [editDueDate, setEditDueDate] = useState('');

    const loadTasks = async () => {
        try {
            const res = await tasksApi.getUserTasks(1, 100);
            setTasks(res.tasks || []);
        } catch (err) { 
            console.error("Ошибка загрузки:", err); 
        }
    };

    useEffect(() => { 
        loadTasks(); 
    }, []);

    const handleCreateTask = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!title.trim()) return;

        if (dueDate) {
            const selectedDate = new Date(dueDate);
            const today = new Date();
            today.setHours(0, 0, 0, 0);

            if (selectedDate < today) {
                alert("Нельзя запланировать задачу на прошедшее число!");
                return;
            }
        }

        try {
            await tasksApi.createTask({
                title,
                priority,
                due_date: dueDate ? new Date(dueDate).toISOString() : undefined
            });
            setTitle('');
            setDueDate('');
            loadTasks();
        } catch (err) {
            alert("Не удалось создать задачу");
        }
    };

    const handleUpdate = async (taskId: string, updates: any) => {
        try {
            if (updates && typeof updates.status === 'string') {
                updates.status = updates.status.trim();
            }
            await tasksApi.updateTask(taskId, updates);
            loadTasks();
        } catch (err) { 
            alert("Ошибка обновления задачи. Возможно, срок задачи находится в прошлом!"); 
        }
    };

    const handleEditClick = (task: Task) => {
        setEditingTask(task);
        setEditTitle(task.title);
        
        const currentTaskDate = task.due_date && task.due_date !== "0001-01-01T00:00:00Z" 
            ? task.due_date.split('T')[0] 
            : "";
        setEditDueDate(currentTaskDate);
    };

    const handleEditSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!editingTask) return;

        const updates: any = { title: editTitle };
        if (editDueDate.trim() !== "") {
            updates.due_date = new Date(editDueDate).toISOString();
        }

        await handleUpdate(editingTask.id, updates);
        setEditingTask(null);
    };

    const handleDelete = async (taskId: string) => {
        if (!confirm("Удалить задачу?")) return;
        try {
            await tasksApi.deleteTask(taskId);
            loadTasks();
        } catch (err) {
            alert("Не удалось удалить задачу");
        }
    };

    const columns: { title: string; status: TaskStatus; color: string }[] = [
        { title: 'Ожидают', status: 'pending', color: 'var(--text-muted)' },
        { title: 'В процессе', status: 'in_progress', color: 'var(--primary)' },
        { title: 'Выполнены', status: 'completed', color: 'var(--success)' },
        { title: 'Отменены', status: 'cancelled', color: 'var(--danger)' }
    ];

    const btnBase: React.CSSProperties = { 
        fontSize: '11px', 
        padding: '6px 8px', 
        cursor: 'pointer', 
        border: 'none', 
        borderRadius: '4px', 
        color: '#fff',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        flex: '1 1 auto',
        minWidth: '32px'
    };

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

    return (
        <div>
            <h2>📋 Мой Трекер Задач</h2>
            {/* Форма быстрого создания */}
            <form onSubmit={handleCreateTask} style={{ display: 'flex', gap: '10px', backgroundColor: 'var(--panel-bg)', padding: '20px', borderRadius: '8px', marginBottom: '30px', alignItems: 'flex-end' }}>
                <div style={{ flex: 2 }}>
                    <label style={{ display: 'block', marginBottom: '5px', fontSize: '12px' }}>Название задачи</label>
                    <input type="text" value={title} onChange={e => setTitle(e.target.value)} placeholder="Что нужно сделать?" required />
                </div>
                <div style={{ flex: 1 }}>
                    <label style={{ display: 'block', marginBottom: '5px', fontSize: '12px' }}>Приоритет</label>
                    <select value={priority} onChange={e => setPriority(e.target.value as TaskPriority)}>
                        <option value="low">Низкий</option>
                        <option value="medium">Средний</option>
                        <option value="high">Высокий</option>
                        <option value="urgent">Срочный</option>
                    </select>
                </div>
                <div style={{ flex: 1 }}>
                    <label style={{ display: 'block', marginBottom: '5px', fontSize: '12px' }}>Срок исполнения</label>
                    <input type="date" value={dueDate} onChange={e => setDueDate(e.target.value)} />
                </div>
                <button type="submit">Добавить</button>
            </form>

            {/* Канбан Сетка */}
            <div style={{ display: 'grid', gridTemplateColumns: 'repeat(4, 1fr)', gap: '15px' }}>
                {columns.map(col => (
                    <div key={col.status} style={{ backgroundColor: 'var(--panel-bg)', padding: '15px', borderRadius: '8px', minHeight: '500px' }}>
                        <h4 style={{ margin: '0 0 15px 0', borderBottom: `2px solid ${col.color}`, paddingBottom: '10px' }}>
                            {col.title} ({tasks.filter(t => t.status === col.status).length})
                        </h4>
                        <div style={{ display: 'flex', flexDirection: 'column', gap: '10px' }}>
                            {tasks.filter(t => t.status === col.status).map(task => {
                                let priorityColor = '#8a8f98';
                                if (task.priority === 'medium') priorityColor = '#4a90e2';
                                if (task.priority === 'high') priorityColor = '#f5a623';
                                if (task.priority === 'urgent') priorityColor = 'var(--danger)';

                                return (
                                    <div key={task.id} style={{ 
                                        backgroundColor: '#202024', 
                                        padding: '12px', 
                                        borderRadius: '6px', 
                                        borderLeft: `4px solid ${priorityColor}`
                                    }}>
                                        <div style={{ marginBottom: '4px', fontWeight: 'bold', color: 'white' }}>{task.title}</div>
                                        {task.due_date && task.due_date !== "0001-01-01T00:00:00Z" && (
                                            <div style={{ fontSize: '11px', color: 'var(--text-muted)', marginBottom: '10px' }}>
                                                📅 до: {new Date(task.due_date).toLocaleDateString()}
                                            </div>
                                        )}
                                        {/* Группа кнопок с исправленной передачей статуса */}
                                        <div style={{ display: 'flex', gap: '4px', marginTop: '10px', width: '100%', justifyContent: 'flex-start' }}>
                                            {task.status !== 'in_progress' && (
                                                <button type="button" title="В работу" style={{...btnBase, backgroundColor: '#3b82f6'}} 
                                                        onClick={() => handleUpdate(task.id, { status: 'in_progress' })}>▶</button>
                                            )}
                                            {task.status !== 'completed' && (
                                                <button type="button" title="Готово" style={{...btnBase, backgroundColor: '#10b981'}} 
                                                        onClick={() => handleUpdate(task.id, { status: 'completed' })}>✓</button>
                                            )}
                                            {task.status !== 'cancelled' && (
                                                <button type="button" title="Отменить" style={{...btnBase, backgroundColor: '#666'}} 
                                                        onClick={() => handleUpdate(task.id, { status: 'cancelled' })}>✕</button>
                                            )}
                                            <button type="button" title="Редактировать" style={{...btnBase, backgroundColor: '#4a90e2'}} onClick={() => handleEditClick(task)}>✎</button>
                                            <button type="button" title="Удалить" style={{...btnBase, backgroundColor: '#d9534f', marginLeft: 'auto'}} onClick={() => handleDelete(task.id)}>🗑️</button>
                                        </div>
                                    </div>
                                );
                            })}
                        </div>
                    </div>
                ))}
            </div>

            {/* Модальное окно редактирования задачи (в стиле Layout) */}
            {editingTask && (
                <div style={{ position: 'fixed', top: 0, left: 0, right: 0, bottom: 0, backgroundColor: 'rgba(0,0,0,0.7)', display: 'flex', alignItems: 'center', justifyContent: 'center', zIndex: 1000 }}>
                    <form onSubmit={handleEditSubmit} style={{ backgroundColor: '#1e1e1e', padding: '30px', borderRadius: '12px', display: 'flex', flexDirection: 'column', gap: '15px', color: 'white', width: '400px' }}>
                        <h3>Редактировать задачу</h3>
                        
                        <div style={{ display: 'flex', flexDirection: 'column', gap: '5px' }}>
                            <label>Название задачи</label>
                            <input 
                                value={editTitle} 
                                onChange={e => setEditTitle(e.target.value)} 
                                required 
                                style={{ padding: '8px' }} 
                            />
                        </div>

                        <div style={{ display: 'flex', flexDirection: 'column', gap: '5px' }}>
                            <label>Срок исполнения</label>
                            <input 
                                type="date" 
                                value={editDueDate} 
                                onChange={e => setEditDueDate(e.target.value)} 
                                style={{ padding: '8px' }} 
                            />
                        </div>

                        <div style={{ display: 'flex', gap: '10px', marginTop: '10px' }}>
                            <button type="submit" style={modalButtonStyle(true)}>Сохранить</button>
                            <button type="button" onClick={() => setEditingTask(null)} style={modalButtonStyle(false)}>Отмена</button>
                        </div>
                    </form>
                </div>
            )}
        </div>
    );
};