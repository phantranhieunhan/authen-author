package redis

// Config redis
type config struct {
	Address  string
	Password string
	Database int
	Timeout  int
}

type Option interface {
	apply(*config)
}

type optionFunc func(*config)

func (o optionFunc) apply(c *config) {
	o(c)
}

func newConfig(opts ...Option) (*config, error) {
	c := &config{}
	for _, opt := range opts {
		opt.apply(c)
	}
	if c.Address == "" {
		c.Address = "127.0.0.1:6379" // set default address
	}
	if c.Timeout == 0 {
		c.Timeout = 10
	}
	if c.Database == 0 {
		c.Database = 0
	}

	return c, nil
}

func WithTimeout(timeout int) Option {
	return optionFunc(func(c *config) {
		c.Timeout = timeout
	})
}

func WithAddress(address string) Option {
	return optionFunc(func(c *config) {
		c.Address = address
	})
}

func WithPassword(password string) Option {
	return optionFunc(func(c *config) {
		c.Password = password
	})
}

func WithDatabase(database int) Option {
	return optionFunc(func(c *config) {
		c.Database = database
	})
}
