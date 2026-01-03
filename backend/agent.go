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

	// Prevent SQLITE_BUSY garbage
	agentDB.Exec("PRAGMA journal_mode=WAL;")
	agentDB.Exec("PRAGMA busy_timeout = 5000;")

	now := time.Now().UTC()

	rows, err := agentDB.Query(`
		SELECT id, name, room_number, monthly_rate, check_in_date, contact
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
		month   int
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

		months := monthsStayed(checkIn, now)
		if months < 1 {
			continue
		}

		toNotify = append(toNotify, candidate{
			id:      id,
			name:    name,
			room:    room,
			rate:    rate,
			month:   months,
			contact: contact,
		})
	}

	for _, g := range toNotify {
		res, err := agentDB.Exec(`
			INSERT OR IGNORE INTO notifications (guest_id, month_number)
			VALUES (?, ?)
		`, g.id, g.month)

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
			"Monthly billing reminder\n\n"+
				"Guest: %s\n"+
				"Room: %s\n"+
				"Month: %d\n"+
				"Rate: $%d\n"+
				"Contact Information: %s\n",
			g.name,
			g.room,
			g.month,
			g.rate,
			g.contact,
		)

		if err := sendEmail(subject, body); err != nil {
			log.Println("Agent error: email send failed")
			continue
		}
	}
}
