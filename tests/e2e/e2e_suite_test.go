package e2e_test

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/netologist/hotel-booking/pkg/client"
)

var (
	baseURL    string
	httpClient *http.Client
	apiClient  *client.ClientWithResponses
)

const (
	defaultBaseURL      = "http://localhost:8080"
	startupTimeout      = 30 * time.Second
	startupPollInterval = 500 * time.Millisecond
)

func TestE2E(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Service E2E Test Suite")
}

var _ = BeforeSuite(func(ctx context.Context) {
	baseURL = os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = defaultBaseURL
	}
	baseURL = strings.TrimRight(baseURL, "/")

	httpClient = &http.Client{
		Timeout: 10 * time.Second,
	}

	var err error
	apiClient, err = client.NewClientWithResponses(baseURL, client.WithHTTPClient(httpClient))
	Expect(err).NotTo(HaveOccurred())

	Eventually(ctx, func(ctx context.Context) error {
		resp, err := apiClient.GetHealthWithResponse(ctx)
		if err != nil {
			return err
		}

		if resp.StatusCode() != http.StatusOK {
			return fmt.Errorf("readiness probe returned status %d", resp.StatusCode())
		}
		return nil
	}).WithTimeout(startupTimeout).WithPolling(startupPollInterval).Should(
		Succeed(),
		"Target server %s did not become healthy within %s",
		baseURL,
		startupTimeout,
	)
})
