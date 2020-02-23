package ecovacs

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	// Frameworks
	home "github.com/djthorpe/mutablehome"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type DeviceState struct {
	values  map[string]*XMPPMessage
	expires map[string]time.Time
	queue   []home.EcovacsEventType

	sync.RWMutex
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Sets a value for key, and returns true if value was added or modified
func (this *DeviceState) Set(value *XMPPMessage, ttl time.Duration) bool {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	// Create data structures
	this.init()

	// Add TTL value
	mapKey := mapKey(value)

	// Check for modified message
	modified := true
	if old, exists := this.values[mapKey]; exists {
		if time.Now().After(this.expires[mapKey]) {
			// Special case where the value has expired already
			modified = true
		} else {
			// Mark as modified if the value is different
			modified = old.Equals(value) == false
		}
	}

	// Update values
	this.expires[mapKey] = time.Now().Add(ttl)
	this.values[mapKey] = value

	// Return modified flag
	return modified
}

// Remove all elements from state
func (this *DeviceState) RemoveAll() {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	// Create data structures
	this.values = nil
	this.expires = nil
	this.queue = nil
}

func (this *DeviceState) AddExpiredKey(key home.EcovacsEventType) {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	// Create data structures
	this.init()

	// Check to see if key has already been added
	for _, elem := range this.queue {
		if elem == key {
			return
		}
	}

	// Append to end
	this.queue = append(this.queue, key)
}

func (this *DeviceState) NextExpiredKey() home.EcovacsEventType {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	// Create data structures
	this.init()

	// Gather all expired keys
	expired := make(map[home.EcovacsEventType]bool, len(this.values))
	for k, v := range this.values {
		if this.exists(k) == false {
			expired[v.Type()] = true
		}
	}

	// Append any from the queue
	for _, elem := range this.queue {
		expired[elem] = true
	}

	// No expired keys, return NONE
	if len(expired) == 0 {
		return home.ECOVACS_EVENT_NONE
	}

	// Return a random key
	expired_keys := []home.EcovacsEventType{}
	for k := range expired {
		expired_keys = append(expired_keys, k)
	}
	i := rand.Int() % len(expired_keys)
	return expired_keys[i]
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *DeviceState) init() {
	if this.values == nil {
		this.values = make(map[string]*XMPPMessage)
	}
	if this.expires == nil {
		this.expires = make(map[string]time.Time)
	}
	if this.queue == nil {
		this.queue = make([]home.EcovacsEventType, 0, 10)
	}
}

// exists returns true if a value exists for key and it's not expired
func (this *DeviceState) exists(key string) bool {
	if _, exists := this.values[key]; exists == false {
		return false
	} else if expires, exists := this.expires[key]; exists == false {
		return false
	} else if time.Now().After(expires) {
		return false
	} else {
		return true
	}
}

func mapKey(value *XMPPMessage) string {
	valueType := value.Type()
	prefix := strings.TrimPrefix(fmt.Sprint(valueType), "ECOVACS_EVENT_")
	if valueType == home.ECOVACS_EVENT_LIFESPAN {
		part, _, _ := value.LifeSpan()
		return prefix + "_" + strings.ToUpper(fmt.Sprint(part))
	} else {
		return prefix
	}
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *DeviceState) String() string {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	str := "<DeviceState"
	for k, v := range this.values {
		if this.exists(k) {
			str += fmt.Sprintf(" %v=%v", k, v.Value())
		} else {
			str += fmt.Sprintf(" %v=%v [EXPIRED]", k, v.Value())
		}
	}
	return str + ">"
}
