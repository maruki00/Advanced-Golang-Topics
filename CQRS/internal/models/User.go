package models

type User struct {
	ID   int    `yaml:"id" json:"id,omitempty"`
	Name string `yaml:"name" json:"name,omitempty"`
}

func (obj User) method() {

}
