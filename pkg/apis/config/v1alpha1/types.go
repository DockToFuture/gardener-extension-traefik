// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TraefikConfigSpec defines the desired state of [TraefikConfig]
type TraefikConfigSpec struct {
	// Image is the Traefik container image to use.
	Image string `json:"image,omitempty"`

	// Replicas is the number of Traefik replicas to deploy.
	// Defaults to 2 if not specified.
	Replicas int32 `json:"replicas,omitempty"`

	// IngressClass is the ingress class name that Traefik will handle.
	// Defaults to "traefik" if not specified.
	// This replaces the deprecated nginx ingress class.
	IngressClass string `json:"ingressClass,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// TraefikConfig is the configuration schema for the Traefik extension.
// This extension deploys Traefik ingress controller to shoot clusters
// as a replacement for the nginx-ingress-controller which is out of maintenance.
type TraefikConfig struct {
	metav1.TypeMeta `json:",inline"`

	// Spec provides the Traefik extension configuration spec.
	Spec TraefikConfigSpec `json:"spec"`
}
