// Package main can be used to build the Vault GitHub secrets plugin.
package main

import (
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/plugin"

	"github.com/martinbaillie/vault-plugin-secrets-github/github"
)

func main() {
	// KEEPSAKE FIX - logging
	f, err := os.OpenFile("/var/log/vault/vault-plugin-secrets-github.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		hclog.New(&hclog.LoggerOptions{}).Error(
			"plugin shutting down",
			"failed to create log file",
			err,
		)
		os.Exit(2)
	}
	devLogger := hclog.New(&hclog.LoggerOptions{
		Name:   "vault-plugin-secrets-github",
		Level:  hclog.LevelFromString("DEBUG"),
		Output: f,
	})
	// END KEEPSAKE FIX

	devLogger.Info("Running main..")
	devLogger.Info("args: ", os.Args[1:])

	devLogger.Info("Setting up plugin client meta")
	apiClientMeta := &api.PluginAPIClientMeta{}

	flags := apiClientMeta.FlagSet()
	devLogger.Info("Parsing flags..")
	if err := flags.Parse(os.Args[1:]); err != nil {
		devLogger.Error("Flags parsing failed", err)
		fatalErr(err)
	}
	devLogger.Info("Parsed flags OK", flags)

	devLogger.Info("Getting TLS config..")
	tlsConfig := apiClientMeta.GetTLSConfig()
	devLogger.Info("Setting up TLS provider with config", tlsConfig)
	tlsProviderFunc := api.VaultPluginTLSProvider(tlsConfig)

	devLogger.Info("Starting to serve..")
	if err := plugin.Serve(&plugin.ServeOpts{
		BackendFactoryFunc: github.Factory,
		TLSProviderFunc:    tlsProviderFunc,
		Logger:             devLogger,
	}); err != nil {
		devLogger.Error("Error when serving", err)
		fatalErr(err)
	}
}

func fatalErr(err error) {
	hclog.New(&hclog.LoggerOptions{}).Error(
		"plugin shutting down",
		"error",
		err,
	)
	os.Exit(1)
}
