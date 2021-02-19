package mongo

import (
	"context"
	"dryan/db"
	logutil "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func getDB(ctx context.Context, dbName string) (*db.MgoDB, error) {
	session, err := db.Mongo(dbName)
	if err != nil {
		return nil, err
	}

	db, err := session.DB()
	if err != nil {
		logutil.Errorf("Opened DB connection failed, DB Name : [%s], error : [%v]", dbName, err)
		return nil, err
	}

	return db, nil
}

func Insert(ctx context.Context, dbName string, collectionName string, data interface{}) error {
	db, err := getDB(ctx, dbName)
	if err != nil {
		logutil.Errorf("Get DB connection failed, DB Name : [%s], Collection Name : [%s], error : [%v]", dbName, collectionName, err)
		return err
	}
	defer db.Close()

	c := db.C(collectionName)

	err = c.Insert(data)
	if err != nil {
		logutil.Errorf("Insert %s failed, source data : [%v], error : [%v]", collectionName, data, err)
		return err
	}

	return nil
}

func BatchInsert(ctx context.Context, dbName string, collectionName string, data []interface{}) error {
	db, err := getDB(ctx, dbName)
	if err != nil {
		logutil.Errorf("Get DB connection failed, DB Name : [%s], Collection Name : [%s], error : [%v]", dbName, collectionName, err)
		return err
	}
	defer db.Close()

	c := db.C(collectionName)

	err = c.Insert(data...)
	if err != nil {
		logutil.Errorf("Insert %s failed, source data : [%v], error : [%v]", collectionName, data, err)
		return err
	}

	return nil
}

func DeleteById(ctx context.Context, dbName string, collectionName string, id bson.ObjectId) error {
	return Delete(ctx, dbName, collectionName, bson.M{"_id": id})
}

func Delete(ctx context.Context, dbName string, collectionName string, selector interface{}) error {
	db, err := getDB(ctx, dbName)
	if err != nil {
		logutil.Errorf("Get DB connection failed, DB Name : [%s], Collection Name : [%s], error : [%v]", dbName, collectionName, err)
		return err
	}
	defer db.Close()

	c := db.C(collectionName)

	_, err = c.RemoveAll(selector)
	if err != nil {
		logutil.Errorf("Delete %s failed, selector : [%v], error : [%v]", collectionName, selector, err)
		return err
	}

	return nil
}

func UpdateById(ctx context.Context, dbName string, collectionName string, id bson.ObjectId, update interface{}) error {
	return Update(ctx, dbName, collectionName, bson.M{"_id": id}, update)
}

func Update(ctx context.Context, dbName string, collectionName string, selector interface{}, update interface{}) error {
	db, err := getDB(ctx, dbName)
	if err != nil {
		logutil.Errorf("Get DB connection failed, DB Name : [%s], Collection Name : [%s], error : [%v]", dbName, collectionName, err)
		return err
	}
	defer db.Close()

	c := db.C(collectionName)

	_, err = c.UpdateAll(selector, update)
	if err != nil {
		logutil.Errorf("Update %s failed, selector : [%v], update : [%v], error : [%v]", collectionName, selector, update, err)
		return err
	}

	return nil
}

func UpdateByIds(ctx context.Context, dbName string, collectionName string, updateDataMap map[bson.ObjectId]interface{}) error {
	db, err := getDB(ctx, dbName)
	if err != nil {
		logutil.Errorf("get DB connection failed, DB Name : [%s], Collection Name : [%s], error : [%v]", dbName, collectionName, err)
		return err
	}
	defer db.Close()

	c := db.C(collectionName)

	for id, updateData := range updateDataMap {
		selector := bson.M{"_id": id}
		_, err = c.UpdateAll(selector, updateData)

		if err != nil {
			logutil.Errorf("update by id failed, id:%s, error:%v", id.Hex(), updateData)
			return err
		}
	}

	return nil
}

func UpsertById(ctx context.Context, dbName string, collectionName string, id bson.ObjectId, update interface{}) error {
	return Upsert(ctx, dbName, collectionName, bson.M{"_id": id}, update)
}

func Upsert(ctx context.Context, dbName string, collectionName string, selector interface{}, update interface{}) error {
	db, err := getDB(ctx, dbName)
	if err != nil {
		logutil.Errorf("Get DB connection failed, DB Name : [%s], Collection Name : [%s], error : [%v]", dbName, collectionName, err)
		return err
	}
	defer db.Close()

	c := db.C(collectionName)

	_, err = c.Upsert(selector, update)
	if err != nil {
		logutil.Errorf("Upsert %s failed, selector : [%v], update : [%v], error : [%v]", collectionName, selector, update, err)
		return err
	}

	return nil
}

func FindById(ctx context.Context, dbName string, collectionName string, id bson.ObjectId, result interface{}) error {
	db, err := getDB(ctx, dbName)
	if err != nil {
		logutil.Errorf("Get DB connection failed, DB Name : [%s], Collection Name : [%s], error : [%v]", dbName, collectionName, err)
		return err
	}
	defer db.Close()

	c := db.C(collectionName)
	err = c.FindId(id).One(result)
	if err == mgo.ErrNotFound {
		return err
	}
	if err != nil {
		logutil.Errorf("Find %s failed, id : [%v], error : [%v]", collectionName, id, err)
		return err
	}

	return nil
}

func Count(ctx context.Context, dbName string, collectionName string, query interface{}) (int, error) {
	db, err := getDB(ctx, dbName)
	if err != nil {
		logutil.Errorf("Get DB connection failed, DB Name : [%s], Collection Name : [%s], error : [%v]", dbName, collectionName, err)
		return -1, err
	}
	defer db.Close()

	c := db.C(collectionName)
	number, err := c.Find(query).Count()
	if err != nil {
		logutil.Errorf("Count %s failed, query: [%v], error : [%v]", collectionName, query, err)
		return -1, err
	}

	return number, nil
}

type FindModel struct {
	CollectionName string
	Query          interface{}
	Fields         interface{}
	Sort           []string
	Cursor         int
	Size           int
	Results        interface{}
}

func FindByModel(ctx context.Context, dbName string, model FindModel) error {
	return Find(ctx, dbName, model.CollectionName, model.Query, model.Fields, model.Sort, model.Cursor, model.Size, model.Results)
}

/**
 * 执行查询
 * @param dbName         string [数据库名称]
 * @param collectionName string [集合名称]
 * @param query          bson.M [查询条件，传nil表示查询不做约束]
 * @param fields         bson.M [投影字段，传nil表示查询出所有字段，1为inclusion模式 指定返回的键，不返回其他键；0为exclusion模式 指定不返回的键,返回其他键]
 * @param sort           bson.M [排序条件，传nil表示不做排序，条件前面加【+/-】表示【升序/降序】，默认为升序]
 * @param cursor         int    [起始位置，从0开始]
 * @param size           int    [查询数量，填0表示查询所有]
 * @param results        interface{} [结果集，数组指针类型]
 *
 * 查询示例
 * var f []do.Favourite（必须为数组类型）
 * mongo.Find(nil, "news", "favourites", bson.M{"person_id": "wuyiwen"}, bson.M{"person_id": 1, "create_time": 1}, []string{"status", "-create_time"}, 0, 10, &f)
 * 等价于sql：SELECT person_id, create_time FROM favourites WHERE person_id = "wuyiwen" SORT BY status, create_time DESC LIMIT 0, 10
 */
func Find(ctx context.Context, dbName string, collectionName string, query interface{}, fields interface{}, sort []string, cursor int, size int, results interface{}) error {
	db, err := getDB(ctx, dbName)
	if err != nil {
		logutil.Errorf("Get DB connection failed, DB Name : [%s], Collection Name : [%s], error : [%v]", dbName, collectionName, err)
		return err
	}
	defer db.Close()
	c := db.C(collectionName)
	err = c.Find(query).Select(fields).Sort(sort...).Skip(cursor).Limit(size).All(results)
	if err != nil {
		logutil.Errorf("Find %s failed, query : [%v], fields : [%v], sort : [%v], cursor : [%v], size : [%v], error : [%v]",
			collectionName, query, fields, sort, cursor, size, err)
		return err
	}

	return nil
}

func FindOneByModel(ctx context.Context, dbName string, model FindModel) error {
	return FindOne(ctx, dbName, model.CollectionName, model.Query, model.Fields, model.Sort, model.Cursor, model.Size, model.Results)
}

func FindOne(ctx context.Context, dbName string, collectionName string, query interface{}, fields interface{}, sort []string, cursor int, size int, result interface{}) error {
	db, err := getDB(ctx, dbName)
	if err != nil {
		logutil.Errorf("Get DB connection failed, DB Name : [%s], Collection Name : [%s], error : [%v]", dbName, collectionName, err)
		return err
	}
	defer db.Close()
	c := db.C(collectionName)
	err = c.Find(query).Select(fields).Sort(sort...).Skip(cursor).Limit(size).One(result)
	if err == mgo.ErrNotFound {
		return err
	}
	if err != nil {
		logutil.Errorf("FindOne %s failed, query : [%v], fields : [%v], sort : [%v], cursor : [%v], size : [%v], error : [%v]",
			collectionName, query, fields, sort, cursor, size, err)
		return err
	}

	return nil
}

func Pipe(ctx context.Context, dbName string, collectionName string, pipeline interface{}, result interface{}) error {
	db, err := getDB(ctx, dbName)
	if err != nil {
		logutil.Errorf("Get DB connection failed, DB Name : [%s], Collection Name : [%s], error : [%v]", dbName, collectionName, err)
		return err
	}
	defer db.Close()
	c := db.C(collectionName)

	err = c.Pipe(pipeline).All(result)

	if err != nil {
		logutil.Errorf("Pipe %s failed, pipe : [%v], error : [%v]",
			collectionName, pipeline, err)
		return err
	}

	return nil
}

func FindByIds(ctx context.Context, dbName string, collectionName string, ids []bson.ObjectId, result interface{}) error {
	db, err := getDB(ctx, dbName)
	if err != nil {
		logutil.Errorf("Get DB connection failed, DB Name : [%s], Collection Name : [%s], error : [%v]", dbName, collectionName, err)
		return err
	}
	defer db.Close()

	c := db.C(collectionName)

	query := bson.M{"_id": bson.M{"$in": ids}}
	err = c.Find(query).All(result)

	if err != nil {
		logutil.Errorf("Find %s failed, id : [%v], error : [%v]", collectionName, ids, err)
		return err
	}

	return nil
}

func AddToSetById(ctx context.Context, dbName string, collectionName string, id bson.ObjectId, fieldName string, data []interface{}) error {
	db, err := getDB(ctx, dbName)
	if err != nil {
		logutil.Errorf("Get DB connection failed, DB Name : [%s], Collection Name : [%s], error : [%v]", dbName, collectionName, err)
		return err
	}
	defer db.Close()

	c := db.C(collectionName)

	updateMap := bson.M{
		"$addToSet": bson.M{fieldName: bson.M{"$each": data}},
	}

	err = c.UpdateId(id, updateMap)

	if err != nil {
		logutil.Errorf("addToSet failed, fieldName:%s, data:%+v, error : [%v]", fieldName, data, err)
		return err
	}

	return nil
}

func AddToSet(ctx context.Context, dbName string, collectionName string, selector interface{}, fieldName string, data []interface{}) error {
	db, err := getDB(ctx, dbName)
	if err != nil {
		logutil.Errorf("Get DB connection failed, DB Name : [%s], Collection Name : [%s], error : [%v]", dbName, collectionName, err)
		return err
	}
	defer db.Close()

	c := db.C(collectionName)

	updateMap := bson.M{
		"$addToSet": bson.M{fieldName: bson.M{"$each": data}},
	}

	_, err = c.UpdateAll(selector, updateMap)

	if err != nil {
		logutil.Errorf("addToSet failed, fieldName:%s, data:%+v, error : [%v]", fieldName, data, err)
		return err
	}

	return nil
}

func PipeUseDisk(ctx context.Context, dbName string, collectionName string, pipeline interface{}, result interface{}) error {
	db, err := getDB(ctx, dbName)
	if err != nil {
		logutil.Errorf("Get DB connection failed, DB Name : [%s], Collection Name : [%s], error : [%v]", dbName, collectionName, err)
		return err
	}
	defer db.Close()
	c := db.C(collectionName)

	err = c.Pipe(pipeline).AllowDiskUse().All(result)

	if err != nil {
		logutil.Errorf("Pipe %s failed, pipe : [%v], error : [%v]", collectionName, pipeline, err)
		return err
	}

	return nil
}

type BulkUpdateItem struct {
	Selector bson.M
	Update   bson.M
}

func BulkUpdate(ctx context.Context, dbName string, collectionName string, bulkUpdateItems []*BulkUpdateItem) error {
	db, err := getDB(ctx, dbName)
	if err != nil {
		logutil.Errorf("Get DB connection failed, DB Name : [%s], Collection Name : [%s], error : [%v]", dbName, collectionName, err)
		return err
	}
	defer db.Close()

	c := db.C(collectionName)

	bulk := c.Bulk()
	bulk.Unordered()

	for _, bulkUpdateItem := range bulkUpdateItems {
		if bulkUpdateItem.Selector == nil || bulkUpdateItem.Update == nil {
			continue
		}
		bulk.Update(bulkUpdateItem.Selector, bulkUpdateItem.Update)
	}

	_, err = bulk.Run()
	if err != nil {
		logutil.Errorf("Bulk update %s failed, size : [%v], error : [%v]", collectionName, len(bulkUpdateItems), err)
		return err
	}

	return nil
}

func BulkInsert(ctx context.Context, dbName string, collectionName string, data []interface{}) error {
	db, err := getDB(ctx, dbName)
	if err != nil {
		logutil.Errorf("Get DB connection failed, DB Name : [%s], Collection Name : [%s], error : [%v]", dbName, collectionName, err)
		return err
	}
	defer db.Close()

	c := db.C(collectionName)

	bulk := c.Bulk()
	bulk.Unordered()

	bulk.Insert(data...)

	_, err = bulk.Run()
	if err != nil {
		logutil.Errorf("Bulk insert %s failed, size : [%v], error : [%v]", collectionName, len(data), err)
		return err
	}

	return nil
}

func BulkUpsert(ctx context.Context, dbName string, collectionName string, bulkUpdateItems []*BulkUpdateItem) error {
	db, err := getDB(ctx, dbName)
	if err != nil {
		logutil.Errorf("Get DB connection failed, DB Name : [%s], Collection Name : [%s], error : [%v]", dbName, collectionName, err)
		return err
	}
	defer db.Close()

	c := db.C(collectionName)

	bulk := c.Bulk()
	bulk.Unordered()

	for _, bulkUpdateItem := range bulkUpdateItems {
		if bulkUpdateItem.Selector == nil || bulkUpdateItem.Update == nil {
			continue
		}
		bulk.Upsert(bulkUpdateItem.Selector, bulkUpdateItem.Update)
	}

	_, err = bulk.Run()
	if err != nil {
		logutil.Errorf("Bulk upsert %s failed, size : [%v], error : [%v]", collectionName, len(bulkUpdateItems), err)
		return err
	}

	return nil
}
