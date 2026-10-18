package e2e_test

import (
	"context"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/netologist/hotel-booking/pkg/client"
)

var _ = Describe("Admin API (/api/v1/seed, /api/v1/reset)", func() {
	Context("Database lifecycle operations", func() {
		It("resets and seeds database cleanly", func(ctx context.Context) {
			// Step 1: Reset database
			resetResp, err := apiClient.ResetDatabaseWithResponse(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(resetResp.StatusCode()).To(Equal(http.StatusOK))
			Expect(resetResp.JSON200).NotTo(BeNil())
			Expect(resetResp.JSON200.Message).To(Equal("Database reset successfully"))

			// Step 2: Verify hotels are empty
			searchResp, err := apiClient.FindHotelsWithResponse(ctx, &client.FindHotelsParams{Name: "Continental"})
			Expect(err).NotTo(HaveOccurred())
			Expect(searchResp.StatusCode()).To(Equal(http.StatusOK))
			Expect(searchResp.JSON200).NotTo(BeNil())
			Expect(searchResp.JSON200.Hotels).NotTo(BeNil())
			Expect(*searchResp.JSON200.Hotels).To(BeEmpty())

			// Step 3: Seed database
			seedResp, err := apiClient.SeedDatabaseWithResponse(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(seedResp.StatusCode()).To(Equal(http.StatusOK))
			Expect(seedResp.JSON200).NotTo(BeNil())
			Expect(seedResp.JSON200.Message).To(Equal("Database seeded successfully"))
			Expect(seedResp.JSON200.HotelsCount).To(Equal(2))
			Expect(seedResp.JSON200.RoomsCount).To(Equal(12))

			// Step 4: Verify seeded hotels are found with their 6 rooms each
			searchContinental, err := apiClient.FindHotelsWithResponse(ctx, &client.FindHotelsParams{Name: "Continental"})
			Expect(err).NotTo(HaveOccurred())
			Expect(searchContinental.StatusCode()).To(Equal(http.StatusOK))
			Expect(searchContinental.JSON200).NotTo(BeNil())
			Expect(searchContinental.JSON200.Hotels).NotTo(BeNil())
			Expect(*searchContinental.JSON200.Hotels).To(HaveLen(1))
			Expect((*searchContinental.JSON200.Hotels)[0].Rooms).NotTo(BeNil())
			Expect(*(*searchContinental.JSON200.Hotels)[0].Rooms).To(HaveLen(6))

			searchBudapest, err := apiClient.FindHotelsWithResponse(ctx, &client.FindHotelsParams{Name: "Budapest"})
			Expect(err).NotTo(HaveOccurred())
			Expect(searchBudapest.StatusCode()).To(Equal(http.StatusOK))
			Expect(searchBudapest.JSON200).NotTo(BeNil())
			Expect(searchBudapest.JSON200.Hotels).NotTo(BeNil())
			Expect(*searchBudapest.JSON200.Hotels).To(HaveLen(1))
			Expect((*searchBudapest.JSON200.Hotels)[0].Rooms).NotTo(BeNil())
			Expect(*(*searchBudapest.JSON200.Hotels)[0].Rooms).To(HaveLen(6))

			// Step 5: Reset database again
			resetResp2, err := apiClient.ResetDatabaseWithResponse(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(resetResp2.StatusCode()).To(Equal(http.StatusOK))

			// Step 6: Verify empty again
			searchAfterReset, err := apiClient.FindHotelsWithResponse(ctx, &client.FindHotelsParams{Name: "Budapest"})
			Expect(err).NotTo(HaveOccurred())
			Expect(searchAfterReset.StatusCode()).To(Equal(http.StatusOK))
			Expect(searchAfterReset.JSON200).NotTo(BeNil())
			Expect(searchAfterReset.JSON200.Hotels).NotTo(BeNil())
			Expect(*searchAfterReset.JSON200.Hotels).To(BeEmpty())
		})
	})
})
