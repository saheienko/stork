package k8s

import (
	"fmt"
	"strings"
)

func validate(name, namespace string) error {
	if err := validateNamespace(namespace); err != nil {
		return err
	}
	return validateName(name)
}

func validateName(name string) error {
	if strings.TrimSpace(name) == nil {
		return fmt.Errorf("name should not be empty")
	}
	return nil
}

func validateNamespace(namespace string) error {
	if strings.TrimSpace(namespace) == "" {
		return fmt.Errorf("namespace should not be empty")
	}
	return nil
}