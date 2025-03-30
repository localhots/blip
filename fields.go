package blip

import (
	"context"
	"sync"
)

// Field is a key-value pair that is used to add structured data to log entries.
type Field struct {
	Key   string
	Value any
}

// F is a convenient alias for a map of fields.
type F map[string]any

func makeFields(ctx context.Context, ff []F) *[]Field {
	var n int
	for _, f := range ff {
		n += len(f)
	}
	if n == 0 {
		return nil
	}

	fields := getFields()
	for k, v := range FromContext(ctx) {
		addField(fields, k, v)
	}
	for _, f := range ff {
		for k, v := range f {
			addField(fields, k, v)
		}
	}
	return fields
}

func addField(f *[]Field, key string, val any) {
	// Update existing entry if exists
	for i := range *f {
		if (*f)[i].Key == key {
			(*f)[i].Value = val
			return
		}
	}

	// Add new entry
	(*f) = append(*f, Field{key, val})
}

func sortFields(f []Field) {
	if len(f) > 1 {
		insertionSort(f)
	}
}

// insertionSort is great for small slices. Using this custom function instead
// of sort.Slice() reduces the number of allocations to zero.
func insertionSort(f []Field) {
	for i := 1; i < len(f); i++ {
		for j := i; j > 0 && f[j].Key < f[j-1].Key; j-- {
			f[j], f[j-1] = f[j-1], f[j]
		}
	}
}

//
// Fields pool
//

// Fields are pooled to reduce allocations.
// Preallocated slices of 20 fields should be enough for most cases, in worst
// case the slice will grow.
var fieldsPool = sync.Pool{
	New: func() any {
		fields := make([]Field, 0, 20)
		return &fields
	},
}

func getFields() *[]Field {
	return fieldsPool.Get().(*[]Field)
}

func putFields(fields *[]Field) {
	*fields = (*fields)[:0] // Reset
	fieldsPool.Put(fields)
}
