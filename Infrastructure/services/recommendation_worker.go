package services

import (
	"blog-api/Domain/interfaces"
	"log"
	"time"
)

type RecommendationWorker struct {
	recommendationUC interfaces.RecommendationUseCase
	interval         time.Duration
	stopChan         chan bool
}

func NewRecommendationWorker(recommendationUC interfaces.RecommendationUseCase) *RecommendationWorker {
	return &RecommendationWorker{
		recommendationUC: recommendationUC,
		interval:         1 * time.Hour, // Run every hour
		stopChan:         make(chan bool),
	}
}

// Start starts the background worker
func (rw *RecommendationWorker) Start() {
	log.Println("Starting recommendation worker...")

	go func() {
		ticker := time.NewTicker(rw.interval)
		defer ticker.Stop()

		// Run initial processing
		rw.processAll()

		for {
			select {
			case <-ticker.C:
				rw.processAll()
			case <-rw.stopChan:
				log.Println("Stopping recommendation worker...")
				return
			}
		}
	}()
}

// Stop stops the background worker
func (rw *RecommendationWorker) Stop() {
	rw.stopChan <- true
}

// processAll runs all background processing tasks
func (rw *RecommendationWorker) processAll() {
	log.Println("Running recommendation background processing...")

	// Process content similarities
	go func() {
		if err := rw.recommendationUC.UpdateContentSimilarities(); err != nil {
			log.Printf("Error updating content similarities: %v", err)
		} else {
			log.Println("Content similarities updated successfully")
		}
	}()

	// Process user recommendations
	go func() {
		if err := rw.recommendationUC.ProcessRecommendations(); err != nil {
			log.Printf("Error processing recommendations: %v", err)
		} else {
			log.Println("User recommendations processed successfully")
		}
	}()

	// Cleanup old data
	go func() {
		if err := rw.recommendationUC.CleanupOldData(); err != nil {
			log.Printf("Error cleaning up old data: %v", err)
		} else {
			log.Println("Old data cleaned up successfully")
		}
	}()
}

// ProcessContentSimilarities processes content similarities immediately
func (rw *RecommendationWorker) ProcessContentSimilarities() error {
	log.Println("Processing content similarities...")
	return rw.recommendationUC.UpdateContentSimilarities()
}

// ProcessUserRecommendations processes user recommendations immediately
func (rw *RecommendationWorker) ProcessUserRecommendations() error {
	log.Println("Processing user recommendations...")
	return rw.recommendationUC.ProcessRecommendations()
}

// CleanupOldData cleans up old data immediately
func (rw *RecommendationWorker) CleanupOldData() error {
	log.Println("Cleaning up old data...")
	return rw.recommendationUC.CleanupOldData()
}
