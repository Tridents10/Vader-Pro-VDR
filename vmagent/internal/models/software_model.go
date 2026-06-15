package models

// Application represents an installed software package
type Application struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Publisher   string `json:"publisher"`
	InstallDate string `json:"install_date,omitempty"`
}
