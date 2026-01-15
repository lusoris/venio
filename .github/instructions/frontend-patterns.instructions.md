---
alwaysApply: true
always_on: true
trigger: always_on
applyTo: "web/**/*.ts,web/**/*.tsx,web/**/*.js,web/**/*.jsx"
description: Frontend Patterns & Best Practices
---

# Frontend Patterns & Best Practices

## Core Principle

**Next.js 15 App Router with React 19**: Server Components by default, Client Components only when needed. Type-safe API calls, consistent error handling.

## Technology Stack

- **Framework**: Next.js 15 (App Router)
- **React**: React 19
- **TypeScript**: Latest stable
- **Styling**: Tailwind CSS
- **State Management**: React Context + useState/useReducer
- **API Client**: Fetch API with type-safe wrappers
- **Authentication**: JWT tokens in httpOnly cookies

## Project Structure

```
web/
├── src/
│   ├── app/                    # Next.js App Router pages
│   │   ├── layout.tsx          # Root layout
│   │   ├── page.tsx            # Home page
│   │   ├── (auth)/             # Auth route group
│   │   │   ├── login/
│   │   │   │   └── page.tsx
│   │   │   └── register/
│   │   │       └── page.tsx
│   │   ├── dashboard/          # Protected pages
│   │   │   └── page.tsx
│   │   └── api/                # API routes (if needed)
│   │       └── auth/
│   │           └── route.ts
│   ├── components/             # React components
│   │   ├── ui/                 # Reusable UI components
│   │   ├── forms/              # Form components
│   │   └── layouts/            # Layout components
│   ├── lib/                    # Utilities and helpers
│   │   ├── api.ts              # API client
│   │   ├── auth.ts             # Auth utilities
│   │   └── validators.ts       # Validation functions
│   ├── contexts/               # React contexts
│   │   └── AuthContext.tsx
│   ├── hooks/                  # Custom hooks
│   │   ├── useAuth.ts
│   │   └── useApi.ts
│   └── types/                  # TypeScript types
│       ├── api.ts
│       └── models.ts
```

## Server vs Client Components

### ✅ CORRECT: Server Components by Default

```tsx
// app/dashboard/page.tsx
// ✅ Server Component (default in App Router)
import { getUserData } from '@/lib/api'

export default async function DashboardPage() {
  // ✅ Data fetching in Server Component
  const user = await getUserData()

  return (
    <div>
      <h1>Welcome, {user.name}</h1>
      <UserStats stats={user.stats} />
    </div>
  )
}
```

### ✅ CORRECT: Client Components When Needed

```tsx
// components/LoginForm.tsx
'use client'  // ✅ Use 'use client' only when needed

import { useState } from 'react'
import { useRouter } from 'next/navigation'
import { login } from '@/lib/api'

export function LoginForm() {
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState<string | null>(null)
  const router = useRouter()

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    
    try {
      await login({ email, password })
      router.push('/dashboard')
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Login failed')
    }
  }

  return (
    <form onSubmit={handleSubmit}>
      {/* Form fields */}
    </form>
  )
}
```

### When to Use Client Components

Use `'use client'` directive when:
- Using React hooks (useState, useEffect, useContext)
- Event handlers (onClick, onChange, onSubmit)
- Browser-only APIs (localStorage, window, document)
- Third-party libraries that use hooks or browser APIs

## API Client

### ✅ CORRECT: Type-Safe API Client

```typescript
// lib/api.ts
import { User, LoginRequest, LoginResponse, ErrorResponse } from '@/types/api'

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:3690'

class ApiError extends Error {
  constructor(
    message: string,
    public status: number,
    public code?: string
  ) {
    super(message)
    this.name = 'ApiError'
  }
}

async function fetchApi<T>(
  endpoint: string,
  options?: RequestInit
): Promise<T> {
  const url = `${API_BASE_URL}${endpoint}`
  
  const response = await fetch(url, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      ...options?.headers,
    },
    credentials: 'include',  // ✅ Include cookies (JWT)
  })

  if (!response.ok) {
    const error: ErrorResponse = await response.json()
    throw new ApiError(
      error.message || 'An error occurred',
      response.status,
      error.code
    )
  }

  return response.json()
}

// Authentication
export async function login(data: LoginRequest): Promise<LoginResponse> {
  return fetchApi<LoginResponse>('/api/v1/auth/login', {
    method: 'POST',
    body: JSON.stringify(data),
  })
}

export async function register(data: RegisterRequest): Promise<User> {
  return fetchApi<User>('/api/v1/auth/register', {
    method: 'POST',
    body: JSON.stringify(data),
  })
}

export async function logout(): Promise<void> {
  return fetchApi<void>('/api/v1/auth/logout', {
    method: 'POST',
  })
}

// Users
export async function getCurrentUser(): Promise<User> {
  return fetchApi<User>('/api/v1/users/me')
}

export async function getUsers(): Promise<User[]> {
  return fetchApi<User[]>('/api/v1/users')
}

export async function updateUser(id: number, data: Partial<User>): Promise<User> {
  return fetchApi<User>(`/api/v1/users/${id}`, {
    method: 'PUT',
    body: JSON.stringify(data),
  })
}
```

## Authentication Context

### ✅ CORRECT: Auth Context with Provider

```tsx
// contexts/AuthContext.tsx
'use client'

import { createContext, useContext, useState, useEffect, ReactNode } from 'react'
import { User } from '@/types/models'
import { getCurrentUser, logout as apiLogout } from '@/lib/api'

interface AuthContextType {
  user: User | null
  loading: boolean
  login: (user: User) => void
  logout: () => Promise<void>
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    // Check if user is authenticated on mount
    getCurrentUser()
      .then(setUser)
      .catch(() => setUser(null))
      .finally(() => setLoading(false))
  }, [])

  const login = (user: User) => {
    setUser(user)
  }

  const logout = async () => {
    await apiLogout()
    setUser(null)
  }

  return (
    <AuthContext.Provider value={{ user, loading, login, logout }}>
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  const context = useContext(AuthContext)
  if (context === undefined) {
    throw new Error('useAuth must be used within AuthProvider')
  }
  return context
}
```

### Root Layout with Provider

```tsx
// app/layout.tsx
import { AuthProvider } from '@/contexts/AuthContext'
import './globals.css'

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="en">
      <body>
        <AuthProvider>
          {children}
        </AuthProvider>
      </body>
    </html>
  )
}
```

## Protected Routes

### ✅ CORRECT: Route Protection Middleware

```typescript
// middleware.ts (root of web/)
import { NextResponse } from 'next/server'
import type { NextRequest } from 'next/server'

export function middleware(request: NextRequest) {
  // Check for auth token in cookies
  const token = request.cookies.get('auth_token')

  // Protected routes
  if (request.nextUrl.pathname.startsWith('/dashboard')) {
    if (!token) {
      return NextResponse.redirect(new URL('/login', request.url))
    }
  }

  // Auth routes (redirect to dashboard if already logged in)
  if (request.nextUrl.pathname.startsWith('/login') || 
      request.nextUrl.pathname.startsWith('/register')) {
    if (token) {
      return NextResponse.redirect(new URL('/dashboard', request.url))
    }
  }

  return NextResponse.next()
}

export const config = {
  matcher: ['/dashboard/:path*', '/login', '/register'],
}
```

## Form Validation

### ✅ CORRECT: Client-Side Validation

```tsx
// components/forms/LoginForm.tsx
'use client'

import { useState } from 'react'
import { useRouter } from 'next/navigation'
import { useAuth } from '@/contexts/AuthContext'
import { login } from '@/lib/api'
import { validateEmail, validatePassword } from '@/lib/validators'

export function LoginForm() {
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [errors, setErrors] = useState<Record<string, string>>({})
  const [isSubmitting, setIsSubmitting] = useState(false)
  const router = useRouter()
  const { login: setUser } = useAuth()

  const validate = () => {
    const newErrors: Record<string, string> = {}

    const emailError = validateEmail(email)
    if (emailError) newErrors.email = emailError

    const passwordError = validatePassword(password)
    if (passwordError) newErrors.password = passwordError

    setErrors(newErrors)
    return Object.keys(newErrors).length === 0
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    if (!validate()) return

    setIsSubmitting(true)
    try {
      const response = await login({ email, password })
      setUser(response.user)
      router.push('/dashboard')
    } catch (err) {
      setErrors({ 
        form: err instanceof Error ? err.message : 'Login failed' 
      })
    } finally {
      setIsSubmitting(false)
    }
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <div>
        <label htmlFor="email" className="block text-sm font-medium">
          Email
        </label>
        <input
          id="email"
          type="email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          className="mt-1 block w-full rounded-md border-gray-300"
          disabled={isSubmitting}
        />
        {errors.email && (
          <p className="mt-1 text-sm text-red-600">{errors.email}</p>
        )}
      </div>

      <div>
        <label htmlFor="password" className="block text-sm font-medium">
          Password
        </label>
        <input
          id="password"
          type="password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          className="mt-1 block w-full rounded-md border-gray-300"
          disabled={isSubmitting}
        />
        {errors.password && (
          <p className="mt-1 text-sm text-red-600">{errors.password}</p>
        )}
      </div>

      {errors.form && (
        <div className="rounded-md bg-red-50 p-4">
          <p className="text-sm text-red-800">{errors.form}</p>
        </div>
      )}

      <button
        type="submit"
        disabled={isSubmitting}
        className="w-full rounded-md bg-blue-600 px-4 py-2 text-white hover:bg-blue-700 disabled:opacity-50"
      >
        {isSubmitting ? 'Logging in...' : 'Log in'}
      </button>
    </form>
  )
}
```

## Error Handling

### ✅ CORRECT: Error Boundaries

```tsx
// components/ErrorBoundary.tsx
'use client'

import { Component, ErrorInfo, ReactNode } from 'react'

interface Props {
  children: ReactNode
  fallback?: ReactNode
}

interface State {
  hasError: boolean
  error?: Error
}

export class ErrorBoundary extends Component<Props, State> {
  constructor(props: Props) {
    super(props)
    this.state = { hasError: false }
  }

  static getDerivedStateFromError(error: Error): State {
    return { hasError: true, error }
  }

  componentDidCatch(error: Error, errorInfo: ErrorInfo) {
    console.error('Error caught by boundary:', error, errorInfo)
  }

  render() {
    if (this.state.hasError) {
      return this.props.fallback || (
        <div className="rounded-md bg-red-50 p-4">
          <h2 className="text-lg font-semibold text-red-800">
            Something went wrong
          </h2>
          <p className="mt-2 text-sm text-red-600">
            {this.state.error?.message || 'An unexpected error occurred'}
          </p>
        </div>
      )
    }

    return this.props.children
  }
}
```

## Loading States

### ✅ CORRECT: Loading UI

```tsx
// app/dashboard/loading.tsx
export default function Loading() {
  return (
    <div className="flex items-center justify-center min-h-screen">
      <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-gray-900"></div>
    </div>
  )
}

// app/dashboard/page.tsx
import { Suspense } from 'react'
import { UserList } from '@/components/UserList'
import Loading from './loading'

export default function DashboardPage() {
  return (
    <div>
      <h1>Dashboard</h1>
      <Suspense fallback={<Loading />}>
        <UserList />
      </Suspense>
    </div>
  )
}
```

## TypeScript Types

### ✅ CORRECT: Shared Types

```typescript
// types/models.ts
export interface User {
  id: number
  email: string
  username: string
  firstName: string
  lastName: string
  isActive: boolean
  createdAt: string
  updatedAt: string
}

export interface Role {
  id: number
  name: string
  description: string
}

// types/api.ts
export interface LoginRequest {
  email: string
  password: string
}

export interface LoginResponse {
  token: string
  user: User
}

export interface ErrorResponse {
  error: string
  message: string
  code?: string
}

export interface PaginatedResponse<T> {
  data: T[]
  total: number
  page: number
  pageSize: number
}
```

## Environment Variables

### ✅ CORRECT: Type-Safe Env Vars

```typescript
// lib/env.ts
function getEnvVar(key: string, defaultValue?: string): string {
  const value = process.env[key]
  if (!value && !defaultValue) {
    throw new Error(`Missing required environment variable: ${key}`)
  }
  return value || defaultValue!
}

export const env = {
  apiUrl: getEnvVar('NEXT_PUBLIC_API_URL', 'http://localhost:3690'),
  appName: getEnvVar('NEXT_PUBLIC_APP_NAME', 'Venio'),
  environment: getEnvVar('NODE_ENV', 'development'),
} as const
```

### .env.local

```bash
# .env.local (development)
NEXT_PUBLIC_API_URL=http://localhost:3690
NEXT_PUBLIC_APP_NAME=Venio

# .env.production (production)
NEXT_PUBLIC_API_URL=https://api.venio.io
NEXT_PUBLIC_APP_NAME=Venio
```

## Common Mistakes

### ❌ DON'T Do This

```tsx
// 1. Using 'use client' unnecessarily
'use client'  // ❌ Not needed if no client features
export default function Page() {
  return <div>Static content</div>
}

// 2. Fetching in useEffect in Server Components
export default function Page() {
  const [data, setData] = useState()
  useEffect(() => {
    fetch('/api/data').then(setData)  // ❌ Use async Server Component
  }, [])
}

// 3. Not handling loading states
const data = await fetchData()  // ❌ No loading UI

// 4. No error handling
const data = await fetchData()  // ❌ What if it fails?

// 5. Exposing secrets in client
const API_KEY = process.env.API_KEY  // ❌ Exposed to client
// Use NEXT_PUBLIC_ prefix only for non-sensitive config
```

---

## Document Version

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0.0 | 2026-01-15 | AI Assistant | Initial frontend patterns guide |

**Remember**: Server Components by default. Use Client Components only when necessary. Type-safe API calls prevent runtime errors.
