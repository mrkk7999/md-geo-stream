package repository

import (
	"md-geo-track/request_response/location"
)

// UpdateLocationStatus updates the status field based on the given ID
func (r *repository) UpdateLocationStatus(req location.LocationReq) error {
	// Update status in the database where ID matches
	if err := r.db.Model(&location.LocationReq{}).
		Where("id = ?", req.ID).
		Update("status", "processed").Error; err != nil {
		return err
	}

	return nil
}
