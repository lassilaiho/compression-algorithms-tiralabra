package lz77

const (
	// The maximum load factor allowed for a dictionary.
	dictMaxLoadFactor = 0.75
	// The multiplier used when expanding dictionary bucket slice.
	dictGrowthMultiplier = 2
)

// dictionary maps keys into byte sequence positions in a data stream.
//
// A key is formed from the initial bytes of the sequence. The number of bytes
// used is specified using the constant dictKeySize. Sequences shorter than this
// can't be stored in the dictionary.
//
// Each key maps to an entry. An entry is a linked list of values. Values must
// be added in ascending order. The addition order is not enforced. It is the
// caller's responsibility to add values in the correct order.
type dictionary struct {
	buckets   []dictBucket
	itemCount int
}

// newDictionary returns creates a new dictionary.
func newDictionary() *dictionary {
	return &dictionary{
		buckets: make([]dictBucket, 1),
	}
}

// add adds value to the end of the entry corresponding to key.
func (d *dictionary) add(key dictKey, value int64) {
	entry := d.getEntry(key)
	if entry == nil {
		entry = &dictEntry{
			key:   key,
			first: &dictValue{value: value},
		}
		entry.last = entry.first
		d.setEntry(entry)
	} else {
		entry.last.next = &dictValue{value: value}
		entry.last = entry.last.next
	}
}

// get returns the first value corresponding to key or nil if the entry
// corresponding to key has no values.
func (d *dictionary) get(key dictKey) *dictValue {
	entry := d.getEntry(key)
	if entry == nil {
		return nil
	}
	return entry.first
}

// removeLesserThan removes all values lesser to value from the entry
// corresponding to key.
func (d *dictionary) removeLesserThan(key dictKey, value int64) {
	entry := d.getEntry(key)
	if entry == nil {
		return
	}
	dictValue := entry.first
	for dictValue != nil && dictValue.value < value {
		dictValue = dictValue.next
	}
	if dictValue == nil {
		d.getBucket(key).remove(key)
	} else {
		entry.first = dictValue
	}
}

// setEntry sets adds entry to the dictionary. If an entry with the same key
// already exists, it is overwritten with the new one.
func (d *dictionary) setEntry(entry *dictEntry) {
	bucket := d.getBucket(entry.key)
	existing := bucket.get(entry.key)
	if existing == nil {
		d.itemCount++
		bucket.add(entry)
		if d.loadFactor() > dictMaxLoadFactor {
			d.expand()
		}
	} else {
		existing.first = entry.first
		existing.last = entry.last
	}
}

// getEntry returns the entry corresponding to key. If no such entry exists, nil
// is returned.
func (d *dictionary) getEntry(key dictKey) *dictEntry {
	entry := d.getBucket(key).get(key)
	if entry == nil {
		return nil
	}
	return entry
}

// getBucket returns the bucket corresponding to key.
func (d *dictionary) getBucket(key dictKey) *dictBucket {
	return &d.buckets[key.hash()%uint(len(d.buckets))]
}

// loadFactor returns the load factor of the dictionary.
func (d *dictionary) loadFactor() float32 {
	return float32(d.itemCount) / float32(len(d.buckets))
}

// expand adds buckets to the dictionary according to dictGrowthMultiplier.
// expand always adds at least one bucket.
func (d *dictionary) expand() {
	oldBuckets := d.buckets
	newLen := dictGrowthMultiplier * len(oldBuckets)
	if newLen == len(oldBuckets) {
		newLen++
	}
	d.buckets = make([]dictBucket, newLen)
	for i := 0; i < len(oldBuckets); i++ {
		bucket := oldBuckets[i]
		for j := 0; j < len(bucket); j++ {
			entry := bucket[j]
			d.getBucket(entry.key).add(entry)
		}
	}
}

const dictKeySize = 2

type dictKey [dictKeySize]byte

// hash returns the hash of the key.
func (k dictKey) hash() uint {
	const shiftSize = 5
	const hashSize = 1024

	hash := uint(0)
	for i := 0; i < len(k); i++ {
		hash = (hash << shiftSize) ^ uint(k[i])
		hash %= hashSize
	}
	return hash
}

type dictValue struct {
	value int64
	next  *dictValue
}

type dictEntry struct {
	first *dictValue
	last  *dictValue
	key   dictKey
}

// dictBucket stores tableEntries with colliding key hashes.
type dictBucket []*dictEntry

// get returns the entry corresponding to key in the bucket or nil if no such
// entry exists.
func (b *dictBucket) get(key dictKey) *dictEntry {
	for i := 0; i < len(*b); i++ {
		entry := (*b)[i]
		if entry.key == key {
			return entry
		}
	}
	return nil
}

// add adds entry to the bucket.
func (b *dictBucket) add(entry *dictEntry) {
	if len(*b) == cap(*b) {
		newBucket := make(dictBucket, len(*b), cap(*b)+1)
		for i := 0; i < len(*b); i++ {
			newBucket[i] = (*b)[i]
		}
		*b = newBucket
	}
	*b = (*b)[:len(*b)+1]
	(*b)[len(*b)-1] = entry
}

// remove removes the entry corresponding to key from the bucket.
func (b *dictBucket) remove(key dictKey) {
	for i := 0; i < len(*b); i++ {
		if (*b)[i].key == key {
			(*b)[i] = (*b)[len(*b)-1]
			(*b)[len(*b)-1] = nil
			*b = (*b)[:len(*b)-1]
			return
		}
	}
}
