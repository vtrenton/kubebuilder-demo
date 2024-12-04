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

package v1

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
	batchv1 "tutorial.kubebuilder.io/project/api/v1"
)

// log is for logging in this package.
var cronjoblog = logf.Log.WithName("cronjob-resource")

// SetupWebhookWithManager will setup the manager to manage the webhooks
func (r *CronJob) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).For(&batchv1.CronJob{}).
		WithValidator(&CronJobCustomValidator{}).
		WithDefaulter(&CronJobCustomDefaulter{
			DefaultConcurrencyPolicy:         batchv1.AllowCurrent,
			DefaultSuspend:                   false,
			DefaultSuccessfulJobHistoryLimit: 3,
			DefaultFailedJobsHistoryLimit:    1,
		}).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// +kubebuilder:webhook:path=/mutate-batch-tutorial-kubebuilder-io-v1-cronjob,mutating=true,failurePolicy=fail,sideEffects=None,groups=batch.tutorial.kubebuilder.io,resources=cronjobs,verbs=create;update,versions=v1,name=mcronjob.kb.io,admissionReviewVersions=v1

type CronJobCustomDefaulter struct {
	DefaultConcurrencyPolicy         batchv1.ConcurrencyPolicy
	DefaultSuspend                   bool
	DefaultSuccessfulJobHistoryLimit int32
	DefaultFailedJobsHistoryLimit    int32
}

var _ webhook.CustomDefaulter = &CronJobCustomDefaulter{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (d *CronJobCustomDefaulter) Default(ctx context.Context, obj runtime.Object) error {
	cronjob, ok := obj.(*batchv1.CronJob)

	if !ok {
		return fmt.Errorf("Expected a cronjob object but got a %T", obj)
	}
	cronjoblog.Info("Defaulting for CronJob", "name", cronjob.GetName())

	d.applyDefaults(cronjob)
	return nil
}

func (d *CronJobCustomDefaulter) applyDefaults(cronJob *batchv1.CronJob) {
	if cronJob.Spec.ConcurrencyPolicy == "" {
		cronJob.Spec.ConcurrencyPolicy = d.DefaultConcurrencyPolicy
	}
	if cronJob.Spec.Suspend == nil {
		cronJob.Spec.Suspend = new(bool)
		*&cronJob.Spec.Suspend = d.DefaultSuspend
	}
	if cronJob.Spec.SuccessfulJobsHistoryLimit == nil {
		cronJob.Spec.SuccessfulJobsHistoryLimit = new(int32)
		*cronJob.Spec.SuccessfulJobsHistoryLimit = d.DefaultSuccessfulJobHistoryLimit
	}
	if cronJob.Spec.FailedJobsHistoryLimit == nil {
		cronJob.Spec.FailedJobsHistoryLimit = new(int32)
		*cronJob.Spec.FailedJobsHistoryLimit = d.DefaultFailedJobsHistoryLimit
	}
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
// NOTE: The 'path' attribute must follow a specific pattern and should not be modified directly here.
// Modifying the path for an invalid path can cause API server errors; failing to locate the webhook.
// +kubebuilder:webhook:path=/validate-batch-tutorial-kubebuilder-io-v1-cronjob,mutating=false,failurePolicy=fail,sideEffects=None,groups=batch.tutorial.kubebuilder.io,resources=cronjobs,verbs=create;update,versions=v1,name=vcronjob.kb.io,admissionReviewVersions=v1

type CronJobCustomValidator struct {
	// TODO
}

var _ webhook.CronJobCustomValidator = &CronJobCustomValidator{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *CronJob) ValidateCreate(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	cronjob, ok := obj.(*batchv1.CronJob)
	if !ok {
		return nil, fmt.Errorf("expected a CronJob object but got %T", obj)
	}
	cronjoblog.Info("Validation for CronJob Creation", "name", cronjob.GetName())

	return nil, validateCronJob(cronjob)
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *CronJob) ValidateUpdate(ctx context.Context, newObj runtime.Object) (admission.Warnings, error) {
	cronjob, ok := newObj.(*batchv1.CronJob)
	if !ok {
		return nil, fmt.Errorf("expected a CronJob for the newObj but got %T", newObj)
	}
	cronjoblog.Info("Validation for cronjob upon update", "name", cronjob.GetName())

	// TODO(user): fill in your validation logic upon object update.
	return nil, validateCronJob(cronjob)
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *CronJob) ValidateDelete(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	cronjob, ok := obj.(*batchv1.CronJob)
	if !ok {
		return nil, fmt.Errorf("expected CronJob object but got %T", obj)
	}
	cronjoblog.Info("Validation for cronjob upon update", "name", cronjob.GetName())

	// TODO(user): fill in your validation logic upon object update.
	return nil, nil
}

func validateCronJob(cronjob *batchv1.CronJob) error {
	var allErrs field.ErrorList
	if err := validateCronJobName(cronjob); err != nil {
		allErrs = append(allErrs, err)
	}
	if err := validateCronJobSpec(cronjob); err != nil {
		allErrs = append(allErrs, err)
	}
	if len(allErrs) == 0 {
		return nil
	}

	return apierrors.NewInvalid(schema.GroupKind{
		Group: "batch.tutorial.kubebuilder.io", Kind: "CronJob"},
		cronjob.Name, allErrs)
}
