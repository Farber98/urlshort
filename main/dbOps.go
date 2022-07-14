package main

import "github.com/boltdb/bolt"

func writeDb(db *bolt.DB) error {
	path1 := []byte("/pcbook")
	url1 := []byte("https://github.com/Farber98/go-grpc-pcbook")
	path2 := []byte("/socket")
	url2 := []byte("https://github.com/Farber98/Datagram-Socket-App")
	// store some data
	err := db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("pathurls"))
		if err != nil {
			return err
		}

		err = bucket.Put(path1, url1)
		if err != nil {
			return err
		}
		err = bucket.Put(path2, url2)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

func readDb(db *bolt.DB) (map[string]string, error) {
	mp := make(map[string]string)
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("pathurls"))

		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			mp[string(k)] = string(v)
		}
		return nil
	})
	return mp, err
}
