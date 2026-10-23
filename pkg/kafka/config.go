package kafka

// Config kafka config
type Config struct {
	Brokers    []string `mapstructure:"brokers"`
	GroupID    string   `mapstructure:"groupID"`
	InitTopics bool     `mapstructure:"initTopics"`
}

// TopicConfig kafka topic config
type TopicConfig struct {
	TopicName         string `mapstructure:"topicName"`
	Partitions        int    `mapstructure:"partitions"`
	ReplicationFactor int    `mapstructure:"replicationFactor"`
	// MaxMessageBytes caps a single produced payload so callers can split oversized
	// batches instead of hitting the broker message.max.bytes (~1MB default). 0 = unset.
	MaxMessageBytes int `mapstructure:"maxMessageBytes"`
}
