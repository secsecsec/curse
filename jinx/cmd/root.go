// Copyright © 2017 Michael Smith <mikejsmitty@gmail.com>
//
// This file is part of jinx.
//
// jinx is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// jinx is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with jinx. If not, see <http://www.gnu.org/licenses/>.
//

package cmd

import (
	"fmt"
	"os"

	"github.com/mikesmitty/curse/jinx/jinxlib"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "jinx",
	Short: "SSH certificate client",
	Long: `JINX is a client to the CURSE SSH certificate authority.
It is used to provide short-lived SSH certificates in place of semi-permanent SSH pubkeys
in authorized_keys files, which are difficult to manage at scale and over long periods
of time.`,
	Run: func(cmd *cobra.Command, args []string) {
		jinxlib.Jinx(args)
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	//RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.jinx.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	//if cfgFile != "" { // enable ability to specify config file via flag
	//	viper.SetConfigFile(cfgFile)
	//}

	viper.SetConfigName("jinx") // name of config file (without extension)
	viper.AddConfigPath("/etc/jinx")
	viper.AddConfigPath("$HOME/.jinx/")

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	viper.SetDefault("autogenkeys", true)
	viper.SetDefault("bastionip", "")
	viper.SetDefault("insecure", false)
	viper.SetDefault("keygenbitsize", 2048)
	viper.SetDefault("keygenpubkey", "$HOME/.ssh/id_jinx.pub")
	viper.SetDefault("keygentype", "ed25519")
	viper.SetDefault("pubkey", "$HOME/.ssh/id_ed25519.pub")
	viper.SetDefault("sshuser", "root") // FIXME Need to revisit this?
	viper.SetDefault("timeout", 30)
	viper.SetDefault("url", "https://localhost/")
}
