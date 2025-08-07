package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	cobra.OnInitialize(initConfig)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// processConfigTemplate processes a config file as a Go template with environment variables
func processConfigTemplate(configPath string) ([]byte, error) {
	// Read the template file
	templateContent, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	// Create template with helper functions
	tmpl := template.New("config").Funcs(template.FuncMap{
		"env": func(key string) string {
			return os.Getenv(key)
		},
		"envDefault": func(key, defaultValue string) string {
			if value := os.Getenv(key); value != "" {
				return value
			}
			return defaultValue
		},
	})

	// Parse the template
	tmpl, err = tmpl.Parse(string(templateContent))
	if err != nil {
		return nil, fmt.Errorf("failed to parse config template: %w", err)
	}

	// Execute template with environment variables
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, nil); err != nil {
		return nil, fmt.Errorf("failed to execute config template: %w", err)
	}

	return buf.Bytes(), nil
}

func initConfig() {
	var configFile string

	// Don't forget to read config either from cfgFile or from home directory!
	if cfgFile != "" {
		// Use config file from the flag.
		configFile = cfgFile
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Look for config file in current directory first, then home directory
		configPaths := []string{
			filepath.Join(".", ".weather.yml"),
			filepath.Join(".", ".weather.yaml"),
			filepath.Join(home, ".weather.yml"),
			filepath.Join(home, ".weather.yaml"),
		}

		for _, path := range configPaths {
			if _, err := os.Stat(path); err == nil {
				configFile = path
				break
			}
		}
	}

	if configFile == "" {
		fmt.Println("No config file found")
		return
	}

	// Process the config file as a template
	processedConfig, err := processConfigTemplate(configFile)
	if err != nil {
		fmt.Printf("Error processing config template: %v\n", err)
		os.Exit(1)
	}

	// Set viper to read from the processed content
	viper.SetConfigType(filepath.Ext(configFile)[1:]) // Remove the dot from extension
	if err := viper.ReadConfig(bytes.NewReader(processedConfig)); err != nil {
		fmt.Printf("Error reading processed config: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Using config file:", configFile)
}
