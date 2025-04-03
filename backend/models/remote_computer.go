package models

import (
	"database/sql"
	"time"
)

type RemoteComputer struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedByID int64     `json:"created_by_id"`
}

type Credential struct {
	ID          int64     `json:"id"`
	Username    string    `json:"username"`
	Password    string    `json:"password"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedByID int64     `json:"created_by_id"`
}

type ComputerCredentialMapping struct {
	ID           int64     `json:"id"`
	ComputerID   int64     `json:"computer_id"`
	CredentialID int64     `json:"credential_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	CreatedByID  int64     `json:"created_by_id"`
}

func (table *RemoteComputer) InitTable(db *sql.DB) error {
	// Create RemoteComputers table
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS remote_computers (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			created_by_id INTEGER NOT NULL,
			FOREIGN KEY(created_by_id) REFERENCES users(id)
		)
	`)

	return err
}

// InitDB initializes the database tables
func (table *Credential) InitTable(db *sql.DB) error {
	// Create Credentials table
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS credentials (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL,
			password TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			created_by_id INTEGER NOT NULL,
			FOREIGN KEY(created_by_id) REFERENCES users(id)
		)
	`)

	return err
}

func (table *ComputerCredentialMapping) InitTable(db *sql.DB) error {
	// Create ComputerCredentialMappings table
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS computer_credential_mappings (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			computer_id INTEGER NOT NULL,
			credential_id INTEGER NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			created_by_id INTEGER NOT NULL,
			FOREIGN KEY(created_by_id) REFERENCES users(id),
			FOREIGN KEY(computer_id) REFERENCES remote_computers(id),
			FOREIGN KEY(credential_id) REFERENCES credentials(id)
		)
	`)
	return err
}

// GetCredentialsByUserID retrieves all credentials created by a specific user
func GetCredentialsByUserID(db *sql.DB, userID int64) ([]Credential, error) {
	rows, err := db.Query(`
		SELECT id, username, password, created_at, updated_at, created_by_id
		FROM credentials
		WHERE created_by_id = ?
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var credentials []Credential
	for rows.Next() {
		var cred Credential
		err := rows.Scan(&cred.ID, &cred.Username, &cred.Password, &cred.CreatedAt, &cred.UpdatedAt, &cred.CreatedByID)
		if err != nil {
			return nil, err
		}
		credentials = append(credentials, cred)
	}
	return credentials, nil
}

// DeleteCredential deletes a credential by its ID
func DeleteCredential(db *sql.DB, credentialID int64) error {
	// First delete any mappings that use this credential
	_, err := db.Exec(`
		DELETE FROM computer_credential_mappings
		WHERE credential_id = ?
	`, credentialID)
	if err != nil {
		return err
	}

	// Then delete the credential itself
	_, err = db.Exec(`
		DELETE FROM credentials
		WHERE id = ?
	`, credentialID)
	return err
}

// UpdateCredentialPassword updates the password of a credential
func UpdateCredentialPassword(db *sql.DB, credentialID int64, newPassword string) error {
	_, err := db.Exec(`
		UPDATE credentials
		SET password = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, newPassword, credentialID)
	return err
}

// UpdateComputerCredentialMapping updates the mapping between a computer and credential
func UpdateComputerCredentialMapping(db *sql.DB, computerCredentialMappingID int64, credentialID int64, userID int64) error {
	_, err := db.Exec(`
		UPDATE computer_credential_mappings
		SET credential_id = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, credentialID, computerCredentialMappingID)
	return err
}
