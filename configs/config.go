package configs

type conf struct {
	DBDriver           string   `mapstructure:"DB_DRIVER"`
	DBHost             string   `mapstructure:"DB_HOST"`
	DBPort             string   `mapstructure:"DB_PORT"`
	DBUser             string   `mapstructure:"DB_USER"`
	DBPassword         string   `mapstructure:"DB_PASSWORD"`
	DBName             string   `mapstructure:"DB_NAME"`
	WebServerPort      string   `mapstructure:"WEB_SERVER_PORT"`
	GRPCServerPort     string   `mapstructure:"GRPC_SERVER_PORT"`
	InitialChatMessage string   `mapstructure:"INITIAL_CHAT_MESSAGE"`
	OpenAIApiKey       string   `mapstructure:"OPENAI_API_KEY"`
	Model              string   `mapstructure:"MODEL"`
	ModelMaxTokens     int      `mapstructure:"MODEL_MAX_TOKENS"`
	Temperature        float64  `mapstructure:"TEMPERATURE"`
	TopP               float64  `mapstructure:"TOP_P"`
	N                  int      `mapstructure:"N"`
	Stop               []string `mapstructure:"STOP"`
	MaxTokens          int      `mapstructure:"MAX_TOKENS"`
	AuthToken          string   `mapstructure:"AUTH_TOKEN"`
}

func LoadConfig(path string) (*conf, error) {
	return &conf{
		DBDriver:           "mysql",
		DBHost:             "mysql",
		DBPort:             "3306",
		DBUser:             "root",
		DBPassword:         "root",
		DBName:             "jarvis",
		WebServerPort:      "8080",
		GRPCServerPort:     "50051",
		InitialChatMessage: "Seu nome é Jarvis. Você é a inteligência artificial de Lucas Gois. Você da suporte a programadores e arquitetos de software",
		OpenAIApiKey:       "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",
		Model:              "gpt-3.5-turbo",
		ModelMaxTokens:     4096,
		Temperature:        0.2,
		TopP:               0.2,
		N:                  1,
		Stop:               []string{"super-end"},
		MaxTokens:          300,
		AuthToken:          "123456",
	}, nil
}
