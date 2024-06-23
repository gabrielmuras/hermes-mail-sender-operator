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
	mailv1 "hermes-mail-sender-operator/api/v1"
	secrets "hermes-mail-sender-operator/internal/k8s"
	providers "hermes-mail-sender-operator/internal/providers"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// EmailReconciler reconciles a Email object
type EmailReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=email.hermes.sender,resources=emails,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=email.hermes.sender,resources=emails/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=email.hermes.sender,resources=emails/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Email object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.18.2/pkg/reconcile
func (r *EmailReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	log.Info("Reconciling Email")

	var email mailv1.Email
	var emailSenderConfig mailv1.EmailSenderConfig

	if err := r.Get(ctx, req.NamespacedName, &email); err != nil {
		log.Error(err, "Unable to fetch Email")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if email.Status.DeliveryStatus == "" {
		log.Info("Email status is empty")

		key := client.ObjectKey{

			Namespace: req.Namespace,
			Name:      email.Spec.SenderConfigRef,
		}

		if err := r.Get(ctx, key, &emailSenderConfig); err != nil {

			log.Error(err, "Unable to fetch EmailSenderConfig. "+email.Spec.SenderConfigRef)
			return ctrl.Result{}, nil

		} else {

			status := emailSenderConfig.Status.Status

			if status == "Error" || status == "Unknown Provider" {
				log.Error(err, "EmailSenderConfig is in error state and cannot be used")
				email.Status.DeliveryStatus = "EmailSenderConfigError"
				return ctrl.Result{}, nil
			}

			//provider := emailSenderConfig.Spec.Provider
			log.Info("Able to fetch EmailSenderConfig " + emailSenderConfig.Spec.Provider)

		}

		apiTokenDecode, err := secrets.GetDecodeApiKey(r.Client, ctx, req, emailSenderConfig.Spec.ApiTokenSecretRef)

		if err != nil {
			log.Error(err, "Unable to decode secret")
		}

		EmailConfig := providers.EmailConfig{

			ApiToken:       apiTokenDecode,
			Subject:        email.Spec.Subject,
			Text:           email.Spec.Body,
			FromEmail:      emailSenderConfig.Spec.SenderEmail,
			RecipientEmail: email.Spec.RecipientEmail,
		}

		switch provider := emailSenderConfig.Spec.Provider; provider {

		case "mailersender":
			if messageID, err := providers.SendEmailMailerSender(EmailConfig); err != nil {

				log.Error(err, "Error sending email")
				email.Status.DeliveryStatus = "Error"
				email.Status.Error = err.Error()

			} else {
				log.Info("Email sent successfully")
				email.Status.DeliveryStatus = "Sent"
				email.Status.MessageId = *messageID
			}

		case "mailgun":
			if _, messageID, err := providers.SendEmailMailgun(EmailConfig); err != nil {

				log.Error(err, "Error sending email")
				email.Status.DeliveryStatus = "Error"
				email.Status.Error = err.Error()

			} else {

				log.Info("Email sent successfully")
				email.Status.DeliveryStatus = "Sent"
				email.Status.MessageId = *messageID

			}

		}
	}

	if err := r.Status().Update(ctx, &email); err != nil {
		log.Error(err, "Unable to update Email status")
		return ctrl.Result{}, err

	} else {
		log.Info("Email status updated successfully")

	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *EmailReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		// Uncomment the following line adding a pointer to an instance of the controlled resource as an argument
		// For().
		For(&mailv1.Email{}).
		Complete(r)
}
