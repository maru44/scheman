package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/maru44/scheman/db"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/volatiletech/sqlboiler/v4/boilingcore"
	"github.com/volatiletech/sqlboiler/v4/drivers"
	"github.com/volatiletech/sqlboiler/v4/importers"
)

var (
	flagConfigFile string
	// state          *boilingcore.State
	state db.SchemanState
)

func main() {
	var rootCmd = &cobra.Command{
		Use:           "scheman [flags] <driver>",
		Short:         "Scheman generates schema table to notion, etc", // @TODO
		Long:          "Scheman generates a schema table to notion",    // @TODO
		Example:       `scheman psql`,
		PreRunE:       setState,
		RunE:          run,
		PostRunE:      postRun,
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&flagConfigFile, "config", "c", "", "Filename of config file to override default lookup")
	rootCmd.PersistentFlags().BoolP("add-enum-types", "", false, "Enable generation of types for enums")
	rootCmd.PersistentFlags().StringP("enum-null-prefix", "", "Null", "Name prefix of nullable enum types")
	rootCmd.PersistentFlags().StringP("output-type", "", "Notion", "Platform table definition")
	rootCmd.PersistentFlags().StringP("output-page-id", "", "Null", "page id")
	rootCmd.PersistentFlags().StringP("token", "", "Null", "Token or id for authentication") // recommend not to use
	rootCmd.PersistentFlags().StringP("password", "", "Null", "Password for authentication") // recommend not to use

	viper.BindPFlags(rootCmd.PersistentFlags())
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func setState(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return errors.New("must provide a driver name")
	}

	driverName, _, err := drivers.RegisterBinaryFromCmdArg(args[0])
	if err != nil {
		return err
	}

	config := &boilingcore.Config{
		DriverName:     driverName,
		AddEnumTypes:   viper.GetBool("add-enum-types"),
		EnumNullPrefix: viper.GetString("enum-null-prefix"),
	}
	config.DriverConfig = map[string]interface{}{
		"whitelist":        viper.GetStringSlice(driverName + ".whitelist"),
		"blacklist":        viper.GetStringSlice(driverName + ".blacklist"),
		"add-enum-types":   config.AddEnumTypes,
		"enum-null-prefix": config.EnumNullPrefix,
	}

	keys := allKeys(driverName)
	for _, key := range keys {
		if key != "blacklist" && key != "whitelist" {
			prefixedKey := fmt.Sprintf("%s.%s", driverName, key)
			config.DriverConfig[key] = viper.Get(prefixedKey)
		}
	}

	config.Imports = importers.NewDefaultImports()

	state.State, err = boilingcore.New(config)
	return err
}

func allKeys(prefix string) []string {
	keys := make(map[string]bool)

	prefix += "."

	for _, e := range os.Environ() {
		splits := strings.SplitN(e, "=", 2)
		key := strings.ReplaceAll(strings.ToLower(splits[0]), "_", ".")

		if strings.HasPrefix(key, prefix) {
			keys[strings.ReplaceAll(key, prefix, "")] = true
		}
	}

	for _, key := range viper.AllKeys() {
		if strings.HasPrefix(key, prefix) {
			keys[strings.ReplaceAll(key, prefix, "")] = true
		}
	}

	keySlice := make([]string, 0, len(keys))
	for k := range keys {
		keySlice = append(keySlice, k)
	}
	return keySlice
}

func initConfig() {
	if len(flagConfigFile) != 0 {
		viper.SetConfigFile(flagConfigFile)
		if err := viper.ReadInConfig(); err != nil {
			color.Red("Can't read config:", err)
			os.Exit(1)
		}
		return
	}

	var err error
	viper.SetConfigName("scheman")

	configHome := os.Getenv("XDG_CONFIG_HOME")
	homePath := os.Getenv("HOME")
	wd, err := os.Getwd()
	if err != nil {
		wd = "."
	}

	configPaths := []string{wd}
	if len(configHome) > 0 {
		configPaths = append(configPaths, filepath.Join(configHome, "scheman"))
	} else {
		configPaths = append(configPaths, filepath.Join(homePath, ".config/scheman"))
	}

	for _, p := range configPaths {
		viper.AddConfigPath(p)
	}

	// Ignore errors here, fallback to other validation methods.
	// Users can use environment variables if a config is not found.
	_ = viper.ReadInConfig()
}

func run(cmd *cobra.Command, args []string) error {
	return state.Run()
}

func postRun(cmd *cobra.Command, args []string) error {
	return state.Cleanup()
}
