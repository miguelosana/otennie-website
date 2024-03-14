package storage

import (
	"context"
	"log"
	"otennie/models"

	"cloud.google.com/go/firestore"
	"github.com/google/uuid"
)

type FirestoreStorage struct {
	client *firestore.Client
}

func NewFirestoreStorage(ctx context.Context, projectId string) *FirestoreStorage {

	client, err := firestore.NewClient(ctx, projectId)
	if err != nil {
		log.Fatal(err)
	}

	return &FirestoreStorage{
		client: client,
	}

}

func (f *FirestoreStorage) collection(collection string) *firestore.CollectionRef {
	return f.client.Collection(collection)
}
func (f *FirestoreStorage) InsertContact(ctx context.Context, c models.ContactForm) error {
	id := uuid.New()
	_, err := f.collection("contact-form").Doc(id.String()).Set(ctx, &c)
	return err

}
func (f *FirestoreStorage) InsertVideoWaitlist(ctx context.Context, v models.VideoWaitlistForm) error {

	id := uuid.New()
	_, err := f.collection("video-waitlist").Doc(id.String()).Set(ctx, &v)
	return err
}

func (f *FirestoreStorage) Close() error {
	return f.client.Close()

}
