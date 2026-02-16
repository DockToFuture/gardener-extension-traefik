// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package shoot

import (
	"context"
	"testing"

	gardencorev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	admissionv1 "k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

func TestShootWebhook(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Shoot Webhook Suite")
}

var _ = Describe("Shoot Validator", func() {
	var (
		validator *Validator
		scheme    *runtime.Scheme
		encoder   runtime.Encoder
	)

	BeforeEach(func() {
		scheme = runtime.NewScheme()
		Expect(gardencorev1beta1.AddToScheme(scheme)).To(Succeed())

		client := fake.NewClientBuilder().WithScheme(scheme).Build()
		logger := zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true))
		validator = NewValidator(client, logger)

		serializer := json.NewSerializer(json.DefaultMetaFactory, scheme, scheme, false)
		encoder = serializer
	})

	encodeShoot := func(shoot *gardencorev1beta1.Shoot) []byte {
		data, err := runtime.Encode(encoder, shoot)
		Expect(err).NotTo(HaveOccurred())

		return data
	}

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

			req := admission.Request{
				AdmissionRequest: admissionv1.AdmissionRequest{
					Name:      shoot.Name,
					Namespace: shoot.Namespace,
					Object:    runtime.RawExtension{Raw: encodeShoot(shoot)},
				},
			}

			response := validator.Handle(context.Background(), req)
			Expect(response.Allowed).To(BeTrue())
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

			req := admission.Request{
				AdmissionRequest: admissionv1.AdmissionRequest{
					Name:      shoot.Name,
					Namespace: shoot.Namespace,
					Object:    runtime.RawExtension{Raw: encodeShoot(shoot)},
				},
			}

			response := validator.Handle(context.Background(), req)
			Expect(response.Allowed).To(BeFalse())
			Expect(response.Result.Message).To(ContainSubstring("evaluation"))
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

			req := admission.Request{
				AdmissionRequest: admissionv1.AdmissionRequest{
					Name:      shoot.Name,
					Namespace: shoot.Namespace,
					Object:    runtime.RawExtension{Raw: encodeShoot(shoot)},
				},
			}

			response := validator.Handle(context.Background(), req)
			Expect(response.Allowed).To(BeFalse())
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

			req := admission.Request{
				AdmissionRequest: admissionv1.AdmissionRequest{
					Name:      shoot.Name,
					Namespace: shoot.Namespace,
					Object:    runtime.RawExtension{Raw: encodeShoot(shoot)},
				},
			}

			response := validator.Handle(context.Background(), req)
			Expect(response.Allowed).To(BeTrue())
		})
	})
})
