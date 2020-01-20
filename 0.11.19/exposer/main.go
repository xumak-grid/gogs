package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	discovery "gitlab.xumak.gt/grid/service-discovery"

	"github.com/spf13/viper"
)

const (
	// NAMESPACE ENV VARIABLE
	NAMESPACE = "K8_NAMESPACE"
)

func exposeBalancer(s discovery.Selector, envVar string) {
	ns := os.Getenv(NAMESPACE)
	if ns == "default" {
		ns = ""
	}

	var services []discovery.Service
	var err error
	for {
		services, err = discovery.Locate(s, ns, false)
		if err != nil {
			fmt.Printf("service/balancer not found for %v and  %s \n", s, envVar)
			fmt.Println(err.Error(), "next try in 5 seconds")
			time.Sleep(time.Second * 5)
			continue
		}
		break
	}

	if len(services) == 0 || services[0].Balancer == "" {
		fmt.Printf("service/balancer not found for %v \n", s)
		os.Exit(1)
	}
	service := services[0]
	fmt.Printf("export %s=%s\n", envVar, service.Balancer)
}
func main() {
	viper.SetConfigFile("/meta/labels.properties")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("Error", err)
		os.Exit(1)
	}
	appName := strings.Trim(viper.GetString("app"), "\"")
	gogsSelector := discovery.Selector{
		"app": appName,
	}
	exposeBalancer(gogsSelector, "ELB_HOSTNAME")
	nexusSelector := discovery.Selector{
		"app": "nexus-server",
	}
	exposeBalancer(nexusSelector, "ELB_NEXUS_HOST")
	serviceDiscovery := discovery.Selector{
		"app": "service-deployer",
	}
	exposeBalancer(serviceDiscovery, "ELB_SERVICE_DISCOVERY_HOST")

	droneSelector := discovery.Selector{
		"app": "drone-server",
	}
	exposeBalancer(droneSelector, "ELB_DRONE_HOST")
}
