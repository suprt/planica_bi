# Frontend API Integration Plan

**Дата:** 2025-11-16  
**Статус:** В процессе  
**Backend API:** http://localhost:8080/api  
**Frontend:** http://localhost:3000

---

## 📊 Текущее состояние

### ✅ Backend (готов на 98%)
- ✅ REST API полностью реализован
- ✅ JWT авторизация работает
- ✅ Все endpoints протестированы
- ✅ Запущен в Docker на порту 8080

### 🔶 Frontend (UI готов, API нет)
- ✅ React 19.2 + TypeScript
- ✅ React Router работает
- ✅ UI компоненты готовы (Dashboard, Statistics)
- ❌ Нет интеграции с API
- ❌ Захардкоженные данные

---

## 📋 TODO List (13 задач)

### Phase 1: Инфраструктура (задачи 1-3)

#### ✅ [PENDING] 1. Создать API клиент (axios)
**Файл:** `frontend/src/services/api/apiClient.ts`

```typescript
import axios from 'axios';

const apiClient = axios.create({
    baseURL: process.env.REACT_APP_API_URL || 'http://localhost:8080/api',
    timeout: 30000,
    headers: {
        'Content-Type': 'application/json',
    },
});

// Request interceptor для добавления JWT токена
apiClient.interceptors.request.use(
    (config) => {
        const token = localStorage.getItem('auth_token');
        if (token) {
            config.headers.Authorization = `Bearer ${token}`;
        }
        return config;
    },
    (error) => Promise.reject(error)
);

// Response interceptor для обработки ошибок
apiClient.interceptors.response.use(
    (response) => response,
    (error) => {
        if (error.response?.status === 401) {
            // Токен истек или невалиден - редирект на login
            localStorage.removeItem('auth_token');
            window.location.href = '/login';
        }
        return Promise.reject(error);
    }
);

export default apiClient;
```

**Зависимости:**
```bash
npm install axios
```

---

#### ✅ [PENDING] 2. Реализовать Auth Service
**Файл:** `frontend/src/services/api/authService.ts`

```typescript
import apiClient from './apiClient';

export interface LoginRequest {
    email: string;
    password: string;
}

export interface RegisterRequest {
    name: string;
    email: string;
    password: string;
}

export interface AuthResponse {
    token: string;
    user: {
        id: number;
        name: string;
        email: string;
        is_active: boolean;
    };
}

export const authService = {
    async login(data: LoginRequest): Promise<AuthResponse> {
        const response = await apiClient.post<AuthResponse>('/auth/login', data);
        return response.data;
    },

    async register(data: RegisterRequest): Promise<AuthResponse> {
        const response = await apiClient.post<AuthResponse>('/auth/register', data);
        return response.data;
    },

    setToken(token: string): void {
        localStorage.setItem('auth_token', token);
    },

    getToken(): string | null {
        return localStorage.getItem('auth_token');
    },

    removeToken(): void {
        localStorage.removeItem('auth_token');
    },

    isAuthenticated(): boolean {
        return !!this.getToken();
    },
};
```

**Backend endpoints:**
- `POST /api/auth/login` → `{ token, user }`
- `POST /api/auth/register` → `{ token, user }`

---

#### ✅ [PENDING] 3. Создать AuthContext
**Файл:** `frontend/src/contexts/AuthContext.tsx`

```typescript
import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import { authService } from '../services/api/authService';

interface User {
    id: number;
    name: string;
    email: string;
    is_active: boolean;
}

interface AuthContextType {
    user: User | null;
    isAuthenticated: boolean;
    isLoading: boolean;
    login: (email: string, password: string) => Promise<void>;
    logout: () => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
    const [user, setUser] = useState<User | null>(null);
    const [isLoading, setIsLoading] = useState(true);

    useEffect(() => {
        // Проверка токена при загрузке
        const token = authService.getToken();
        if (token) {
            // TODO: Добавить проверку валидности токена через API
            // Пока просто считаем что токен валиден
            setIsLoading(false);
        } else {
            setIsLoading(false);
        }
    }, []);

    const login = async (email: string, password: string) => {
        const response = await authService.login({ email, password });
        authService.setToken(response.token);
        setUser(response.user);
    };

    const logout = () => {
        authService.removeToken();
        setUser(null);
    };

    return (
        <AuthContext.Provider
            value={{
                user,
                isAuthenticated: !!user,
                isLoading,
                login,
                logout,
            }}
        >
            {children}
        </AuthContext.Provider>
    );
};

export const useAuth = () => {
    const context = useContext(AuthContext);
    if (!context) {
        throw new Error('useAuth must be used within AuthProvider');
    }
    return context;
};
```

---

### Phase 2: Авторизация (задачи 4-5)

#### ✅ [PENDING] 4. Интегрировать Login страницу
**Файл:** `frontend/src/pages/Login/Login.tsx`

Заменить `handleSubmit`:
```typescript
const { login } = useAuth();
const [error, setError] = useState<string>('');
const [loading, setLoading] = useState(false);

const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
        await login(login, password);
        navigate('/dashboard');
    } catch (err: any) {
        setError(err.response?.data?.error || 'Ошибка авторизации');
    } finally {
        setLoading(false);
    }
};
```

---

#### ✅ [PENDING] 5. Добавить ProtectedRoute
**Файл:** `frontend/src/components/ProtectedRoute.tsx`

```typescript
import React from 'react';
import { Navigate } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';

interface ProtectedRouteProps {
    children: React.ReactElement;
}

const ProtectedRoute: React.FC<ProtectedRouteProps> = ({ children }) => {
    const { isAuthenticated, isLoading } = useAuth();

    if (isLoading) {
        return <div>Загрузка...</div>; // TODO: Spinner
    }

    if (!isAuthenticated) {
        return <Navigate to="/login" replace />;
    }

    return children;
};

export default ProtectedRoute;
```

**Обновить App.tsx:**
```typescript
<Route path="/dashboard" element={<ProtectedRoute><Dashboard /></ProtectedRoute>}>
    <Route path="statistics" element={<Statistics />} />
    <Route index element={<Placeholder />} />
</Route>
```

---

### Phase 3: API сервисы (задачи 6-7)

#### ✅ [PENDING] 6. Projects Service
**Файл:** `frontend/src/services/api/projectsService.ts`

```typescript
import apiClient from './apiClient';

export interface Project {
    id: number;
    name: string;
    slug: string;
    timezone: string;
    currency: string;
    is_active: boolean;
    created_at: string;
    updated_at: string;
}

export const projectsService = {
    async getAll(): Promise<Project[]> {
        const response = await apiClient.get<Project[]>('/projects');
        return response.data;
    },

    async getById(id: number): Promise<Project> {
        const response = await apiClient.get<Project>(`/projects/${id}`);
        return response.data;
    },
};
```

**Backend endpoint:**
- `GET /api/projects` → `[{ id, name, slug, ... }]`

---

#### ✅ [PENDING] 7. Reports Service
**Файл:** `frontend/src/services/api/reportsService.ts`

```typescript
import apiClient from './apiClient';

export interface MetricsSummary {
    month: string;
    visits: number;
    users: number;
    bounce: number;
    avgSec: number;
    conv: number;
    dynamics?: {
        visits: number;
        users: number;
        bounce: number;
        avgSec: number;
        conv: number;
    };
}

export interface AgeMetrics {
    month: string;
    age: string;
    visits: number;
    users: number;
    bounce: number;
    avgSec: number;
}

export interface DirectCampaign {
    campaignId: number;
    name: string;
    rows: Array<{
        month: string;
        impressions: number;
        clicks: number;
        ctr: number;
        cpc: number;
        conv?: number;
        cpa?: number;
        cost: number;
    }>;
}

export interface Report {
    projectId: number;
    periods: string[];
    metrica: {
        summary: MetricsSummary[];
        age: AgeMetrics[];
    };
    direct: {
        totals: Array<{
            month: string;
            impressions: number;
            clicks: number;
            ctr: number;
            cpc: number;
            conv?: number;
            cpa?: number;
            cost: number;
        }>;
        campaigns: DirectCampaign[];
    };
    seo: {
        summary: Array<{
            month: string;
            visitors: number;
            conv: number;
        }>;
        queries: Array<{
            month: string;
            query: string;
            position: number;
            url?: string;
        }>;
    };
    ai_insights?: {
        summary: string;
        recommendations: string[];
    };
}

export const reportsService = {
    async getReport(projectId: number): Promise<Report> {
        const response = await apiClient.get<Report>(`/report/${projectId}`);
        return response.data;
    },

    async getChannelMetrics(projectId: number, periods: string[]): Promise<any> {
        const periodsParam = periods.join(',');
        const response = await apiClient.get(`/channel-metrics/${projectId}?periods=${periodsParam}`);
        return response.data;
    },
};
```

**Backend endpoints:**
- `GET /api/report/:id` → полный отчет
- `GET /api/channel-metrics/:id?periods=...` → метрики по каналам

---

### Phase 4: UI компоненты (задачи 8-9)

#### ✅ [PENDING] 8. ProjectsList компонент
**Файл:** `frontend/src/components/ProjectsList/ProjectsList.tsx`

```typescript
import React, { useEffect, useState } from 'react';
import { projectsService, Project } from '../../services/api/projectsService';

const ProjectsList: React.FC = () => {
    const [projects, setProjects] = useState<Project[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string>('');

    useEffect(() => {
        const fetchProjects = async () => {
            try {
                const data = await projectsService.getAll();
                setProjects(data);
            } catch (err: any) {
                setError(err.message || 'Ошибка загрузки проектов');
            } finally {
                setLoading(false);
            }
        };

        fetchProjects();
    }, []);

    if (loading) return <div>Загрузка проектов...</div>;
    if (error) return <div>Ошибка: {error}</div>;

    return (
        <div className="projects-list">
            <h2>Мои проекты</h2>
            <div className="projects-grid">
                {projects.map(project => (
                    <div key={project.id} className="project-card">
                        <h3>{project.name}</h3>
                        <p>{project.slug}</p>
                        <button>Открыть отчет</button>
                    </div>
                ))}
            </div>
        </div>
    );
};

export default ProjectsList;
```

---

#### ✅ [PENDING] 9. Интеграция Statistics с API
**Файл:** `frontend/src/pages/Statistics/Statistics.tsx`

Заменить захардкоженный массив `metrics`:
```typescript
const [report, setReport] = useState<Report | null>(null);
const [loading, setLoading] = useState(true);
const projectId = 1; // TODO: Получать из роута или контекста

useEffect(() => {
    const fetchReport = async () => {
        try {
            const data = await reportsService.getReport(projectId);
            setReport(data);
        } catch (err) {
            console.error('Error fetching report:', err);
        } finally {
            setLoading(false);
        }
    };

    fetchReport();
}, [projectId]);

// Преобразовать данные из API в формат для таблицы
const metrics = report ? [
    {
        name: 'Посетители, кол-во',
        october: report.metrica.summary[0]?.users || 0,
        september: report.metrica.summary[1]?.users || 0,
        august: report.metrica.summary[2]?.users || 0,
        efficiency: report.metrica.summary[0]?.dynamics?.users || 0,
    },
    // ... остальные метрики
] : [];
```

---

### Phase 5: Обработка ошибок (задачи 10-12)

#### ✅ [PENDING] 10. Loading & Error states
**Файл:** `frontend/src/components/LoadingSpinner/LoadingSpinner.tsx`

```typescript
const LoadingSpinner: React.FC = () => (
    <div className="loading-spinner">
        <div className="spinner"></div>
        <p>Загрузка...</p>
    </div>
);
```

**Файл:** `frontend/src/components/ErrorMessage/ErrorMessage.tsx`

```typescript
interface ErrorMessageProps {
    message: string;
    onRetry?: () => void;
}

const ErrorMessage: React.FC<ErrorMessageProps> = ({ message, onRetry }) => (
    <div className="error-message">
        <p>❌ {message}</p>
        {onRetry && <button onClick={onRetry}>Повторить</button>}
    </div>
);
```

---

#### ✅ [PENDING] 11. JWT Interceptor (уже в apiClient)
Проверить что работает автоматическое добавление токена в headers.

---

#### ✅ [PENDING] 12. Реализовать Logout
**Файл:** `frontend/src/components/Dashboard/Dashboard.tsx`

Обновить `handleLogout`:
```typescript
const { logout } = useAuth();

const handleLogout = useCallback(() => {
    logout();
    navigate('/login');
}, [logout, navigate]);
```

---

### Phase 6: Тестирование (задача 13)

#### ✅ [PENDING] 13. E2E тестирование

**Сценарий:**
1. Открыть http://localhost:3000
2. Должен быть редирект на /login (если не авторизован)
3. Ввести email: `admin@example.com`, password: `password123`
4. Клик "Войти" → редирект на /dashboard
5. Проверить что отображаются реальные проекты из API
6. Кликнуть на проект → открыть страницу статистики
7. Проверить что данные загружены из API (не захардкожены)
8. Клик "Выйти" → редирект на /login, токен удален

**Данные для тестирования:**
- Admin user: `admin@example.com` / `password123`
- Backend API: http://localhost:8080/api
- Test project ID: 1

---

## 🔧 Установка зависимостей

```bash
cd frontend
npm install axios
```

---

## 📝 Примечания

### Environment Variables
В `frontend/.env`:
```env
REACT_APP_API_URL=http://localhost:8080/api
REACT_APP_NAME=Planica BI
```

### CORS
Backend уже настроен на CORS (middleware.CORS() в Echo).

### Token Storage
- Используем `localStorage` для хранения JWT токена
- Ключ: `auth_token`
- Токен добавляется автоматически через axios interceptor

### Error Handling
- 401 Unauthorized → автоматический редирект на /login
- Другие ошибки → отображение в UI

---

## 🎯 Текущий статус: Phase 1-2 завершены! 🎉

**Прогресс:** 10/13 задач ✅

### ✅ Завершено:
1. ✅ API клиент (axios) с interceptors
2. ✅ Auth Service (login, register, tokens)
3. ✅ AuthContext (глобальное состояние)
4. ✅ Login интеграция с backend
5. ✅ ProtectedRoute компонент
6. ✅ Projects Service
7. ✅ Reports Service (полная типизация)
8. ✅ Loading/Error компоненты
9. ✅ Logout функционал
10. ✅ JWT автоматически добавляется в headers

### 🔄 Осталось:
- ProjectsList компонент (опционально)
- Интеграция Statistics с API (опционально)
- **E2E тестирование** (критично)

**Следующий шаг:** Протестировать login → dashboard flow

