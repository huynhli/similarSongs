package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Data exposed by your maintenance endpoint.
type MaintenanceData struct {
	IsActive bool   `json:"is_active"`
	Message  string `json:"message"`
	Severity string `json:"severity"`
}

// Represents a single entry returned by SpurtCMS GraphQL.
type ChannelEntry struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	IsActive    int    `json:"isActive"`
}

// Represents the channelEntriesList field returned by GraphQL.
type ChannelEntriesDetails struct {
	ChannelEntriesList []ChannelEntry `json:"channelEntriesList"`
}

// Represents the complete GraphQL response.
type GraphQLResponse struct {
	Data struct {
		ChannelEntriesList ChannelEntriesDetails `json:"channelEntriesList"`
	} `json:"data"`
}

func GetMaintenanceStatus(c *fiber.Ctx) error {
	data, err := fetchFromCMS()
	if err != nil {
		fmt.Printf("fetchFromCMS error: %v\n", err)

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve maintenance status",
		})
	}

	return c.JSON(data)
}

func fetchFromCMS() (*MaintenanceData, error) {
	cmsURL := "http://spurtcms-graphql:9001/query"

	queryData := map[string]string{
		"query": `
            {
                channelEntriesList(channelId: 1, limit: 1, offset: 0) {
                    channelEntriesList {
                        id
                        title
                        slug
                        description
                        isActive
                    }
                }
            }
        `,
	}

	jsonBody, err := json.Marshal(queryData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal GraphQL request: %w", err)
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Post(
		cmsURL,
		"application/json",
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		return nil, fmt.Errorf("CMS request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed reading CMS response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"CMS returned status %d: %s",
			resp.StatusCode,
			string(body),
		)
	}

	fmt.Printf("CMS response: %s\n", string(body))

	var graphqlResponse GraphQLResponse

	if err := json.Unmarshal(body, &graphqlResponse); err != nil {
		return nil, fmt.Errorf(
			"failed decoding CMS response: %w",
			err,
		)
	}

	entries := graphqlResponse.Data.ChannelEntriesList.ChannelEntriesList

	if len(entries) == 0 {
		return nil, fmt.Errorf("no maintenance entry found in channel")
	}

	entry := entries[0]

	parsedData := &MaintenanceData{
		IsActive: entry.IsActive == 1,
		Message:  entry.Description,
	}

	if parsedData.IsActive {
		fmt.Printf(
			"BROADCAST: Maintenance is active! Message: %s\n",
			parsedData.Message,
		)
	}

	return parsedData, nil
}
