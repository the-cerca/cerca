package models

import (
	"database/sql"
	"time"
)

type Transaction struct {
	ID          string
	Emitter     User
	Beneficiary User
	Amount float64
	Date time.Time
	Status string 
}

type Payment struct {
	ID string
	IDTransaction string
	Amount float64
	Date time.Time
}
type TransactionManager  struct {
	DB *sql.DB
}

const(
	waiting string = "waiting"
	done string = "done"
	suspend string = "suspend"
 )
func (tm *TransactionManager)CreateTransaction(emmitter_id, beneficiary_id, status string, )(*Transaction, error)  {
	var id string 
	if err := tm.DB.QueryRow(`INSERT INTO transactions (emmitter_id, beneficiary_id, amount, status) VALUES ($1, $2, $3, $4,) RETURNING id `,emmitter_id, beneficiary_id, suspend ).Scan(&id) ; err != nil {
		return nil, err 
	}
	var transaction Transaction
	if err := tm.DB.QueryRow(`select transactions.id, users.username, transactions.amount, transactions.date, 
		transactions.status from transactions inner join users on transactions.emmiter_id = users.id 
		inner join users on transactions.beneficiary_id = users.id Where transactions.id = $1`, id).Scan(&transaction.ID,&transaction.Emitter.Username, &transaction.Beneficiary.Username, &transaction.Amount, &transaction.Date, &transaction.Status); err != nil {
			return nil, err 
		}
	
	return &transaction, nil 
}

func (tm *TransactionManager)ValidateTransaction(id string)(*Transaction, error)  {
	var transaction Transaction
	if err := tm.DB.QueryRow(`UPDATE transactions SET status=$1 where id =$2 transactions.id, users.username, transactions.amount, transactions.date, 
	transactions.status from transactions inner join users on transactions.emmiter_id = users.id 
	inner join users on transactions.beneficiary_id = users.id Where transactions.id = $1`, done, id).Scan(&transaction.ID,&transaction.Emitter.Username, &transaction.Beneficiary.Username, &transaction.Amount, &transaction.Date, &transaction.Status); err != nil {
		return nil ,err
	}
	return &transaction, nil
}


