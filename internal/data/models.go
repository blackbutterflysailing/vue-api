package data

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

const dbTimeout = time.Second * 3

var db *sql.DB

func New(dbPool *sql.DB) Models {
	db = dbPool

	// initialize db
	c, err := Init()
	if err != nil {
		fmt.Printf("models:New: err=%v \n", err)
		return Models{
			Calculated: Calculated{},
		}
	}
	fmt.Printf("models:New Sucess c=%v \n", *c)

	return Models{
		Calculated: *c,
	}
}

type Models struct {
	Calculated Calculated
}

type Calculated struct {
	ID        int       `json:"id"`
	Result    int       `json:"result"`
	CreatedAt time.Time `json:created_at"`
	UpdatedAt time.Time `json:updated_at"`
}

func (c *Calculated) GetCalculated() (*Calculated, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, result, created_at, updated_at from calculated`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var calculated Calculated

	for rows.Next() {
		fmt.Println(rows)
		err := rows.Scan(
			&calculated.ID,
			&calculated.Result,
			&calculated.CreatedAt,
			&calculated.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
	}
	return &calculated, nil
}

func (c *Calculated) Update() error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `update calculated set
		result = $1,
		updated_at = $2
		where id = $3
	`

	_, err := db.ExecContext(ctx, stmt,
		c.Result,
		time.Now(),
		c.ID,
	)

	if err != nil {
		fmt.Printf("Update(); Error updating record %v \t ID %d result %d \n", c, c.ID, c.Result)
		return err
	}

	return nil
}

func Init() (*Calculated, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	calc := &Calculated{
		ID:        1,
		Result:    0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	stmt := `insert into calculated(id, result, created_at, updated_at) values($1, $2, $3, $4)`

	_, err := db.ExecContext(ctx, stmt,
		calc.ID,
		calc.Result,
		calc.CreatedAt,
		calc.UpdatedAt,
	)

	if err != nil {
		fmt.Println("Init(): Error initializing database \n")
		return nil, err
	}

	return calc, nil
}
