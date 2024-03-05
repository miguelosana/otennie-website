package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	bolt "go.etcd.io/bbolt"
)

type Db struct {
	db *bolt.DB
}

func NewDB() *Db {
	db, err := bolt.Open("ottenie.db", 0660, &bolt.Options{Timeout: 1 * time.Second})

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
	return &Db{db: db}
}

func (d *Db) InsertContact(c ContactForm) error {

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

func (d *Db) InserVideoWaitlist(v VideoWaitlistForm) error {
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

func (d *Db) Close() {
	d.db.Close()
}
