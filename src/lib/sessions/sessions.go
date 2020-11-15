package sessions

import (
	"config"
	"fmt"
	. "interfaces"
	"lib/accounts"
	"lib/players"
	lib "lib/world"
	"logger"
	"sync"
	"time"
	"users/player"
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
			EndSession(s)
			break
		}
	}
}

func EndSession(s SessionInterface) {
	// TODO: Maybe add LogOut() methods to player/account structs?

	// Remove player from the room it is in
	room, _ := lib.GetRoom(s.Player().GetLocation())
	err := players.SyncPlayer(s.Player().(*player.Player))
	if err != nil {
		log.Err("failed to sync player (%s) data: %v", s.Player().GetName(), err)
	}
	room.RemovePlayer(s.Player().GetName())
	// #######
	// TODO: if players library is for ALL players and not just logged in players, remove this ##
	players.RemovePlayer(s.Player().GetName())
	// #######
	acct, err := accounts.GetAccount(s.Player().AccountName())
	if err != nil {
		log.Err("failed to get account while ending session for player (%s)", s.Player().GetName())
	}
	acct.SetLoggedInStatus(false)
	RemoveSessionByPlayerName(s.Player().GetName())
	s.Client().Logout()
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
