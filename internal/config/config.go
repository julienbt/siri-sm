package config

type ConfigCheckStatus struct {
	SubscriberRef   string `required:"true" split_words:"true"`
	SupplierAddress string `required:"true" split_words:"true"` // CanalBox endpoint for SIRI-ET subscription
}

type ConfigSubscribe struct {
	SubscriberRef   string `required:"true" split_words:"true"`
	ConsumerAddress string `required:"true" split_words:"true"`
}
