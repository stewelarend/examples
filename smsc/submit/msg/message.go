package msg

import (
	"fmt"
	"time"
)

type Message struct {
	From    Address
	To      []Address
	Content Content
	After   time.Time
}

func (m Message) Validate() error {
	if err := m.From.Validate(); err != nil {
		return fmt.Errorf("invalid from address(%v): %v", m.From, err)
	}
	if len(m.To) < 1 {
		return fmt.Errorf("missing to address")
	}
	for _, to := range m.To {
		if err := to.Validate(); err != nil {
			return fmt.Errorf("invalid to address(%v): %v", to, err)
		}
	}
	if err := m.Content.Validate(); err != nil {
		return fmt.Errorf("invalid content: %v", err)
	}
	return nil
}
