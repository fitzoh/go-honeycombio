package honeycombio

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTriggers(t *testing.T) {
	var trigger *Trigger
	var err error

	c := newTestClient(t)

	t.Run("Create", func(t *testing.T) {
		data := &Trigger{
			Name:        fmt.Sprintf("Test trigger created at %v", time.Now()),
			Description: "Some description",
			Disabled:    true,
			Query: &QuerySpec{
				Breakdowns: nil,
				Calculations: []CalculationSpec{
					{
						Op:     CalculateOpP99,
						Column: &[]string{"duration_ms"}[0],
					},
				},
				Filters:           nil,
				FilterCombination: nil,
			},
			Frequency: 300,
			Threshold: &TriggerThreshold{
				Op:    TriggerThresholdOpGreaterThan,
				Value: &[]float64{10000}[0],
			},
			Recipients: []TriggerRecipient{
				{
					Type:   "email",
					Target: "hello@example.com",
				},
			},
		}
		trigger, err = c.Triggers.Create(data)
		if err != nil {
			t.Fatal(err)
		}

		assert.NotNil(t, trigger.ID)

		data.ID = trigger.ID
		assert.Equal(t, data, trigger)
	})

	t.Run("List", func(t *testing.T) {
		triggers, err := c.Triggers.List()
		if err != nil {
			t.Fatal(err)
		}

		var createdTrigger *Trigger

		for _, tr := range triggers {
			if trigger.ID == tr.ID {
				createdTrigger = &tr
			}
		}
		if createdTrigger == nil {
			t.Fatalf("could not find newly created trigger with ID = %s", trigger.ID)
		}

		assert.Equal(t, *trigger, *createdTrigger)
	})

	t.Run("Get", func(t *testing.T) {
		getTrigger, err := c.Triggers.Get(trigger.ID)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, *trigger, *getTrigger)
	})

	t.Run("Delete", func(t *testing.T) {
		err = c.Triggers.Delete(trigger.ID)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Get_unexistingID", func(t *testing.T) {
		_, err := c.Markers.Get(trigger.ID)
		assert.Equal(t, ErrNotFound, err)
	})
}
