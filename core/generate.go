package core

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

const buildGradleTemplate = `
plugins {
    // Apply the java-library plugin to add support for Java Library
    id 'java-library'
}

repositories {
    // Use jcenter for resolving dependencies.
    // You can declare any Maven/Ivy/file repository here.
    jcenter()
}

dependencies {
    compile '{{.Dependency}}'
}

sourceCompatibility = '{{.SourceCompatibility}}'
targetCompatibility = '{{.TargetCompatibility}}'

jar {
    dependsOn configurations.runtimeClasspath

    from {
        configurations.runtimeClasspath.collect { it.isDirectory() ? it : zipTree(it) }
    }
}
`

const settingsGradleTemplate = `
rootProject.name = '{{.ProjectName}}'
`

func Generate(dependency, javaVersion string, dst io.Writer) error {
	projectName := "lib"

	if javaVersion == "" {
		javaVersion = "1.8"
	}

	tmpDir, err := os.MkdirTemp(os.TempDir(), "*")
	if err != nil {
		return fmt.Errorf("creating temp dir: %w", err)
	}
	log.Println("Created temp directory:", tmpDir)

	defer func() {
		err = os.RemoveAll(tmpDir)
		if err != nil {
			log.Printf("deleting temp directory: %v\n", err)
		}
		log.Println("Deleted temp directory:", tmpDir)
	}()

	filesTemplates := map[string]string{
		"build.gradle":    buildGradleTemplate,
		"settings.gradle": settingsGradleTemplate,
	}

	data := map[string]interface{}{
		"Dependency":          dependency,
		"SourceCompatibility": javaVersion,
		"TargetCompatibility": javaVersion,
		"ProjectName":         projectName,
	}

	for filename, filetemplate := range filesTemplates {
		buildGradleFile, err := os.OpenFile(filepath.Join(tmpDir, filename), os.O_WRONLY|os.O_CREATE, os.ModePerm)
		if err != nil {
			return fmt.Errorf("opening %s file: %w", filename, err)
		}
		defer buildGradleFile.Close()
		log.Println("Opened temp file:", buildGradleFile.Name())

		tpl, err := template.New(filename).Parse(filetemplate)
		if err != nil {
			return fmt.Errorf("parsing %s template: %w", filename, err)
		}
		log.Println("Parsed template:", filename)

		err = tpl.Execute(buildGradleFile, data)
		if err != nil {
			return fmt.Errorf("executing %s template: %w", filename, err)
		}
		log.Println("Executed template:", filename)
	}

	cmd := exec.Command("gradle", "--project-dir", tmpDir, "jar")
	// cmd := exec.Command("gradle", "--project-dir", tmpDir, "tasks")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("executing command: %s - %w", cmd.String(), err)
	}

	jarpath := filepath.Join(tmpDir, "build/libs", projectName+".jar")
	jarFile, err := os.Open(jarpath)
	if err != nil {
		return fmt.Errorf("opening jar file: %s - %w", jarpath, err)
	}
	defer jarFile.Close()

	_, err = io.Copy(dst, jarFile)
	if err != nil {
		return fmt.Errorf("could copy jar file: %w", err)
	}

	return nil
}
