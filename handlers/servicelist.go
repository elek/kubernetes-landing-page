package handlers

import (
	"github.com/gin-gonic/gin"
	"k8s.io/client-go/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"net/http"
)

type ServiceList struct {
	Namespace string
	Services  []ServiceLink
}

type ServiceLinkType string

type ServiceLink struct {
	Url  string
	Name string
	Type ServiceLinkType
}

const (
	NodePort  ServiceLinkType = "NodePort"
	ClusterIp ServiceLinkType = "ClusterIp"
)

//func GetKubernetesApi(){
//	config, err := rest.InClusterConfig()
//	if err != nil {
//		log.Info("Can't initalize in-cluster configuration. Trying to load config file from home dir.")
//	}
//
//
//
//	// use the current context in kubeconfig
//	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
//	if err != nil {
//		panic(err.Error())
//	}
//}
//

func ListServices(clientset *kubernetes.Clientset, listInterenal bool) func(*gin.Context) {
	return func(c *gin.Context) {

		namespaces, err := clientset.CoreV1().Namespaces().List(metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}

		nodes, err := clientset.CoreV1().Nodes().List(metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		firstNode := nodes.Items[0]

		serviceCatalog := make([]ServiceList, 0, 10)
		for _, namespace := range namespaces.Items {
			services, err := clientset.CoreV1().Services(namespace.Name).List(metav1.ListOptions{})
			if err != nil {
				panic(err.Error())
			}
			serviceLinks := make([]ServiceLink, 0, 10)
			for _, service := range services.Items {
				if service.Spec.Type == corev1.ServiceTypeNodePort {
					serviceLinks = append(serviceLinks, ServiceLinkFromNodeType(firstNode, service))
				} else if service.Spec.Type == corev1.ServiceTypeClusterIP && listInterenal {
					serviceLinks = append(serviceLinks, ServiceLinkFromClusterIp(service))
				}
			}
			serviceCatalog = append(serviceCatalog, ServiceList{Namespace: namespace.Name, Services: serviceLinks})

		}

		c.HTML(http.StatusOK, "index", gin.H{
			"ServiceCatalog": serviceCatalog,
		})

	}
}
func ServiceLinkFromClusterIp(service corev1.Service) ServiceLink {
	url := fmt.Sprintf("http://localhost:8001/api/v1/namespaces/%s/services/%s:/proxy/",
		service.Namespace, service.Name)
	return ServiceLink{Url: url, Name: service.Name, Type: ClusterIp}
}
func ServiceLinkFromNodeType(firstNode corev1.Node, service corev1.Service) ServiceLink {
	url := fmt.Sprintf("http://%s:%d", firstNode.Status.Addresses[0].Address, service.Spec.Ports[0].NodePort)
	return ServiceLink{Url: url, Name: service.Name, Type: NodePort}
}
