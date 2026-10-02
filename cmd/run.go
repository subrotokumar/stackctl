package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

type projectKind int

const (
	kindUnknown projectKind = iota
	kindSpring
	kindQuarkus
)

func (k projectKind) String() string {
	switch k {
	case kindSpring:
		return "Spring Boot"
	case kindQuarkus:
		return "Quarkus"
	}
	return "unknown"
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// detectProject inspects the build files to decide whether this is a
// Quarkus or Spring Boot project.
func detectProject() projectKind {
	for _, f := range []string{"pom.xml", "build.gradle", "build.gradle.kts"} {
		b, err := os.ReadFile(f)
		if err != nil {
			continue
		}
		s := string(b)
		switch {
		case strings.Contains(s, "io.quarkus"):
			return kindQuarkus
		case strings.Contains(s, "org.springframework.boot"):
			return kindSpring
		}
	}
	// Fallback for Spring projects whose build file we couldn't read.
	if fileExists("src/main/resources/application.properties") ||
		fileExists("src/main/resources/application.yml") {
		return kindSpring
	}
	return kindUnknown
}

func wrapperPath(unix, windows string) string {
	if runtime.GOOS == "windows" {
		return ".\\" + windows
	}
	return "./" + unix
}

func detectBuildTool(kind projectKind) (tool string, args []string, err error) {
	mavenGoal, gradleTask := "spring-boot:run", "bootRun"
	if kind == kindQuarkus {
		mavenGoal, gradleTask = "quarkus:dev", "quarkusDev"
	}

	hasMaven := fileExists("pom.xml")
	hasGradle := fileExists("build.gradle") || fileExists("build.gradle.kts")

	if hasMaven || !hasGradle {
		if w := wrapperPath("mvnw", "mvnw.cmd"); fileExists(w) {
			return w, []string{mavenGoal}, nil
		}
		if _, err := exec.LookPath("mvn"); err == nil {
			return "mvn", []string{mavenGoal}, nil
		}
	}
	if hasGradle || !hasMaven {
		if w := wrapperPath("gradlew", "gradlew.bat"); fileExists(w) {
			return w, []string{gradleTask}, nil
		}
		if _, err := exec.LookPath("gradle"); err == nil {
			return "gradle", []string{gradleTask}, nil
		}
	}
	return "", nil, fmt.Errorf("no Maven/Gradle build tool found")
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a Spring Boot or Quarkus project",
	Run: func(cmd *cobra.Command, args []string) {
		kind := detectProject()
		if kind == kindUnknown {
			fmt.Println("❌ Not a Spring Boot or Quarkus project (no pom.xml/build.gradle with Spring Boot or Quarkus found)")
			os.Exit(1)
		}

		tool, toolArgs, err := detectBuildTool(kind)
		if err != nil {
			fmt.Printf("❌ %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("🚀 Starting %s using: %s %s\n", kind, tool, strings.Join(toolArgs, " "))

		c := exec.Command(tool, toolArgs...)
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		c.Stdin = os.Stdin

		if err := c.Run(); err != nil {
			fmt.Printf("❌ Error running %s project: %v\n", kind, err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
