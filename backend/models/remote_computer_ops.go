package models

import (
	"database/sql"
	"fmt"
	"time"
)

// CreateRemoteComputer creates a new remote computer record
func CreateRemoteComputer(db *sql.DB, name string, createdByID int64) (*RemoteComputer, error) {
	fmt.Printf("Creating remote computer: %s\n", name)
	fmt.Printf("Created by user ID: %d\n", createdByID)
	result, err := db.Exec(`
		INSERT INTO remote_computers (name, created_by_id)
		VALUES (?, ?)
	`, name, createdByID)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &RemoteComputer{
		ID:          id,
		Name:        name,
		CreatedByID: createdByID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil
}

// CreateCredential creates a new credential record
func CreateCredential(db *sql.DB, username, password string, userID int64) (*Credential, error) {

	result, err := db.Exec(`
		INSERT INTO credentials (username, password, created_by_id)
		VALUES (?, ?, ?)
	`, username, password, userID)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &Credential{
		ID:        id,
		Username:  username,
		Password:  password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

// MapComputerToCredential creates a mapping between a computer and a credential
func MapComputerToCredential(db *sql.DB, computerID, credentialID int64, createdByID int64) error {
	_, err := db.Exec(`
		INSERT INTO computer_credential_mappings (computer_id, credential_id, created_by_id)
		VALUES (?, ?, ?)
	`, computerID, credentialID, createdByID)
	return err
}

// GetComputersByUserID retrieves all computers created by a specific user
func GetComputersByUserID(db *sql.DB, userID int64) ([]RemoteComputer, error) {
	rows, err := db.Query(`
		SELECT id, name, created_at, updated_at, created_by_id
		FROM remote_computers
		WHERE created_by_id = ?
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var computers []RemoteComputer
	for rows.Next() {
		var computer RemoteComputer
		err := rows.Scan(&computer.ID, &computer.Name, &computer.CreatedAt, &computer.UpdatedAt, &computer.CreatedByID)
		if err != nil {
			return nil, err
		}
		computers = append(computers, computer)
	}
	return computers, nil
}

// GetCredentialsByComputerID retrieves all credentials associated with a computer
func GetCredentialsByComputerID(db *sql.DB, computerID int64) ([]Credential, error) {
	rows, err := db.Query(`
		SELECT c.id, c.username, c.password, c.created_at, c.updated_at
		FROM credentials c
		JOIN computer_credential_mappings m ON c.id = m.credential_id
		WHERE m.computer_id = ?
	`, computerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var credentials []Credential
	for rows.Next() {
		var cred Credential
		err := rows.Scan(&cred.ID, &cred.Username, &cred.Password, &cred.CreatedAt, &cred.UpdatedAt)
		if err != nil {
			return nil, err
		}
		credentials = append(credentials, cred)
	}
	return credentials, nil
}

// DeleteComputerCredentialMapping removes a mapping between a computer and a credential
func DeleteComputerCredentialMapping(db *sql.DB, computerID, credentialID int64) error {
	_, err := db.Exec(`
		DELETE FROM computer_credential_mappings
		WHERE computer_id = ? AND credential_id = ?
	`, computerID, credentialID)
	return err
}

// DeleteComputer deletes a remote computer and all its credential mappings
func DeleteComputer(db *sql.DB, computerID int64) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Delete mappings first due to foreign key constraints
	_, err = tx.Exec(`
		DELETE FROM computer_credential_mappings
		WHERE computer_id = ?
	`, computerID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Delete the computer
	_, err = tx.Exec(`
		DELETE FROM remote_computers
		WHERE id = ?
	`, computerID)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// CheckComputerOwnership checks if the computer was created by the specified user
func CheckComputerOwnership(db *sql.DB, computerID, userID int64) (bool, error) {
	var count int
	err := db.QueryRow(`
		SELECT COUNT(*)
		FROM remote_computers
		WHERE id = ? AND created_by_id = ?
	`, computerID, userID).Scan(&count)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// CheckCredentialOwnership checks if the credential was created by the specified user
func CheckCredentialOwnership(db *sql.DB, credentialID, userID int64) (bool, error) {
	var count int
	err := db.QueryRow(`
		SELECT COUNT(*)
		FROM credentials
		WHERE id = ? AND created_by_id = ?
	`, credentialID, userID).Scan(&count)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// CheckComputerCredentialMappingOwnership checks if the computer credential mapping was created by the specified user
func CheckComputerCredentialMappingOwnership(db *sql.DB, computerCredentialMappingID, userID int64) (bool, error) {
	var count int
	err := db.QueryRow(`
		SELECT COUNT(*)
		FROM computer_credential_mappings
		WHERE id = ? AND created_by_id = ?
	`, computerCredentialMappingID, userID).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
