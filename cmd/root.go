/*
Copyright © 2022 Edgar Ramírez-Mondragón

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"meltano.com/target-jsonl-blob/target"

	"gocloud.dev/blob"
	// _ "gocloud.dev/blob/gcsblob"
	_ "gocloud.dev/blob/fileblob"
	_ "gocloud.dev/blob/s3blob"
)

var (
	configFile string
	inputFile  string
	C          target.Config
)

var rootCmd = &cobra.Command{
	Use:   "target-jsonl-blob",
	Short: "JSONLines Singer target for blob storages",
	Run:   func(cmd *cobra.Command, args []string) {},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}

	readConfig(configFile)

	if C.Bucket == "" {
		log.Fatalf("Bucket is required")
	}

	if C.KeyTemplate == "" {
		C.KeyTemplate = target.DEFAULT_KEY_TEMPLATE
	}

	ctx := context.Background()
	bucket, err := blob.OpenBucket(ctx, C.Bucket)
	if err != nil {
		log.Fatalf("Unable to open bucket, %v", err)
	}
	defer bucket.Close()

	var lines io.Reader

	if inputFile == "" {
		lines = os.Stdin
	} else {
		lines, err = os.Open(inputFile)
		if err != nil {
			fmt.Print(err)
		}
	}

	target.ProcessLines(lines, C, ctx, bucket)

	if inputFile == "" {
		target.ProcessLines(os.Stdin, C, ctx, bucket)
	} else {
		inputLines, err := os.Open(inputFile)
		if err != nil {
			fmt.Print(err)
		}
		target.ProcessLines(inputLines, C, ctx, bucket)
	}

}

func init() {
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "Config file")
	rootCmd.PersistentFlags().StringVarP(&inputFile, "input", "i", "", "Input file")
	rootCmd.MarkPersistentFlagRequired("config")
}

// initConfig reads in config file and ENV variables if set.
func readConfig(file string) {
	viper.SetConfigFile(file)

	if err := viper.ReadInConfig(); err == nil {
		log.Printf("Using config file %s", viper.ConfigFileUsed())
	} else {
		fmt.Fprintln(os.Stderr, "Config file is not valid")
		os.Exit(1)
	}

	if err := viper.Unmarshal(&C); err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}
}
