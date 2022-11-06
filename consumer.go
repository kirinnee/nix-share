package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"time"
)

func cleanOutdated(cachePath string, now int64, timeout int, c chan int) {
	val := <-c
	log.Println("Clearing up from goroutine... now: ", now)
	defer func() {
		log.Println("Exiting clear up from goroutine...")
		c <- val
	}()
	store, err := getCacheValue(cachePath)
	if err != nil {
		log.Println("Failed to read cache path")
		log.Println(err)
		return
	}
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Println("Failed to obtain home directory")
		log.Println(err)
		return
	}
	remove := make([]StorePacket, 0)
	updated := make(map[string]StorePacketData)
	for k, v := range store {
		if now-v.TimeStamp > int64(timeout) {
			remove = append(remove, StorePacket{
				Host:   v.Host,
				PubKey: v.PubKey,
			})
		} else {
			updated[k] = v
		}
	}

	nixPath := filepath.Join(dirname, ".config/nix/nix.conf")

	for _, v := range remove {
		errRemove := removeNix(nixPath, v)
		if errRemove != nil {
			log.Println("Failed to remove nix sub and pub key")
			log.Println(err)
			return
		}
	}

	updateErr := setCacheValue(cachePath, updated)
	if updateErr != nil {
		log.Println("Failed to write cache path")
		log.Println(err)
		return
	}
}

func update(original map[string]StorePacketData, packet StorePacket, now int64, timeout int) ([]StorePacket, []StorePacket, map[string]StorePacketData) {

	remove := make([]StorePacket, 0)
	add := make([]StorePacket, 0)
	updated := make(map[string]StorePacketData)
	match := false
	for k, v := range original {
		if packet.PubKey == k {
			updated[k] = StorePacketData{
				Host:      v.Host,
				PubKey:    v.PubKey,
				TimeStamp: now,
			}
			match = true
		} else {
			// decide to add or remove
			if now-v.TimeStamp > int64(timeout) {
				remove = append(remove, StorePacket{
					Host:   v.Host,
					PubKey: v.PubKey,
				})
			} else {
				updated[k] = v
			}
		}
	}

	if !match {
		updated[packet.Host] = StorePacketData{
			Host:      packet.Host,
			PubKey:    packet.PubKey,
			TimeStamp: now,
		}
		add = append(add, packet)

	}

	return remove, add, updated
}

func receive(port, cachePath string, timeout int, c chan int) {

	pc, err := net.ListenPacket("udp4", ":"+port)
	if err != nil {
		panic(err)
	}
	defer pc.Close()

	buf := make([]byte, 4096*4)
	n, addr, err := pc.ReadFrom(buf)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s sent this: %s\n", addr, buf[:n])

	// lock
	val := <-c
	defer func() {
		c <- val
	}()

	var packet StorePacket
	content := string(buf[:n])
	e := json.Unmarshal([]byte(content), &packet)
	if e != nil {
		log.Println("Failed to unmarshal packet")
		log.Println(err)
		return
	}

	store, err := getCacheValue(cachePath)
	if err != nil {
		log.Println("Failed to read cache path")
		log.Println(err)
		return
	}

	remove, add, updated := update(store, packet, time.Now().Unix(), timeout)

	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Println("Failed to obtain home directory")
		log.Println(err)
		return
	}

	nixPath := filepath.Join(dirname, ".config/nix/nix.conf")

	for _, v := range remove {
		errRemove := removeNix(nixPath, v)
		if errRemove != nil {
			log.Println("Failed to remove nix sub and pub key")
			log.Println(err)
			return
		}
	}

	for _, v := range add {
		errAdd := addNix(nixPath, v)
		if errAdd != nil {
			log.Println("Failed to add nix sub and pub key")
			log.Println(err)
			return
		}
	}

	updateErr := setCacheValue(cachePath, updated)
	if updateErr != nil {
		log.Println("Failed to write cache path")
		log.Println(err)
		return
	}

}
