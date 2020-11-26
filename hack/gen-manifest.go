package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/argoproj/gitops-engine/pkg/utils/kube"
	extensionsobj "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
)

func getCRD(yamlCrdBytes []byte) error {

	objs, err := kube.SplitYAML(yamlCrdBytes)
	if err != nil {
		return err
	}

	crds := []extensionsobj.CustomResourceDefinition{}

	for _, obj := range objs {
		b, err := obj.MarshalJSON()
		if err != nil {
			log.Fatalln(err)
		}

		var crd extensionsobj.CustomResourceDefinition
		err = json.Unmarshal(b, &crd)
		if err != nil {
			return err
		}

		crds = append(crds, crd)
	}

	for _, crd := range crds {
		fmt.Println(crd)
	}

	return nil
}

func main() {

	b, err := ioutil.ReadFile("manifests/circle-crd.yaml")
	if err != nil {
		log.Fatalln(err)
	}

	if err := getCRD(b); err != nil {
		log.Fatalln(err)
	}
}
