package main

import "time"

type Guest struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	RoomNumber  string    `json:"room_number"`
	DailyRate   int       `json:"daily_rate"`
	CheckInDate time.Time `json:"check_in_date"`
	Contact     string    `json:"contact"`
	WeeksStayed int       `json:"weeks_stayed"`
}
