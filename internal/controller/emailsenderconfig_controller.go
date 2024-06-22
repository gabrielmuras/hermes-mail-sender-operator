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

	//b64 "encoding/base64"
	"fmt"
	mailv1 "hermes-mail-sender-operator/api/v1"
	providers "hermes-mail-sender-operator/internal/providers"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// EmailSenderConfigReconciler reconciles a EmailSenderConfig object
type EmailSenderConfigReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=email.hermes.sender,resources=emailsenderconfigs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=email.hermes.sender,resources=emailsenderconfigs/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=email.hermes.sender,resources=emailsenderconfigs/finalizers,verbs=update
// +kubebuilder:rbac:groups=email.hermes.sender,resources=secrets,verbs=get;list;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the EmailSenderConfig object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.18.2/pkg/reconcile
func (r *EmailSenderConfigReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	log.Info("Reconciling EmailSenderConfig")

	var emailSenderConfig mailv1.EmailSenderConfig

	if err := r.Get(ctx, req.NamespacedName, &emailSenderConfig); err != nil {
		log.Error(err, "Unable to fetch EmailSenderConfig")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Fetch the referenced secret
	var apiTokenDecode string
	secret := &corev1.Secret{}
	secretName := types.NamespacedName{
		Name:      emailSenderConfig.Spec.ApiTokenSecretRef,
		Namespace: req.Namespace,
	}

	if err := r.Get(ctx, secretName, secret); err != nil {
		log.Error(err, "Failed to get Secret", "SecretName", secretName)
		return ctrl.Result{}, err
	}

	// Retrieve the API token from the secret
	apiToken, exists := secret.Data["apiToken"]
	if !exists {
		log.Error(nil, "Secret does not contain key 'apiToken'")

	} else {
		apiTokenDecode = string(apiToken)
		fmt.Println("API Token: ", apiTokenDecode)
	}

	//apiToken, subject, text, fromEmail, recipientEmail
	if emailSenderConfig.Status.Status == "" || emailSenderConfig.Status.Status == "Error" || emailSenderConfig.Status.Status == "Unknown Provider" {

		dummyValidateConfig := providers.EmailConfig{
			ApiToken:       apiTokenDecode,
			Subject:        "test",
			Text:           "test",
			FromEmail:      emailSenderConfig.Spec.SenderEmail,
			RecipientEmail: emailSenderConfig.Spec.SenderEmail,
		}

		switch provider := emailSenderConfig.Spec.Provider; provider {

		case "mailersender":
			if _, err := providers.SendEmailMailerSender(dummyValidateConfig); err != nil {

				log.Error(err, "Unable to verify emailSenderConfig")
				emailSenderConfig.Status.Status = "Error"

			} else {
				log.Info("EmailSenderConfig verified successfully")
				emailSenderConfig.Status.Status = "Ok"
			}

		case "mailgun":
			//paid email verification not using providers.validateDomainMailGun function
			emailSenderConfig.Status.Status = "Ok"

		default:
			log.Error(nil, "Invalid provider. Please use mailersender or mailgun.")
			emailSenderConfig.Status.Status = "Unknown Provider"

		}
		if err := r.Status().Update(ctx, &emailSenderConfig); err != nil {
			log.Error(err, "Unable to create EmailSenderConfig status")
			return ctrl.Result{}, err
		} else {
			log.Info("EmailSenderConfig status created successfully")

		}
	}

	if err := r.Status().Update(ctx, &emailSenderConfig); err != nil {
		log.Error(err, "unable to update EmailSenderConfig status")
		return ctrl.Result{}, err
	} else {
		log.Info("EmailSenderConfig status updated successfully")
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *EmailSenderConfigReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		// Uncomment the following line adding a pointer to an instance of the controlled resource as an argument
		// For().
		For(&mailv1.EmailSenderConfig{}).
		Complete(r)
}
