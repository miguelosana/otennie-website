package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"otennie/models"
	"time"

	"github.com/google/uuid"
	bolt "go.etcd.io/bbolt"
)

type BoltStorage struct {
	db *bolt.DB
}

func NewBoltStorage(dbFile string) *BoltStorage {
	db, err := bolt.Open(dbFile, 0660, &bolt.Options{Timeout: 1 * time.Second})

	if err != nil {
		log.Fatal(err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("CONTACT"))
		if err != nil {
			return fmt.Errorf("could not create contact bucket: %v", err)
		}

		_, err = tx.CreateBucketIfNotExists([]byte("VIDEO_WAITLIST"))
		if err != nil {
			return fmt.Errorf("could not create video_waitlist bucket: %v", err)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	return &BoltStorage{db: db}
}

func (d *BoltStorage) InsertContact(ctx context.Context, c models.ContactForm) error {

	encoded, err := json.Marshal(c)
	if err != nil {
		return err
	}
	u := uuid.New()
	err = d.db.Update(func(tx *bolt.Tx) error {
		err = tx.Bucket([]byte("CONTACT")).Put([]byte(u.String()), encoded)
		if err != nil {
			return fmt.Errorf("could not set contact %v", err)
		}
		return nil
	})
	return err
}

func (d *BoltStorage) InsertVideoWaitlist(ctx context.Context, v models.VideoWaitlistForm) error {
	encoded, err := json.Marshal(v)
	if err != nil {
		return err
	}
	u := uuid.New()
	err = d.db.Update(func(tx *bolt.Tx) error {
		err = tx.Bucket([]byte("VIDEO_WAITLIST")).Put([]byte(u.String()), encoded)
		if err != nil {
			return fmt.Errorf("could not set waitlist %v", err)
		}
		return nil
	})
	return err
}

func (d *BoltStorage) Close() error {
	return d.db.Close()

}
