package main

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func deleteGuest(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		_, err := db.Exec("DELETE FROM notifications WHERE guest_id = ?", id)
		if err != nil {
			c.JSON(500, gin.H{"error": "db error"})
			return
		}

		res, err := db.Exec("DELETE FROM guests WHERE id = ?", id)
		if err != nil {
			c.JSON(500, gin.H{"error": "db error"})
			return
		}

		rows, _ := res.RowsAffected()
		if rows == 0 {
			c.JSON(404, gin.H{"error": "guest not found"})
			return
		}

		c.Status(204)
	}
}

func getGuests(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		rows, err := db.Query(`
			SELECT g.id, g.name, g.room_number, g.monthly_rate,
			       g.check_in_date, g.contact
			FROM guests g
			LEFT JOIN notifications n ON g.id = n.guest_id
			GROUP BY g.id
		`)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
			return
		}
		defer rows.Close()

		now := time.Now().UTC()
		var guests []Guest

		for rows.Next() {
			var g Guest

			err := rows.Scan(
				&g.ID, &g.Name, &g.RoomNumber, &g.MonthlyRate,
				&g.CheckInDate, &g.Contact,
			)
			if err != nil {
				continue
			}

			g.MonthsStayed = monthsStayed(g.CheckInDate, now)

			guests = append(guests, g)
		}

		c.JSON(http.StatusOK, guests)

	}
}

func createGuest(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			Name        string `json:"name"`
			Contact     string `json:"contact"`
			RoomNumber  string `json:"room_number"`
			MonthlyRate int    `json:"monthly_rate"`
			CheckInDate string `json:"check_in_date"`
		}

		if err := c.BindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
			return
		}

		checkIn, err := time.Parse("2006-01-02", input.CheckInDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date"})
			return
		}

		_, err = db.Exec(
			`INSERT INTO guests (name, room_number, monthly_rate, check_in_date, contact)
			 VALUES (?, ?, ?, ?, ?)`,
			input.Name,
			input.RoomNumber,
			input.MonthlyRate,
			checkIn.UTC(),
			input.Contact,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
			return
		}

		c.Status(http.StatusCreated)
		runAgent()
	}
}

func monthsStayed(checkIn time.Time, now time.Time) int {
	yearDiff := now.Year() - checkIn.Year()
	monthDiff := int(now.Month()) - int(checkIn.Month())

	months := yearDiff*12 + monthDiff

	if now.Day() < checkIn.Day() {
		months--
	}

	if months < 0 {
		return 0
	}

	return months
}
