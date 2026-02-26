// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"os"

	"github.com/urfave/cli/v3"
	ctrl "sigs.k8s.io/controller-runtime"
	ctrllog "sigs.k8s.io/controller-runtime/pkg/log"

	managercmd "github.com/gardener/gardener-extension-shoot-traefik/cmd/extension-traefik/internal/manager"
	webhookcmd "github.com/gardener/gardener-extension-shoot-traefik/cmd/extension-traefik/webhook"
	"github.com/gardener/gardener-extension-shoot-traefik/pkg/version"
)

func main() {
	app := &cli.Command{
		Name:                  "gardener-extension-shoot-traefik",
		Version:               version.Version,
		EnableShellCompletion: true,
		Usage:                 "traefik extension for Gardener",
		Commands: []*cli.Command{
			managercmd.New(),
			webhookcmd.New(),
		},
	}

	ctx := ctrl.SetupSignalHandler()
	if err := app.Run(ctx, os.Args); err != nil {
		ctrllog.Log.Error(err, "failed to start traefik extension")
		os.Exit(1)
	}
}
