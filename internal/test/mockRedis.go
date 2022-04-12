package test

type mockRedis struct {
	Host     string
	Password string
	DB       int
	PoolSize int
}
