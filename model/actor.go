package model

import (
	"database/sql"
	log "github.com/sirupsen/logrus"
)

type Actor struct {
	Id       int64
	Code     string `json:"code"`
	FullName string `json:"full_name"`
}
type ActorCollection struct {
	Actors []Actor `json:"actors"`
}

type ActorModel struct {
	DB *sql.DB
}

func (a *ActorModel) Add(actor *Actor) error {
	Query := `INSERT INTO actors 
        (code, name) 
        VALUES (?, ?)
        `

	stmt, err := a.DB.Prepare(Query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(
		actor.Code,
		actor.FullName,
	)

	if err != nil {
		return err
	}

	actor.Id, _ = result.LastInsertId()

	return nil
}

func (a *ActorModel) Update(actor *Actor) error {
	Query := `UPDATE 
	actors a
       SET a.name = ?,
       WHERE a.id = ?
       `

	stmt, err := a.DB.Prepare(Query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		actor.FullName,
	)

	if err != nil {
		return err
	}

	return nil
}

func (a *ActorModel) Find(code string) (*Actor, error) {
	Query := `SELECT
	a.id,
	a.code,
	a.name,
FROM
	actors a
WHERE 
	a.code = ?
	AND d.deleted_at IS NULL`

	stmt, err := a.DB.Prepare(Query)
	defer stmt.Close()
	if err != nil {
		log.Errorf(err.Error())

		return nil, err
	}

	actor := &Actor{}
	err = stmt.QueryRow(code).Scan(
		&actor.Id,
		&actor.Code,
		&actor.FullName,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		log.Println(err.Error())

		return nil, err
	}

	return actor, nil
}

func (a *ActorModel) FindAll() (*ActorCollection, error) {
	Query := `SELECT
	a.id,
	a.code,
	a.name,
FROM
	actors a
WHERE 
	d.deleted_at IS NULL`

	collection := ActorCollection{}

	stmt, err := a.DB.Prepare(Query)
	defer stmt.Close()
	if err != nil {
		return &collection, err
	}

	rows, err := stmt.Query()
	if err != nil {
		return &collection, err
	}
	defer rows.Close()

	for rows.Next() {
		actor := Actor{}
		err := rows.Scan(
			&actor.Id,
			&actor.Code,
			&actor.FullName,
		)

		if err != nil {
			return &collection, err
		}

		collection.Actors = append(collection.Actors, actor)
	}

	err = rows.Err()
	if err != nil {
		return &collection, err
	}

	return &collection, nil
}
