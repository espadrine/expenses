package expense

import (
	"crypto/rand"
	"encoding/base32"
	"strings"
)

func (s *Store) genBase32ID(table string) (string, error) {
	foundUniqueID := false
	var id string

	for !foundUniqueID {
		id = genRandBase32ID()
		// If there is already an entity with this ID, try again.
		row := s.db.QueryRow("select count(*) from "+table+" where id = ? limit 1", id)
		var count int
		err := row.Scan(&count)
		if err != nil {
			return "", err
		}
		foundUniqueID = count == 0
	}

	return id, nil
}

func genRandBase32ID() string {
	// Generate a random byte slice. 16×8 = 128 bits.
	b := make([]byte, 16)
	rand.Read(b)

	// Encode the byte slice to lowercase base32 with no padding.
	rawId := base32.StdEncoding.EncodeToString(b)
	id := strings.ToLower(strings.TrimRight(rawId, "="))

	return id
}
