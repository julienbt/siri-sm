package config

type Config struct {
	SiriSm struct {
		//CheckStatusSchedule string   `required:"true" split_words:"true"`
		SubscriberRef string `default:"NAVITIA" split_words:"true"`
		// WsUrlPublic         string   `required:"true" split_words:"true"` // Kisio endpoint for SIRI-ET notifications
		SupplierAddress string `required:"true" split_words:"true"` // CanalBox endpoint for SIRI-ET subscription
		// ChangeBeforeUpdates string   `default:"PT1M" split_words:"true"`
		// SubscribedLines     []string `required:"true" split_words:"true"`
	}
	// Redis struct {
	// 	Host     string `required:"true" split_words:"true"`
	// 	Port     string `required:"true" split_words:"true"`
	// 	Password string `required:"true" split_words:"true"`
	// }
}
