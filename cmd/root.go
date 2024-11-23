/*
Copyright © 2024 Mukund

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/MukundSinghRajput/wtr/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var version = "0.0.1"

var rootCmd = &cobra.Command{
	Use:   "wtr",
	Short: "A basic CLI tool to get weather of your city",
	Long: `Copyright © 2024 Mukund
	
A basic CLI tool to get weather of your city 
developed in hurry xD btw to report any error visit (https://t.me/TreatEveryoneEqually)`,
	Run: func(cmd *cobra.Command, args []string) {
		if versionFlag, _ := cmd.Flags().GetBool("version"); versionFlag {
			fmt.Printf("wtr CLI version %s\n", version)
			return
		}

		if cmd.Flags().Changed("set") {
			city, _ := cmd.Flags().GetString("set")
			viper.Set("default", city)
			err := viper.WriteConfig()
			if err != nil {
				fmt.Println("Failed to save the default city:", err)
				return
			}
			fmt.Printf("Default city set to: %s\n", city)
			return
		}

		v := viper.GetString("default")
		if v == "" {
			fmt.Println("Set default city using `wtr --set`")
			return
		}

		client := api.NewClient("https://wttr.in/")
		s, err := client.Today(v)
		if errors.Is(err, api.ErrCityNotFound) {
			fmt.Println("Your default city probably doesn't exist")
			return
		} else if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Println(s)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.wtr.yaml)")
	rootCmd.Flags().BoolP("version", "v", false, "Get the current version of wtr CLI")
	rootCmd.Flags().StringP("set", "s", "", "Set the default city for wtr CLI")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".wtr")
		configFilePath := fmt.Sprintf("%s/.wtr.yaml", home)
		if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
			file, err := os.Create(configFilePath)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error creating config file:", err)
				os.Exit(1)
			}
			defer file.Close()
			fmt.Fprintln(os.Stderr, "Created config file:", configFilePath)
		}
	}
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err == nil {
		return
	}
}
