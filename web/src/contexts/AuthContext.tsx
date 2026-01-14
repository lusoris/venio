'use client';

import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import { apiClient, User, LoginRequest, RegisterRequest } from '@/lib/api';
import { useRouter } from 'next/navigation';

interface AuthContextType {
  user: User | null;
  loading: boolean;
  token: string | null;
  login: (data: LoginRequest) => Promise<void>;
  register: (data: RegisterRequest) => Promise<void>;
  logout: () => void;
  isAuthenticated: boolean;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);
  const [token, setToken] = useState<string | null>(null);
  const router = useRouter();

  useEffect(() => {
    // Check if user is already logged in
    const checkAuth = async () => {
      const storedToken = localStorage.getItem('access_token');
      if (storedToken) {
        try {
          setToken(storedToken);
          // Try to fetch current user or validate token
          // For now, we'll just set loading to false
          setLoading(false);
        } catch (error) {
          console.error('Auth check failed:', error);
          apiClient.logout();
          setToken(null);
          setLoading(false);
        }
      } else {
        setLoading(false);
      }
    };

    checkAuth();
  }, []);

  const login = async (data: LoginRequest) => {
    try {
      const response = await apiClient.login(data);
      setUser(response.user);
      const storedToken = localStorage.getItem('access_token');
      setToken(storedToken);
      router.push('/dashboard');
    } catch (error) {
      throw error;
    }
  };

  const register = async (data: RegisterRequest) => {
    try {
      await apiClient.register(data);
      // Auto-login after registration
      await login({ email: data.email, password: data.password });
    } catch (error) {
      throw error;
    }
  };

  const logout = () => {
    apiClient.logout();
    setUser(null);
    setToken(null);
    router.push('/login');
  };

  return (
    <AuthContext.Provider
      value={{
        user,
        loading,
        token,
        login,
        register,
        logout,
        isAuthenticated: !!user,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
}
