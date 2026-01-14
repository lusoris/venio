// API Client Configuration
const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:3690/api/v1';

export interface User {
  id: number;
  email: string;
  username: string;
  first_name: string;
  last_name: string;
  avatar?: string;
  is_active: boolean;
  roles?: string[];
  created_at: string;
  updated_at: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  email: string;
  username: string;
  first_name: string;
  last_name: string;
  password: string;
}

export interface LoginResponse {
  access_token: string;
  refresh_token: string;
  user: User;
}

export interface ErrorResponse {
  error: string;
  message: string;
}

class ApiClient {
  private baseUrl: string;
  private accessToken: string | null = null;

  constructor(baseUrl: string) {
    this.baseUrl = baseUrl;
    // Load token from localStorage on client-side
    if (typeof window !== 'undefined') {
      this.accessToken = localStorage.getItem('access_token');
    }
  }

  setAccessToken(token: string) {
    this.accessToken = token;
    if (typeof window !== 'undefined') {
      localStorage.setItem('access_token', token);
    }
  }

  clearTokens() {
    this.accessToken = null;
    if (typeof window !== 'undefined') {
      localStorage.removeItem('access_token');
      localStorage.removeItem('refresh_token');
    }
  }

  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<T> {
    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
      ...(options.headers as Record<string, string>),
    };

    if (this.accessToken) {
      headers['Authorization'] = `Bearer ${this.accessToken}`;
    }

    const response = await fetch(`${this.baseUrl}${endpoint}`, {
      ...options,
      headers,
    });

    if (!response.ok) {
      const error: ErrorResponse = await response.json();
      throw new Error(error.message || 'API request failed');
    }

    return response.json();
  }

  // Auth endpoints
  async register(data: RegisterRequest): Promise<User> {
    return this.request<User>('/auth/register', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  async login(data: LoginRequest): Promise<LoginResponse> {
    const response = await this.request<LoginResponse>('/auth/login', {
      method: 'POST',
      body: JSON.stringify(data),
    });

    this.setAccessToken(response.access_token);
    if (typeof window !== 'undefined') {
      localStorage.setItem('refresh_token', response.refresh_token);
    }

    return response;
  }

  async refreshToken(): Promise<{ access_token: string }> {
    const refreshToken =
      typeof window !== 'undefined'
        ? localStorage.getItem('refresh_token')
        : null;

    if (!refreshToken) {
      throw new Error('No refresh token available');
    }

    const response = await this.request<{ access_token: string }>(
      '/auth/refresh',
      {
        method: 'POST',
        body: JSON.stringify({ refresh_token: refreshToken }),
      }
    );

    this.setAccessToken(response.access_token);
    return response;
  }

  logout() {
    this.clearTokens();
  }

  // User endpoints
  async getUsers(limit = 10, offset = 0): Promise<User[]> {
    return this.request<User[]>(`/users?limit=${limit}&offset=${offset}`);
  }

  async getUser(id: number): Promise<User> {
    return this.request<User>(`/users/${id}`);
  }

  async updateUser(
    id: number,
    data: Partial<Omit<User, 'id' | 'created_at' | 'updated_at'>>
  ): Promise<User> {
    return this.request<User>(`/users/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    });
  }

  async deleteUser(id: number): Promise<void> {
    await this.request<void>(`/users/${id}`, {
      method: 'DELETE',
    });
  }

  // Health check
  async healthCheck(): Promise<{ app: string; status: string; version: string }> {
    const response = await fetch(`${this.baseUrl.replace('/api/v1', '')}/health`);
    return response.json();
  }
}

export const apiClient = new ApiClient(API_BASE_URL);
