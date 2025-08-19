package grpcservice

var ENV struct {
	Server struct {
		Port string `env:"PORT" envDefault:"8080"`
	} `envPrefix:"SERVER_"`
	DB struct {
		Host     string `env:"HOST" envDefault:"localhost"`
		User     string `env:"USER" envDefault:"root"`
		Password string `env:"PASSWORD" envDefault:"root"`
		Name     string `env:"NAME" envDefault:"grpc-golang"`
		Port     string `env:"PORT" envDefault:"3306"`
	} `envPrefix:"DB_"`
}
