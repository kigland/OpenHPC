package dockerProv

func getBaseUrl() string {
	return "http://localhost:8080"
}

func (ops StartContainerOptions) WithBaseURL(envVar string, baseURL string) StartContainerOptions {
	if envVar == "" || baseURL == "" {
		return ops
	}

	ops.Env = append(ops.Env, envVar+"="+baseURL)
	return ops
}
