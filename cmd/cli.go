/*
Copyright © 2020 Michael Gruener & Simon Fuhrer

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
	"strings"

	"github.com/bedag/kusible/pkg/printer"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Cli is the core command line interface
// Wrapping the cli in its own public type might allow
// it to be integrated in other programs, creating
// something like a composite application
type Cli struct {
	RootCommand *cobra.Command
	viper       *viper.Viper
}

// NewCli creates a
func NewCli() *Cli {
	v := viper.GetViper()
	v.SetEnvPrefix(appName)
	dashReplacer := strings.NewReplacer("-", "_", ".", "_")
	v.SetEnvKeyReplacer(dashReplacer)
	v.AutomaticEnv()
	cli := &Cli{
		viper: v,
	}
	cli.RootCommand = NewRootCommand(cli)
	return cli
}

func (c *Cli) bindAllFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().VisitAll(func(f *pflag.Flag) {
		if err := c.viper.BindPFlag(f.Name, cmd.PersistentFlags().Lookup(f.Name)); err != nil {
			panic(err) // Should never happen
		}
	})
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if err := c.viper.BindPFlag(f.Name, cmd.Flags().Lookup(f.Name)); err != nil {
			panic(err) // Should never happen
		}
	})
}

// wrapper func to bind all flags with viper on command execution
// and to perform global post-command execution steps (if necessary)
func (c *Cli) wrap(f func(*Cli, *cobra.Command, []string) error) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		c.setupLogger()
		c.bindAllFlags(cmd)
		return f(c, cmd, args)
	}
}

func (c *Cli) setupLogger() {
	if c.viper.GetBool("log-json") {
		log.SetFormatter(&log.JSONFormatter{})
	}

	logLevel, err := log.ParseLevel(c.viper.GetString("log-level"))
	if err != nil {
		log.Fatal(err.Error())
	}

	// According to the logrus documentation this is very costly, but
	// if we need debug tracing or more, this is also very helpful
	// See https://github.com/sirupsen/logrus/blob/d417be0fe654de640a82370515129985b407c7e3/README.md#logging-method-name
	if c.viper.GetBool("log-functions") {
		log.SetReportCaller(true)
	}

	log.SetLevel(logLevel)
}

func (c *Cli) output(queue printer.Queue) error {
	printerFormat := c.viper.GetString("format")
	printerFields := c.viper.GetStringSlice("fields")

	format, err := printer.ParseFormat(printerFormat)
	if err != nil {
		log.WithFields(log.Fields{
			"format": printerFormat,
		}).Error("Unknown printer format")
		return err
	}

	options := printer.Options{
		ListWrapSingleItem: c.viper.GetBool("list-wrap-single"),
	}
	printer, err := printer.New(format, printerFields, queue, options)
	if err != nil {
		log.WithFields(log.Fields{
			"format": printerFormat,
			"error":  err,
		}).Error("Failed to create printer")
		return err
	}
	if !c.viper.GetBool("quiet") {
		printer.Print()
	}
	return nil
}
