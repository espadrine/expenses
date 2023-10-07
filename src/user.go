package expense

type User struct {
	id   string
	name string
}

func (s *Store) getUsers() (users []User, err error) {
	rows, err := s.db.Query("select * from users")
	if err != nil {
		return []User{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var id, name string
		err = rows.Scan(&id, &name)
		users = append(users, User{id, name})
		if err != nil {
			return []User{}, err
		}
	}
	return users, nil
}

func (s *Store) getUser(userID string) (*User, error) {
	row := s.db.QueryRow("select id, name from users where id = ? limit 1", userID)
	var id, name string
	err := row.Scan(&id, &name)
	if err != nil {
		return nil, err
	}
	return &User{id, name}, nil
}

func (s *Store) createUser(username string) (*User, error) {
	id, err := s.genBase32ID("users")
	if err != nil {
		return nil, err
	}
	_, err = s.db.Exec("insert into users (id, name) values (?, ?)", id, username)
	if err != nil {
		return nil, err
	}
	return &User{id, username}, nil
}
