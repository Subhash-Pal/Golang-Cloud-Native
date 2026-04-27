package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
     "strings"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	rootCmd = &cobra.Command{
		Use:   "myapp",
		Short: "Advanced Cobra CLI application",
		Long: `A fully-featured CLI application demonstrating:
- Viper configuration management
- Persistent and local flags
- Custom logging
- Subcommand inheritance
- Automatic config file creation`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// Initialize config before any command runs
			initConfig()
			setupLogger()
		},
	}

	verbose bool
	cfgFile string
)

// Application configuration struct
type Config struct {
	APIKey      string        `mapstructure:"api_key"`
	Timeout     time.Duration `mapstructure:"timeout"`
	MaxRetries  int           `mapstructure:"max_retries"`
	LogLevel    string        `mapstructure:"log_level"`
	OutputDir   string        `mapstructure:"output_dir"`
	FeatureFlag bool          `mapstructure:"feature_flag"`
}

var appConfig Config

func init() {
	// Persistent flags (available to all commands)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.myapp.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().StringP("output-dir", "o", "", "output directory (overrides config)")

	// Bind persistent flags to Viper
	viper.BindPFlag("output_dir", rootCmd.PersistentFlags().Lookup("output-dir"))
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))

	// Set default config values
	viper.SetDefault("timeout", "30s")
	viper.SetDefault("max_retries", 3)
	viper.SetDefault("log_level", "info")
	viper.SetDefault("feature_flag", false)

	// Add subcommands
	rootCmd.AddCommand(greetCmd)
	rootCmd.AddCommand(processCmd)
	rootCmd.AddCommand(configCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func initConfig() {
	// Find home directory
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	// Set config name and type
	viper.SetConfigName(".myapp")
	viper.SetConfigType("yaml")

	// Add possible config paths
	viper.AddConfigPath(".")
	viper.AddConfigPath(filepath.Join(home))

	// Use specified config file if provided
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	}

	// Read config (ignore error if config doesn't exist)
	if err := viper.ReadInConfig(); err == nil {
		log.Printf("Using config file: %s", viper.ConfigFileUsed())
	} else {
		// Create default config if none exists
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			createDefaultConfig(home)
		}
	}

	// Unmarshal config into struct
	if err := viper.Unmarshal(&appConfig); err != nil {
		log.Fatalf("unable to decode config into struct: %v", err)
	}

	// Handle duration conversion
	timeoutStr := viper.GetString("timeout")
	if timeout, err := time.ParseDuration(timeoutStr); err == nil {
		appConfig.Timeout = timeout
	} else {
		log.Fatalf("invalid timeout format: %v", err)
	}
}

func createDefaultConfig(home string) {
	configPath := filepath.Join(home, ".myapp.yaml")
	defaultConfig := []byte(`api_key: "your_api_key_here"
timeout: "30s"
max_retries: 3
log_level: "info"
output_dir: "./output"
feature_flag: false
`)

	if err := os.WriteFile(configPath, defaultConfig, 0644); err != nil {
		log.Fatalf("failed to create default config: %v", err)
	}
	log.Printf("Created default config: %s", configPath)
	viper.SetConfigFile(configPath)
	viper.ReadInConfig()
}

func setupLogger() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)
	if appConfig.LogLevel == "debug" || verbose {
		log.SetOutput(os.Stdout)
	} else {
		log.SetOutput(os.Stderr)
	}
}

var greetCmd = &cobra.Command{
	Use:   "greet [name]",
	Short: "Send a greeting",
	Long:  "Send a personalized greeting with optional enthusiasm level",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		enthusiasm, _ := cmd.Flags().GetInt("enthusiasm")
		name := args[0]

		msg := fmt.Sprintf("Hello, %s!", name)
		if enthusiasm > 0 {
			msg += strings.Repeat("!", enthusiasm)
		}

		if verbose || appConfig.LogLevel == "debug" {
			log.Printf("Greeting generated (enthusiasm=%d)", enthusiasm)
		}

		fmt.Println(msg)
	},
}

func init() {
	greetCmd.Flags().IntP("enthusiasm", "e", 0, "level of enthusiasm (0-5)")
	greetCmd.Flags().Int("max-enthusiasm", 5, "maximum allowed enthusiasm")
	greetCmd.MarkFlagRequired("max-enthusiasm")
}

var processCmd = &cobra.Command{
	Use:   "process",
	Short: "Process data files",
	Long: `Process data files with advanced configuration:
- Uses output directory from config
- Respects timeout and retry settings
- Supports dry-run mode`,
	Run: func(cmd *cobra.Command, args []string) {
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		files, _ := cmd.Flags().GetStringSlice("files")

		if len(files) == 0 {
			log.Fatal("no files specified")
		}

		log.Printf("Starting processing (timeout: %v, retries: %d)", 
			appConfig.Timeout, appConfig.MaxRetries)
		
		if dryRun {
			log.Printf("DRY RUN: Would process %d files", len(files))
			return
		}

		// Actual processing would happen here
		fmt.Printf("Processing %d files with output to: %s\n", 
			len(files), appConfig.OutputDir)
	},
}

func init() {
	processCmd.Flags().StringSliceP("files", "f", []string{}, "files to process")
	processCmd.Flags().Bool("dry-run", false, "perform a trial run")
	processCmd.MarkFlagRequired("files")
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration",
	Long:  "View or reset application configuration",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Current configuration:\n")
		fmt.Printf("API Key: %s\n", maskAPIKey(appConfig.APIKey))
		fmt.Printf("Timeout: %v\n", appConfig.Timeout)
		fmt.Printf("Max Retries: %d\n", appConfig.MaxRetries)
		fmt.Printf("Log Level: %s\n", appConfig.LogLevel)
		fmt.Printf("Output Dir: %s\n", appConfig.OutputDir)
		fmt.Printf("Feature Flag: %v\n", appConfig.FeatureFlag)
	},
}

func maskAPIKey(key string) string {
	if len(key) < 8 {
		return "****"
	}
	return key[:4] + "****" + key[len(key)-4:]
}

// Helper function to join strings (since we can't import strings in single file)
func stringsRepeat(s string, count int) string {
	if count <= 0 {
		return ""
	}
	b := make([]byte, len(s)*count)
	bp := copy(b, s)
	for bp < len(b) {
		copy(b[bp:], b[:bp])
		bp *= 2
	}
	return string(b)
}
/*
# Initialize default config
./myapp

# Greet with enthusiasm
./myapp greet Alice -e 3

# Process files with dry run
./myapp process -f file1.txt file2.csv --dry-run

# View current configuration
./myapp config

# Override config with flags
./myapp process -f data.log -o /tmp/results --timeout 60s
*/