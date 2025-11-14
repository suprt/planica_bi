package services

// SyncService handles data synchronization with Yandex APIs
type SyncService struct {
	// TODO: add repositories and integrations
}

// NewSyncService creates a new sync service
func NewSyncService() *SyncService {
	return &SyncService{}
}

// SyncProject synchronizes data for a specific project
func (s *SyncService) SyncProject(projectID uint) error {
	// TODO: implement
	// 1. Fetch data from Yandex.Metrica
	// 2. Fetch data from Yandex.Direct
	// 3. Aggregate and save to database
	return nil
}

// SyncAllProjects synchronizes data for all active projects
func (s *SyncService) SyncAllProjects() error {
	// TODO: implement
	return nil
}

// FinalizeMonth finalizes data for the previous month
func (s *SyncService) FinalizeMonth() error {
	// TODO: implement
	// Called on 1st of each month at 07:00 MSK
	return nil
}

