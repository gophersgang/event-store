package main

import (
	"fmt"

	"k8s.io/client-go/pkg/api/resource"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/util/intstr"
)

func buildAppContainer(config MicroserviceConfig, e EnvironmentConfig) v1.Container {
	// Secrets
	env := []v1.EnvVar{
		{
			Name: "GKE_PODNAME",
			ValueFrom: &v1.EnvVarSource{
				FieldRef: &v1.ObjectFieldSelector{
					FieldPath: "metadata.name",
				},
			},
		},
		{
			Name: "GKE_NAMESPACE",
			ValueFrom: &v1.EnvVarSource{
				FieldRef: &v1.ObjectFieldSelector{
					FieldPath: "metadata.namespace",
				},
			},
		},
		{
			Name:  "ENVIRONMENT",
			Value: e.Name,
		},
		{
			Name:  "TLS_KEY_FILE",
			Value: "/etc/certs/server-key.pem",
		},
		{
			Name:  "TLS_CERT_FILE",
			Value: "/etc/certs/server-cert.pem",
		},
		{
			Name:  "DEV_CA_FILE",
			Value: "/etc/certs/ca-cert.pem",
		},
		{
			Name: "ROOT_CA_FILE",
			ValueFrom: &v1.EnvVarSource{
				ConfigMapKeyRef: &v1.ConfigMapKeySelector{
					LocalObjectReference: v1.LocalObjectReference{
						Name: fmt.Sprintf("%s-config", config.Name),
					},
					Key: "root_ca_file",
				},
			},
		},
	}

	if config.Apps.Redis != nil {
		env = append(env, v1.EnvVar{Name: "REDIS_HOST", Value: fmt.Sprintf("redis-%s.%s.svc.cluster.local:6379", config.Name, e.K8sNamespace)})
		env = append(env, v1.EnvVar{Name: "REDIS_PASSWORD", Value: config.Apps.Redis.Password})
	}

	for _, customEnv := range e.PodConfig.PodEnv {
		env = append(env, v1.EnvVar{Name: customEnv.Key, Value: customEnv.Value})
	}

	// Volume Mounts
	vm := []v1.VolumeMount{
		{
			Name:      "certs",
			MountPath: "/etc/certs",
		},
	}
	for _, secret := range e.PodConfig.Secrets {
		vm = append(vm, v1.VolumeMount{Name: secret.Name, MountPath: secret.MountPath})
	}

	if e.Name == "local" {
		vm = append(vm,
			v1.VolumeMount{Name: "local-auth", MountPath: "/etc/local-auth/"},
			v1.VolumeMount{Name: "local-app-creds", MountPath: "/etc/local-app-creds/"},
		)
		env = append(env, v1.EnvVar{
			Name:  "GOOGLE_APPLICATION_CREDENTIALS",
			Value: "/etc/local-auth/application_default_credentials.json",
		})
	}

	var probe *v1.Probe
	if e.EndpointsVersion == "" {
		probe = &v1.Probe{
			Handler: v1.Handler{
				HTTPGet: &v1.HTTPGetAction{
					Path: "/healthz",
					Port: intstr.FromInt(11001),
				},
			},
		}
	}

	rvalue := v1.Container{
		Name:            config.Name,
		Image:           fmt.Sprintf("gcr.io/repcore-prod/%s:%s", config.Name, config.Version),
		ImagePullPolicy: v1.PullAlways,
		Ports: []v1.ContainerPort{
			{
				Name:          "http",
				ContainerPort: 11001,
			},
			{
				Name:          "grpc",
				ContainerPort: 11000,
			},
		},
		Env: env,
		Resources: v1.ResourceRequirements{
			Limits: v1.ResourceList{
				v1.ResourceCPU:    resource.MustParse(e.Resources.CpuLimit),
				v1.ResourceMemory: resource.MustParse(e.Resources.MemoryLimit),
			},
			Requests: v1.ResourceList{
				v1.ResourceCPU:    resource.MustParse(e.Resources.CpuRequest),
				v1.ResourceMemory: resource.MustParse(e.Resources.MemoryRequest),
			},
		},
		VolumeMounts: vm,
		LivenessProbe: probe,
		ReadinessProbe: probe,
	}

	if e.Name == "local" {
		rvalue.ImagePullPolicy = v1.PullNever
	}
	return rvalue
}

func buildGoogleAuthContainer(e EnvironmentConfig) v1.Container {
	cpuLimits := resource.MustParse("50m")
	cpuRequests := resource.MustParse("25m")
	memoryLimits := resource.MustParse("32Mi")
	memoryRequests := resource.MustParse("16Mi")
	if e.Name == "prod" {
		cpuLimits = resource.MustParse("100m")
		cpuRequests = resource.MustParse("50m")
		memoryLimits = resource.MustParse("128Mi")
		memoryRequests = resource.MustParse("64Mi")
	}

	rvalue := v1.Container{
		Name:            "auth-proxy",
		Image:           "gcr.io/repcore-prod/google_auth_proxy:v11",
		ImagePullPolicy: v1.PullAlways,
		Ports: []v1.ContainerPort{
			{
				Name:          "https",
				ContainerPort: int32(11002),
			},
		},
		Env: []v1.EnvVar{
			{
				Name:  "REDIRECT_URL",
				Value: fmt.Sprintf("https://%s/oauth2/callback", e.Network.HTTPSHost),
			},
			{
				Name:  "EMAIL_DOMAIN",
				Value: "vendasta.com",
			},
			{
				Name:  "UPSTREAM_URL",
				Value: "http://127.0.0.1:11001",
			},
			{
				Name:  "HTTPS_ADDRESS",
				Value: "0.0.0.0:11002",
			},
			{
				Name:  "SECURE_COOKIE",
				Value: "true",
			},
			{
				Name:  "CLIENT_ID",
				Value: "999898651218.apps.googleusercontent.com",
			},
			{
				Name:  "CLIENT_SECRET",
				Value: "HxfjhCvesf9cmymiUIptoT88",
			},
			{
				Name:  "COOKIE_SECRET",
				Value: "ZDRhMTlkNTczNDYyNDY2ZGJhMDFiYWQ2M2YyM2IxMTYK",
			},
			{
				Name:  "SERVER_CERT",
				Value: "/etc/vendasta-internal/tls.crt",
			},
			{
				Name:  "SERVER_KEY",
				Value: "/etc/vendasta-internal/tls.key",
			},
		},
		Resources: v1.ResourceRequirements{
			Limits: v1.ResourceList{
				v1.ResourceCPU:    cpuLimits,
				v1.ResourceMemory: memoryLimits,
			},
			Requests: v1.ResourceList{
				v1.ResourceCPU:    cpuRequests,
				v1.ResourceMemory: memoryRequests,
			},
		},
		VolumeMounts: []v1.VolumeMount{
			{
				Name:      "vendasta-internal",
				MountPath: "/etc/vendasta-internal",
			},
		},
	}

	if e.Name == "local" {
		rvalue.VolumeMounts = append(rvalue.VolumeMounts,
			v1.VolumeMount{
				Name:      "local-app-creds",
				MountPath: "/etc/local-app-creds/",
			})
	}
	return rvalue
}

func buildEndpointsContainer(e EnvironmentConfig) v1.Container {
	fmt.Printf("Environment:\n----\n%#v\n----\n", e)
	cpuLimits := resource.MustParse("50m")
	cpuRequests := resource.MustParse("25m")
	memoryLimits := resource.MustParse("32Mi")
	memoryRequests := resource.MustParse("16Mi")
	if e.Name == "prod" {
		cpuLimits = resource.MustParse("1000m")
		cpuRequests = resource.MustParse("500m")
		memoryLimits = resource.MustParse("128Mi")
		memoryRequests = resource.MustParse("64Mi")
	}

	env := []v1.EnvVar{}
	vm := []v1.VolumeMount{

	}

	// Configure Volume Mounts
	if e.Name == "local" {
		vm = append(vm,
			v1.VolumeMount{Name: "local-auth", MountPath: "/etc/local-auth/"},
			v1.VolumeMount{Name: "local-app-creds", MountPath: "/etc/local-app-creds/"},
			v1.VolumeMount{Name: "local-app-creds", MountPath: "/etc/nginx/ssl"},
		)
		env = append(env, v1.EnvVar{
			Name:  "GOOGLE_APPLICATION_CREDENTIALS",
			Value: "/etc/local-auth/application_default_credentials.json",
		})
	} else {
		vm = append(vm, v1.VolumeMount{
			Name:      "vendasta-internal",
			MountPath: "/etc/nginx/ssl",
			ReadOnly:  true,
		})
	}

	// Endpoints Container Args
	args := []string{
		"-s", e.GRPCHost,
		"-v", e.EndpointsVersion,
		"-a", "grpc://127.0.0.1:11000",
		"-p", "11003",
		"-P", "11005",
		"-S", "11006",
		"-z", "healthz",
	}
	if e.Name == "local" {
		args = append(args, "-k", "/etc/local-app-creds/key.json")
	}

	var probe *v1.Probe
	if e.EndpointsVersion == "" {
		probe = &v1.Probe{
			Handler: v1.Handler{
				HTTPGet: &v1.HTTPGetAction{
					Path: "/healthz",
					Port: intstr.FromInt(11001),
				},
			},
		}
	} else {
		probe = &v1.Probe{
			Handler: v1.Handler{
				HTTPGet: &v1.HTTPGetAction{
					Path: "/healthz",
					Port: intstr.FromInt(11003),
				},
			},
		}
	}

	rvalue :=
		v1.Container{
			Name:  "endpoints",
			Image: "b.gcr.io/endpoints/endpoints-runtime:1",
			Args: args,
			ImagePullPolicy: v1.PullAlways,
			LivenessProbe: probe,
			ReadinessProbe: probe,
			Ports: []v1.ContainerPort{
				{
					ContainerPort: 11003, //< Serve GRPC + REST
				},
				{
					ContainerPort: 10005, //< Serve Debug info
				},
				{
					ContainerPort: 11006, //< Optionally terminate SSL for GRPC + REST
				},
			},
			Resources: v1.ResourceRequirements{
				Limits: v1.ResourceList{
					v1.ResourceCPU:    cpuLimits,
					v1.ResourceMemory: memoryLimits,
				},
				Requests: v1.ResourceList{
					v1.ResourceCPU:    cpuRequests,
					v1.ResourceMemory: memoryRequests,
				},
			},
			VolumeMounts: vm,
			Env: env,
		}
	return rvalue
}

func buildLocalProxyContainer() v1.Container {

	cpuLimits := resource.MustParse("50m")
	cpuRequests := resource.MustParse("25m")
	memoryLimits := resource.MustParse("32Mi")
	memoryRequests := resource.MustParse("16Mi")

	env := []v1.EnvVar{}
	vm := []v1.VolumeMount{
		{Name: "vendasta-local-secret", MountPath: "/etc/local-proxy", ReadOnly: true},
	}
	rvalue :=
		v1.Container{
			Name:  "local-proxy",
			Image: "vendasta/local-proxy",
			ImagePullPolicy: v1.PullAlways,
			Ports: []v1.ContainerPort{
				{
					ContainerPort: 443,
				},
			},
			Resources: v1.ResourceRequirements{
				Limits: v1.ResourceList{
					v1.ResourceCPU:    cpuLimits,
					v1.ResourceMemory: memoryLimits,
				},
				Requests: v1.ResourceList{
					v1.ResourceCPU:    cpuRequests,
					v1.ResourceMemory: memoryRequests,
				},
			},
			VolumeMounts: vm,
			Env: env,
		}
	return rvalue
}

func buildRedisContainer(msConfig MicroserviceConfig) v1.Container {

	cpuLimits := resource.MustParse("50m")
	cpuRequests := resource.MustParse("25m")
	memoryLimits := resource.MustParse("32Mi")
	memoryRequests := resource.MustParse("16Mi")
	maxMemory := "16mb"

	if msConfig.GetEnvironment().Name == "prod" {
		memoryLimits = resource.MustParse("256Mi")
		memoryRequests = resource.MustParse("128Mi")
		maxMemory = "200mb"
	}

	env := []v1.EnvVar{}
	vm := []v1.VolumeMount{}
	rvalue :=
		v1.Container{
			Name:  "redis",
			Image: "redis",
			Args: []string{
				"redis-server", fmt.Sprintf("--requirepass %s", msConfig.Apps.Redis.Password), fmt.Sprintf("--maxmemory %s", maxMemory),
			},
			ImagePullPolicy: v1.PullAlways,
			Ports: []v1.ContainerPort{
				{
					ContainerPort: 6379,
				},
			},
			Resources: v1.ResourceRequirements{
				Limits: v1.ResourceList{
					v1.ResourceCPU:    cpuLimits,
					v1.ResourceMemory: memoryLimits,
				},
				Requests: v1.ResourceList{
					v1.ResourceCPU:    cpuRequests,
					v1.ResourceMemory: memoryRequests,
				},
			},
			VolumeMounts: vm,
			Env: env,
		}
	return rvalue
}
