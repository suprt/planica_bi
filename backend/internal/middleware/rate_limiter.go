package middleware

import (
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/time/rate"
)

// RateLimiterConfig конфигурация rate limiter
type RateLimiterConfig struct {
	RequestsPerSecond float64       // Количество запросов в секунду
	BurstSize         int           // Максимальный размер "пачки" запросов
	KeyExpiry         time.Duration // Время жизни ключа в памяти
}

// DefaultRateLimiterConfig конфигурация по умолчанию
func DefaultRateLimiterConfig() RateLimiterConfig {
	return RateLimiterConfig{
		RequestsPerSecond: 10,  // 10 запросов в секунду
		BurstSize:         20,  // Пачка до 20 запросов
		KeyExpiry:         10 * time.Minute,
	}
}

// RateLimiter middleware для ограничения запросов
type RateLimiter struct {
	config  RateLimiterConfig
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
	stopCh   chan struct{}
}

// NewRateLimiter создаёт новый rate limiter middleware
func NewRateLimiter(config RateLimiterConfig) *RateLimiter {
	rl := &RateLimiter{
		config:   config,
		limiters: make(map[string]*rate.Limiter),
		stopCh:   make(chan struct{}),
	}

	// Запускаем очистку старых ключей
	go rl.cleanup()

	return rl
}

// Middleware создаёт Echo middleware
func (rl *RateLimiter) Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Получаем ключ (IP адрес)
			key := c.RealIP()

			// Получаем или создаём limiter для этого ключа
			limiter := rl.getLimiter(key)

			// Проверяем, можем ли обработать запрос
			if !limiter.Allow() {
				return echo.NewHTTPError(429, "Too many requests, please try again later")
			}

			return next(c)
		}
	}
}

// getLimiter получает или создаёт limiter для ключа
func (rl *RateLimiter) getLimiter(key string) *rate.Limiter {
	rl.mu.RLock()
	if limiter, exists := rl.limiters[key]; exists {
		rl.mu.RUnlock()
		return limiter
	}
	rl.mu.RUnlock()

	rl.mu.Lock()
	defer rl.mu.Unlock()

	// Проверяем ещё раз после получения write lock
	if limiter, exists := rl.limiters[key]; exists {
		return limiter
	}

	// Создаём новый limiter
	limiter := rate.NewLimiter(rate.Limit(rl.config.RequestsPerSecond), rl.config.BurstSize)
	rl.limiters[key] = limiter

	return limiter
}

// cleanup периодически удаляет старые ключи
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			rl.mu.Lock()
			for key, limiter := range rl.limiters {
				if limiter.AllowN(time.Now(), rl.config.BurstSize) {
					delete(rl.limiters, key)
				}
			}
			rl.mu.Unlock()
		case <-rl.stopCh:
			return
		}
	}
}

// Stop останавливает очистку
func (rl *RateLimiter) Stop() {
	close(rl.stopCh)
}

// StrictRateLimiterConfig для строгих endpoints (например, auth)
func StrictRateLimiterConfig() RateLimiterConfig {
	return RateLimiterConfig{
		RequestsPerSecond: 2,   // 2 запроса в секунду
		BurstSize:         5,   // Пачка до 5 запросов
		KeyExpiry:         5 * time.Minute,
	}
}

// PermissiveRateLimiterConfig для лёгких endpoints
func PermissiveRateLimiterConfig() RateLimiterConfig {
	return RateLimiterConfig{
		RequestsPerSecond: 50,  // 50 запросов в секунду
		BurstSize:         100, // Пачка до 100 запросов
		KeyExpiry:         10 * time.Minute,
	}
}
