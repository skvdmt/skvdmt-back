package internal

type Option func(data *OptionData) error

type OptionData struct {
	parseFlags bool
}

// WithoutParseFlags without processing app launch parameters
func WithoutParseFlags() Option {
	return func(data *OptionData) error {
		data.parseFlags = false
		return nil
	}
}
