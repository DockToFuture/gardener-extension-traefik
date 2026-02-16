// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

// Package shoot provides admission webhook handlers for Shoot resources.
package shoot

import (
	"context"
	"fmt"
	"net/http"

	gardencorev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// Validator validates Shoot resources for the Traefik extension.
type Validator struct {
	client  client.Client
	decoder admission.Decoder
	logger  logr.Logger
}

// NewValidator creates a new Validator for Shoot resources.
func NewValidator(c client.Client, logger logr.Logger) *Validator {
	return &Validator{
		client:  c,
		decoder: admission.NewDecoder(c.Scheme()),
		logger:  logger.WithName("shoot-validator"),
	}
}

// Handle validates the Shoot resource to ensure Traefik extension is only
// enabled for shoots with purpose "evaluation".
func (v *Validator) Handle(ctx context.Context, req admission.Request) admission.Response {
	v.logger.V(1).Info("validating shoot", "name", req.Name, "namespace", req.Namespace)

	shoot := &gardencorev1beta1.Shoot{}
	if err := v.decoder.Decode(req, shoot); err != nil {
		v.logger.Error(err, "failed to decode shoot")

		return admission.Errored(http.StatusBadRequest, err)
	}

	// Check if the Traefik extension is being added
	hasTraefikExtension := false
	for _, ext := range shoot.Spec.Extensions {
		if ext.Type == "traefik" {
			hasTraefikExtension = true

			break
		}
	}

	// If no Traefik extension, allow the request
	if !hasTraefikExtension {
		return admission.Allowed("no traefik extension configured")
	}

	// Validate that the shoot purpose is "evaluation"
	if shoot.Spec.Purpose == nil || *shoot.Spec.Purpose != gardencorev1beta1.ShootPurposeEvaluation {
		purposeStr := "nil"
		if shoot.Spec.Purpose != nil {
			purposeStr = string(*shoot.Spec.Purpose)
		}
		v.logger.Info("denying shoot with traefik extension: purpose is not 'evaluation'",
			"name", shoot.Name,
			"namespace", shoot.Namespace,
			"purpose", purposeStr,
		)

		return admission.Denied(fmt.Sprintf(
			"Traefik extension can only be enabled for shoots with purpose 'evaluation'. "+
				"Current purpose: %s. Traefik acts as a replacement for the nginx ingress controller "+
				"and is only supported for evaluation clusters.",
			purposeStr,
		))
	}

	v.logger.Info("allowing shoot with traefik extension: purpose is 'evaluation'",
		"name", shoot.Name,
		"namespace", shoot.Namespace,
	)

	return admission.Allowed("shoot purpose is evaluation, traefik extension allowed")
}
