/*
Copyright 2024.

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

package controller

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	v1 "k8s.io/client-go/applyconfigurations/core/v1"
	"k8s.io/kubernetes/pkg/controller/cronjob"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	batchv1 "tutorial.kubebuilder.io/project/api/v1"
)

var _ = Describe("CronJob Controller", func() {
	const (
		CronJobName = "test-cronjob"
		CronJobNamespace = "default"
		JobName = "test-job"

		time = time.Second * 10
		duration = time.Second * 10
		interval = time.Millisecond * 250
	)

	Context("When updating CronJob status", func() {
		It("Should increase CronJob status.Active count when new Jobs are created", func() {
			By("Creating a new CronJob")
			ctx := context.Background()
			cronJob := &cronjobv1.CronJob{
				TypeMeta: metav1.TypeMeta {
					APIVersion: "batch.tutorial.kubebuilder.io/v1",
					Kind: "CronJob",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: CronJobName,
					Namespace: CronJobNamespace,
				},
				Spec: cronjobv1.CronJobSpec {
					Schedule: "1 * * * *",
					JobTemplate: batchv1.JobTemplateSpec{
						Spec: batchv1.JobSpec{
							Template: v1.PodTemplateSpec{
								Spec: v1.PodSpec{
									Containers: []v1.Container{
										{
											Name: "test-container",
											Image: "test-image",
										},
									},
									RestartPolicy: v1.RestartPolicyOnFailure,
								},
							},
						},
					},
				},
			}
			Expect(k8sClient.Create(ctx, CronJob)).To(Succeed())
		},
	})

	Context("When reconciling a resource", func() {
		const resourceName = "test-resource"

		ctx := context.Background()

		typeNamespacedName := types.NamespacedName{
			Name:      resourceName,
			Namespace: "default", // TODO(user):Modify as needed
		}
		cronjob := &batchv1.CronJob{}

		BeforeEach(func() {
			By("creating the custom resource for the Kind CronJob")
			err := k8sClient.Get(ctx, typeNamespacedName, cronjob)
			if err != nil && errors.IsNotFound(err) {
				resource := &batchv1.CronJob{
					ObjectMeta: metav1.ObjectMeta{
						Name:      resourceName,
						Namespace: "default",
					},
					// TODO(user): Specify other spec details if needed.
				}
				Expect(k8sClient.Create(ctx, resource)).To(Succeed())
			}
		})

		AfterEach(func() {
			// TODO(user): Cleanup logic after each test, like removing the resource instance.
			resource := &batchv1.CronJob{}
			err := k8sClient.Get(ctx, typeNamespacedName, resource)
			Expect(err).NotTo(HaveOccurred())

			By("Cleanup the specific resource instance CronJob")
			Expect(k8sClient.Delete(ctx, resource)).To(Succeed())
		})
		It("should successfully reconcile the resource", func() {
			By("Reconciling the created resource")
			controllerReconciler := &CronJobReconciler{
				Client: k8sClient,
				Scheme: k8sClient.Scheme(),
			}

			_, err := controllerReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: typeNamespacedName,
			})
			Expect(err).NotTo(HaveOccurred())
			// TODO(user): Add more specific assertions depending on your controller's reconciliation logic.
			// Example: If you expect a certain status condition after reconciliation, verify it here.
		})
	})
})
