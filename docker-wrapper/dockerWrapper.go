package dockerWrapper

import (
	"fmt"
	"os"

	"github.com/reivaj05/GoJSON"
	"github.com/reivaj05/GoRequester"
)

// TODO: Refactor

var dockerUsernameKey = "DOCKER_USERNAME"
var dockerPasswordKey = "DOCKER_PASSWORD"
var dockerRegistryHostKey = "DOCKER_REGISTRY_HOST"

var repositoriesEndpoint = "%s/v2/repositories/%s/"
var loginEndpoint = "%s/v2/users/login/"

type dockerRegistryManager struct {
	client   *requester.Requester
	token    string
	host     string
	username string
}

func NewDockerRegistryManager() (*dockerRegistryManager, error) {
	if err := checkDockerCredentials(); err != nil {
		return nil, err
	}
	host, username := os.Getenv(dockerRegistryHostKey), os.Getenv(dockerUsernameKey)
	manager := &dockerRegistryManager{host: host, username: username, client: requester.New()}
	if err := manager.setToken(); err != nil {
		return nil, err
	}
	return manager, nil
}

func checkDockerCredentials() error {
	for _, key := range []string{dockerUsernameKey, dockerPasswordKey, dockerRegistryHostKey} {
		if value := os.Getenv(key); value == "" {
			return fmt.Errorf("Env var %s not set", key)
		}
	}
	return nil
}

func (manager *dockerRegistryManager) setToken() error {
	token, err := manager.getToken()
	if err != nil {
		return err
	}
	manager.token = token
	return nil
}

func (manager *dockerRegistryManager) getToken() (string, error) {
	username, password := os.Getenv(dockerUsernameKey), os.Getenv(dockerPasswordKey)
	data := fmt.Sprintf(`{"username":"%s","password":"%s"}`, username, password)
	return manager.login(data)
}

func (manager *dockerRegistryManager) login(data string) (string, error) {
	response, status, err := manager.client.MakeRequest(manager.createLoginRequestConfig(data))
	if err != nil || status >= 400 {
		return "", fmt.Errorf("Error %v: Got status %d", status, err)
	}
	return manager.parseLoginResponse(response)
}

func (manager *dockerRegistryManager) createLoginRequestConfig(data string) *requester.RequestConfig {
	return &requester.RequestConfig{
		URL:     fmt.Sprintf(loginEndpoint, manager.host),
		Method:  "POST",
		Body:    []byte(data),
		Headers: map[string]string{"Content-Type": "application/json"},
	}
}

func (manager *dockerRegistryManager) parseLoginResponse(response string) (string, error) {
	jsonString, _ := GoJSON.New(response)
	value, _ := jsonString.GetStringFromPath("token")
	return value, nil
}

func (manager *dockerRegistryManager) FilterByExistingRepos(dataToFilter []string) []string {
	if reposResponse, err := manager.SearchRepos(); err == nil {
		dataToFilter = manager.filterExistingRepos(dataToFilter, reposResponse)
	}
	return dataToFilter
}

func (manager *dockerRegistryManager) filterExistingRepos(
	dataToFilter []string, reposResponse *GoJSON.JSONWrapper) (results []string) {

	repos := reposResponse.GetArrayFromPath("results")
	for _, item := range dataToFilter {
		if manager.isItemInDockerRegistry(repos, item) {
			results = append(results, item)
		}
	}
	return results
}

func (manager *dockerRegistryManager) isItemInDockerRegistry(
	repos []*GoJSON.JSONWrapper, item string) bool {

	for _, repo := range repos {
		if repo.HasPath("name") {
			if name, ok := repo.GetStringFromPath("name"); ok && name == item {
				return true
			}
		}
	}
	return false
}

func (manager *dockerRegistryManager) SearchRepos() (*GoJSON.JSONWrapper, error) {
	response, status, err := manager.client.MakeRequest(manager.createSearchReposRequestConfig())
	if err != nil || status >= 400 {
		return nil, fmt.Errorf("Error %v: Got status %d", status, err)
	}
	return manager.toJSON(response)
}

func (manager *dockerRegistryManager) createSearchReposRequestConfig() *requester.RequestConfig {
	return &requester.RequestConfig{
		Method:  "GET",
		URL:     fmt.Sprintf(repositoriesEndpoint, manager.host, manager.username),
		Headers: map[string]string{"Authorization": "JWT " + manager.token},
	}
}

func (manager *dockerRegistryManager) toJSON(data string) (*GoJSON.JSONWrapper, error) {
	jsonData, err := GoJSON.New(data)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}
