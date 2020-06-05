package payment

import (
	"database/sql"
	"fmt"
	"github.com/AgileBits/go-redis-queue/redisqueue"
	"github.com/google/uuid"
	"strings"
	"time"
)

const SessionLifeTime time.Duration = 30 * time.Minute

type Payment struct {
	ID        int64     `json:"-"`
	UID       uuid.UUID `json:"uid"`
	Amount    float32   `json:"amount"`
	Purpose   string    `json:"purpose"`
	CreatedAt time.Time `json:"created_at"`
}

func (p *Payment) TimeLeft() float64 {
	return p.CreatedAt.Add(SessionLifeTime).Sub(time.Now()).Seconds()
}

func NewPayment(amount float32, purpose string) (*Payment, error) {
	uid, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	return &Payment{
		UID:       uid,
		Amount:    amount,
		Purpose:   purpose,
		CreatedAt: time.Now(),
	}, nil
}

// --------------

type Repository struct {
	db    *sql.DB
	queue *redisqueue.Queue
}

func NewRepository(db *sql.DB, queue *redisqueue.Queue) *Repository {
	return &Repository{
		db:    db,
		queue: queue,
	}
}

func (r *Repository) Create(amount float32, purpose string) (*Payment, error) {
	p, err := NewPayment(amount, purpose)
	if err != nil {
		return nil, err
	}
	row, err := r.db.Exec("INSERT INTO `payments` (`uid`, `amount`, `purpose`, `created_at`) VALUES (?, ?, ?, ?)", p.UID, p.Amount, p.Purpose, p.CreatedAt)
	if err != nil {
		return nil, err
	}
	id, err := row.LastInsertId()
	if err != nil {
		return nil, err
	}
	if _, err = r.db.Exec("INSERT INTO `payment_status` (`payment_id`, `status_id`) VALUES (?, ?)", id, StatusExpected); err != nil {
		return nil, err
	}
	now := time.Now().Add(SessionLifeTime)
	if _, err = r.queue.Schedule(fmt.Sprintf("payment_timeout:%d", id), now); err != nil {
		return nil, err
	}

	return p, nil
}

func (r *Repository) GetByID(id int64) (*Payment, error) {
	if id <= 0 {
		return nil, nil
	}
	p := &Payment{}
	row := r.db.QueryRow("SELECT id, uid, amount, purpose, created_at FROM payments WHERE id = ? AND NOT EXISTS (SELECT * FROM payment_status WHERE payment_id = payments.id AND (status_id = ? OR status_id = ?))", id, StatusPaid, StatusError)
	err := row.Scan(&p.ID, &p.UID, &p.Amount, &p.Purpose, &p.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (r *Repository) GetByUID(uid string) (*Payment, error) {
	if len(uid) == 0 {
		return nil, nil
	}
	p := &Payment{}
	row := r.db.QueryRow("SELECT id, uid, amount, purpose, created_at FROM payments WHERE uid = ? AND NOT EXISTS (SELECT * FROM payment_status WHERE payment_id = payments.id AND (status_id = ? OR status_id = ?))", uid, StatusPaid, StatusError)
	err := row.Scan(&p.ID, &p.UID, &p.Amount, &p.Purpose, &p.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (r *Repository) SetStatus(p *Payment, statusID int64) error {
	_, err := r.db.Exec("INSERT INTO `payment_status` (`payment_id`, `status_id`) VALUES (?, ?)", p.ID, statusID)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) All(from, to string) ([]*Payment, error) {
	var fromTime, toTime time.Time
	var err error

	args := make([]interface{}, 0, 3)
	args = append(args, StatusPaid)
	queryTime := make([]string, 0, 2)

	from = strings.ToUpper(from)
	if len(from) > 0 {
		fromTime, err = time.Parse(time.RFC3339, from)
		if err != nil {
			return nil, err
		}
		queryTime = append(queryTime, "created_at >= ?")
		args = append(args, fromTime)
	}

	to = strings.ToUpper(to)
	if len(to) > 0 {
		toTime, err = time.Parse(time.RFC3339, to)
		if err != nil {
			return nil, err
		}
		queryTime = append(queryTime, "created_at <= ?")
		args = append(args, toTime)
	}

	query := "SELECT id, uid, amount, purpose, created_at FROM payments WHERE EXISTS (SELECT * FROM payment_status WHERE payment_id = payments.id AND status_id = ?)"
	if len(queryTime) > 0 {
		query += " AND " + strings.Join(queryTime, " AND ")
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	payments := make([]*Payment, 0, 5)
	for rows.Next() {
		p := &Payment{}
		_ = rows.Scan(&p.ID, &p.UID, &p.Amount, &p.Purpose, &p.CreatedAt)
		payments = append(payments, p)
	}

	return payments, nil
}

func (r *Repository) Cancell(p *Payment) (bool, error) {
	_, err := r.db.Exec("INSERT INTO `payment_status` (`payment_id`, `status_id`) VALUES (?, ?)", p.ID, StatusError)
	if err != nil {
		return false, err
	}
	return true, err
}
