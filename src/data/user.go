package data

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/labstack/echo/v4"
)

// User is a user record
type User struct {
	UID  uint64 `json:"uid"`
	Name string `json:"name"`
}

func userFromIncompleteWithUID(uid uint64, incomplete IncompleteUser) User {
	return User{uid, incomplete.Name}
}

// IncompleteUser denotes user data without associated ID
type IncompleteUser struct {
	Name string `json:"name"`
}

// Users storage for instances of User
type Users struct {
	Users    []User          `json:"users"`
	LastUID  uint64          `json:"lastuid"`
	UsersMap map[uint64]User `json:"-"`
	lock     sync.RWMutex    `json:"-"`
}

func (store *Users) commitToStorage() error {
	jsonFile, err := os.Create("users.json")
	if err != nil {
		return err
	}
	jsonBytes, err := json.MarshalIndent(store, "", "\t")
	_, writeErr := jsonFile.Write(jsonBytes)
	if writeErr != nil {
		return writeErr
	}
	return nil
}

// Get fetches a user matching specified UID
func (store *Users) Get(uid uint64) *User {
	store.lock.RLock()
	defer store.lock.RUnlock()
	user, existed := store.UsersMap[uid]
	if existed {
		return &user
	}
	return nil
}

// Insert tries to create a new user
func (store *Users) Insert(incomplete IncompleteUser) (*User, error) {
	store.lock.Lock()
	store.LastUID++
	uid := store.LastUID
	user := userFromIncompleteWithUID(uid, incomplete)
	store.Users = append(store.Users, user)
	store.UsersMap[uid] = user
	err := store.commitToStorage()
	defer store.lock.Unlock()
	return &user, err
}

// Update tries to update a user matching specified UID
func (store *Users) Update(uid uint64, incomplete IncompleteUser) (*User, error) {
	store.lock.Lock()
	defer store.lock.Unlock()
	user := new(User)
	_, existed := store.UsersMap[uid]
	if existed {
		*user = userFromIncompleteWithUID(uid, incomplete)
		store.UsersMap[uid] = *user
		store.updateList()
		err := store.commitToStorage()
		if err != nil {
			return nil, err
		}
		return user, nil
	}
	return nil, nil
}

// Delete tries to delete a user matching specified UID
func (store *Users) Delete(uid uint64) error {
	store.lock.Lock()
	defer store.lock.Unlock()
	//fmt.Printf("\n  %v\n", store.UsersMap)
	delete(store.UsersMap, uid)
	store.updateList()
	err := store.commitToStorage()
	if err != nil {
		return err
	}
	return nil
}

func (store *Users) updateList() {
	store.Users = make([]User, 0)
	for _, v := range store.UsersMap {
		store.Users = append(store.Users, v)
	}
}

// UserContext echo.Context implementor providing access to users store
type UserContext struct {
	echo.Context
	Users *Users
}
