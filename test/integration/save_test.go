//go:build integration

package integration

import (
	"github.com/derivatan/si"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSaveToCreate(t *testing.T) {
	db := DB(t)
	ids := Seed(db, []Artist{
		{Name: "Roger Waters"},
	})
	c := Contact{
		Email:          "roger@waters.com",
		Phone:          1357924680,
		RadioFrequency: 21.789,
		LastCall:       time.Date(1943, time.September, 6, 5, 4, 3, 0, time.Local),
		OnSocialMedia:  true,
		ArtistID:       ids[0],
	}
	err := si.Save(db, &c)
	c2 := si.Query[Contact]().MustFind(db)

	assert.NoError(t, err)
	assert.NotNil(t, c2)
	assert.Equal(t, c.Email, c2.Email)
	assert.Equal(t, c.Phone, c2.Phone)
	assert.Equal(t, c.RadioFrequency, c2.RadioFrequency)
	assert.Equal(t, c.LastCall, c2.LastCall.Local())
	assert.Equal(t, c.OnSocialMedia, c2.OnSocialMedia)
	assert.Equal(t, c.ArtistID, c2.ArtistID)
}

func TestSaveToUpdate(t *testing.T) {
	db := DB(t)
	ids := Seed(db, []Artist{
		{Name: "Timbuktu"},
	})
	Seed(db, []Contact{
		{

			// TODO: Nullable values

			Email:          "Ett brev",
			Phone:          192837465,
			RadioFrequency: 98.54,
			LastCall:       time.Now(),
			OnSocialMedia:  false,
			ArtistID:       ids[0],
		},
	})

	email := "Det Löser sig"
	phone := 7592836
	radio := 73.11
	lastCall := time.Date(1234, 5, 6, 7, 8, 9, 0, time.Local)
	onSM := true
	c := si.Query[Contact]().MustFind(db)
	c.Email = email
	c.Phone = phone
	c.RadioFrequency = radio
	c.LastCall = lastCall
	c.OnSocialMedia = onSM
	err := si.Save(db, c)

	c2 := si.Query[Contact]().MustFind(db)

	assert.NoError(t, err)
	assert.Equal(t, email, c2.Email)
	assert.Equal(t, phone, c2.Phone)
	assert.Equal(t, radio, c2.RadioFrequency)
	assert.Equal(t, lastCall, c2.LastCall.Local())
	assert.Equal(t, onSM, c2.OnSocialMedia)
}

func TestUpdateWhenNotExists(t *testing.T) {
	db := DB(t)

	err := si.Update[Artist](db, &Artist{
		Name:     "Whatever",
		Nickname: "Who cares",
	}, []string{"name", "nickname"})

	assert.Error(t, err)
	assert.Equal(t, si.ResourceNotFoundError, err)
}

func TestUpdate(t *testing.T) {
	db := DB(t)
	artist := &Artist{
		Name:     "Aleks Christensen",
		Nickname: "Alex",
	}
	err := si.Save[Artist](db, artist)

	artist.Name = "Alex Christensen"
	artist.Nickname = "Aleks" // This should not update.
	err2 := si.Update[Artist](db, artist, []string{"name"})

	result, err3 := si.Query[Artist]().First(db)

	assert.NoError(t, err)
	assert.NoError(t, err2)
	assert.NoError(t, err3)
	assert.Equal(t, "Alex", result.Nickname)
	assert.Equal(t, "Alex Christensen", result.Name)
}

func TestInsertWithID(t *testing.T) {
	db := DB(t)
	artistID := uuid.MustParse("00001111-2222-3333-4444-555566667777")
	artist := &Artist{
		Model: si.Model{
			ID: &artistID,
		},
		Name:     "System of a Down",
		Nickname: "soad",
	}
	err := si.Insert[Artist](db, artist)
	result, err2 := si.Query[Artist]().First(db)

	assert.NoError(t, err)
	assert.NoError(t, err2)
	assert.Equal(t, artistID, *result.Model.ID)
}
