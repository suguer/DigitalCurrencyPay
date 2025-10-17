package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type StorageConfig struct {
	Path string `yaml:"path"`
	Log  string `yaml:"log"`
}

type EthConfig struct {
	Name        string
	Chain       string
	Enable      int
	Key         string
	Node        string
	Precision   int
	Logger      bool
	GrpcAddress string `yaml:"grpc_address"`
	ChainId     int    `yaml:"chain_id"`
	// Contract []ContractConfig
}
type ContractConfig struct {
	Token   string
	Address string
	Decimal int
}

func (c *EthConfig) ContractMap() map[string]ContractConfig {
	m := make(map[string]ContractConfig)
	// for _, v := range c.Contract {
	// 	m[v.Address] = v
	// }
	return m
}

type RedisConf struct {
	Host     string
	Port     int
	Database int
	Password string
}

// Config 配置结构体
type Config struct {
	Salt   string `yaml:"salt"`
	Server struct {
		Port int    `yaml:"port"`
		Mode string `yaml:"mode"`
	}
	BlockChain struct {
		Tron       EthConfig `yaml:"tron"`
		TronShasta EthConfig `yaml:"tron-shasta"`
		Arbitrum   EthConfig `yaml:"arbitrum"`
		Avax       EthConfig `yaml:"avax"`
		Matic      EthConfig `yaml:"matic"`
		Base       EthConfig `yaml:"base"`
		Op         EthConfig `yaml:"op"`
		Ethereum   EthConfig `yaml:"ethereum"`
	} `yaml:"blockchain"`

	Redis    RedisConf
	Database DatabaseConf `yaml:"database"`
	Queue    struct {
		Driver string `yaml:"driver"`
	} `yaml:"queue"`
	Storage StorageConfig `yaml:"storage"`
}

var (
	Conf *Config
)

// Load 加载配置文件
func Load(file string) (*Config, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	// fmt.Printf("cfg: %+v\n", cfg)
	Conf = &cfg
	return &cfg, nil
}
