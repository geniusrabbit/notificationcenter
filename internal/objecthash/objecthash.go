package objecthash

import (
	"crypto/sha256"
	"fmt"
)

type objectHash interface {
	ObjectHash() string
}

// Hash value from the object
func Hash(obj interface{}) string {
	if ohasher, ok := obj.(objectHash); ok {
		return ohasher.ObjectHash()
	}
	h := sha256.New()
	// Ignore error check...
	_, _ = h.Write([]byte(fmt.Sprintf("%v", obj)))
	return fmt.Sprintf("%x", h.Sum(nil))
}
