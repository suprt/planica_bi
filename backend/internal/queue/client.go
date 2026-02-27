package queue

import (
	"time"

	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
	"github.com/suprt/planica_bi/backend/internal/config"
)

// Client wraps asynq client for task enqueueing
type Client struct {
	client *asynq.Client
}

// NewClient creates a new queue client
func NewClient(cfg *config.Config) (*Client, error) {
	redisOpt := asynq.RedisClientOpt{
		Addr:     cfg.RedisHost + ":" + cfg.RedisPort,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	}

	client := asynq.NewClient(redisOpt)
	return &Client{client: client}, nil
}

// Close closes the queue client
func (c *Client) Close() error {
	return c.client.Close()
}

// EnqueueSyncMetricaTask enqueues a task to sync Metrica data
func (c *Client) EnqueueSyncMetricaTask(projectID uint, year, month int) (*asynq.TaskInfo, error) {
	task := NewSyncMetricaTask(projectID, year, month)
	return c.client.Enqueue(task,
		asynq.MaxRetry(3),
		asynq.Timeout(10*60*time.Second), // 10 minutes timeout
		asynq.Queue("default"),
	)
}

// EnqueueSyncDirectTask enqueues a task to sync Direct data
func (c *Client) EnqueueSyncDirectTask(projectID uint, year, month int) (*asynq.TaskInfo, error) {
	task := NewSyncDirectTask(projectID, year, month)
	return c.client.Enqueue(task,
		asynq.MaxRetry(3),
		asynq.Timeout(10*60*time.Second), // 10 minutes timeout
		asynq.Queue("default"),
	)
}

// EnqueueSyncProjectTask enqueues a task to sync entire project
func (c *Client) EnqueueSyncProjectTask(projectID uint) (*asynq.TaskInfo, error) {
	task := NewSyncProjectTask(projectID)
	return c.client.Enqueue(task,
		asynq.MaxRetry(3),
		asynq.Timeout(15*60*time.Second), // 15 minutes timeout
		asynq.Queue("default"),
	)
}

// EnqueueAnalyzeMetricsTask enqueues a task to analyze channel metrics using AI
func (c *Client) EnqueueAnalyzeMetricsTask(projectID uint, periods []string) (*asynq.TaskInfo, error) {
	task := NewAnalyzeMetricsTask(projectID, periods)
	return c.client.Enqueue(task,
		asynq.MaxRetry(2),               // Fewer retries for AI analysis
		asynq.Timeout(2*60*time.Second), // 2 minutes timeout
		asynq.Queue("low"),              // Low priority queue for AI analysis
	)
}

// EnqueueGenerateReportTask enqueues a task to generate report
func (c *Client) EnqueueGenerateReportTask(projectID uint) (*asynq.TaskInfo, error) {
	task := NewGenerateReportTask(projectID)
	return c.client.Enqueue(task,
		asynq.MaxRetry(2),
		asynq.Timeout(5*60*time.Second), // 5 minutes timeout for report generation
		asynq.Queue("default"),
	)
}

// GetRedisClient returns underlying Redis client (for worker)
func GetRedisClient(cfg *config.Config) redis.UniversalClient {
	return redis.NewClient(&redis.Options{
		Addr:     cfg.RedisHost + ":" + cfg.RedisPort,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})
}
