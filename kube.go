package main

import (
	"errors"
	"os"
	"sort"
	"strconv"

	"github.com/Luke-Vear/binnacle/dnsprovider"

	"k8s.io/client-go/kubernetes"
	rest "k8s.io/client-go/rest"
)

// getKubeIngressIPs polls the kubernetes API server for the current IP addresses
// of the ingress controllers.
func getKubeIngressIPs(dp dns.Provider) ([]string, error) {

	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	// get ingress information from API server
	ing, err := clientset.
		ExtensionsV1beta1().
		Ingresses(dp.NameSpace()).
		Get(dp.IngressName())
	if err != nil {
		return nil, err
	}

	if len(ing.Status.LoadBalancer.Ingress) == 0 {
		return nil, errors.New("no ingress IPs")
	}

	// build up ingress slice sort and return
	ingSlice := []string{}

	for _, lbing := range ing.Status.LoadBalancer.Ingress {
		ingSlice = append(ingSlice, lbing.IP)
	}

	sort.Strings(ingSlice)
	return ingSlice, nil
}

// getKubeConfigMap reads the kubernetes configmap attached to the pod and fills out
// the dns.Data with the configuration.
// TODO error handling on this
func getKubeConfigMap() (dns.Conf, error) {
	if ttl, err := strconv.ParseInt(os.Getenv("BINNACLE_TTL"), 10, 64); err == nil {
		return dns.NewConf(
			os.Getenv("BINNACLE_PROVIDER"),
			os.Getenv("BINNACLE_INGRESSNAME"),
			os.Getenv("BINNACLE_NAMESPACE"),
			os.Getenv("BINNACLE_RECORDNAME"),
			os.Getenv("BINNACLE_ZONEID"),
			ttl,
		), nil
	}
	return dns.Conf{}, errors.New("problem getting kube ConfigMap")
}
