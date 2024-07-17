//go:build integration

package integration

import (
	"github.com/derivatan/si"
	"github.com/google/uuid"
	"time"
)

type Contact struct {
	si.Model

	Email          string
	Phone          int
	RadioFrequency float64
	LastCall       time.Time
	OnSocialMedia  bool
	ArtistID       uuid.UUID

	artist si.RelationData[Artist]
}

func (c Contact) GetModel() si.Model {
	return c.Model
}

func (c Contact) GetTable() string {
	return "contacts"
}

func (c Contact) Artist() *si.Relation[Contact, Artist] {
	return si.BelongsTo[Contact, Artist](c, "ArtistID", "artist", func(c *Contact) *si.RelationData[Artist] {
		return &c.artist
	})
}

type Artist struct {
	si.Model

	Name string

	contact si.RelationData[Contact]
	albums  si.RelationData[Album]
}

func (a Artist) GetModel() si.Model {
	return a.Model
}

func (a Artist) GetTable() string {
	return "artists"
}

func (a Artist) Contact() *si.Relation[Artist, Contact] {
	return si.HasOne[Artist, Contact](a, "ArtistID", "contact", func(a *Artist) *si.RelationData[Contact] {
		return &a.contact
	})
}

func (a Artist) Albums() *si.Relation[Artist, Album] {
	return si.HasMany[Artist, Album](a, "ArtistID", "albums", func(a *Artist) *si.RelationData[Album] {
		return &a.albums
	})
}

type Album struct {
	si.Model

	Name     string
	Year     int
	ArtistID uuid.UUID

	artist si.RelationData[Artist]
}

func (a Album) GetModel() si.Model {
	return a.Model
}

func (a Album) GetTable() string {
	return "albums"
}

func (a Album) Artist() *si.Relation[Album, Artist] {
	return si.BelongsTo[Album, Artist](a, "ArtistID", "artist", func(a *Album) *si.RelationData[Artist] {
		return &a.artist
	})
}
