
# Hermes Mail Sender Operator

![image](https://github.com/gabrielmuras/hermes-mail-sender-operator-helm-charts/assets/62755656/b8a29704-d9be-41f2-8872-43d8777a0411)

Hermes came back from mount olympus and instead of deliverying letters he is now sending emails through a kubernetes cluster. 

He is a kubernetes operator that was designed to manage and send emails using the mailersend and mailgun providers.

## Architeture and Components

![image](https://github.com/gabrielmuras/hermes-mail-sender-operator-helm-charts/assets/62755656/dd9e8798-7109-4d8e-8f35-e89a5768e190)

The project was developed using the Kubebuilder operator framework https://kubebuilder.io which speeds up the development process in comparisson to create everything from scratch.

The email providers were configured using the go sdk libraries that they provide 
- https://github.com/mailgun/mailgun-go
- https://github.com/mailersend/mailersend-go

Link to the helm charts repository

- https://github.com/gabrielmuras/hermes-mail-sender-operator-helm-charts


## Installation

You can install it by different ways.

### Helm

Add the helm repository

```bash
  helm repo add hermes-mail-sender-operator-charts https://gabrielmuras.github.io/hermes-mail-sender-operator-helm-charts
```
Update

```bash
  helm repo update
```

You can check the available charts

```bash
  helm search repo hermes-mail-sender-operator-charts
```

Install the Operator with the default values or use a custom values that suits your needs

```bash
helm install email-operator hermes-mail-sender-operator-charts/hermes-mail-sender-operator
```

Install the EmailSenderConfig using your values because it doesn't provide a default out of the box default values

```bash
helm install <name> hermes-mail-sender-operator-charts/email-sender-config -f <Your-Custom-Yaml>.yaml
```
Install the EmailSenderConfig using your values because it doesn't provide a default out of the box default values

```bash
helm install <name> hermes-mail-sender-operator-charts/email -f <Your-Custom-Yaml>.yaml
```

### Make

Clone this repository and then execute

```bash
  make; make manifests; make install; make run
```

### Manifests

Use the manifests present in `config/ ` and apply the contents in those folders:

- `config/crds`
- `config/rbac`
- `config/manager`

The `EmailSenderConfig`and `Email`will be on the `config/samples` and also with some examples.

Take note that when using the kind secret for the EmailSenderConfig you should pass it in base64 format before the kubectl apply

```bash
  echo -n yourApiToken | base64
```

## Usage

