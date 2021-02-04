package cmd

import (
	"fmt"
	"markit/utils"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

func processCmd(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&cfgFile, "config", "", "默认配置文件(将从指定路径及上层路径递归查找.markit.toml,如果未能找到将尝试查询$HOME/.markit.toml)")
}

func loadConfig(path string) {
	viper.SetConfigName(".markit")
	viper.SetConfigType("toml")
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		path := utils.ResolveConfigPath(path)
		if path == "" {
			var err error
			path, err = homedir.Dir()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
		viper.AddConfigPath(path)
	}

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
