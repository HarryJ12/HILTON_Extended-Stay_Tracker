package main

import "time"

type Guest struct {
	ID int `json:"id"`

	Name         string     `json:"name"`
	RoomNumber   string     `json:"room_number"`
	MonthlyRate  int        `json:"monthly_rate"`
	CheckInDate  time.Time  `json:"check_in_date"`
	Contact      string     `json:"contact"`
	LastReminder *time.Time `json:"last_reminder"`

	MonthsStayed int `json:"months_stayed"`
}
