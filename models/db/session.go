package db

import (
	"github.com/billyboar/battleships/models"
)

// GetEvents returns all events for a session stream
func (store *Store) GetEvents(sessionID string) ([]*models.Event, error) {
	events, err := store.connection.XRange(sessionID, "-", "+").Result()
	if err != nil {
		return nil, err
	}

	deserializedEvents := make([]*models.Event, len(events))
	for i, event := range events {
		deserializedEvents[i] = models.DeserializeRedisStream(event)
	}

	return deserializedEvents, nil
}

// AppendEvent adds new event to stream
func (store *Store) AppendEvent(sessionID string, event *models.Event) error {
	return store.connection.XAdd(event.SerializeRedisStream()).Err()
}
