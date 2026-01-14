package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func runAgent() {
	agentDB, err := sql.Open("sqlite3", "data.db")
	if err != nil {
		log.Println("Agent error: cannot open db")
		return
	}
	defer agentDB.Close()

	agentDB.Exec("PRAGMA journal_mode=WAL;")
	agentDB.Exec("PRAGMA busy_timeout = 5000;")

	now := time.Now().UTC()

	rows, err := agentDB.Query(`
		SELECT id, name, room_number, daily_rate, check_in_date, contact
		FROM guests
	`)
	if err != nil {
		log.Println("Agent error: read failed")
		return
	}
	defer rows.Close()

	type candidate struct {
		id      int
		name    string
		room    string
		rate    int
		week    int
		contact string
	}

	var toNotify []candidate

	for rows.Next() {
		var (
			id      int
			name    string
			room    string
			rate    int
			checkIn time.Time
			contact string
		)

		if err := rows.Scan(
			&id,
			&name,
			&room,
			&rate,
			&checkIn,
			&contact,
		); err != nil {
			continue
		}

		weeks := weeksStayed(checkIn, now)
		if weeks < 1 {
			continue
		}

		toNotify = append(toNotify, candidate{
			id:      id,
			name:    name,
			room:    room,
			rate:    rate,
			week:    weeks,
			contact: contact,
		})
	}

	for _, g := range toNotify {
		res, err := agentDB.Exec(`
			INSERT OR IGNORE INTO notifications (guest_id, period_number)
			VALUES (?, ?)
		`, g.id, g.week)

		if err != nil {
			log.Println("Agent error: insert failed")
			continue
		}

		affected, _ := res.RowsAffected()
		if affected != 1 {
			continue
		}

		subject := "Extended-Stay Guest Billing Reminder"

		body := fmt.Sprintf(
			"Weekly Billing Reminder for: \n"+"%s\n\n"+
				"Room: %s\n"+
				"Weeks Stayed: %d\n"+
				"Daily Rate: $%d\n"+
				"Contact Information: %s\n",
			g.name,
			g.room,
			g.week,
			g.rate,
			g.contact,
		)

		if err := sendEmail(subject, body); err != nil {
			log.Println("Agent error: email send failed")
			continue
		}
	}
}
