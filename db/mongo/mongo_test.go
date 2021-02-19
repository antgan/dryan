package mongo

import (
	"context"
	_ "dryan/db"
	logutil "github.com/sirupsen/logrus"
	"testing"
)

type Person struct {
	Name string `bson:"name"`
}

func TestMongo(t *testing.T) {
	err := Insert(context.Background(), "dryan", "c", &Person{
		Name: "Ant",
	})
	logutil.Info(err)
}
