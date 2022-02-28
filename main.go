package tasks_db

import (
	"encoding/binary"
	"time"

	"github.com/boltdb/bolt"
)

type Task struct {
	Key   int
	Value string
}

// Creates a database file if it hasn't been created yet.
// Takes the database path string as input
func Init(dbPath string) ([]byte, *bolt.DB, error) {
	var taskBucket = []byte("tasks")
	db, err := bolt.Open(dbPath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, nil, err
	}
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(taskBucket)
		return err
	})

	if err != nil {
		return nil, nil, err
	}
	return taskBucket, db, nil
}

// Creates a New task and save it in the database
func CreateTask(task string, taskBucket []byte, db *bolt.DB) (int, error) {
	createdAt := time.Now().Format(time.RFC1123)
	var id int
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(taskBucket)
		id64, _ := b.NextSequence()
		id = int(id64)
		key := itob(id)
		return b.Put(key, []byte(task+"\n"+createdAt+"\n"+"false"))
	})
	if err != nil {
		return -1, err
	}
	return id, nil
}

// Updates a task status completed or not
func UpdateTask(task string, id int, taskBucket []byte, db *bolt.DB) (int, error) {
	createdAt := time.Now().Format(time.RFC1123)
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(taskBucket)
		key := itob(id)
		return b.Put(key, []byte(task+"\n"+createdAt))
	})
	if err != nil {
		return -1, err
	}
	return id, nil
}

// returns all the tasks currently stored in the database
func AllTasks(taskBucket []byte, db *bolt.DB) ([]Task, error) {
	var tasks []Task
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(taskBucket)
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			tasks = append(tasks, Task{
				Key:   btoi(k),
				Value: string(v),
			})
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

// Deletes a certain task
func DeleteTask(key int, taskBucket []byte, db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(taskBucket)
		return b.Delete(itob(key))
	})
}

/* func main() {
	home, _ := homedir.Dir()
	dbPath := filepath.Join(home, "munna_go_tasks.db")

	fmt.Println(dbPath)
	taskBucket, db, err := Init(dbPath)
	// _, _, err := Init(dbPath)

	if err != nil {
		fmt.Println("Can't Initialize Database")
	} else {
		fmt.Println("Successfully Initialized Database")
	}

	CreateTask("New Time", taskBucket, db)

	tasks, err := AllTasks(taskBucket, db)

	if err != nil {
		fmt.Println("Can't Access")
	}

	for _, task := range tasks {
		fmt.Println(task.Key)
		fmt.Println(task.Value)
		fmt.Println()
		fmt.Println()
		UpdateTask(task.Value, task.Key, taskBucket, db)
	}

	for _, task := range tasks {
		fmt.Println(task.Key)
		fmt.Println(task.Value)
		fmt.Println()
		fmt.Println()
	}

	// for _, task := range tasks {
	// 	DeleteTask(task.Key, taskBucket, db)
	// }

} */

func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

func btoi(b []byte) int {
	return int(binary.BigEndian.Uint64(b))
}
