package db

import "time"

type User struct {
	ID           string    `json:"id"`        // The user's ID.
	SID          string    `json:"-"`         // The user's session ID.
	Name         string    `json:"name"`      // The user's display name.
	PasswordHash string    `json:"-"`         // The user's password, hashed.
	Score        int       `json:"score"`     // Total upvotes the user received minutes total downvotes.
	Bio          string    `json:"bio"`       // The user's biography.
	Joined       time.Time `json:"joined"`    // The time when the user joined.
	IsAdmin      bool      `json:"is_admin"`  // If the user is an instance administrator.
	IsBanned     bool      `json:"is_banned"` // If the user is banned.
	Items        []string  `json:"items"`     // A list of the items the user has created.
}

type Item struct {
	ID        string    `json:"id"`        // The ID of the item.
	Creator   string    `json:"creator"`   // The ID of the User who made the Item.
	Points    int       `json:"points"`    // The amount of upvotes minus the amount of downvotes.
	Title     string    `json:"title"`     // The title of the Item. Only applicable for posts.
	Content   string    `json:"content"`   // The text in the item.
	Parent    string    `json:"parent"`    // The ID of the parent of the Item. Only applicable for comments.
	Children  []string  `json:"children"`  // The IDs of the replies to the Item.
	Published time.Time `json:"published"` // The time when the Item was published.
}
