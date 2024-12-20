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
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ConcurrencyPolicy describes how the job will be handled.
// Only one of the following concurrent policies may be specified.
// If none of the following policies is specified, the default one
// is AllowConcurrent.
// +kubebuilder:validation:Enum=Allow;Forbid;Replace
type ConcurrencyPolicy string

const (
	// AllowConcurrent allows concurrent jobs
	AllowConcurrent ConcurrencyPolicy = "Allow"

	// ForbidConcurrent forbids concurrent runs
	ForbidConcurrent ConcurrencyPolicy = "Forbid"

	// ReplaceConcurrent replaces the current run
	ReplaceConcurrent ConcurrencyPolicy = "Replace"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// CronJobSpec defines the desired state of CronJob
type CronJobSpec struct {
	// +kubebuilder:validation:MinLength=0

	// The schedule in the CronJob format
	Schedule string `json:"schedule"`

	// +kubebuilder:validation:Minimum=0

	// Optional deadline in seconds for starting the job
	// if it misses scheduling the job for any reason.
	// missed jobs will be counted as failed ones
	// +optional
	StartingDeadlineSeconds *int64 `json:"startingDeadlineSeconds,omitempty"`

	// Specifies how to treat concurrent exexcutions of a job.
	// Valid values are:
	// - "Allow" (default): allows CronJobs to run concurrently
	// - "Forbid": forbids concurrent runs. Skips next one if the current has not finished yet
	// - "Replace": Cancels the current job and replaces it with a the new one.
	// +optional
	ConcurrencyPolicy ConcurrencyPolicy `json:"concurrencyPolicy,omitempty"`

	// This flag tells the controller to suppend further executions
	// it does not apply to currently running jobs. defaults to false.
	// +optional
	Suspend *bool `json:"suspend,omitempty"`

	// Specifies the job that will be created when executing the CronJob.
	JobTemplate batchv1.JobTemplateSpec `json:"jobTemplate,omitempty"`

	// kubebuilder:validation:Minimum=0

	// The number of finished jobs to retain
	// This is a pointer to distinguish between explicit zero and not specified
	// +optional
	SucessfulJobsHistoryLimit *int32 `json:"sucessfulJobsHistoryLimit,omitempty"`

	// +kubebuilder:validation:Minimum=0

	// The number of fialed jobs to retain
	// This is a pointer to distinguish between explicit zero and not specified
	// +optional
	FailedJobsHistoryLimit *int32 `json:"failedJobsHistoryLimit,omitempty"`
}

// CronJobStatus defines the observed state of CronJob
type CronJobStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// A list of pointers to currently running jobs
	// +optional
	Active []corev1.ObjectReference `json:"active:omitempty"`

	// Information as to the last time the job was sucessfully scheduled
	// +optional
	LastScheduledTime *metav1.Time `json:"lastScheduledTime,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// CronJob is the Schema for the cronjobs API
type CronJob struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CronJobSpec   `json:"spec,omitempty"`
	Status CronJobStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// CronJobList contains a list of CronJob
type CronJobList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CronJob `json:"items"`
}

func init() {
	SchemeBuilder.Register(&CronJob{}, &CronJobList{})
}
