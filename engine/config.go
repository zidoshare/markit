package engine

import "github.com/spf13/viper"

const (
	AUTO_SPACE_NAME = "autoSpace"
	FIX_TERM_TYPE   = "fixTermType"
)

type Options struct {
	AutoSpace   bool
	FixTermTypo bool
}

func NewOptions() Options {
	return Options{
		AutoSpace:   viper.GetBool(AUTO_SPACE_NAME),
		FixTermTypo: viper.GetBool(FIX_TERM_TYPE),
	}
}

func init() {
	viper.SetDefault(AUTO_SPACE_NAME, true)
	viper.SetDefault(FIX_TERM_TYPE, true)
}
