package main

import (
	"fmt"
	"log"
	"os"
	"os/user"

	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/apis/extensions/v1beta1"
	"k8s.io/client-go/pkg/util/intstr"
	"k8s.io/client-go/pkg/util/json"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	autoScalingV1 "k8s.io/client-go/pkg/apis/autoscaling/v1"
	"k8s.io/client-go/pkg/api/errors"
)

func getKubeConfigFile() string {
	var u *user.User
	var err error

	if u, err = user.Current(); err != nil {
		log.Fatalf("Error getting current user: %s", err.Error())
	}

	kubeFile := fmt.Sprintf("%s/.kube/config", u.HomeDir)
	if _, err = os.Stat(kubeFile); err != nil {
		log.Fatalf("Could not find kubernetes config at %s", kubeFile)
	}
	return kubeFile
}

func GetK8sClientSet(config MicroserviceConfig) *kubernetes.Clientset {
	clusterConfig, err := rest.InClusterConfig()
	if err != nil {
		envConfig := config.GetEnvironment()
		cfg := clientcmd.GetConfigFromFileOrDie(getKubeConfigFile())
		cfg2 := clientcmd.NewNonInteractiveClientConfig(*cfg, envConfig.K8sContext, &clientcmd.ConfigOverrides{}, nil)
		cfg3, err := cfg2.ClientConfig()
		if err != nil {
			log.Fatalf("Error creating config: %s", err.Error())
		}
		// creates the clientset
		var clientset *kubernetes.Clientset
		//var err error
		if clientset, err = kubernetes.NewForConfig(cfg3); err != nil {
			log.Fatalf("Error generating kubernetes clientset: %s", err.Error())
		}

		return clientset
	}
	clientset, err := kubernetes.NewForConfig(clusterConfig)
	if err != nil {
		panic(err)
	}
	return clientset
}

func CreateConfigMap(config MicroserviceConfig, clientset *kubernetes.Clientset) {
	log.Printf("Creating ConfigMap\n")
	envConfig := config.GetEnvironment()

	cmName := fmt.Sprintf("%s-config", config.Name)
	cmExists := false

	var err error
	var list *v1.ConfigMapList
	if list, err = clientset.Core().ConfigMaps(envConfig.K8sNamespace).List(v1.ListOptions{}); err != nil {
		log.Fatalf("Error getting configmap list: %s", err.Error())
	}
	log.Printf("%d ConfigMaps:\n", len(list.Items))
	for i := range list.Items {
		name := list.Items[i].ObjectMeta.Name
		log.Printf(" - %s\n", name)
		if name == cmName {
			cmExists = true
		}
	}

	if !cmExists {
		log.Printf("Creating ConfigMap...\n")
		//Create a new ConfigMap
		cm := v1.ConfigMap{
			ObjectMeta: v1.ObjectMeta{
				Name:      cmName,
				Namespace: envConfig.K8sNamespace,
			},
			Data: map[string]string{
				"environment":  envConfig.Name,
				"https-url":    fmt.Sprintf("https://%s:%d", envConfig.Network.HTTPSHost, envConfig.Network.HTTPSPort),
				"grpc-host":    fmt.Sprintf("%s:%d", envConfig.Network.GRPCHost, envConfig.Network.GRPCPort),
				"root_ca_file": "",
			},
		}
		if envConfig.Name == "local" {
			cm.Data["root_ca_file"] = "yes"
		}
		if _, err = clientset.Core().ConfigMaps(envConfig.K8sNamespace).Create(&cm); err != nil {
			log.Fatalf("Error creating configMap: %s", err.Error())
		}
	} else {
		log.Printf("ConfigMap %s exists, skipping creation...\n", cmName)
	}
}

func CreateDeployment(config MicroserviceConfig, clientset *kubernetes.Clientset) {
	log.Printf("Creating Deployment\n")
	envConfig := config.GetEnvironment()

	depName := config.Name
	depExists := false

	var err error
	var list *v1beta1.DeploymentList
	if list, err = clientset.Extensions().Deployments(envConfig.K8sNamespace).List(v1.ListOptions{}); err != nil {
		log.Fatalf("Error getting deployments list: %s", err.Error())
	}
	log.Printf("%d Deployments:\n", len(list.Items))
	for i := range list.Items {
		name := list.Items[i].ObjectMeta.Name
		log.Printf(" - %s\n", name)
		if name == depName {
			depExists = true
		}
	}
	replicas := int32(1)
	revisionHistoryLimit := int32(5)
	maxUnavailable := intstr.FromString("25%")
	maxSurge := intstr.FromString("25%")
	containers := []v1.Container{
		buildAppContainer(config, envConfig),
		buildGoogleAuthContainer(envConfig),
	}
	if envConfig.AppConfig.EndpointsVersion != "" {
		containers = append(containers, buildEndpointsContainer(envConfig))
	}
	volumes := []v1.Volume{
		{
			Name: "certs",
			VolumeSource: v1.VolumeSource{
				Secret: &v1.SecretVolumeSource{
					SecretName: fmt.Sprintf("%s-secret", config.Name),
				},
			},
		},
		{
			Name: "vendasta-internal",
			VolumeSource: v1.VolumeSource{
				Secret: &v1.SecretVolumeSource{
					SecretName: "vendasta-internal-secret",
				},
			},
		},
	}
	if envConfig.Name == "local" {
		volumes = append(volumes,
			v1.Volume{
				Name: "local-auth",
				VolumeSource: v1.VolumeSource{
					HostPath: &v1.HostPathVolumeSource{
						Path: wellKnownFile(),
					},
				},
			},
			v1.Volume{
				Name: "local-app-creds",
				VolumeSource: v1.VolumeSource{
					Secret: &v1.SecretVolumeSource{
						SecretName: "vendasta-local-secret",
					},
				},
			},
		)
	}
	for _, secret := range envConfig.PodConfig.Secrets {
		volumes = append(volumes, v1.Volume{
			Name: secret.Name,
			VolumeSource: v1.VolumeSource{
				Secret: &v1.SecretVolumeSource{
					SecretName: secret.Name,
				},
			},
		})
	}
	dep := v1beta1.Deployment{
		ObjectMeta: v1.ObjectMeta{
			Name:      depName,
			Namespace: envConfig.K8sNamespace,
		},
		Spec: v1beta1.DeploymentSpec{
			Replicas:             &replicas,
			RevisionHistoryLimit: &revisionHistoryLimit,
			Strategy: v1beta1.DeploymentStrategy{
				Type: v1beta1.RollingUpdateDeploymentStrategyType,
				RollingUpdate: &v1beta1.RollingUpdateDeployment{
					MaxUnavailable: &maxUnavailable,
					MaxSurge:       &maxSurge,
				},
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: v1.ObjectMeta{
					Labels: map[string]string{
						"app":         config.Name,
						"environment": config.Environment,
					},
				},
				Spec: v1.PodSpec{
					Containers: containers,
					Volumes:    volumes,
				},
			},
		},
	}
	if envConfig.Name == "local" {
		dep.Spec.Template.Spec.ImagePullSecrets = append(dep.Spec.Template.Spec.ImagePullSecrets, v1.LocalObjectReference{"vendasta-local-gcr"})
	}
	o, _ := json.Marshal(dep)
	log.Printf("Deployment %s\n", string(o))
	if !depExists {
		log.Printf("Creating Deployment...\n")
		if _, err = clientset.Extensions().Deployments(envConfig.K8sNamespace).Create(&dep); err != nil {
			log.Fatalf("Error creating deployment: %s", err.Error())
		}
	} else {
		log.Printf("Deployment %s exists, updating...\n", depName)
		if _, err = clientset.Extensions().Deployments(envConfig.K8sNamespace).Update(&dep); err != nil {
			log.Fatalf("Error updating deployment: %s", err.Error())
		}
	}
}

func CreateLocalProxyDeployment(clientset *kubernetes.Clientset) {
	log.Printf("Creating Local Proxy Deployment\n")
	depExists := true

	if _, err := clientset.Extensions().Deployments("default").Get("local-proxy"); err != nil {
		if !errors.IsNotFound(err) {
			log.Fatalf("Error getting local proxy: %s", err.Error())
		}
		depExists = false
	}
	replicas := int32(1)
	revisionHistoryLimit := int32(5)
	maxUnavailable := intstr.FromString("25%")
	maxSurge := intstr.FromString("25%")
	containers := []v1.Container{
		buildLocalProxyContainer(),
	}
	volumes := []v1.Volume{
		{
			Name: "vendasta-local-secret",
			VolumeSource: v1.VolumeSource{
				Secret: &v1.SecretVolumeSource{
					SecretName: "vendasta-local-secret",
				},
			},
		},
	}
	dep := v1beta1.Deployment{
		ObjectMeta: v1.ObjectMeta{
			Name:      "local-proxy",
			Namespace: "default",
		},
		Spec: v1beta1.DeploymentSpec{
			Replicas:             &replicas,
			RevisionHistoryLimit: &revisionHistoryLimit,
			Strategy: v1beta1.DeploymentStrategy{
				Type: v1beta1.RollingUpdateDeploymentStrategyType,
				RollingUpdate: &v1beta1.RollingUpdateDeployment{
					MaxUnavailable: &maxUnavailable,
					MaxSurge:       &maxSurge,
				},
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: v1.ObjectMeta{
					Labels: map[string]string{
						"app": "local-proxy",
					},
				},
				Spec: v1.PodSpec{
					Containers: containers,
					Volumes:    volumes,
				},
			},
		},
	}
	if !depExists {
		log.Printf("Creating Local Proxy Deployment...\n")
		if _, err := clientset.Extensions().Deployments("default").Create(&dep); err != nil {
			log.Fatalf("Error creating deployment: %s", err.Error())
		}
	} else {
		log.Printf("Local proxy Deployment exists, updating...\n")
		if _, err := clientset.Extensions().Deployments("default").Update(&dep); err != nil {
			log.Fatalf("Error updating deployment: %s", err.Error())
		}
	}
}

func CreateAppsDeployment(config MicroserviceConfig, clientset *kubernetes.Clientset) {
	log.Printf("Creating Apps Deployment\n")
	if config.Apps.Redis != nil {
		createRedisCacheDeployment(config, clientset)
		createRedisCacheService(config, clientset)
	}
}

func createRedisCacheDeployment(config MicroserviceConfig, clientset *kubernetes.Clientset) {
	log.Printf("Creating Redis Deployment\n")
	if config.Apps.Redis.Password == "" {
		log.Fatal("Redis password must be supplied.")
	}
	depExists := true

	envConfig := config.GetEnvironment()
	depName := fmt.Sprintf("redis-%s", config.Name)

	if _, err := clientset.Extensions().Deployments(envConfig.K8sNamespace).Get(depName); err != nil {
		if !errors.IsNotFound(err) {
			log.Fatalf("Error getting redis deployment: %s", err.Error())
		}
		depExists = false
	}
	replicas := int32(1)
	revisionHistoryLimit := int32(5)
	containers := []v1.Container{
		buildRedisContainer(config),
	}
	dep := v1beta1.Deployment{
		ObjectMeta: v1.ObjectMeta{
			Name:      depName,
			Namespace: envConfig.K8sNamespace,
		},
		Spec: v1beta1.DeploymentSpec{
			Replicas:             &replicas,
			RevisionHistoryLimit: &revisionHistoryLimit,
			Template: v1.PodTemplateSpec{
				ObjectMeta: v1.ObjectMeta{
					Labels: map[string]string{
						"app": depName,
					},
				},
				Spec: v1.PodSpec{
					Containers: containers,
				},
			},
		},
	}

	if !depExists {
		log.Println("Creating redis cache Deployment...")
		if _, err := clientset.Extensions().Deployments(envConfig.K8sNamespace).Create(&dep); err != nil {
			log.Fatalf("Error creating deployment: %s", err.Error())
		}
	} else {
		log.Println("Redis cache Deployment exists, updating...")
		if _, err := clientset.Extensions().Deployments(envConfig.K8sNamespace).Update(&dep); err != nil {
			log.Fatalf("Error updating deployment: %s", err.Error())
		}
	}
}

func guessUnixHomeDir() string {
	// Prefer $HOME over user.Current due to glibc bug: golang.org/issue/13470
	if v := os.Getenv("HOME"); v != "" {
		return v
	}
	// Else, fall back to user.Current:
	if u, err := user.Current(); err == nil {
		return u.HomeDir
	}
	return ""
}

func wellKnownFile() string {
	return filepath.Join(guessUnixHomeDir(), ".config", "gcloud")
}

func CreateSecret(config MicroserviceConfig, clientset *kubernetes.Clientset) {
	log.Printf("Creating Secret, shhhhhhh\n")
	envConfig := config.GetEnvironment()

	secretName := fmt.Sprintf("%s-secret", config.Name)
	if !DoesSecretExist(envConfig, secretName, clientset) {
		log.Printf("Generating certs...\n")
		var certs *TlsPems
		if config.Environment == "local" {
			certs = GenerateLocalCerts()
		} else {
			certs = &TlsPems{}
		}

		log.Printf("Creating Secret...\n")
		secret := v1.Secret{
			ObjectMeta: v1.ObjectMeta{
				Name:      secretName,
				Namespace: envConfig.K8sNamespace,
			},
			Data: map[string][]byte{
				"server-key.pem":  certs.ServerKeyPem,
				"server-cert.pem": certs.ServerCertPem,
				"ca-cert.pem":     certs.CaCertPem,
			},
		}
		log.Printf("Secret: %+v\n", secret)
		if _, err := clientset.Core().Secrets(envConfig.K8sNamespace).Create(&secret); err != nil {
			log.Fatalf("Error creating secret: %s", err.Error())
		}
	} else {
		log.Printf("Secret %s exists, skipping creation...\n", secretName)
	}

	createVendastaInternalSecret(envConfig, clientset)
	if envConfig.Name == "local" {
		createVendastaLocalSecret(envConfig, clientset)
		createLocalPullSecret(envConfig, clientset)
	}
}

func createVendastaInternalSecret(envConfig EnvironmentConfig, clientset *kubernetes.Clientset) {
	log.Printf("Creating Vendasta internal secret...\n")

	if !DoesSecretExist(envConfig, "vendasta-internal-secret", clientset) {
		secret := v1.Secret{
			ObjectMeta: v1.ObjectMeta{
				Name:      "vendasta-internal-secret",
				Namespace: envConfig.K8sNamespace,
			},
			Data: map[string][]byte{
				// Kubernetes Ingress format
				"tls.crt": []byte(VENDASTA_INTERNAL_CERT),
				"tls.key": []byte(VENDASTA_INTERNAL_KEY),
				// NGINX format
				"nginx.crt": []byte(VENDASTA_INTERNAL_CERT),
				"nginx.key": []byte(VENDASTA_INTERNAL_KEY),
			},
		}

		_, err := clientset.Core().Secrets(envConfig.K8sNamespace).Create(&secret)
		if err != nil {
			log.Fatalf("Error creating vendasta internal secret: %s\n", err.Error())
		}
	}
}

type dockerConfig struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	//Auth     string `json:"auth"`
}

func createVendastaLocalSecret(envConfig EnvironmentConfig, clientset *kubernetes.Clientset) {
	log.Printf("Creating Vendasta Local secret...\n")

	//auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", "_json_key", VENDASTA_LOCAL_JSON_KEY)))
	dc := dockerConfig{
		Username: "_json_key",
		Password: VENDASTA_LOCAL_JSON_KEY,
		Email: "123@3456.com",
		//Auth: auth,
	}
	dockerConfigBytes, err := json.Marshal(map[string]interface{}{"https://gcr.io": dc}); if err != nil {
		log.Fatalf("Error creating json key for docker gcr.io auth %s", err.Error())
	}

	if !DoesSecretExist(envConfig, "vendasta-local-secret", clientset) {
		secret := v1.Secret{
			ObjectMeta: v1.ObjectMeta{
				Name:      "vendasta-local-secret",
				Namespace: envConfig.K8sNamespace,
			},
			Data: map[string][]byte{
				// vendasta local service account key
				"key.json": []byte(VENDASTA_LOCAL_JSON_KEY),

				// Kubernetes Ingress format
				"tls.crt": []byte(VENDASTA_LOCAL_CERT),
				"tls.key": []byte(VENDASTA_LOCAL_KEY),

				// NGINX format
				"nginx.crt": []byte(VENDASTA_LOCAL_CERT),
				"nginx.key": []byte(VENDASTA_LOCAL_KEY),
			},
		}
		_, err := clientset.Core().Secrets(envConfig.K8sNamespace).Create(&secret)
		if err != nil {
			log.Fatalf("Error creating vendasta internal secret: %s\n", err.Error())
		}
	}
	if !DoesSecretExist(envConfig, "vendasta-local-gcr", clientset) {
		secret := v1.Secret{
			ObjectMeta: v1.ObjectMeta{
				Name:      "vendasta-local-gcr",
				Namespace: envConfig.K8sNamespace,
			},
			Data: map[string][]byte{
				// GCR.IO Pull Secrets
				".dockercfg": dockerConfigBytes,
			},
			Type: v1.SecretTypeDockercfg,
		}
		_, err := clientset.Core().Secrets(envConfig.K8sNamespace).Create(&secret)
		if err != nil {
			log.Fatalf("Error creating vendasta internal secret: %s\n", err.Error())
		}
	}
}

func createLocalPullSecret(envConfig EnvironmentConfig, clientset *kubernetes.Clientset) {
	log.Printf("Creating Vendasta Local Pull Secret...\n")
	sai := clientset.ServiceAccounts(envConfig.K8sNamespace)
	defaultSA, err := sai.Get("default"); if err != nil {
		log.Fatalf("Error getting default service account %s", err.Error())
	}
	if len(defaultSA.ImagePullSecrets) >= 1 {
		return
	}
	defaultSA.ImagePullSecrets = append(defaultSA.ImagePullSecrets,
		v1.LocalObjectReference{Name: "vendasta-local-gcr"},
	)
	_, err = sai.Update(defaultSA); if err != nil {
		log.Fatalf("Error updating default service account. %s", err.Error())
	}
}

func DoesSecretExist(envConfig EnvironmentConfig, secretName string, clientset *kubernetes.Clientset) bool {
	secretExists := false

	var err error
	var list *v1.SecretList
	if list, err = clientset.Core().Secrets(envConfig.K8sNamespace).List(v1.ListOptions{}); err != nil {
		log.Fatalf("Error getting secrets list: %s", err.Error())
	}
	log.Printf("%d Secrets:\n", len(list.Items))
	for i := range list.Items {
		name := list.Items[i].ObjectMeta.Name
		log.Printf(" - %s\n", name)
		if name == secretName {
			secretExists = true
		}
	}
	return secretExists
}

func CreateHorizontalPodAutoscaler(msConfig MicroserviceConfig, clientset *kubernetes.Clientset) {
	log.Println("Creating Horizontal Pod Autoscaler...")
	envConfig := msConfig.GetEnvironment()
	hpa := clientset.HorizontalPodAutoscalers(envConfig.K8sNamespace)

	hpaExists := true

	_, err := hpa.Get(msConfig.Name); if err != nil {
		if !errors.IsNotFound(err) {
			log.Fatalf("Error getting HPA %s", err.Error())
		}
		hpaExists = false
	}
	if envConfig.Scaling.MinReplicas <= 0 {
		envConfig.Scaling.MinReplicas = 1
	}
	if envConfig.Scaling.MaxReplicas <= 0 {
		envConfig.Scaling.MaxReplicas = 3
	}
	if envConfig.Scaling.TargetCPU <= 0 {
		envConfig.Scaling.TargetCPU = 50
	}

	hpaSpec := &autoScalingV1.HorizontalPodAutoscaler{
		ObjectMeta: v1.ObjectMeta{
			Name:      msConfig.Name,
			Namespace: envConfig.K8sNamespace,
		},
		Spec: autoScalingV1.HorizontalPodAutoscalerSpec{
			ScaleTargetRef: autoScalingV1.CrossVersionObjectReference{
				Kind: "Deployment",
				Name: msConfig.Name,
			},
			MinReplicas: &envConfig.Scaling.MinReplicas,
			MaxReplicas: envConfig.Scaling.MaxReplicas,
			TargetCPUUtilizationPercentage: &envConfig.Scaling.TargetCPU,
		},
	}
	if !hpaExists {
		_, err = hpa.Create(hpaSpec)
	} else {
		_, err = hpa.Update(hpaSpec)
	}
	log.Printf("HPA %s", hpaSpec)
	if err != nil {
		log.Fatalf("Unable to create/update the HPA. %s", err.Error())
	}
}
