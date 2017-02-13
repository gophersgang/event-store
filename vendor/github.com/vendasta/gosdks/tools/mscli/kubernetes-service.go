package main

import (
	"fmt"
	"log"

	"encoding/json"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/errors"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/util/intstr"
)

func CreateServices(config MicroserviceConfig, clientset *kubernetes.Clientset) {
	fmt.Printf("Creating Services\n")

	createGrpcService(config, clientset)
	createHttpService(config, clientset)
}

func createGrpcService(config MicroserviceConfig, clientset *kubernetes.Clientset) {
	envConfig := config.GetEnvironment()
	grpcServiceName := fmt.Sprintf("%s-grpc-svc", config.Name)

	svcs := clientset.Core().Services(envConfig.K8sNamespace)
	exists := true
	_, err := svcs.Get(grpcServiceName)
	if err != nil {
		if !errors.IsNotFound(err) {
			log.Fatalf("Error getting service list: %s", err.Error())
		}
		exists = false
	}
	svc := &v1.Service{
		ObjectMeta: v1.ObjectMeta{
			Name:      grpcServiceName,
			Namespace: envConfig.K8sNamespace,
		},
		Spec: v1.ServiceSpec{
			Ports: []v1.ServicePort{
				{
					Name:       "grpc",
					Port:       int32(443),
					TargetPort: intstr.FromInt(11006),
					Protocol:   v1.ProtocolTCP,
				},
			},
			Type:           v1.ServiceTypeLoadBalancer,
			LoadBalancerIP: envConfig.Network.LoadBalancerIP,
			Selector: map[string]string{
				"app":         config.Name,
				"environment": config.Environment,
			},
		},
	}
	o, _ := json.MarshalIndent(svc, "", "  ")
	fmt.Printf("gRPC Service: %s", o)

	if exists {
		log.Printf("gRPC Service already exists...")
		return
	}

	if config.Environment == "local" {
		svc.Spec.LoadBalancerIP = ""
		svc.Spec.Ports[0].Port = int32(31957)
		svc.Spec.Ports[0].TargetPort = intstr.FromInt(11003)
		svc.ObjectMeta.Annotations = map[string]string{
			"vendasta-local.com/domain": envConfig.GRPCHost,
			"vendasta-local.com/port":   "31957",
		}
	}
	if _, err := svcs.Create(svc); err != nil {
		log.Fatalf("Error creating svc %s: %s", grpcServiceName, err.Error())
	}
}

func createHttpService(config MicroserviceConfig, clientset *kubernetes.Clientset) {
	envConfig := config.GetEnvironment()
	httpsServiceName := fmt.Sprintf("%s-https-svc", config.Name)

	svcs := clientset.Core().Services(envConfig.K8sNamespace)
	exists := true
	_, err := svcs.Get(httpsServiceName)
	if err != nil {
		if !errors.IsNotFound(err) {
			log.Fatalf("Error getting service list: %s", err.Error())
		}
		exists = false
	}
	svc := &v1.Service{
		ObjectMeta: v1.ObjectMeta{
			Name:      httpsServiceName,
			Namespace: envConfig.K8sNamespace,
			Annotations: map[string]string{
				"prometheus.io/scrape": "true",
			},
		},
		Spec: v1.ServiceSpec{
			Ports: []v1.ServicePort{
				{
					Name:       "https",
					Port:       int32(443),
					TargetPort: intstr.FromInt(11002),
					Protocol:   v1.ProtocolTCP,
				},
			},
			Type:           v1.ServiceTypeLoadBalancer,
			LoadBalancerIP: envConfig.Network.LoadBalancerIP,
			Selector: map[string]string{
				"app":         config.Name,
				"environment": config.Environment,
			},
		},
	}
	o, _ := json.MarshalIndent(svc, "", "  ")
	fmt.Printf("HTTPS Service: %s", o)

	if exists {
		log.Printf("HTTPS Service already exists...")
		return
	}

	if config.Environment == "local" {
		svc.Spec.LoadBalancerIP = ""
	}
	if _, err := svcs.Create(svc); err != nil {
		log.Fatalf("Error creating svc %s: %s", httpsServiceName, err.Error())
	}
}

func CreateLocalProxyService(clientset *kubernetes.Clientset) {
	svcs := clientset.Core().Services("default")
	exists := true
	_, err := svcs.Get("local-proxy")
	if err != nil {
		if !errors.IsNotFound(err) {
			log.Fatalf("Error getting service list: %s", err.Error())
		}
		exists = false
	}
	svc := &v1.Service{
		ObjectMeta: v1.ObjectMeta{
			Name:      "local-proxy",
			Namespace: "default",
		},
		Spec: v1.ServiceSpec{
			Ports: []v1.ServicePort{
				{
					Name:       "tcp",
					Port:       int32(443),
					TargetPort: intstr.FromInt(443),
					Protocol:   v1.ProtocolTCP,
					NodePort:   32000,
				},
			},
			Type: v1.ServiceTypeLoadBalancer,
			Selector: map[string]string{
				"app": "local-proxy",
			},
		},
	}
	if exists {
		log.Println("Local Proxy Service already exists...")
		return
	}

	if _, err := svcs.Create(svc); err != nil {
		log.Fatalf("Error creating svc local proxy %s", err.Error())
	}
}

func createRedisCacheService(config MicroserviceConfig, clientset *kubernetes.Clientset) {
	envConfig := config.GetEnvironment()
	svcs := clientset.Core().Services(envConfig.K8sNamespace)
	exists := true
	svcName := fmt.Sprintf("redis-%s", config.Name)
	_, err := svcs.Get(svcName)
	if err != nil {
		if !errors.IsNotFound(err) {
			log.Fatalf("Error getting : %s", err.Error())
		}
		exists = false
	}
	svc := &v1.Service{
		ObjectMeta: v1.ObjectMeta{
			Name:      svcName,
			Namespace: envConfig.K8sNamespace,
		},
		Spec: v1.ServiceSpec{
			Ports: []v1.ServicePort{
				{
					Name:       "tcp",
					Port:       int32(6379),
					TargetPort: intstr.FromInt(6379),
					Protocol:   v1.ProtocolTCP,
				},
			},
			Type: v1.ServiceTypeLoadBalancer,
			Selector: map[string]string{
				"app": svcName,
			},
		},
	}
	if exists {
		log.Println("Redis service already exists...")
		return
	}

	if _, err := svcs.Create(svc); err != nil {
		log.Fatalf("Error creating redis service: %s", err.Error())
	}
}
