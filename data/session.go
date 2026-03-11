package data

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"sourcecraft.dev/organization-shipmonitor/ship-cloud-auth/db"
	"sourcecraft.dev/organization-shipmonitor/ship-cloud-auth/keyval"
)

type Session struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"userId"`
	UserBlocked bool      `json:"userBlocked"`
	Permissions []string  `json:"permissions"`
}

type SessionRecord struct {
	*bun.BaseModel   `bun:"table:sessions"`
	*TimestampsModel `bun:",embed"`

	SessionID uuid.UUID `bun:",type:uuid,pk" json:"sessionId"`
	UserID    uuid.UUID `bun:",type:uuid,notnull" json:"userId"`
}

func newSessionRecord(session *Session) (*SessionRecord, error) {
	dbSession := SessionRecord{
		SessionID: session.ID,
		UserID:    session.UserID,
	}
	_, err := db.DB.NewInsert().Model(&dbSession).Exec(context.TODO())
	return &dbSession, err
}

const sessionTTL = time.Hour * 24 * 30

func NewSession(userID uuid.UUID) (*Session, error) {
	session := Session{
		ID:          uuid.New(),
		UserID:      userID,
		UserBlocked: false,      // TODO: get blocked status from database
		Permissions: []string{}, // TODO: get permissions from database
	}

	_, err := newSessionRecord(&session)
	if err != nil {
		return nil, err
	}
	data, err := json.Marshal(session)
	if err != nil {
		return nil, fmt.Errorf("failed marshal Session to JSON: %s", err)
	}
	err = keyval.RDB.Set(context.TODO(), session.ID.String(), data, sessionTTL).Err()
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func GetSession(sessionID uuid.UUID) (*Session, error) {
	sessionJSON, err := keyval.RDB.Get(context.TODO(), sessionID.String()).Result()
	if err != nil {
		return nil, err
	}

	var session Session
	err = json.Unmarshal([]byte(sessionJSON), &session)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func getSessionRecords(userID uuid.UUID) ([]SessionRecord, error) {
	var records []SessionRecord
	// TODO: add where statement to filter by TTL
	err := db.DB.NewSelect().Model(&records).Where("user_id = ?", userID).Scan(context.TODO())
	if err != nil {
		return nil, err
	}
	return records, nil
}

func GetActiveSessions(userID uuid.UUID) ([]*Session, error) {
	records, err := getSessionRecords(userID)
	if err != nil {
		return nil, fmt.Errorf("failed get session records: %s", err)
	}

	var allErr error

	sessions := []*Session{}
	for _, record := range records {
		s, err := GetSession(record.SessionID)
		if err != nil {
			allErr = errors.Join(allErr, err)
			continue
		}
		sessions = append(sessions, s)
	}

	return sessions, nil
}

func DeleteSession(sessionID uuid.UUID) error {
	return keyval.RDB.Del(context.TODO(), sessionID.String()).Err()
}
