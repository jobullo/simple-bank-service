package database

/** SampleModel for storing name and description */
type SampleEntity struct {
	ID          int    `json:"ID,string"`
	Name        string `json:"Name"`
	Description string `json:"Description"`
}
