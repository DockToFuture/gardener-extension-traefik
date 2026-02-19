// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

// Package webhook provides webhook handlers for the Traefik extension.
package webhook

import (
	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/webhook"

	"github.com/gardener/gardener-extension-shoot-traefik/pkg/webhook/shoot"
)

const (
	// WebhookPath is the path for the shoot validation webhook.
	WebhookPath = "/validate-shoot-traefik"
)

// AddToManager adds the webhook handlers to the manager.
func AddToManager(mgr manager.Manager, logger logr.Logger) error {
	hookServer := mgr.GetWebhookServer()

	shootValidator := shoot.NewValidator(mgr.GetClient(), logger)
	hookServer.Register(WebhookPath, &webhook.Admission{Handler: shootValidator})

	logger.Info("registered shoot validation webhook", "path", WebhookPath)

	return nil
}
