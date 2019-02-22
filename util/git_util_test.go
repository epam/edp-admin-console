package util

import (
	"log"
	"os"
	"testing"
)

func TestCheckExistingPublicRepository(t *testing.T) {
	gitUrl := "https://github.com/epmd-edp/java-maven-springboot.git"
	result := IsGitRepoAvailable(gitUrl, "", "")

	if result != true {
		t.Error("Expected true")
	}
}

func TestCheckNonExistingPublicRepository(t *testing.T) {
	gitUrl := "https://github.com/epmd-edp/java-maven-springboot-fake.git"
	result := IsGitRepoAvailable(gitUrl, "", "")

	if result != false {
		t.Error("Expected false")
	}
}

func TestCheckExistingPrivateRepository(t *testing.T) {
	gitUrl := "https://git.epam.com/epmd-edp/examples/basic/edp-auto-tests-simple-example.git"
	gitUser := lookupEnv("GIT_USER")
	gitPass := lookupEnv("GIT_PASSWORD")
	result := IsGitRepoAvailable(gitUrl, gitUser, gitPass)

	if result != true {
		t.Error("Expected true")
	}
}

func TestCheckNonExistingPrivateRepository(t *testing.T) {
	gitUrl := "https://git.epam.com/epmd-edp/examples/basic/edp-auto-tests-simple-example-fake.git"
	gitUser := lookupEnv("GIT_USER")
	gitPass := lookupEnv("GIT_PASSWORD")
	result := IsGitRepoAvailable(gitUrl, gitUser, gitPass)

	if result != false {
		t.Error("Expected false")
	}
}

func TestCheckExistingPrivateRepositoryWithFakeCredentials(t *testing.T) {
	gitUrl := "https://git.epam.com/epmd-edp/examples/basic/edp-auto-tests-simple-example.git"
	result := IsGitRepoAvailable(gitUrl, "fake", "fake")

	if result != false {
		t.Error("Expected false")
	}
}

func lookupEnv(key string) string {
	value, isPresented := os.LookupEnv(key)
	if !isPresented {
		log.Fatalf("required env variable by key %s is not presented", key)
	}
	return value
}
