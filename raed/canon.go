package raed

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"sort"
	"strings"
)

func digestBytes(b []byte) string {
	h := sha256.Sum256(b)
	return "sha256:" + hex.EncodeToString(h[:])
}

func canonicalJSON(v any) ([]byte, error) {
	return json.Marshal(v)
}

func MintID(namespace, normalizedRef, objectType string) (Identity, error) {
	namespace = strings.TrimSpace(namespace)
	normalizedRef = strings.TrimSpace(normalizedRef)
	objectType = strings.TrimSpace(objectType)
	if namespace == "" || normalizedRef == "" || objectType == "" {
		return Identity{}, errors.New("namespace, normalized reference, and object type are required")
	}
	seed := []byte("raed:v0.1\x00" + namespace + "\x00" + normalizedRef + "\x00" + objectType)
	return Identity{ID: digestBytes(seed), Namespace: namespace, NormalizedRef: normalizedRef, ObjectType: objectType}, nil
}

func relationDigest(r Relation) (string, error) {
	cp := r
	cp.RelationID = ""
	cp.Signature = ""
	b, err := canonicalJSON(cp)
	if err != nil {
		return "", err
	}
	return digestBytes(b), nil
}

func sortedCopy(in []string) []string {
	out := append([]string(nil), in...)
	sort.Strings(out)
	return out
}
