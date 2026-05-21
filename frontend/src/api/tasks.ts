import api from './axios';
import type { Task, CreateTaskRequest, UpdateTaskRequest, PaginatedTasksResponse } from '../types';

export const tasksApi = {
    getUserTasks: async (page: number, limit: number): Promise<PaginatedTasksResponse> => {
        const response = await api.get<PaginatedTasksResponse>('/tasks', {
            params: { page, limit }
        });
        return response.data;
    },

    getTaskById: async (id: string): Promise<Task> => {
        const response = await api.get<Task>(`/tasks/${id}`);
        return response.data;
    },

    createTask: async (data: CreateTaskRequest): Promise<Task> => {
        const response = await api.post<Task>('/tasks', data);
        return response.data;
    },

    updateTask: async (id: string, data: UpdateTaskRequest): Promise<Task> => {
        const response = await api.put<Task>(`/tasks/${id}`, data);
        return response.data;
    },

    deleteTask: async (id: string): Promise<void> => {
        await api.delete(`/tasks/${id}`);
    },

    adminGetAllTasks: async (page: number, limit: number): Promise<PaginatedTasksResponse> => {
        const response = await api.get<PaginatedTasksResponse>('/admin/tasks', {
            params: { page, limit }
        });
        return response.data;
    },

    adminUpdateTask: async (id: string, data: UpdateTaskRequest): Promise<Task> => {
        const response = await api.put<Task>(`/admin/tasks/${id}`, data);
        return response.data;
    },

    adminDeleteTask: async (id: string): Promise<void> => {
        await api.delete(`/admin/tasks/${id}`);
    }
};