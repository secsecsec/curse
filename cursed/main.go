package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/spf13/viper"
)

type config struct {
	//cfgFile    string
	caSigner   ssh.Signer
	caKey      string
	duration   int
	dur        time.Duration
	extensions []string
	exts       map[string]string
	forceCmd   bool
	port       int
}

func main() {
	var conf config
	err := viper.Unmarshal(&conf)
	if err != nil {
		log.Fatal("Unable to read config into struct: ", err)
	}

	//caKey = viper.GetString("cakey")
	//durInt := viper.GetInt("duration")
	//confExts := viper.GetStringSlice("extensions")
	//forceCmd = viper.GetBool("force-command")
	//port = viper.GetInt("port")

	// Check our extensions for validity
	var errSlice []error
	conf.exts, errSlice = validateExtensions(conf.extensions)
	if len(errSlice) > 0 {
		for _, err := range errSlice {
			log.Printf("%v", err)
		}
	}

	// Convert our cert validity duration from int to time.Duration
	conf.dur = time.Duration(conf.duration) * time.Second

	fmt.Printf("%v", conf) // DEBUG

	conf.caSigner, err = loadCAKey(conf.caKey)
	if err != nil {
		log.Fatal(err)
	}

	// Start web service
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		webHandler(w, r, &conf)
	})

	log.Printf("Starting HTTP server on %d", conf.port)
	err = http.ListenAndServe(":"+strconv.Itoa(conf.port), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func init() {
	//if cfgFile != "" { // enable ability to specify config file via flag
	//	viper.SetConfigFile(cfgFile)
	//}

	viper.SetConfigName(".cursed") // name of config file (without extension)
	viper.AddConfigPath("$HOME")   // adding home directory as first search path
	viper.ReadInConfig()

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Printf("Using config file: %s", viper.ConfigFileUsed())
	}

	viper.SetDefault("cakey", "test_keys/user_ca")
	viper.SetDefault("duration", 2*60)
	viper.SetDefault("extensions", []string{"permit-pty"})
	viper.SetDefault("forcecmd", false)
	viper.SetDefault("port", 8000)
}

func validateExtensions(confExts []string) (map[string]string, []error) {
	validExts := []string{"permit-X11-forwarding", "permit-agent-forwarding",
		"permit-port-forwarding", "permit-pty", "permit-user-rc"}
	exts := make(map[string]string)
	errSlice := make([]error, 0)

	// Compare each of the config items from our config file against our known-good list, and
	// add them as a key in a map[string]string with empty value, as SSH expects
	for i := range confExts {
		valid := false
		for j := range validExts {
			if confExts[i] == validExts[j] {
				name := confExts[i]
				exts[name] = ""
				valid = true
				break
			}
		}
		if !valid {
			err := fmt.Errorf("Invalid extension in config: %s", confExts[i])
			errSlice = append(errSlice, err)
		}
	}

	return exts, errSlice
}
