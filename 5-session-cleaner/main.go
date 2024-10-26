//////////////////////////////////////////////////////////////////////
//
// Given is a SessionManager that stores session information in
// memory. The SessionManager itself is working, however, since we
// keep on adding new sessions to the manager our program will
// eventually run out of memory.
//
// Your task is to implement a session cleaner routine that runs
// concurrently in the background and cleans every session that
// hasn't been updated for more than 5 seconds (of course usually
// session times are much longer).
//
// Note that we expect the session to be removed anytime between 5 and
// 7 seconds after the last update. Also, note that you have to be
// very careful in order to prevent race conditions.
//

package main

import (
	"container/list"
	"errors"
	"log"
	"sync"
	"time"
)

// SessionManager keeps track of all sessions from creation, updating
// to destroying.
type SessionManager struct {
	sessions             map[string]*list.Element
	recentlyUsedSessions *list.List
	mu                   *sync.Mutex
}

// Session stores the session's data
type Session struct {
	ID       string
	Data     map[string]interface{}
	lastUsed time.Time
}

// NewSessionManager creates a new sessionManager
func NewSessionManager() *SessionManager {
	m := &SessionManager{
		sessions:             make(map[string]*list.Element),
		mu:                   &sync.Mutex{},
		recentlyUsedSessions: list.New(),
	}

	// start the background cleaner for the session manager
	go backgroundCleaner(m)
	return m
}

// CreateSession creates a new session and returns the sessionID
func (m *SessionManager) CreateSession() (string, error) {
	sessionID, err := MakeSessionID()
	if err != nil {
		return "", err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	session := Session{
		ID:       sessionID,
		Data:     make(map[string]interface{}),
		lastUsed: time.Now(),
	}
	// add the session to session manager map and list
	m.sessions[sessionID] = &list.Element{Value: &session}
	m.recentlyUsedSessions.PushFront(&session)

	return sessionID, nil
}

// ErrSessionNotFound returned when sessionID not listed in
// SessionManager
var ErrSessionNotFound = errors.New("SessionID does not exists")

// GetSessionData returns data related to session if sessionID is
// found, errors otherwise
func (m *SessionManager) GetSessionData(sessionID string) (map[string]interface{}, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	session, ok := m.sessions[sessionID]
	if !ok {
		return nil, ErrSessionNotFound
	}
	// // update the position of the session in list and update the lastUsed time
	// m.sessions[sessionID].Value.(*Session).lastUsed = time.Now()
	// m.recentlyUsedSessions.MoveToFront(m.sessions[sessionID])

	return session.Value.(*Session).Data, nil
}

// UpdateSessionData overwrites the old session data with the new one
func (m *SessionManager) UpdateSessionData(sessionID string, data map[string]interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	_, ok := m.sessions[sessionID]
	if !ok {
		return ErrSessionNotFound
	}

	// Hint: you should renew expiry of the session here
	m.sessions[sessionID].Value.(*Session).lastUsed = time.Now()

	m.recentlyUsedSessions.MoveToFront(m.sessions[sessionID])

	return nil
}

func backgroundCleaner(sm *SessionManager) {
	for {
		sm.mu.Lock()
		if sm.recentlyUsedSessions != nil && sm.recentlyUsedSessions.Len() > 0 && time.Since(sm.recentlyUsedSessions.Back().Value.(*Session).lastUsed) >= 5*time.Second {
			session := sm.recentlyUsedSessions.Back().Value.(*Session)
			println("removing the session with id:", session.ID, " which is used:", time.Since(session.lastUsed))

			// remove the element from list
			sm.recentlyUsedSessions.Remove(sm.recentlyUsedSessions.Back())
			// remove the session from sessionManagerMap
			delete(sm.sessions, session.ID)

		} else if sm.recentlyUsedSessions != nil && sm.recentlyUsedSessions.Len() > 0 && time.Since(sm.recentlyUsedSessions.Back().Value.(*Session).lastUsed) < 5*time.Second {
			sm.recentlyUsedSessions.MoveToFront(sm.recentlyUsedSessions.Back())
		}
		sm.mu.Unlock()
	}
}

func main() {
	println("main function is called")

	// Create new sessionManager and new session
	m := NewSessionManager()
	sID, err := m.CreateSession()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Created new session with ID", sID)

	// Update session data
	data := make(map[string]interface{})
	data["website"] = "longhoang.de"

	err = m.UpdateSessionData(sID, data)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Update session data, set website to longhoang.de")

	// Retrieve data from manager again
	updatedData, err := m.GetSessionData(sID)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Get session data:", updatedData)
}
