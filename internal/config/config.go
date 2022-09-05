package config

type ConfigCheckStatus struct {
	SupplierAddress string `required:"true" split_words:"true"` // CanalBox endpoint for SIRI-ET subscription
	SubscriberRef   string `required:"true" split_words:"true"`
}

type ConfigSubscribe struct {
	SupplierAddress string `required:"true" split_words:"true"` // CanalBox endpoint for SIRI-ET subscription
	SubscriberRef   string `required:"true" split_words:"true"`
	ProducerRef     string `required:"true" split_words:"true"`
	ConsumerAddress string `required:"true" split_words:"true"`
}
