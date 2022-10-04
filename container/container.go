package container

import (
	"echoapp/config"
	"echoapp/logger"
	"echoapp/repo"
	"echoapp/repository"
	"echoapp/session"
)

// Container represents a interface for accessing the data which sharing in overall application.
type Container interface {
	GetRepository() repository.Repository
	GetRepo() repo.Repo
	GetConfig() *config.Config
	GetLogger() logger.Logger
	GetEnv() string
}

// container struct is for sharing data which such as database setting, the setting of application and logger in overall this application.
type container struct {
	rep     repository.Repository
	repo    repo.Repo
	session session.Session
	config  *config.Config
	logger  logger.Logger
	env     string
}

// NewContainer is constructor.
func NewContainer(rep repository.Repository, repo repo.Repo, config *config.Config, logger logger.Logger, env string) Container {
	return &container{rep: rep, repo: repo, config: config, logger: logger, env: env}
}

// GetRepository returns the object of repository.
func (c *container) GetRepository() repository.Repository {
	return c.rep
}

// GetReporeturns the object of repo.
func (c *container) GetRepo() repo.Repo {
	return c.repo
}

// GetSession returns the object of session.
func (c *container) GetSession() session.Session {
	return c.session
}

// GetConfig returns the object of configuration.
func (c *container) GetConfig() *config.Config {
	return c.config
}

// GetLogger returns the object of logger.
func (c *container) GetLogger() logger.Logger {
	return c.logger
}

// GetEnv returns the running environment.
func (c *container) GetEnv() string {
	return c.env
}
