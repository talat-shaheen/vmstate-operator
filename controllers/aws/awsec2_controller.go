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

	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	ctrllog "sigs.k8s.io/controller-runtime/pkg/log"
)

// AWSEC2Reconciler reconciles a AWSEC2 object
type AWSEC2Reconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=aws.xyzcompany.com,resources=awsec2s,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=aws.xyzcompany.com,resources=awsec2s/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=aws.xyzcompany.com,resources=awsec2s/finalizers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=batch,resources=jobs/finalizers;jobs,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the AWSEC2 object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile

const AWSEC2Finalizer = "aws.xyzcompany.com/finalizer"

func (r *AWSEC2Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	//_ = log.FromContext(ctx)
	log := ctrllog.FromContext(ctx)
	//log := r.Log.WithValues("AWSEC2", req.NamespacedName)
	log.Info("Reconciling AWSEC2s CRs")

	// Fetch the AWSEC2 CR
	//awsEC2, err := services.FetchAWSEC2CR(req.Name, req.Namespace)

	// Fetch the AWSEC2 instance
	awsEC2 := &awsv1.AWSEC2{}
	//ctrl.SetControllerReference(awsEC2, awsEC2, r.Scheme)
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

	// Add const values for mandatory specs ( if left blank)
	// log.Info("Adding awsEC2 mandatory specs")
	// utils.AddBackupMandatorySpecs(awsEC2)
	// Check if the jobJob already exists, if not create a new one

	found := &batchv1.Job{}
	err = r.Client.Get(ctx, types.NamespacedName{Name: awsEC2.Name + "create", Namespace: awsEC2.Namespace}, found)
	//log.Info(*found.)
	if err != nil && errors.IsNotFound(err) {
		// Define a new Job
		job := r.JobForAWSEC2(awsEC2, "create")
		log.Info("Creating a new Job", "job.Namespace", job.Namespace, "job.Name", job.Name)
		err = r.Client.Create(ctx, job)
		if err != nil {
			log.Error(err, "Failed to create new Job", "job.Namespace", job.Namespace, "job.Name", job.Name)
			return ctrl.Result{}, err
		}
		// job created successfully - return and requeue
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
	// Check if the AWSEC2 instance is marked to be deleted, which is
	// indicated by the deletion timestamp being set.
	isAWSEC2MarkedToBeDeleted := awsEC2.GetDeletionTimestamp() != nil
	if isAWSEC2MarkedToBeDeleted {
		if controllerutil.ContainsFinalizer(awsEC2, AWSEC2Finalizer) {
			// Run finalization logic for AWSEC2Finalizer. If the
			// finalization logic fails, don't remove the finalizer so
			// that we can retry during the next reconciliation.
			log.Info(awsEC2.Name)
			log.Info("CR is marked for deletion")
			if err := r.finalizeAWSEC2(ctx, awsEC2); err != nil {
				return ctrl.Result{}, err
			}

			// Remove AWSEC2Finalizer. Once all finalizers have been
			// removed, the object will be deleted.
			controllerutil.RemoveFinalizer(awsEC2, AWSEC2Finalizer)
			err := r.Client.Update(ctx, awsEC2)
			if err != nil {
				return ctrl.Result{}, err
			}
			log.Info("Finalizer removed")
			log.Info(awsEC2.Name)
		}
		return ctrl.Result{}, nil
	}

	// Add finalizer for this CR
	if !controllerutil.ContainsFinalizer(awsEC2, AWSEC2Finalizer) {
		log.Info("Finalizer added again")
		log.Info(awsEC2.Name)
		controllerutil.AddFinalizer(awsEC2, AWSEC2Finalizer)
		err = r.Client.Update(ctx, awsEC2)
		if err != nil {
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}
func (r *AWSEC2Reconciler) finalizeAWSEC2(ctx context.Context, awsEC2 *awsv1.AWSEC2) error {
	// TODO(user): Add the cleanup steps that the operator
	// needs to do before the CR can be deleted. Examples
	// of finalizers include performing backups and deleting
	// resources that are not owned by this CR, like a PVC.
	log := ctrllog.FromContext(ctx)
	log.Info("Successfully finalized AWSEC2")
	found := &batchv1.Job{}
	err := r.Client.Get(ctx, types.NamespacedName{Name: awsEC2.Name + "delete", Namespace: awsEC2.Namespace}, found)
	//log.Info(*found.)
	if err != nil && errors.IsNotFound(err) {
		// Define a new job
		job := r.JobForAWSEC2(awsEC2, "delete")
		log.Info("Creating a new Job", "job.Namespace", job.Namespace, "job.Name", job.Name)
		err = r.Client.Create(ctx, job)
		if err != nil {
			log.Error(err, "Failed to create new Job", "job.Namespace", job.Namespace, "job.Name", job.Name)
			return err
		}
		// job created successfully - return and requeue
		return nil
	} else if err != nil {
		log.Error(err, "Failed to get job")
		return err
	}
	return nil
}

// Job Spec
func (r *AWSEC2Reconciler) JobForAWSEC2(awsEC2 *awsv1.AWSEC2, command string) *batchv1.Job {
	jobName := awsEC2.Name + command
	job := &batchv1.Job{
		ObjectMeta: v1.ObjectMeta{
			Name:      jobName,
			Namespace: awsEC2.Namespace,
			Labels:    AWSEC2Labels(awsEC2, "awsEC2"),
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Name:  awsEC2.Name,
						Image: awsEC2.Spec.Image,
						Env: []corev1.EnvVar{
							{
								Name:  "ec2_command",
								Value: command,
							},
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
