export interface User {
    id: string;
    email: string;
    username: string;
    user_type: 'regular' | 'admin';
    created_at?: string;
    updated_at?: string;
}

export interface RegisterRequest {
    email: string;
    username: string;
    password: string;
}

export interface LoginRequest {
    email: string;
    password: string;
}

export interface AuthResponse {
    token: string;
    user: User;
}

export interface UpdateUserInfoRequest {
    email?: string;
    username?: string;
}

export interface UpdateUserPasswordRequest {
    password?: string;
}

export type TaskPriority = 'low' | 'medium' | 'high' | 'urgent';
export type TaskStatus = 'pending' | 'in_progress' | 'completed' | 'cancelled';

export interface Task {
    id: string;
    user_id: string;
    title: string;
    description: string;
    status: TaskStatus;
    priority: TaskPriority;
    due_date?: string;
    created_at: string;
    updated_at: string;
}

export interface PaginatedTasksResponse {
    tasks: Task[];
    total_count: number;
}

export interface CreateTaskRequest {
    title: string;
    description?: string;
    priority: TaskPriority;
    due_date?: string;
}

export interface UpdateTaskRequest {
    title?: string;
    description?: string;
    priority?: TaskPriority;
    status?: TaskStatus;
    due_date?: string;
}

export type EventType = 'formal' | 'celebration' | 'important' | 'personal' | 'holiday' | 'casual' | 'business';

export interface Event {
    id: string;
    user_id: string;
    title: string;
    description: string;
    start_time: string;
    end_time: string;
    event_type: EventType;
    created_at: string;
    updated_at: string;
}

export interface PaginatedEventsResponse {
    events: Event[];
    total_count: number;
}

export interface CreateEventRequest {
    title: string;
    description?: string;
    start_time: string;
    end_time: string;
    event_type: EventType;
}

export interface UpdateEventRequest {
    title?: string;
    description?: string;
    start_time?: string;
    end_time?: string;
    event_type?: EventType;
}