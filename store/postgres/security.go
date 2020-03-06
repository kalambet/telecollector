package postgres

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/kalambet/telecollector/telegram"

	"github.com/kalambet/telecollector/telecollector"
)

const (
	createAllowances = `
create table allowances(
    chat_id bigint, 
    modified date,
    follow bool, 
    primary key(chat_id)
);`

	queryAllowances = `select chat_id, follow from allowances;`

	insertAllowance = `
insert into 
    allowances (chat_id, follow, modified) 
    values ($1, $2, $3) 
        on conflict (chat_id) 
        do update set follow = $2, modified = $3;`
)

type credentialsService struct {
	Allowances map[int64]*telecollector.Allowance
	Admin      map[int64]bool
}

func NewCredentialService() (telecollector.CredentialService, error) {
	err := gracefulCreateTable("allowances", createAllowances)
	if err != nil {
		return nil, err
	}

	cs := &credentialsService{}
	err = cs.loadAllowances()
	if err != nil {
		return nil, err
	}

	err = cs.loadAdmins()
	if err != nil {
		return nil, err
	}

	return cs, nil
}

func (cs *credentialsService) loadAllowances() error {
	rows, err := db.Query(queryAllowances)

	if err != nil {
		return err
	}

	cs.Allowances = make(map[int64]*telecollector.Allowance)
	for rows.Next() {
		a := telecollector.Allowance{}
		if err := rows.Scan(&a.ChatID, &a.Follow); err != nil {
			log.Printf("postgres: error unmarshaling allowance query result: %s", err.Error())
			continue
		}
		cs.Allowances[a.ChatID] = &a
	}

	return rows.Close()
}

func (cs *credentialsService) loadAdmins() error {
	adminStr := os.Getenv("BOT_ADMINS")
	if len(adminStr) == 0 {
		return nil
	}

	admins := strings.Split(adminStr, ",")

	if len(admins) == 0 {
		return nil
	}

	cs.Admin = make(map[int64]bool)
	for _, a := range admins {
		id, err := strconv.Atoi(a)
		if err != nil {
			log.Printf("postgres: parse admin error: %s", err.Error())
			continue
		}
		cs.Admin[int64(id)] = true
	}

	return nil
}

func (cs *credentialsService) CheckAdmin(authorID int64) bool {
	exist, ok := cs.Admin[authorID]
	return ok && exist
}

func (cs *credentialsService) CheckChat(chatID int64) bool {
	a, ok := cs.Allowances[chatID]
	return ok && a.Follow
}

func (cs *credentialsService) FollowChat(chat *telegram.Chat, follow bool) error {
	_, ok := cs.Allowances[chat.ID]
	if ok {
		cs.Allowances[chat.ID].Follow = follow
	} else {
		cs.Allowances[chat.ID] = &telecollector.Allowance{
			ChatID: chat.ID,
			Follow: follow,
		}
	}

	_, err := db.Exec(insertAllowance, &chat.ID, &follow, time.Now())
	if err != nil {
		return err
	}

	return nil
}
