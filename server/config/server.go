package config

var Values ServerConfig

// Model that links to config.yml file
type ServerConfig struct {
	Server struct {
		ApiPath            string   `yaml:"api-path"  env:"API_PATH" env-description:"API base path"`
		ApiVersion         string   `yaml:"api-version"  env:"API_VERSION" env-description:"API Version"`
		CorsAllowedClients []string `yaml:"cors-allowed-clients" env:"CORS_ALLOWED_CLIENTS"  env-description:"List of allowed CORS Clients"`
		Environment        string   `yaml:"environment" env:"SERVER_ENVIRONMENT"  env-description:"server environment"`

		Host              string `yaml:"host"  env:"SERVER_HOST" env-description:"server host"`
		Port              string `yaml:"port" env:"SERVER_PORT"  env-description:"server port"`
		Protocol          string `yaml:"protocol" env:"SERVER_PROTOCOL"  env-description:"server protocol"`
		MaxBulkUploadSize int    `yaml:"max-bulk-upload-size"  env:"MAX_BULK_UPLOAD_SIZE" env-description:"max total bulk upload size"`
		MaxUploadFileSize int    `yaml:"max-upload-file-size"  env:"MAX_UPLOAD_FILE_SIZE" env-description:"max upload size of a single file"`
		UploadFolder      string `yaml:"upload-folder"  env:"UPLOAD_FOLDER" env-description:"folder where to store bulk upload files"`
	} `yaml:"server"`
}
