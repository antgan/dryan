package util

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gopkg.in/mgo.v2/bson"
	"strings"
)

func StringToObjectId(id string) (bson.ObjectId, error) {
	d, err := hex.DecodeString(id)
	if err != nil || len(d) != 12 {
		err := errors.New(fmt.Sprintf("invalid input to ObjectIdHex: %q", id))
		return "", err
	}

	return bson.ObjectId(d), nil
}

func StringsToObjectIds(ids []string) []bson.ObjectId {
	objIds := make([]bson.ObjectId, 0)
	for _, id := range ids {
		objId, err := StringToObjectId(id)
		if err == nil {
			objIds = append(objIds, objId)
		}
	}
	return objIds
}

func NewUUID() string {
	return strings.Replace(uuid.New().String(), "-", "", -1)
}

func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
