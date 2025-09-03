package provider

import (
	"fmt"
	"github.com/charmbracelet/crush/internal/config"
	"github.com/charmbracelet/crush/internal/log"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/azure"
	"github.com/openai/openai-go/option"
)

type azureClient struct {
	*openaiClient
}

type AzureClient ProviderClient

func newAzureClient(opts providerClientOptions) AzureClient {
	apiVersion := opts.extraParams["apiVersion"]
	if apiVersion == "" {
		apiVersion = "2025-01-01-preview"
	}

	model := opts.model(opts.modelType)
	deploymentID := model.ID

	// Construct the full base URL without the API version query parameter.
	// The openai-go client will append "/chat/completions" to this base URL.
	fullBaseURL := fmt.Sprintf("%s/openai/deployments/%s", opts.baseURL, deploymentID)

	reqOpts := []option.RequestOption{
		option.WithBaseURL(fullBaseURL),
		azure.WithAPIKey(opts.apiKey),
		option.WithQuery("api-version", apiVersion), // Add API version as a separate query parameter
	}

	if config.Get().Options.Debug {
		httpClient := log.NewHTTPClient()
		reqOpts = append(reqOpts, option.WithHTTPClient(httpClient))
	}

	client := openai.NewClient(reqOpts...)

	base := &openaiClient{
		providerOptions: opts,
		client:          client,
	}

	return &azureClient{openaiClient: base}
}
