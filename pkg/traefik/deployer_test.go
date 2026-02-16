// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package traefik

import (
	"testing"

	"github.com/gardener/gardener/pkg/utils/imagevector"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestDeployment_ImageOverride(t *testing.T) {
	tests := []struct {
		name          string
		configImage   string
		imageVector   imagevector.ImageVector
		expectedImage string
		expectError   bool
		errorContains string
	}{
		{
			name:          "use config image when specified",
			configImage:   "custom.registry.io/traefik:v2.0",
			imageVector:   nil, // Should not even be consulted
			expectedImage: "custom.registry.io/traefik:v2.0",
			expectError:   false,
		},
		{
			name:        "use image vector when config empty",
			configImage: "",
			imageVector: imagevector.ImageVector{
				{
					Name:       "traefik",
					Repository: strPtr("docker.io/library/traefik"),
					Tag:        strPtr("v3.6.7"),
				},
			},
			expectedImage: "docker.io/library/traefik:v3.6.7",
			expectError:   false,
		},
		{
			name:          "fail when config empty and image not in vector",
			configImage:   "",
			imageVector:   imagevector.ImageVector{}, // Empty vector
			expectedImage: "",
			expectError:   true,
			errorContains: "failed to find traefik image",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a fake client
			scheme := runtime.NewScheme()
			client := fake.NewClientBuilder().WithScheme(scheme).Build()

			config := Config{
				Image:        tt.configImage,
				Replicas:     2,
				IngressClass: "traefik",
			}

			deployer := NewDeployer(client, logr.Discard(), config, tt.imageVector)

			deployment, err := deployer.deployment()

			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
				} else if tt.errorContains != "" && !contains(err.Error(), tt.errorContains) {
					t.Errorf("expected error to contain %q, got: %v", tt.errorContains, err)
				}

				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)

				return
			}

			if deployment == nil {
				t.Error("expected deployment but got nil")

				return
			}

			actualImage := deployment.Spec.Template.Spec.Containers[0].Image
			if actualImage != tt.expectedImage {
				t.Errorf("expected image %q, got %q", tt.expectedImage, actualImage)
			}
		})
	}
}

func strPtr(s string) *string {
	return &s
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || (len(s) > 0 && len(substr) > 0 && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}

	return false
}
