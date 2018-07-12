package entity

// NotificationPreference contains flags indicating whether to send SMS / Email to the user,
// and determine what kind of messages will be sent.
type NotificationPreference struct {
	EnableSMS   bool `json:"enableSMS" bson:"enableSMS"`
	EnableEmail bool `json:"enableEmail" bson:"enableEmail"`
}
