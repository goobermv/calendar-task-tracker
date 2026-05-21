import React, { useEffect, useState } from 'react';
import { eventsApi } from '../api/events';
import type { Event, EventType } from '../types';

const EVENT_TYPE_DETAILS: Record<EventType, { label: string; color: string }> = {
    personal: { label: 'Личное', color: '#3b82f6' },
    business: { label: 'Бизнес / Работа', color: '#10b981' },
    formal: { label: 'Официальное', color: '#6366f1' },
    important: { label: 'Важное', color: '#ef4444' },
    celebration: { label: 'Праздник / Торжество', color: '#f59e0b' },
    holiday: { label: 'Выходной', color: '#ec4899' },
    casual: { label: 'Повседневное', color: '#6b7280' }
};

export const CalendarPage: React.FC = () => {
    const [events, setEvents] = useState<Event[]>([]);
    const [title, setTitle] = useState('');
    const [eventType, setEventType] = useState<EventType>('personal');
    const [startTime, setStartTime] = useState('')
    const [endTime, setEndTime] = useState('');

    const [editingEventId, setEditingEventId] = useState<string | null>(null);

    const loadEvents = async () => {
        try {
            const res = await eventsApi.getUserEvents(1, 100);
            const eventsArray = Array.isArray(res) ? res : ((res as any).events || []);
            setEvents(eventsArray);
        } catch (err) {
            console.error("Ошибка загрузки событий", err);
        }
    };

    useEffect(() => {
        loadEvents();
    }, []);

    const resetForm = () => {
        setTitle('');
        setEventType('personal');
        setStartTime('');
        setEndTime('');
        setEditingEventId(null);
    };

    const handleFormSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!title || !startTime || !endTime) return;

    const start = new Date(startTime);
    const end = new Date(endTime);

    if (start >= end) {
        alert("Ошибка: Дата и время начала не могут быть позже или равны времени окончания события!");
        return;
    }

    const eventData = {
        title,
        event_type: eventType,
        start_time: start.toISOString(),
        end_time: end.toISOString()
    };

    try {
        if (editingEventId) {
            await eventsApi.updateEvent(editingEventId, eventData);
        } else {
            await eventsApi.createEvent(eventData);
        }
        resetForm();
        loadEvents();
    } catch (err) {
        alert(editingEventId ? "Не удалось обновить событие" : "Не удалось создать событие");
    }
};

    const handleEditStart = (ev: Event) => {
        const id = ev.id || (ev as any).ID;
        const startRaw = ev.start_time || (ev as any).StartTime;
        const endRaw = ev.end_time || (ev as any).EndTime;
        const currentType = ev.event_type || (ev as any).EventType || 'personal';

        setEditingEventId(id);
        setTitle(ev.title || (ev as any).Title || '');
        setEventType(currentType as EventType);

        if (startRaw) setStartTime(new Date(startRaw).toISOString().slice(0, 16));
        if (endRaw) setEndTime(new Date(endRaw).toISOString().slice(0, 16));
    };

    const handleDeleteEvent = async (id: string) => {
        if (!window.confirm("Вы уверены, что хотите удалить это событие?")) return;

        try {
            await eventsApi.deleteEvent(id);
            if (editingEventId === id) resetForm();
            loadEvents();
        } catch (err) {
            alert("Не удалось удалить событие");
        }
    };

    return (
        <div>
            <h2>📅 Календарное расписание</h2>

            <div style={{ display: 'grid', gridTemplateColumns: '1fr 2fr', gap: '30px', marginTop: '20px' }}>
                {/* Форма создания / Редактирования */}
                <form onSubmit={handleFormSubmit} style={{ backgroundColor: 'var(--panel-bg)', padding: '24px', borderRadius: '8px', height: 'fit-content', display: 'flex', flexDirection: 'column', gap: '15px' }}>
                    <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                        <h4>{editingEventId ? 'Редактирование события' : 'Новое событие'}</h4>
                        {editingEventId && (
                            <button type="button" onClick={resetForm} style={{ background: 'none', border: 'none', color: '#ef4444', cursor: 'pointer', fontSize: '13px' }}>
                                Отмена
                            </button>
                        )}
                    </div>
                    <input type="text" placeholder="Название встречи/события" value={title} onChange={e => setTitle(e.target.value)} required />
                    <select value={eventType} onChange={e => setEventType(e.target.value as EventType)}>
                        {Object.entries(EVENT_TYPE_DETAILS).map(([key, value]) => (
                            <option key={key} value={key}>{value.label}</option>
                        ))}
                    </select>

                    <div>
                        <label style={{ fontSize: '12px', color: 'var(--text-muted)' }}>Начало</label>
                        <input type="datetime-local" value={startTime} onChange={e => setStartTime(e.target.value)} required />
                    </div>
                    <div>
                        <label style={{ fontSize: '12px', color: 'var(--text-muted)' }}>Окончание</label>
                        <input type="datetime-local" value={endTime} onChange={e => setEndTime(e.target.value)} required />
                    </div>
                    <button type="submit" style={{ backgroundColor: editingEventId ? '#d97706' : 'var(--primary)' }}>
                        {editingEventId ? 'Сохранить изменения' : 'Запланировать'}
                    </button>
                </form>

                {/* Список событий */}
                <div style={{ backgroundColor: 'var(--panel-bg)', padding: '24px', borderRadius: '8px' }}>
                    <h4>Предстоящие события ({events.length})</h4>
                    <div style={{ display: 'flex', flexDirection: 'column', gap: '15px', marginTop: '15px' }}>
                        {events.length === 0 ? <p style={{ color: 'var(--text-muted)' }}>Событий пока нет</p> : 
                        events.map(ev => {
                            const currentId = ev.id || (ev as any).ID;
                            const currentType = (ev.event_type || (ev as any).EventType || "personal") as EventType;
                            const typeDetails = EVENT_TYPE_DETAILS[currentType] || EVENT_TYPE_DETAILS.personal;

                            return (
                                <div 
                                    key={currentId} 
                                    style={{ 
                                        padding: '15px', 
                                        backgroundColor: '#202024', 
                                        borderRadius: '6px', 
                                        borderLeft: `4px solid ${typeDetails.color}`,
                                        display: 'flex',
                                        justifyContent: 'space-between',
                                        alignItems: 'flex-start'
                                    }}
                                >
                                    <div style={{ flex: 1, marginRight: '15px' }}>
                                        <div style={{ display: 'flex', alignItems: 'center', gap: '10px', flexWrap: 'wrap' }}>
                                            <strong style={{ color: 'white' }}>
                                                {ev.title || (ev as any).Title || "Без названия"}
                                            </strong>
                                            <span style={{ fontSize: '11px', padding: '2px 8px', borderRadius: '10px', background: `${typeDetails.color}22`, color: typeDetails.color, fontWeight: 'medium' }}>
                                                {typeDetails.label}
                                            </span>
                                        </div>
                                        <div style={{ fontSize: '12px', color: 'var(--text-muted)', marginTop: '8px' }}>
                                            {(() => {
                                                const startRaw = ev.start_time || (ev as any).StartTime || (ev as any).start;
                                                const endRaw = ev.end_time || (ev as any).EndTime || (ev as any).end;

                                                const start = startRaw ? new Date(startRaw) : null;
                                                const end = endRaw ? new Date(endRaw) : null;
                                                if (!start || isNaN(start.getTime()) || !end || isNaN(end.getTime())) {
                                                    return <span>⏰ Время не указано</span>;
                                                }
                                                return `⏰ ${start.toLocaleString()} — ${end.toLocaleString()}`;
                                            })()}
                                        </div>
                                    </div>

                                    {/* Блок управления кнопками (Редактировать / Удалить) */}
                                    <div style={{ display: 'flex', gap: '8px' }}>
                                        <button 
                                            type="button" 
                                            onClick={() => handleEditStart(ev)}
                                            style={{ background: 'none', border: 'none', color: 'var(--text-muted)', cursor: 'pointer', fontSize: '12px', padding: '4px' }}
                                            onMouseEnter={(e) => e.currentTarget.style.color = '#f59e0b'}
                                            onMouseLeave={(e) => e.currentTarget.style.color = 'var(--text-muted)'}
                                            title="Редактировать"
                                        >
                                            ✏️
                                        </button>
                                        <button 
                                            type="button" 
                                            onClick={() => handleDeleteEvent(currentId)}
                                            style={{ background: 'none', border: 'none', color: 'var(--text-muted)', cursor: 'pointer', fontSize: '12px', padding: '4px' }}
                                            onMouseEnter={(e) => e.currentTarget.style.color = '#ef4444'}
                                            onMouseLeave={(e) => e.currentTarget.style.color = 'var(--text-muted)'}
                                            title="Удалить"
                                        >
                                            🗑️
                                        </button>
                                    </div>
                                </div>
                            );
                        })}
                    </div>
                </div>
            </div>
        </div>
    );
}; 