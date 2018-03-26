package main

import (
	"fmt"
	"github.com/elek/gin-template"
	"github.com/gin-gonic/gin"
	"github.com/shokunin/contrib/ginrus"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"kubernetes-landing-page/handlers"
	"os"
	"time"
)

func main() {
	var kubeConfig string

	var rootCmd = &cobra.Command{
		Use:   "kubernetes-landing-page",
		Short: "Landing page for kubernetes clusters.",
		Long: `A generic landing page which displays links to all
		the available services/ingress`,
		Run: func(cmd *cobra.Command, args []string) {
			kubernetes := initKubeApi(kubeConfig)
			router := gin.New()
			logrus.SetFormatter(&logrus.JSONFormatter{})
			router.Use(ginrus.Ginrus(logrus.StandardLogger(), time.RFC3339, true, "kubernetes-landing-page"))
			router.HTMLRender = gintemplate.New(gintemplate.TemplateConfig{
				Root:         "views",
				Extension:    ".html",
				Master:       "layouts/master",
				Partials:     []string{},
				DisableCache: true,
				AssetFunction: func(name string) ([]byte, error) {
					return handlers.Asset("views/" + name)
				},
			})

			// Start routes
			router.GET("/health", handlers.HealthCheck)
			router.GET("/", handlers.ListServices(kubernetes, false))

			// RUN rabit run
			router.Run() // listen and serve on 0.0.0.0:8080
		},
	}

	rootCmd.PersistentFlags().StringVarP(&kubeConfig, "kubeconfig", "c",
		"", "Kubernetes config file.")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
func initKubeApi(configFile string) *kubernetes.Clientset {
	var config *rest.Config
	var err error
	if configFile != "" {
		config, err = clientcmd.BuildConfigFromFlags("", configFile)
		if err != nil {
			panic(err.Error())
		}
	} else {
		config, err = rest.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}
	}
	// creates the in-cluster config

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return clientset
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
