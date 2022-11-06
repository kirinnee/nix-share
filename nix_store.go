package main

import (
	"github.com/samber/lo"
	"log"
	"os"
	"strings"
)

func parseRemove(search string) func(line, add string) (string, bool) {
	return func(line, remove string) (string, bool) {
		if strings.HasPrefix(line, search) {
			subs := strings.Split(line, " ")
			subs = lo.Drop[string](subs, 2)
			subs = lo.Filter[string](subs, func(v string, k int) bool {
				return v != remove
			})
			if len(subs) == 0 {
				return "", true
			}
			ret := strings.Join(subs, " ")
			return search + " = " + ret, true
		}
		return line, false
	}
}

func parseAdd(search string) func(line, add string) (string, bool) {
	return func(line, add string) (string, bool) {
		if strings.HasPrefix(line, search) {
			subs := strings.Split(line, " ")
			subs = lo.Drop[string](subs, 2)
			subs = append(subs, add)
			subs = lo.Uniq[string](subs)
			ret := strings.Join(subs, " ")
			return search + " = " + ret, true
		}
		return line, false
	}
}

func addNix(path string, p StorePacket) error {
	b, err := os.ReadFile(path)

	if err != nil {
		log.Println("Failed to read nix config file")
		return err
	}

	content := string(b)
	lines := strings.Split(content, "\n")
	outContent := ""
	subExist := false
	pubKeyExist := false

	addSub := parseAdd("substituters")
	addPubKey := parseAdd("trusted-public-keys")

	for _, l := range lines {
		l1, se := addSub(l, p.Host)
		subExist = subExist || se
		l2, pke := addPubKey(l1, p.PubKey)
		pubKeyExist = pubKeyExist || pke
		outContent += l2 + "\n"
	}
	if !subExist {
		outContent += "substituters = " + p.Host + "\n"
	}

	if !pubKeyExist {
		outContent += "trusted-public-keys = " + p.PubKey + "\n"
	}

	return os.WriteFile(path, []byte(outContent), 0644)
}

func removeNix(path string, p StorePacket) error {
	b, err := os.ReadFile(path)

	if err != nil {
		log.Println("Failed to read nix config file")
		return err
	}

	content := string(b)
	lines := strings.Split(content, "\n")

	outContent := ""
	rmSub := parseRemove("substituters")
	rmPubKey := parseRemove("trusted-public-keys")

	for _, l := range lines {
		l, _ = rmSub(l, p.Host)
		l, _ = rmPubKey(l, p.PubKey)
		outContent += l + "\n"
	}

	return os.WriteFile(path, []byte(outContent), 0644)
}
