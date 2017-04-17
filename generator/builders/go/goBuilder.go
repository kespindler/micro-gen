package goBuilder

// TODO: Update/Fix glide package path

import (
	"fmt"
	"os"

	"github.com/reivaj05/micro-gen/generator/utils"
)

type data struct {
	ServiceName string
}

func Build(serviceName string) error {
	fmt.Println("TODO: Implement go generator")
	if err := createServiceDir(serviceName); err != nil {
		return err
	}
	return createService(serviceName)
}

func createServiceDir(serviceName string) error {
	dst := fmt.Sprintf("./%s", serviceName)
	return os.MkdirAll(dst, os.ModePerm)
}

func createService(serviceName string) error {
	return generateFiles(serviceName)
}

func generateFiles(serviceName string) error {
	if err := generateGoFiles(serviceName); err != nil {
		return err
	}
	if err := generateConfigFile(serviceName); err != nil {
		return err
	}
	return generateBuildFiles(serviceName)
}

func generateGoFiles(serviceName string) error {
	if err := generateMainFile(serviceName); err != nil {
		return err
	}
	if err := generateEndpointsFile(serviceName); err != nil {
		return err
	}
	return generateControllersFile(serviceName)
}

func generateMainFile(serviceName string) error {
	return utils.GenerateFile(serviceName, "main", "go",
		"main.gen", &data{ServiceName: serviceName})
}

func generateEndpointsFile(serviceName string) error {
	return utils.GenerateFile(serviceName, "endpoints", "go",
		"endpoints.gen", nil)
}

func generateControllersFile(serviceName string) error {
	return utils.GenerateFile(serviceName, "controllers", "go",
		"controllers.gen", nil)
}

func generateConfigFile(serviceName string) error {
	return utils.GenerateFile(serviceName, "config", "json",
		"config.gen", nil)
}

func generateBuildFiles(serviceName string) error {
	if err := generateDockerFile(serviceName); err != nil {
		return err
	}
	if err := generateMakeFile(serviceName); err != nil {
		return err
	}
	if err := generateGlideFile(serviceName); err != nil {
		return err
	}
	if err := generateGitIgnoreFile(serviceName); err != nil {
		return err
	}
	if err := generateDockerIgnoreFile(serviceName); err != nil {
		return err
	}
	return nil
}

func generateDockerFile(serviceName string) error {
	return utils.GenerateFile(serviceName, "Dockerfile", "",
		"Dockerfile.gen", &data{ServiceName: serviceName})
}

func generateMakeFile(serviceName string) error {
	return utils.GenerateFile(serviceName, "Makefile", "",
		"Makefile.gen", nil)
}

func generateGlideFile(serviceName string) error {
	return utils.GenerateFile(serviceName, "glide", "yaml",
		"glide.gen", &data{ServiceName: serviceName})
}

func generateGitIgnoreFile(serviceName string) error {
	return utils.GenerateFile(serviceName, ".gignore", "",
		"ignore.gen", nil)
}

func generateDockerIgnoreFile(serviceName string) error {
	return utils.GenerateFile(serviceName, ".dockerignore", "",
		"ignore.gen", nil)
}
