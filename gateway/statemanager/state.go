package statemanager

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

type User struct {
	ID   string       `json:"user_id"`
	File []Fileheader `json:"file"`
}

type Fileheader struct {
	File_id     string    `json:"file_id"`
	Filename    string    `json:"filename"`
	Uploaded    time.Time `json:"uploaded"`
	StoragePath string
	Size        int       `json:"size"`
	LastUpdated time.Time `json:"lastupdated"`
}

type SystemState struct {
	User map[string]*User
}

type Manager struct {
	Mu       sync.RWMutex
	Filepath string
	State    *SystemState
}

func NewManger(filepath string) (*Manager, error) {
	mgr := &Manager{
		Filepath: filepath,
		State: &SystemState{
			User: make(map[string]*User),
		},
	}
	data, err := os.ReadFile(filepath)
	if err != nil {

		if os.IsNotExist(err) {
			return mgr, nil
		}
		return nil, err
	}
	if err := json.Unmarshal(data, mgr.State); err != nil {
		return nil, err
	}
	return mgr, nil

}
func (m *Manager) Save() error {
	data, err := json.MarshalIndent(m.State, "", " ")
	if err != nil {
		return err
	}
	return os.WriteFile(m.Filepath, data, 0644)
}

func (m *Manager) Addfile(User_id string, fileheader Fileheader) error {

	m.Mu.Lock()
	defer m.Mu.Unlock()
	user, exists := m.State.User[User_id]
	if !exists {
		user = &User{
			ID:   User_id,
			File: []Fileheader{},
		}
		m.State.User[User_id] = user
	}
	user.File = append(user.File, fileheader)
	return m.Save()
}

func (m *Manager) GetUserFile(user_id string) []Fileheader {
	m.Mu.Lock()
	defer m.Mu.Unlock()
	user, exists := m.State.User[user_id]
	if !exists {
		return []Fileheader{}
	}
	copiedfile := make([]Fileheader, len(user.File))
	copy(copiedfile, user.File)
	return copiedfile
}
