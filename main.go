package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/maru44/scheman/core"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/volatiletech/sqlboiler/v4/boilingcore"
	"github.com/volatiletech/sqlboiler/v4/drivers"
	"github.com/volatiletech/sqlboiler/v4/importers"
)

const (
	version          = "1.0.0"
	sqlBoilerVersion = "4.8.6"
)

var (
	flagConfigFile string
	state          *core.SchemanState
)

func main() {
	for _, arg := range os.Args {
		if arg == "--version" {
			fmt.Println("SQLBoiler v" + version)
			return
		}
	}

	var rootCmd = &cobra.Command{
		Use:   "scheman [flags] <driver>",
		Short: "Scheman will write database schema in notion. I'm going to add output destination like spreadsheets.",
		Long: "Scheman will write database schema in notion. I'm going to add output destination like spreadsheets.\n" +
			"Complete documentation is available at https://github.com/maru44/scheman\n" +
			"This package depends on https://github.com/volatiletech/sqlboiler. Thank you.",
		Example:       `scheman psql`,
		PreRunE:       setState,
		RunE:          run,
		PostRunE:      postRun,
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&flagConfigFile, "config", "c", "", "Filename of config file to override default lookup")

	rootCmd.PersistentFlags().StringArray("services", []string{string(core.ServiceNotion)}, "Service table definition")
	rootCmd.PersistentFlags().StringP("notion-page-id", "", "", "Page id for notion")
	rootCmd.PersistentFlags().StringP("notion-token", "", "", "Notion integration token")
	rootCmd.PersistentFlags().StringP("notion-table-index", "", "", "Table List to refer table name and its definition database id")

	rootCmd.PersistentFlags().BoolP("disable-views", "", false, "If true, views are ignored and only tables are written")
	rootCmd.PersistentFlags().StringArray("attr-ignore", []string{},
		"List of attributes that should be ignored. ('Data Type', 'Default', 'PK', 'Auto Generate', 'Unique', 'Null', 'Enum', 'Comment', 'Free Entry')",
	)

	if err := viper.BindPFlags(rootCmd.PersistentFlags()); err != nil {
		exit(err)
	}
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := rootCmd.Execute(); err != nil {
		exit(err)
	}
}

func exit(err error) {
	color.Red("%v", err)
	os.Exit(1)
}

func setState(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return errors.New("must provide a driver name")
	}

	driverName, _, err := drivers.RegisterBinaryFromCmdArg(args[0])
	if err != nil {
		switch args[0] {
		case "psql", "mysql", "mssql", "sqlite3":
			return fmt.Errorf("Driver is not installed. Please run following command.\ngo install github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-%s@v%s", args[0], sqlBoilerVersion)
		case "crdb":
			return fmt.Errorf("Driver is not installed. Please run following command.\ngo install github.com/glerchundi/sqlboiler-crdb")
		default:
			return fmt.Errorf("Driver that name is %s is not supported.", args[0])
		}
	}

	config := &boilingcore.Config{
		DriverName: driverName,
	}
	config.DriverConfig = map[string]interface{}{
		"whitelist": viper.GetStringSlice(driverName + ".whitelist"),
		"blacklist": viper.GetStringSlice(driverName + ".blacklist"),
	}

	keys := allKeys(driverName)
	for _, key := range keys {
		if key != "blacklist" && key != "whitelist" {
			prefixedKey := fmt.Sprintf("%s.%s", driverName, key)
			config.DriverConfig[key] = viper.Get(prefixedKey)
		}
	}

	config.Imports = importers.NewDefaultImports()

	state, err = core.New(config)
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
			color.Red("cannot read config:", err)
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
