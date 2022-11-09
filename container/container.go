package container

import (
	"echoapp/config"
	"echoapp/logger"
	"echoapp/repo"

	"echoapp/session"

	"github.com/allegro/bigcache/v3"
)

// Container represents a interface for accessing the data which sharing in overall application.
type Container interface {
	GetRepository() repo.Repository
	GetConfig() config.Config
	GetLogger() logger.Logger
	GetEnv() string
	GetBigCache() *bigcache.BigCache
}

// container struct is for sharing data which such as database setting, the setting of application and logger in overall this application.
type container struct {
	repository *repo.Repository
	session    session.Session
	config     *config.Config
	bigCache   *bigcache.BigCache
	logger     logger.Logger
	env        string
}

// NewContainer is constructor.
func NewContainer(repository *repo.Repository, config *config.Config, bigCache *bigcache.BigCache, logger logger.Logger, env string) Container {
	return &container{repository: repository, config: config, logger: logger, bigCache: bigCache, env: env}
}

// GetRepository returns the object of repo.
func (c *container) GetRepository() repo.Repository {
	return *c.repository
}

// GetSession returns the object of session.
func (c *container) GetSession() session.Session {
	return c.session
}

// GetConfig returns the object of configuration.
func (c *container) GetConfig() config.Config {
	return *c.config
}

// GetLogger returns the object of logger.
func (c *container) GetLogger() logger.Logger {
	return c.logger
}

// GetEnv returns the running environment.
func (c *container) GetEnv() string {
	return c.env
}

// GetBigCache returns the object of caching.
func (c *container) GetBigCache() *bigcache.BigCache {
	return c.bigCache
}
