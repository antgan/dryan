package mongo

import (
	"context"

	"gopkg.in/mgo.v2/bson"
)

func NewMongoOp(dbName, collectionName string) *MongoOp {
	return &MongoOp{
		dbName:         dbName,
		collectionName: collectionName,
	}
}

type MongoOp struct {
	dbName         string
	collectionName string
}

func (op *MongoOp) Insert(ctx context.Context, data interface{}) error {
	return Insert(ctx, op.dbName, op.collectionName, data)
}

func (op *MongoOp) BatchInsert(ctx context.Context, data []interface{}) error {
	return BatchInsert(ctx, op.dbName, op.collectionName, data)
}

func (op *MongoOp) DeleteById(ctx context.Context, id bson.ObjectId) error {
	return DeleteById(ctx, op.dbName, op.collectionName, id)
}

func (op *MongoOp) Delete(ctx context.Context, selector interface{}) error {
	return Delete(ctx, op.dbName, op.collectionName, selector)
}

func (op *MongoOp) Update(ctx context.Context, selector interface{}, updateOp interface{}) error {
	return Update(ctx, op.dbName, op.collectionName, selector, updateOp)
}

func (op *MongoOp) UpdateById(ctx context.Context, id bson.ObjectId, updateOp interface{}) error {
	return UpdateById(ctx, op.dbName, op.collectionName, id, updateOp)
}

func (op *MongoOp) Count(ctx context.Context, query interface{}) (int, error) {
	return Count(ctx, op.dbName, op.collectionName, query)
}

func (op *MongoOp) FindById(ctx context.Context, id bson.ObjectId, result interface{}) error {
	return FindById(ctx, op.dbName, op.collectionName, id, result)
}

func (op *MongoOp) Find(ctx context.Context, results, query interface{}, sort []string, fields bson.M, cursor, size int) error {
	return Find(ctx, op.dbName, op.collectionName, query, fields, sort, cursor, size, results)
}

func (op *MongoOp) FindOne(ctx context.Context, result, query interface{}, sort []string, fields bson.M, cursor int) error {
	return FindOne(ctx, op.dbName, op.collectionName, query, fields, sort, cursor, 1, result)
}

func (op *MongoOp) Pipe(ctx context.Context, pipeline interface{}, result interface{}) error {
	return Pipe(ctx, op.dbName, op.collectionName, pipeline, result)
}

func (op *MongoOp) Upsert(ctx context.Context, selector interface{}, updateOp interface{}) error {
	return Upsert(ctx, op.dbName, op.collectionName, selector, updateOp)
}

func (op *MongoOp) BulkUpdate(ctx context.Context, bulkUpdateItems []*BulkUpdateItem) error {
	return BulkUpdate(ctx, op.dbName, op.collectionName, bulkUpdateItems)
}

func (op *MongoOp) BulkUpsert(ctx context.Context, bulkUpdateItems []*BulkUpdateItem) error {
	return BulkUpsert(ctx, op.dbName, op.collectionName, bulkUpdateItems)
}
