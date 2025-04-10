package models

type Schedule struct {
	ID                            int  `db:"id"`
	IsOfflineNotificationsAllowed bool `db:"is_offline_notifications_allowed"`
}
