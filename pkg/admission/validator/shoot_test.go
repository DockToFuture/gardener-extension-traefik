// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package validator

import (
	"context"
	"testing"

	gardencorev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestValidator(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Admission Validator Suite")
}

var _ = Describe("Shoot Validator", func() {
	var (
		validator *shootValidator
		scheme    *runtime.Scheme
	)

	BeforeEach(func() {
		scheme = runtime.NewScheme()
		Expect(gardencorev1beta1.AddToScheme(scheme)).To(Succeed())

		client := fake.NewClientBuilder().WithScheme(scheme).Build()
		decoder := serializer.NewCodecFactory(scheme, serializer.EnableStrict).UniversalDecoder()
		validator = &shootValidator{
			client:  client,
			decoder: decoder,
		}
	})

	Context("when shoot has traefik extension", func() {
		It("should allow shoot with purpose 'evaluation'", func() {
			purpose := gardencorev1beta1.ShootPurposeEvaluation
			shoot := &gardencorev1beta1.Shoot{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "core.gardener.cloud/v1beta1",
					Kind:       "Shoot",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-shoot",
					Namespace: "garden-test",
				},
				Spec: gardencorev1beta1.ShootSpec{
					Purpose: &purpose,
					Extensions: []gardencorev1beta1.Extension{
						{Type: "traefik"},
					},
				},
			}

			err := validator.Validate(context.Background(), shoot, nil)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should deny shoot with purpose 'production'", func() {
			purpose := gardencorev1beta1.ShootPurposeProduction
			shoot := &gardencorev1beta1.Shoot{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "core.gardener.cloud/v1beta1",
					Kind:       "Shoot",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-shoot",
					Namespace: "garden-test",
				},
				Spec: gardencorev1beta1.ShootSpec{
					Purpose: &purpose,
					Extensions: []gardencorev1beta1.Extension{
						{Type: "traefik"},
					},
				},
			}

			err := validator.Validate(context.Background(), shoot, nil)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("evaluation"))
		})

		It("should deny shoot with nil purpose", func() {
			shoot := &gardencorev1beta1.Shoot{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "core.gardener.cloud/v1beta1",
					Kind:       "Shoot",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-shoot",
					Namespace: "garden-test",
				},
				Spec: gardencorev1beta1.ShootSpec{
					Purpose: nil,
					Extensions: []gardencorev1beta1.Extension{
						{Type: "traefik"},
					},
				},
			}

			err := validator.Validate(context.Background(), shoot, nil)
			Expect(err).To(HaveOccurred())
		})

		It("should deny shoot with purpose 'development'", func() {
			purpose := gardencorev1beta1.ShootPurposeDevelopment
			shoot := &gardencorev1beta1.Shoot{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "core.gardener.cloud/v1beta1",
					Kind:       "Shoot",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-shoot",
					Namespace: "garden-test",
				},
				Spec: gardencorev1beta1.ShootSpec{
					Purpose: &purpose,
					Extensions: []gardencorev1beta1.Extension{
						{Type: "traefik"},
					},
				},
			}

			err := validator.Validate(context.Background(), shoot, nil)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("evaluation"))
		})
	})

	Context("when shoot does not have traefik extension", func() {
		It("should allow shoot without traefik extension regardless of purpose", func() {
			purpose := gardencorev1beta1.ShootPurposeProduction
			shoot := &gardencorev1beta1.Shoot{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "core.gardener.cloud/v1beta1",
					Kind:       "Shoot",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-shoot",
					Namespace: "garden-test",
				},
				Spec: gardencorev1beta1.ShootSpec{
					Purpose: &purpose,
					Extensions: []gardencorev1beta1.Extension{
						{Type: "other-extension"},
					},
				},
			}

			err := validator.Validate(context.Background(), shoot, nil)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should allow shoot with no extensions", func() {
			purpose := gardencorev1beta1.ShootPurposeProduction
			shoot := &gardencorev1beta1.Shoot{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "core.gardener.cloud/v1beta1",
					Kind:       "Shoot",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-shoot",
					Namespace: "garden-test",
				},
				Spec: gardencorev1beta1.ShootSpec{
					Purpose:    &purpose,
					Extensions: nil,
				},
			}

			err := validator.Validate(context.Background(), shoot, nil)
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
