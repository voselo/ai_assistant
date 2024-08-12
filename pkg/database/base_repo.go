package db

import (
	"errors"
	"fmt"
	"strings"

	"github.com/doug-martin/goqu/v9"
	"github.com/jmoiron/sqlx"
)

const (
	pgErrDuplicateCode = "SQLSTATE 23505"
)

var ErrDuplicate = errors.New("duplicate record")

type BaseRepo struct {
	Db    *sqlx.DB
	Table string
}

func (repo *BaseRepo) Create(newModel interface{}) (string, error) {
	query, _, err := goqu.Insert(repo.Table).Rows(
		newModel,
	).Returning("id").ToSQL()
	if err != nil {
		return "", fmt.Errorf("error build query create: %w", err)
	}

	var id string
	if err := repo.Db.QueryRowx(query).Scan(&id); err != nil {
		if strings.Contains(err.Error(), pgErrDuplicateCode) {
			return "", ErrDuplicate
		}
		return "", fmt.Errorf("can't create: %w", err)
	}

	return id, nil
}

func (rerpo *BaseRepo) GetQuery(id string) (string, error) {
	qb := goqu.From(rerpo.Table)
	qb = qb.Where(goqu.I("id").Eq(id))

	query, _, err := qb.ToSQL()
	if err != nil {
		return "", fmt.Errorf("can't build query to get: %w", err)
	}

	return query, nil
}

func (repo *BaseRepo) Delete(id string) error {
	ds := goqu.From(repo.Table).Delete().
		Where(goqu.I("id").Eq(id))

	sql, _, err := ds.ToSQL()
	if err != nil {
		return err
	}

	if _, err := repo.Db.Exec(sql); err != nil {
		return fmt.Errorf("error can't delete %w", err)
	}

	return nil
}

func (repo *BaseRepo) Update(updateModel interface{}, id string) error {
	qu, _, err := goqu.Update(repo.Table).Set(
		updateModel,
	).Returning("id").Where(goqu.I("id").Eq(id)).ToSQL()
	if err != nil {
		return err
	}
	if _, err := repo.Db.Exec(qu); err != nil {
		if strings.Contains(err.Error(), pgErrDuplicateCode) {
			return ErrDuplicate
		}
		return fmt.Errorf("error can't update updateModel %w", err)
	}

	return nil
}

func (repo *BaseRepo) GetID(id string, getModel interface{}) error {
	q, err := repo.GetQuery(id)
	if err != nil {
		return err
	}
	if err := repo.Db.Get(getModel, q); err != nil {
		return fmt.Errorf("can't get model: %w", err)
	}

	return nil
}
