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

package awsmanager

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	ctrllog "sigs.k8s.io/controller-runtime/pkg/log"

	awsmanagerv1 "github.com/talat-shaheen/vmstate-operator/apis/awsmanager/v1"
)

// AWSManagerReconciler reconciles a AWSManager object
type AWSManagerReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=awsmanager.xyzcompany.com,resources=awsmanagers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=awsmanager.xyzcompany.com,resources=awsmanagers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=awsmanager.xyzcompany.com,resources=awsmanagers/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=configmaps;secrets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=batch,resources=jobs,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the AWSManager object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *AWSManagerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	//_ = log.FromContext(ctx)
	log := ctrllog.FromContext(ctx)
	log.Info("Reconciling AWS EC2 state manager ")
	// Fetch the awsManager instance
	awsManager := &awsmanagerv1.AWSManager{}

	//log.Info(req.NamespacedName.Name)

	err := r.Client.Get(ctx, req.NamespacedName, awsManager)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			log.Info("awsManager resource not found. Ignoring since object must be deleted.")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		log.Error(err, "Failed to get awsManager.")
		return ctrl.Result{}, err
	}

	log.Info(awsManager.Name, awsManager.Namespace, awsManager.Spec.Image)
	// Add const values for mandatory specs ( if left blank)
	log.Info("Checking awsManager mandatory specs")

	// Check if the Deployment already exists, if not create a new one

	found := &appsv1.Deployment{}
	err = r.Client.Get(ctx, types.NamespacedName{Name: awsManager.Name, Namespace: awsManager.Namespace}, found)
	//log.Info(*found.)
	if err != nil && errors.IsNotFound(err) {
		// Define a new DeploymentJob
		Deployment := r.DeploymentForAWSManager(ctx, req, awsManager)
		log.Info("Tried Creating a new Deployment", "Deployment.Namespace", Deployment.Namespace, "Deployment.Name", Deployment.Name)
		err = r.Client.Create(ctx, Deployment)
		if err != nil {
			log.Error(err, "Failed to create new Deployment", "Deployment.Namespace", Deployment.Namespace, "Deployment.Name", Deployment.Name)
			return ctrl.Result{}, err
		}
		// Deploymentjob created successfully - return and requeue
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		log.Error(err, "Failed to get Deployment")
		return ctrl.Result{}, err
	}

	/*
		newImage := awsManager.Spec.Image
		log.Info(newImage)

		var currentImage string = ""

		// Check existing image
		if found.Spec.JobTemplate.Spec.Template.Spec.Containers != nil {
			currentImage = found.Spec.JobTemplate.Spec.Template.Spec.Containers[0].Image
		}

		log.Info(currentImage)

		if applyChange {
			log.Info(strconv.FormatBool(applyChange))
			err = r.Client.Update(ctx, found)
			if err != nil {
				log.Error(err, "Failed to update DeploymentJob", "DeploymentJob.Namespace", found.Namespace, "DeploymentJob.Name", found.Name)
				return ctrl.Result{}, err
			}
			// Spec updated - return and requeue
			return ctrl.Result{Requeue: true}, nil
		}*/

	// Update the awsManager status
	// TODO: Define what needs to be added in status. Currently adding just instanceIds
	//if !reflect.DeepEqual(currentInstanceIds, awsManager.Status.VMStartStatus) ||
	//	!reflect.DeepEqual(currentInstanceIds, awsManager.Status.VMStopStatus) {
	//awsManager.Status = "Running"
	//awsManager.Status.VMStopStatus = currentInstanceIds
	/*err := r.Client.Status().Update(ctx, awsManager)
	if err != nil {
		log.Error(err, "Failed to update awsManager status")
		return ctrl.Result{}, err
	}
	*/

	return ctrl.Result{}, nil
}

// Deployment Spec
func (r *AWSManagerReconciler) DeploymentForAWSManager(ctx context.Context, req ctrl.Request, awsManager *awsmanagerv1.AWSManager) *appsv1.Deployment {
	var replicas int32 = 1
	var labels = map[string]string{
		"app": req.NamespacedName.Name,
	}

	//var trueValue = true

	log := ctrllog.FromContext(ctx)
	log.Info("Inside DeploymentForAWSManager")
	//configMapData := make(map[string]string, 0)
	//configMapData["config.json"] = "{}"
	//fmt.Println("Details", awsManager.Name, awsManager.Namespace, awsManager.Spec.Image)
	Deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      awsManager.Name,
			Namespace: awsManager.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					/*Volumes: []corev1.Volume{{
						Name: "config",
						VolumeSource: corev1.VolumeSource{
							ConfigMap: &corev1.ConfigMapVolumeSource{
								LocalObjectReference: corev1.LocalObjectReference{
									Name: "config",
								},
							},
						},
					}},*/
					Containers: []corev1.Container{{
						Name:  awsManager.Name,
						Image: awsManager.Spec.Image,
						/*VolumeMounts: []corev1.VolumeMount{{
							Name:      "config",
							MountPath: "/opt/config",
						}},*/
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
						//RestartPolicy: "Always",
					}}, // Container
				}, // PodSec
			}, // PodTemplateSpec
		}, // Spec
	} // Deployment
	// Set AWSManager instance as the owner and controller
	ctrl.SetControllerReference(awsManager, Deployment, r.Scheme)
	return Deployment
}

// SetupWithManager sets up the controller with the Manager.
func (r *AWSManagerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&awsmanagerv1.AWSManager{}).
		Owns(&appsv1.Deployment{}).
		Complete(r)
}
