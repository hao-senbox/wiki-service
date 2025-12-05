package app

import (
	"wiki-service/internal/domain/repository"
	"wiki-service/internal/domain/usecase"
	"wiki-service/internal/infrastructure/database"
	infrastructureRepository "wiki-service/internal/infrastructure/repository"
	httpInterface "wiki-service/internal/interface/http"
	"wiki-service/internal/interface/http/handler"
	"wiki-service/internal/interface/middleware"
	"wiki-service/pkg/config"
	"wiki-service/pkg/consul"
	"wiki-service/pkg/gateway"
	"wiki-service/pkg/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/hashicorp/consul/api"
	"github.com/hung-senbox/senbox-cache-service/pkg/cache"
	"github.com/hung-senbox/senbox-cache-service/pkg/cache/cached"
	"github.com/hung-senbox/senbox-cache-service/pkg/redis"
	"go.mongodb.org/mongo-driver/mongo"
)

// Container holds all application dependencies
type Container struct {
	Config            *config.Config
	Logger            *logger.Logger
	MongoDB           *mongo.Database
	AuditMiddleware   *middleware.AuditMiddleware
	WikiRepository    repository.WikiRepository
	WikiUseCase       usecase.WikiUseCase
	WikiHandler       *handler.WikiHandler
	App               *fiber.App
	UserGateway       gateway.UserGateway
	FileGateway       gateway.FileGateway
	MediaGateway      gateway.MediaGateway
	Consul            *api.Client
	ConsulConn        consul.Client
	CacheClientRedis  *cache.RedisCache
	CachedMainGateway cached.CachedMainGateway
}

// NewContainer initializes all application dependencies
func NewContainer() (*Container, error) {
	c := &Container{}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}
	c.Config = cfg

	// Initialize logger
	appLogger, err := logger.NewLogger("logs/wiki-service")
	if err != nil {
		return nil, err
	}
	c.Logger = appLogger
	c.Logger.Info("Logger initialized successfully")

	// Initialize database
	if err := c.initDatabase(); err != nil {
		return nil, err
	}

	// Initialize Redis cache
	if err := c.initCache(); err != nil {
		return nil, err
	}

	// Initialize Consul
	if err := c.initConsul(); err != nil {
		return nil, err
	}

	// Initialize repositories
	c.initRepositories()

	// Initialize gateway
	c.initGateway()

	// Initialize media gateway
	c.initMediaGateway()

	// Initialize file gateway
	c.initFileGateway()

	// Initialize use cases
	c.initUseCases()

	// Initialize handlers
	c.initHandlers()

	// Initialize middlewares
	c.initMiddlewares()

	// Setup router
	c.setupRouter()

	return c, nil
}

// initDatabase initializes MongoDB connection
func (c *Container) initDatabase() error {
	c.Logger.Info("Connecting to MongoDB database")
	mongoDB, err := database.NewMongoConnection(database.MongoConfig{
		Host:     c.Config.MongoDB.Host,
		Port:     c.Config.MongoDB.Port,
		User:     c.Config.MongoDB.User,
		Password: c.Config.MongoDB.Password,
		DBName:   c.Config.MongoDB.DBName,
	})
	if err != nil {
		return err
	}
	c.MongoDB = mongoDB
	c.Logger.Info("MongoDB connection established")
	return nil
}

// initRepositories initializes all repositories
func (c *Container) initRepositories() {
	c.WikiRepository = infrastructureRepository.NewWikiRepositoryMongo(c.MongoDB)
}

// initUseCases initializes all use cases
func (c *Container) initUseCases() {
	c.WikiUseCase = usecase.NewWikiUseCase(c.WikiRepository, c.FileGateway, c.UserGateway, c.MediaGateway)
}

// initHandlers initializes all HTTP handlers
func (c *Container) initHandlers() {
	c.WikiHandler = handler.NewWikiHandler(c.WikiUseCase)
}

// initMiddlewares initializes all middlewares
func (c *Container) initMiddlewares() {
	c.AuditMiddleware = middleware.NewAuditMiddleware(c.Logger)
}

// setupRouter sets up the Fiber application with routes
func (c *Container) setupRouter() {
	c.App = httpInterface.SetupRouter(
		c.WikiHandler,
		c.AuditMiddleware,
		c.UserGateway,
	)
}

func (c *Container) initCache() error {
	// redis cache
	cacheClientRedis, err := redis.InitRedisCache(c.Config.Database.RedisCache.Host, c.Config.Database.RedisCache.Port, c.Config.Database.RedisCache.Password, c.Config.Database.RedisCache.DB)
	if err != nil {
		return err
	}
	c.CacheClientRedis = cacheClientRedis
	c.Logger.Info("Redis cache initialized successfully")
	c.CachedMainGateway = cached.NewCachedMainGateway(cacheClientRedis)
	return nil
}

func (c *Container) initGateway() {
	c.UserGateway = gateway.NewUserGateway("go-main-service", c.Consul, c.CachedMainGateway, c.Logger)
	c.Logger.Info("Gateway initialized successfully")
}

func (c *Container) initFileGateway() {
	c.FileGateway = gateway.NewFileGateway("go-main-service", c.Consul, c.Logger)
	c.Logger.Info("File gateway initialized successfully")
}

func (c *Container) initMediaGateway() {
	c.MediaGateway = gateway.NewMediaGateway("media-service", c.Consul, c.Logger)
	c.Logger.Info("Media gateway initialized successfully")
}

func (c *Container) initConsul() error {
	consulConn := consul.NewConsulConn(c.Logger, c.Config)
	c.Consul = consulConn.Connect()
	c.ConsulConn = consulConn
	c.Logger.Info("Consul initialized successfully")
	return nil
}
