package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

type Message struct {
	text   string
	expire time.Duration
}

type MessageStore struct {
	messages    map[string]Message
	maxMessages int
	mutex       sync.Mutex
}

func (store *MessageStore) Add(message Message) (string, bool) {

	if len(store.messages) >= store.maxMessages {
		return "", false
	}

	id := strconv.Itoa(len(store.messages) + 1) // Change this to be more unique

	store.mutex.Lock()
	defer store.mutex.Unlock()

	if _, ok := store.messages[id]; ok {
		return "", false
	}

	store.messages[id] = message

	// Delete after expiration - TODO: what if the message is read before it expires?
	go func() {
		time.Sleep(message.expire)
		fmt.Println("Deleting message", id)
		store.mutex.Lock()
		defer store.mutex.Unlock()
		delete(store.messages, id)
	}()

	return id, true
}

func (store *MessageStore) Pop(id string) (Message, bool) {

	store.mutex.Lock()
	defer store.mutex.Unlock()

	msg, ok := store.messages[id]
	if ok {
		delete(store.messages, id)
	}

	return msg, ok
}
