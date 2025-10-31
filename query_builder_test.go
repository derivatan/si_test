//go:build integration

package integration

import (
	"testing"
	"time"

	"github.com/derivatan/si"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	db := DB(t)
	name := "Pink Floyd"
	Seed(db, []Artist{
		{Name: name},
	})

	list, err := si.Query[Artist]().Get(db)

	assert.NoError(t, err)
	assert.Len(t, list, 1)
	assert.Equal(t, name, list[0].Name)
}

func TestFirst(t *testing.T) {
	db := DB(t)
	name := "Ray Charles"
	Seed(db, []Artist{
		{Name: name},
		{Name: "Stevie Wonder"},
	})

	artist, err := si.Query[Artist]().OrderBy("name", true).First(db)

	assert.NoError(t, err)
	assert.NotNil(t, artist)
	assert.Equal(t, name, artist.Name)
}

func TestFind(t *testing.T) {
	db := DB(t)
	name := "Portishead"
	Seed(db, []Artist{
		{Name: name},
	})

	artist, err := si.Query[Artist]().Where("name", "=", name).Find(db)

	assert.NoError(t, err)
	assert.NotNil(t, artist)
	assert.Equal(t, name, artist.Name)
}

func TestFindWithID(t *testing.T) {
	db := DB(t)
	name := "Rammstein"
	ids := Seed(db, []Artist{
		{Name: name},
		{Name: "Dream Theater"},
	})

	artist, err := si.Query[Artist]().Find(db, ids[0])

	assert.NoError(t, err)
	assert.NotNil(t, artist)
	assert.Equal(t, name, artist.Name)
}

func TestWithWrongNumberOfResults(t *testing.T) {
	db := DB(t)
	Seed(db, []Artist{
		{Name: "Eminem"},
		{Name: "The Beatles"},
	})

	artist, err := si.Query[Artist]().Find(db)

	assert.Nil(t, artist)
	assert.Error(t, err)
	assert.ErrorIs(t, err, si.ResourceNotFoundError{})
}

// Builder functions

func TestSelect(t *testing.T) {
	db := DB(t)
	Seed(db, []Artist{
		{Name: "314"},
		{Name: "141"},
		{Name: "271"},
	})

	var numResults, max, min int
	selects := []string{
		"COUNT(1)",
		"MIN(name)",
		"MAX(name)",
	}
	_, err := si.Query[Artist]().Select(selects, func(scan func(...any)) {
		scan(&numResults, &min, &max)
	}).Get(db)

	assert.NoError(t, err)
	assert.Equal(t, 3, numResults)
	assert.Equal(t, 141, min)
	assert.Equal(t, 314, max)
}

func TestWhere(t *testing.T) {
	db := DB(t)
	wantedName := "Second"
	Seed(db, []Artist{
		{Name: "First"},
		{Name: wantedName},
		{Name: "Third"},
	})

	list, err := si.Query[Artist]().Where("name", "=", wantedName).Get(db)

	assert.NoError(t, err)
	assert.Len(t, list, 1)
	assert.Equal(t, list[0].Name, wantedName)
}
func TestWhereContains(t *testing.T) {
	db := DB(t)
	name := "Beethoven"
	Seed(db, []Artist{
		{Name: name},
		{Name: "Mozart"},
	})

	list, err := si.Query[Artist]().Where("name", "LIKE", "%ee%").Get(db)

	assert.NoError(t, err)
	assert.Len(t, list, 1)
	assert.Equal(t, list[0].Name, name)
}

func TestOrWhere(t *testing.T) {
	db := DB(t)
	name1 := "Prince"
	name2 := "Queen"
	Seed(db, []Artist{
		{Name: name1},
		{Name: name2},
		{Name: "Michael Jackson"},
	})

	list, err := si.Query[Artist]().
		Where("name", "=", name1).
		OrWhere("name", "=", name2).
		OrderBy("name", true).
		Get(db)

	assert.NoError(t, err)
	assert.Len(t, list, 2)
	assert.Equal(t, name1, list[0].Name)
	assert.Equal(t, name2, list[1].Name)
}

func TestWhereF(t *testing.T) {
	db := DB(t)
	Seed(db, []Artist{
		{Name: "Danny Elfman"},
		{Name: "Hans Zimmer"},
		{Name: "John Williams"},
	})

	// WHERE a AND (b OR c)
	list, err := si.Query[Artist]().Where("name", "ILIKE", "%m%").WhereF(func(q *si.Q[Artist]) *si.Q[Artist] {
		return q.Where("name", "ILIKE", "%Zi%").OrWhere("name", "ILIKE", "%Wi%")
	}).Get(db)

	assert.NoError(t, err)
	assert.Len(t, list, 2)
}

func TestOrWhereF(t *testing.T) {
	db := DB(t)
	Seed(db, []Artist{
		{Name: "Björk"},
		{Name: "Daft Punk"},
		{Name: "The Knife"},
	})

	// WHERE a OR (b AND c)
	list, err := si.Query[Artist]().Where("name", "ILIKE", "%knife%").OrWhereF(func(q *si.Q[Artist]) *si.Q[Artist] {
		return q.Where("name", "ILIKE", "%daft%").Where("name", "ILIKE", "%punk%")
	}).Get(db)

	assert.NoError(t, err)
	assert.Len(t, list, 2)
}

func TestOrderBy(t *testing.T) {
	db := DB(t)
	nameA := "Avalanches, The"
	nameB := "Basement Jaxx"
	nameC := "Cure, The"
	nameD := "Deep Purple"
	Seed(db, []Artist{
		{Name: nameB},
		{Name: nameC},
		{Name: nameA},
		{Name: nameD},
	})

	listAsc, errAsc := si.Query[Artist]().OrderBy("name", true).Get(db)
	listDesc, errDesc := si.Query[Artist]().OrderBy("name", false).Get(db)

	assert.NoError(t, errAsc)
	assert.NoError(t, errDesc)
	assert.Equal(t, nameA, listAsc[0].Name)
	assert.Equal(t, nameB, listAsc[1].Name)
	assert.Equal(t, nameC, listAsc[2].Name)
	assert.Equal(t, nameD, listAsc[3].Name)
	assert.Equal(t, nameD, listDesc[0].Name)
	assert.Equal(t, nameC, listDesc[1].Name)
	assert.Equal(t, nameB, listDesc[2].Name)
	assert.Equal(t, nameA, listDesc[3].Name)
}

func TestTakeAndSkip(t *testing.T) {
	db := DB(t)
	name1 := "Detektivbyrån"
	name2 := "Trazan & Banarne"
	name3 := "Electric Banana Band"
	Seed(db, []Artist{
		{Name: name1},
		{Name: name2},
		{Name: name3},
	})

	listTake, errTake := si.Query[Artist]().OrderBy("name", true).Take(2).Get(db)
	listSkip, errSkip := si.Query[Artist]().OrderBy("name", true).Skip(1).Get(db)

	assert.NoError(t, errTake)
	assert.NoError(t, errSkip)
	assert.Len(t, listTake, 2)
	assert.Equal(t, name1, listTake[0].Name)
	assert.Equal(t, name3, listTake[1].Name)
	assert.Len(t, listSkip, 2)
	assert.Equal(t, name3, listSkip[0].Name)
	assert.Equal(t, name2, listSkip[1].Name)
}

func TestGroupBy(t *testing.T) {
	db := DB(t)
	Seed(db, []Contact{
		{Email: "info@email.com", Phone: 101},
		{Email: "info@email.com", Phone: 103},
		{Email: "support@email.com", Phone: 107},
		{Email: "support@email.com", Phone: 109},
	})

	type result struct {
		Email string
		Sum   int
	}
	var results []result
	_, err := si.Query[Contact]().Select(
		[]string{"email", "SUM(phone)"},
		func(scan func(...any)) {
			var r result
			scan(&r.Email, &r.Sum)
			results = append(results, r)
		},
	).GroupBy("email").OrderBy("email", true).Get(db)

	var havingResult []result
	_, havingErr := si.Query[Contact]().Select(
		[]string{"email", "SUM(phone)"},
		func(scan func(...any)) {
			var r result
			scan(&r.Email, &r.Sum)
			havingResult = append(havingResult, r)
		},
	).GroupBy("email").OrderBy("email", true).Having("SUM(phone)", ">", 210).Get(db)

	assert.NoError(t, err)
	assert.Len(t, results, 2)
	assert.Equal(t, "info@email.com", results[0].Email)
	assert.Equal(t, 204, results[0].Sum)
	assert.Equal(t, "support@email.com", results[1].Email)
	assert.Equal(t, 216, results[1].Sum)

	assert.NoError(t, havingErr)
	assert.Len(t, havingResult, 1)
	assert.Equal(t, "support@email.com", havingResult[0].Email)
	assert.Equal(t, 216, havingResult[0].Sum)
}

func TestRelationHasOne(t *testing.T) {
	db := DB(t)
	name := "Yann Tiersen"
	email := "yann@tiersen.com"
	artistID := Seed(db, []Artist{
		{Name: name},
	})
	Seed(db, []Contact{
		{Email: email, Phone: 123, ArtistID: artistID[0]},
	})

	artist, err := si.Query[Artist]().First(db)
	contact, contactErr := artist.Contact().First(db)

	assert.NoError(t, err)
	assert.NoError(t, contactErr)
	assert.NotNil(t, artist)
	assert.NotNil(t, contact)
	assert.Equal(t, artist.Name, name)
	assert.Equal(t, email, contact.Email)
}

func TestRelationBelongsTo(t *testing.T) {
	db := DB(t)
	albumName := "Brand New Day"
	artistName := "Sting"
	ids := Seed(db, []Artist{
		{Name: artistName},
	})
	Seed(db, []Album{
		{Name: albumName, ArtistID: ids[0]},
	})

	album, err := si.Query[Album]().Find(db)
	artist, artistErr := album.Artist().Find(db)

	assert.NoError(t, err)
	assert.NoError(t, artistErr)
	assert.NotNil(t, album)
	assert.NotNil(t, artist)
	assert.Equal(t, albumName, album.Name)
	assert.Equal(t, artistName, artist.Name)
}

func TestRelationHasMany(t *testing.T) {
	db := DB(t)
	ids := Seed(db, []Artist{
		{Name: "Muse"},
		{Name: "Xploding Plastix"},
	})
	Seed(db, []Album{
		{Name: "The Resistance", ArtistID: ids[0]},
		{Name: "Black Holes And Revelations", ArtistID: ids[0]},
		{Name: "Amateur Girlfriends", ArtistID: ids[1]},
	})

	artist, err := si.Query[Artist]().Find(db, ids[0])
	albums, albErr := artist.Albums().Get(db)

	assert.NoError(t, err)
	assert.NoError(t, albErr)
	assert.Len(t, albums, 2)
}

func TestRelationWithHasOne(t *testing.T) {
	db := DB(t)
	name := "Kraftwerk"
	ids := Seed(db, []Artist{
		{Name: name},
	})
	email := "robots@autobahn"
	Seed(db, []Contact{
		{Email: email, ArtistID: ids[0]},
	})

	artist, err := si.Query[Artist]().With(func(m Artist, r []Artist) error {
		return m.Contact().Execute(db, r)
	}).First(db)
	// db is not needed here since it's already loaded during the 'with' above.
	contact := artist.Contact().MustFind(nil)

	assert.NoError(t, err)
	assert.Equal(t, name, artist.Name)
	assert.Equal(t, email, contact.Email)
}

func TestRelationWithBelongsTo(t *testing.T) {
	db := DB(t)
	name := "Dire staits"
	albName := "Sultans of Swing"
	ids := Seed(db, []Artist{
		{Name: name},
	})
	Seed(db, []Album{
		{Name: albName, ArtistID: ids[0]},
	})

	album, err := si.Query[Album]().With(func(m Album, r []Album) error {
		return m.Artist().Execute(db, r)
	}).First(db)
	// db is not needed here since it's already loaded during the 'with' above.
	artist := album.Artist().MustFirst(nil)

	assert.NoError(t, err)
	assert.Equal(t, albName, album.Name)
	assert.Equal(t, name, artist.Name)
}

func TestRelationWithHasMany(t *testing.T) {
	db := DB(t)
	ids := Seed(db, []Artist{
		{Name: "Metallica"},
	})
	Seed(db, []Album{
		{Name: "Master of puppets", ArtistID: ids[0]},
		{Name: "Ride the Lightning", ArtistID: ids[0]},
	})

	artist, err := si.Query[Artist]().With(func(m Artist, r []Artist) error {
		return m.Albums().Execute(db, r)
	}).First(db)
	// db is not needed here since it's already loaded during the 'with' above.
	albums := artist.Albums().MustGet(nil)

	assert.NoError(t, err)
	assert.Len(t, albums, 2)
}

func TestLoaded(t *testing.T) {
	db := DB(t)
	ids := Seed(db, []Artist{
		{Name: "Vivaldi"},
	})
	albumName := "Le quattro stagioni"
	Seed(db, []Album{
		{Name: albumName, ArtistID: ids[0]},
	})

	artist1, err1 := si.Query[Artist]().First(db)
	artist2, err2 := si.Query[Artist]().With(func(m Artist, r []Artist) error {
		return m.Albums().Execute(db, r)
	}).First(db)

	album := artist2.Albums().MustFirst(nil)

	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.False(t, artist1.Albums().Loaded())
	assert.True(t, artist2.Albums().Loaded())
	assert.Equal(t, album.Name, albumName)
}

func TestJoin(t *testing.T) {
	db := DB(t)
	ids := Seed(db, []Artist{
		{Name: "The Ark"},
		{Name: "The Cranberries"},
		{Name: "Earth, Wind & Fire"},
	})
	Seed(db, []Album{
		{Name: "We Are The Ark", ArtistID: ids[0]},
		{Name: "No Need To Argue", ArtistID: ids[1]},
		{Name: "All 'N All", ArtistID: ids[2]},
	})

	albums, err := si.Query[Album]().Join(func(t Album) *si.JoinConf {
		return t.Artist().Join(si.INNER)
	}).Where("artists.name", "ILIKE", "%The%").Get(db)

	assert.NoError(t, err)
	assert.Len(t, albums, 2)
}

func TestWithDeleted(t *testing.T) {
	db := DB(t)
	si.UseDeletedAt(true)
	name := "Jean-Michel Jarre"
	now := time.Now()
	Seed(db, []Artist{
		{Name: "Kate Bush", Model: si.Model{DeletedAt: &now}},
		{Name: name},
	})
	list, err1 := si.Query[Artist]().Get(db)
	list2, err2 := si.Query[Artist]().WithDeleted().Get(db)

	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.Len(t, list, 1)
	assert.Equal(t, name, list[0].Name)
	assert.Len(t, list2, 2)
}

func TestJoinWithDeleted(t *testing.T) {
	db := DB(t)
	si.UseDeletedAt(true)
	now := time.Now()
	name := "Thousand Sun Sky"
	ids := Seed(db, []Artist{
		{Name: "Infected Mushroom"},
		{Name: name},
	})
	Seed(db, []Album{
		{Name: "B.P.Empire", ArtistID: ids[0]},
		{Name: "Head of NASA and the 2 Amish Boys", ArtistID: ids[0], Model: si.Model{DeletedAt: &now}},
		{Name: "The Aurora Complex", ArtistID: ids[1]},
		{Name: "Passengers", ArtistID: ids[1]},
	})

	list, err := si.Query[Artist]().Join(func(t Artist) *si.JoinConf {
		return t.Albums().Join(si.INNER)
	}).Where("albums.name", "ILIKE", "%the%").Get(db)

	assert.NoError(t, err)
	assert.Len(t, list, 1)
	assert.Equal(t, name, list[0].Name)
}
