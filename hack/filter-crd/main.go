package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

var removeKeys = []string{}

func main() {
	loadVariant()

	if len(flag.Args()) < 1 {
		log.Fatal("Usage: filter-crd <CRD YAML file>")
	}

	f, err := os.Open(flag.Args()[0])
	if err != nil {
		log.Fatal("Error opening file: ", err)
	}

	decoder := yaml.NewDecoder(f)
	var d map[interface{}]interface{}
	output := []string{}

	for decoder.Decode(&d) == nil {

		if len(d) == 0 {
			continue
		}

		checkChain(d, []string{})

		fileOut, err := yaml.Marshal(d)
		if err != nil {
			log.Fatal("Error marshaling output: ", err)
		}

		output = append(output, string(fileOut))

	}

	fmt.Println(strings.Join(output, "---\n"))
}

func checkChain(d map[interface{}]interface{}, chain []string) {
	for k, v := range d {
		if key, ok := k.(string); ok {
			chain = append(chain, key)

			for _, removeKey := range removeKeys {
				if strings.Join(chain, "/") == removeKey {
					delete(d, key)
				}
			}

			if value, ok := v.(map[interface{}]interface{}); ok {
				checkChain(value, chain)
			}
			chain = chain[:len(chain)-1] // we're done with this key, remove it from the chain
		}
	}
}

func loadVariant() {
	variant := ""
	flag.StringVar(&variant, "variant", "", "variant of remove rules")
	flag.Parse()

	if variant == "cert-manager-openshift" {
		// These are the keys that the script will remove for OpenShift compatibility
		removeKeys = []string{
			"spec/preserveUnknownFields",
			"spec/validation/openAPIV3Schema/type",
		}
	}
}
