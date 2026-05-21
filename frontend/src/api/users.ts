import api from './axios';
import type { 
    User, 
    RegisterRequest, 
    LoginRequest, 
    AuthResponse, 
    UpdateUserInfoRequest, 
    UpdateUserPasswordRequest 
} from '../types';

export const usersApi = {
    register: async (data: RegisterRequest): Promise<AuthResponse> => {
        const response = await api.post<AuthResponse>('/auth/register', data);
        return response.data;
    },

    login: async (data: LoginRequest): Promise<AuthResponse> => {
        const response = await api.post<AuthResponse>('/auth/login', data);
        return response.data;
    },

    getProfile: async (): Promise<User> => {
        const response = await api.get<User>('/users/me');
        return response.data;
    },

    updateInfo: async (userId: string, data: UpdateUserInfoRequest): Promise<User> => {
        const response = await api.patch<User>(`/users/${userId}`, data);
        return response.data;
    },

    updatePassword: async (userId: string, data: UpdateUserPasswordRequest): Promise<void> => {
        await api.patch(`/users/${userId}/password`, data);
    },

    deleteAccount: async (userId: string): Promise<void> => {
        await api.delete(`/users/${userId}`);
    },

    promoteToAdmin: async (userId: string): Promise<void> => {
        await api.post(`/admin/users/${userId}/promote`);
    },

    demoteToUser: async (userId: string): Promise<void> => {
        await api.post(`/admin/users/${userId}/demote`);
    },

    getAllUsers: async () => {
        const res = await api.get('/admin/users');
        return res.data;
    }
};