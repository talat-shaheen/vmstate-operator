/*
Copyright 2023.

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

package aws

import (
	"context"

	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	awsv1 "github.com/talat-shaheen/vmstate-operator/apis/aws/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	ctrllog "sigs.k8s.io/controller-runtime/pkg/log"
)

// AWSEC2Reconciler reconciles a AWSEC2 object
type AWSEC2Reconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=aws.xyzcompany.com,resources=awsec2s,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=aws.xyzcompany.com,resources=awsec2s/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=aws.xyzcompany.com,resources=awsec2s/finalizers,verbs=update
//+kubebuilder:rbac:groups=batch,resources=jobjobs;jobs,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the AWSEC2 object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *AWSEC2Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	//_ = log.FromContext(ctx)
	log := ctrllog.FromContext(ctx)
	//log := r.Log.WithValues("AWSEC2", req.NamespacedName)
	log.Info("Reconciling AWSEC2")

	// Fetch the AWSEC2 CR
	//awsEC2, err := services.FetchAWSEC2CR(req.Name, req.Namespace)

	// Fetch the AWSEC2 instance
	awsEC2 := &awsv1.AWSEC2{}

	log.Info(req.NamespacedName.Name)

	err := r.Client.Get(ctx, req.NamespacedName, awsEC2)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			log.Info("awsEC2 resource not found. Ignoring since object must be deleted.")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		log.Error(err, "Failed to get awsEC2.")
		return ctrl.Result{}, err
	}

	log.Info(awsEC2.Name)
	// Add const values for mandatory specs ( if left blank)
	// log.Info("Adding awsEC2 mandatory specs")
	// utils.AddBackupMandatorySpecs(awsEC2)
	// Check if the jobJob already exists, if not create a new one

	found := &batchv1.Job{}
	err = r.Client.Get(ctx, types.NamespacedName{Name: awsEC2.Name, Namespace: awsEC2.Namespace}, found)
	//log.Info(*found.)
	if err != nil && errors.IsNotFound(err) {
		// Define a new jobJob
		job := r.JobForAWSEC2(awsEC2)
		log.Info("Creating a new Job", "job.Namespace", job.Namespace, "job.Name", job.Name)
		err = r.Client.Create(ctx, job)
		if err != nil {
			log.Error(err, "Failed to create new Job", "job.Namespace", job.Namespace, "job.Name", job.Name)
			return ctrl.Result{}, err
		}
		// jobjob created successfully - return and requeue
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		log.Error(err, "Failed to get job")
		return ctrl.Result{}, err
	}

	// Check for any updates for redeployment
	/*applyChange := false

	// Ensure image name is correct, update image if required
	newInstanceIds := awsEC2.Spec.InstanceIds
	log.Info(newInstanceIds)

	newStartSchedule := awsEC2.Spec.StartSchedule
	log.Info(newStartSchedule)

	newImage := awsEC2.Spec.Image
	log.Info(newImage)

	var currentImage string = ""
	var currentStartSchedule string = ""
	var currentInstanceIds string = ""

	// Check existing schedule
	if found.Spec.Schedule != "" {
		currentStartSchedule = found.Spec.Schedule
	}

	if newStartSchedule != currentStartSchedule {
		found.Spec.Schedule = newStartSchedule
		applyChange = true
	}

	// Check existing image
	if found.Spec.JobTemplate.Spec.Template.Spec.Containers != nil {
		currentImage = found.Spec.JobTemplate.Spec.Template.Spec.Containers[0].Image
	}

	if newImage != currentImage {
		found.Spec.JobTemplate.Spec.Template.Spec.Containers[0].Image = newImage
		applyChange = true
	}

	// Check instanceIds
	if found.Spec.JobTemplate.Spec.Template.Spec.Containers != nil {
		currentInstanceIds = found.Spec.JobTemplate.Spec.Template.Spec.Containers[0].Env[0].Value
		log.Info(currentInstanceIds)
	}

	if newInstanceIds != currentInstanceIds {
		found.Spec.JobTemplate.Spec.Template.Spec.Containers[0].Env[0].Value = newInstanceIds
		applyChange = true
	}

	log.Info(currentInstanceIds)
	log.Info(currentImage)
	log.Info(currentStartSchedule)

	log.Info(strconv.FormatBool(applyChange))

	if applyChange {
		log.Info(strconv.FormatBool(applyChange))
		err = r.Client.Update(ctx, found)
		if err != nil {
			log.Error(err, "Failed to update jobJob", "jobJob.Namespace", found.Namespace, "jobJob.Name", found.Name)
			return ctrl.Result{}, err
		}
		// Spec updated - return and requeue
		return ctrl.Result{Requeue: true}, nil
	}*/

	// Update the AWSEC2 status
	// TODO: Define what needs to be added in status. Currently adding just instanceIds
	/*if !reflect.DeepEqual(currentInstanceIds, awsEC2.Status.VMStartStatus) ||
		!reflect.DeepEqual(currentInstanceIds, awsEC2.Status.VMStopStatus) {
		awsEC2.Status.VMStartStatus = currentInstanceIds
		awsEC2.Status.VMStopStatus = currentInstanceIds
		err := r.Client.Status().Update(ctx, awsEC2)
		if err != nil {
			log.Error(err, "Failed to update awsEC2 status")
			return ctrl.Result{}, err
		}
	}*/

	return ctrl.Result{}, nil
}

// Job Spec
func (r *AWSEC2Reconciler) JobForAWSEC2(awsEC2 *awsv1.AWSEC2) *batchv1.Job {

	job := &batchv1.Job{
		ObjectMeta: v1.ObjectMeta{
			Name:      awsEC2.Name,
			Namespace: awsEC2.Namespace,
			Labels:    AWSEC2Labels(awsEC2, "awsEC2"),
		},
		Spec: batchv1.JobSpec{
			//JobTemplate: batchv1.JobTemplateSpec{
			//	Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Name:  awsEC2.Name,
						Image: awsEC2.Spec.Image,
						Env: []corev1.EnvVar{
							{
								Name: "AWS_ACCESS_KEY_ID",
								ValueFrom: &corev1.EnvVarSource{
									SecretKeyRef: &corev1.SecretKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: "aws-secret",
										},
										Key: "aws-access-key-id",
									},
								},
							},
							{
								Name: "AWS_SECRET_ACCESS_KEY",
								ValueFrom: &corev1.EnvVarSource{
									SecretKeyRef: &corev1.SecretKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: "aws-secret",
										},
										Key: "aws-secret-access-key",
									},
								},
							},
							{
								Name: "AWS_DEFAULT_REGION",
								ValueFrom: &corev1.EnvVarSource{
									SecretKeyRef: &corev1.SecretKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: "aws-secret",
										},
										Key: "aws-default-region",
									},
								},
							}},
					}},
					RestartPolicy: "OnFailure",
				},
				//},
				//},
			},
		},
	}
	// Set awsEC2 instance as the owner and controller
	ctrl.SetControllerReference(awsEC2, job, r.Scheme)
	return job
}

func AWSEC2Labels(v *awsv1.AWSEC2, tier string) map[string]string {
	return map[string]string{
		"app":       "AWSEC2",
		"AWSEC2_cr": v.Name,
		"tier":      tier,
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *AWSEC2Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&awsv1.AWSEC2{}).
		Owns(&batchv1.Job{}).
		Complete(r)
}
