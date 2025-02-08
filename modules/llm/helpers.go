package llm

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/YourTechBud/inferix/modules/llm/types"
)

// fetchModelsFromBackend fetches models from a single backend with timeout
func (llm *LLM) fetchModelsFromBackend(ctx context.Context, backendID string, backend types.Backend) ([]types.ModelObject, error) {
	models, err := backend.GetModels(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch models: %w", err)
	}

	// Set the ownedBy field for all models
	for i := range models {
		models[i].OwnedBy = backendID
	}

	return models, nil
}

// fetchAllModels concurrently fetches models from all backends
func (llm *LLM) fetchAllModels() []types.ModelObject {
	var wg sync.WaitGroup
	modelsChannel := make(chan []types.ModelObject)

	// Start fetching models concurrently for each backend
	for backendID, backend := range llm.backends.GetBackends() {
		wg.Add(1)
		go func(backendID string, backend types.Backend) {
			defer wg.Done()

			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			models, err := llm.fetchModelsFromBackend(ctx, backendID, backend)
			if err != nil {
				fmt.Printf("Error fetching models for backend %s: %v\n", backendID, err)
				return
			}

			modelsChannel <- models
		}(backendID, backend)
	}

	// Close channel after all goroutines complete
	go func() {
		wg.Wait()
		close(modelsChannel)
	}()

	// Collect all models
	var allModels []types.ModelObject
	for models := range modelsChannel {
		allModels = append(allModels, models...)
	}

	return allModels
}

// startModelPolling starts a goroutine that polls models for all backends periodically
func (llm *LLM) startModelPolling() {
	llm.wg.Add(1)
	go func() {
		defer llm.wg.Done()

		// Run immediately first
		log.Default().Println("Fetching model list from backends")
		allModels := llm.fetchAllModels()
		llm.models.MergeModels(allModels)

		ticker := time.NewTicker(5 * time.Minute) // TODO: Make this configurable
		defer ticker.Stop()

		for {
			select {
			case <-llm.done:
				return
			case <-ticker.C:
				log.Default().Println("Fetching model list from backends")
				allModels := llm.fetchAllModels()
				llm.models.MergeModels(allModels)
			}
		}
	}()
}
