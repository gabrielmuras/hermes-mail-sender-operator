package secrets

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func GetDecodeApiKey(r client.Client, ctx context.Context, req ctrl.Request, emailSenderName string) (string, error) {
	var apiTokenDecode string
	log := log.FromContext(ctx)

	// Fetch the referenced secret
	secret := &corev1.Secret{}
	secretName := types.NamespacedName{
		Name:      emailSenderName,
		Namespace: req.Namespace,
	}

	if err := r.Get(ctx, secretName, secret); err != nil {
		log.Error(err, "Unable to get secret", "SecretName", secretName)
		return "", err
	}

	// Retrieve the API token from the secret
	apiToken, exists := secret.Data["apiToken"]
	if !exists {
		log.Error(nil, "secret does not contain key 'apiToken'")
		return "", nil
	} else {
		apiTokenDecode = string(apiToken)
	}

	return apiTokenDecode, nil

}
