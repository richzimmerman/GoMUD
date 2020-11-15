package sessions

import (
	"config"
	"fmt"
	. "interfaces"
	"logger"
	"sync"
	"time"
	"utils"
)

var log = logger.NewLogger()

const (
	afkTimerConfig = "AFKTIMER"
)

var Sessions = &sessions{
	sessions: make([]*Session, 0),
}

type sessions struct {
	mu       sync.Mutex
	sessions []*Session
}

type Session struct {
	client    ClientInterface
	player    PlayerInterface
	lastInput int64
}

func (s *Session) Player() PlayerInterface {
	return s.player
}

func (s *Session) Client() ClientInterface {
	return s.client
}

func (s *Session) LastInput() int64 {
	return s.lastInput
}

func (s *Session) InputReceived(timestamp int64) {
	s.lastInput = timestamp
}

func (s *Session) EndSession() {
	return
}

func (s *Session) startAfkTimer() {
	timeoutMinutes, err := config.GetIntegerValue(afkTimerConfig)
	if err != nil {
		log.Err("%v", err)
		return
	}
	timeout := time.Minute * time.Duration(timeoutMinutes)
	for {
		last := time.Unix(s.lastInput, 0)
		now := time.Now()
		diff := now.Sub(last)
		if diff >= timeout {
			// Diff duration between now and last input should not be greater than Duration of timeout minutes
			//  if so, end the session
			s.Client().Out("You've been AFK for too long. Good bye.")
			s.EndSession()
		}
	}
}

func CreateSession(c ClientInterface, p PlayerInterface) {
	Sessions.mu.Lock()
	now := time.Now()
	s := &Session{
		client:    c,
		player:    p,
		lastInput: now.Unix(),
	}
	go s.startAfkTimer()
	Sessions.sessions = append(Sessions.sessions, s)
	Sessions.mu.Unlock()
}

func RemoveSessionByIpAddress(ipaddr string) error {
	_, index, err := GetSessionByIpAddress(ipaddr)
	if err != nil {
		return utils.Error(err)
	}
	removeSessionByIndex(index)
	return nil
}

func RemoveSessionByPlayerName(name string) error {
	_, index, err := GetSessionByPlayerName(name)
	if err != nil {
		return utils.Error(err)
	}
	removeSessionByIndex(index)
	return nil
}

func removeSessionByIndex(i int) {
	// Effectively we're just copying the last element to the desired index to remove, and truncating the slice to exclude
	// the last index. Since order doesn't matter, this is fine and seems to be more efficient
	Sessions.mu.Lock()
	l := len(Sessions.sessions) - 1
	Sessions.sessions[i] = Sessions.sessions[l]
	Sessions.sessions[l] = nil
	Sessions.sessions = Sessions.sessions[:l]
	Sessions.mu.Unlock()
}

func GetSessionByPlayerName(name string) (*Session, int, error) {
	for i, s := range Sessions.sessions {
		if s.player.GetName() == name {
			return s, i, nil
		}
	}
	return nil, -1, fmt.Errorf("unable to find session for player (%s)", name)
}

//  **  Not sure if we want this, if we allow multiple simultaneous connections from the same account  **
// func GetSessionByAccountName(name string) (*session, int, error) {
// 	for i, s := range Sessions.sessions {
// 		if s.player.AccountName() == name {
// 			return s, i, nil
// 		}
// 	}
// 	return nil, -1, fmt.Errorf("unable to find session by account name (%s)", name)
// }

func GetSessionByIpAddress(ipaddr string) (*Session, int, error) {
	for i, s := range Sessions.sessions {
		if s.client.GetRemoteAddress() == ipaddr {
			return s, i, nil
		}
	}
	return nil, -1, fmt.Errorf("unable to find session for remote address (%s)", ipaddr)
}
