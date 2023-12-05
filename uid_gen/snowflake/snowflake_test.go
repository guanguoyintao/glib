package snowflake

import (
	"testing"
	"time"
)

type MockTimeProvider struct {
	currentTimestamp func() int64
}

func (m *MockTimeProvider) CurrentTimestamp() int64 {
	return m.currentTimestamp()
}

func TestSnowflake_GenerateID(t *testing.T) {
	type fields struct {
		datacenterID int64
		machineID    int64
	}
	tests := []struct {
		name              string
		fields            fields
		currentTimestamps []int64
	}{
		{
			name:   "First ID",
			fields: fields{datacenterID: 1, machineID: 1},
			currentTimestamps: []int64{
				1623658170000, // Timestamp 1
				1623658171000, // Timestamp 2
				1623658171000, // Timestamp 3 (collision, same timestamp)
				1623658172000, // Timestamp 4
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTimeProvider := &MockTimeProvider{
				currentTimestamp: func() int64 {
					return tt.currentTimestamps[0]
				},
			}

			s := NewSnowflake(tt.fields.datacenterID, tt.fields.machineID)
			s.timeProvider = mockTimeProvider.CurrentTimestamp

			var prevID int64

			for _, ts := range tt.currentTimestamps {
				mockTimeProvider.currentTimestamp = func() int64 {
					return ts
				}

				id := s.NewInt()

				if prevID > 0 && id <= prevID {
					t.Errorf("GenerateID() does not guarantee ID increment")
				}

				prevID = id
			}
		})
	}
}

func TestSnowflake_WaitNextMillis(t *testing.T) {
	s := NewSnowflake(1, 1)
	// Set lastTimestamp to a future timestamp
	s.lastTimestamp = time.Now().UnixNano() / int64(time.Millisecond)
	lastTimestamp := s.lastTimestamp

	// Call waitNextMillis and check if the returned timestamp is greater than lastTimestamp
	nextTimestamp := s.waitNextMillis(lastTimestamp)

	if nextTimestamp <= lastTimestamp {
		t.Errorf("waitNextMillis() does not guarantee obtaining a greater timestamp")
	}
}
