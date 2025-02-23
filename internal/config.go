package internal

type Config struct {
	config map[string]interface{}
}

func NewConfig() *Config {
	return &Config{
		config: make(map[string]interface{}),
	}
}

func (c *Config) Set(key string, value interface{}) {
	c.config[key] = value
}

func (c *Config) Get(key string) interface{} {
	return c.config[key]
}

func (c *Config) Has(key string) bool {
	_, ok := c.config[key]
	return ok
}

func (c *Config) Delete(key string) {
	delete(c.config, key)
}

func (c *Config) Clear() {
	c.config = make(map[string]interface{})
}
