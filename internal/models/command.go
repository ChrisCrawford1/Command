package models

import (
	"context"
	"database/sql"
	uuid "github.com/satori/go.uuid"
	"log"
	"time"
)

type CommandModel struct {
	DB *sql.DB
}

type CommandCreationRequest struct {
	Name        string `json:"name" validate:"required"`
	Language    string `json:"language" validate:"required"`
	Description string `json:"description"`
	Syntax      string `json:"syntax" validate:"required"`
}

type Command struct {
	ID          int
	UUID        uuid.UUID
	Name        string
	Language    string
	Description string
	Syntax      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type CommandDTO struct {
	UUID        uuid.UUID `json:"uuid"`
	Name        string    `json:"name"`
	Language    string    `json:"language"`
	Description string    `json:"description"`
	Syntax      string    `json:"syntax"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type CommandDtoList struct {
	Commands []CommandDTO `json:"commands"`
}

func (command *Command) AsDTO() CommandDTO {
	return CommandDTO{
		UUID:        command.UUID,
		Name:        command.Name,
		Language:    command.Language,
		Description: command.Description,
		Syntax:      command.Syntax,
		CreatedAt:   command.CreatedAt,
		UpdatedAt:   command.UpdatedAt,
	}
}

func (m CommandModel) CreateCommand(creationRequest CommandCreationRequest) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
			INSERT INTO commands (uuid, name, language, description, syntax, created_at, updated_at) 
			VALUES ($1,$2,$3,$4,$5,$6,$7)`

	newUuid := uuid.NewV4()
	_, err := m.DB.ExecContext(
		ctx,
		stmt,
		newUuid.String(),
		creationRequest.Name,
		creationRequest.Language,
		creationRequest.Description,
		creationRequest.Syntax,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return false, err
	}

	return true, nil
}

func (m CommandModel) GetByUUID(uuid string) (Command, error) {
	stmt, err := m.DB.Prepare("SELECT * FROM commands WHERE uuid = $1")

	if err != nil {
		log.Fatal(err)
	}

	var command Command
	err = stmt.QueryRow(uuid).Scan(
		&command.ID,
		&command.UUID,
		&command.Name,
		&command.Language,
		&command.Description,
		&command.Syntax,
		&command.CreatedAt,
		&command.UpdatedAt,
	)

	if err != nil {
		return Command{}, err
	}

	return command, nil
}

func (m CommandModel) GetAll() []*Command {
	defaultLimit := 20
	commands := make([]*Command, 0)
	stmt, err := m.DB.Prepare("SELECT * FROM commands LIMIT $1")

	if err != nil {
		log.Fatal(err)
	}

	rows, _ := stmt.Query(defaultLimit)

	for rows.Next() {
		command := new(Command)
		if scanErr := rows.Scan(
			&command.ID,
			&command.UUID,
			&command.Name,
			&command.Description,
			&command.Language,
			&command.Syntax,
			&command.CreatedAt,
			&command.UpdatedAt,
		); scanErr != nil {
			log.Fatal(err)
		}
		commands = append(commands, command)
	}

	return commands
}
