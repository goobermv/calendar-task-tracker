import api from './axios';
import type { Event, CreateEventRequest, UpdateEventRequest, PaginatedEventsResponse } from '../types';

export const eventsApi = {
    getUserEvents: async (page: number, limit: number): Promise<PaginatedEventsResponse> => {
        const response = await api.get<PaginatedEventsResponse>('/events', {
            params: { page, limit }
        });
        return response.data;
    },

    getEventById: async (id: string): Promise<Event> => {
        const response = await api.get<Event>(`/events/${id}`);
        return response.data;
    },

    createEvent: async (data: CreateEventRequest): Promise<Event> => {
        const response = await api.post<Event>('/events', data);
        return response.data;
    },

    updateEvent: async (id: string, data: UpdateEventRequest): Promise<Event> => {
        const response = await api.put<Event>(`/events/${id}`, data);
        return response.data;
    },

    deleteEvent: async (id: string): Promise<void> => {
        await api.delete(`/events/${id}`);
    },

    adminGetAllEvents: async (page: number, limit: number): Promise<PaginatedEventsResponse> => {
        const response = await api.get<PaginatedEventsResponse>('/admin/events', {
            params: { page, limit }
        });
        return response.data;
    },

    adminUpdateEvent: async (id: string, data: UpdateEventRequest): Promise<Event> => {
        const response = await api.put<Event>(`/admin/events/${id}`, data);
        return response.data;
    },

    adminDeleteEvent: async (id: string): Promise<void> => {
        await api.delete(`/admin/events/${id}`);
    }
};