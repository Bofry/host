package test

type MockRedis struct {
	Host     string
	Password string
	DB       int
	PoolSize int
}
