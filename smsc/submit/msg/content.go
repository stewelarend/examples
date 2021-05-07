package msg

type Content struct {
	Text string `json:"string"`
}

func (c Content) Validate() error {
	return nil
}
