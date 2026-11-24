package core

import (
	"fmt"
	"io"
	"sort"
)

// Changes contains the cumulative edits relative to the reader passed to OpenMetadata.
// An absent map entry is untouched. A present empty list removes that field.
type Changes struct {
	Fields       map[Field][]Value
	Custom       map[string][]string
	ClearUnknown bool
}

func (c Changes) Empty() bool { return len(c.Fields) == 0 && len(c.Custom) == 0 && !c.ClearUnknown }

// Editor holds metadata and pending edits. It does not own the source or
// destination and is not safe for concurrent use. The caller must keep the
// source open and unmodified for the life of an Editor.
type Editor struct {
	fields  map[Field][]Value
	custom  map[string][]string
	changes Changes
	save    func(io.Writer, Changes) error
}

func NewEditor(fields map[Field][]Value, custom map[string][]string, save func(io.Writer, Changes) error) *Editor {
	return &Editor{
		fields: cloneMap(fields), custom: cloneMap(custom), save: save,
		changes: Changes{Fields: make(map[Field][]Value), Custom: make(map[string][]string)},
	}
}

func cloneMap[K comparable, V any](original map[K][]V) map[K][]V {
	result := make(map[K][]V, len(original))
	for key, values := range original {
		result[key] = append([]V(nil), values...)
	}
	return result
}

// Get returns the field's values in source order. The result belongs to the caller.
func (e *Editor) Get(field Field) ([]Value, bool) {
	values := e.fields[field]
	return append([]Value(nil), values...), len(values) > 0
}

// GetField returns the first value. Use Get for the entire ordered field.
func (e *Editor) GetField(field Field) (Value, bool) {
	values := e.fields[field]
	if len(values) == 0 {
		return Value{}, false
	}
	return values[0], true
}

// SetField replaces a complete field. Omitting values removes the field.
func (e *Editor) SetField(field Field, values ...Value) error {
	d, ok := Lookup(field)
	if !ok {
		return fmt.Errorf("%w: %s", ErrInvalidField, field)
	}
	if !d.Mutable {
		return fmt.Errorf("%w: %s", ErrReadOnly, field)
	}
	if err := Validate(field, values...); err != nil {
		return err
	}
	if len(values) == 0 {
		delete(e.fields, field)
		e.changes.Fields[field] = nil
		return nil
	}
	copied := append([]Value(nil), values...)
	e.fields[field] = copied
	e.changes.Fields[field] = copied
	return nil
}

func (e *Editor) AppendField(field Field, value Value) error {
	existing, _ := e.Get(field)
	return e.SetField(field, append(existing, value)...)
}

func (e *Editor) RemoveField(field Field) error { return e.SetField(field) }

// RemoveValue removes the first equal occurrence and preserves repeated values.
func (e *Editor) RemoveValue(field Field, value Value) error {
	if definition, exists := Lookup(field); exists && !definition.Mutable {
		return fmt.Errorf("%w: %s", ErrReadOnly, field)
	}
	if err := Validate(field, value); err != nil {
		return err
	}
	existing, _ := e.Get(field)
	for i, v := range existing {
		if v.Equal(value) {
			return e.SetField(field, append(existing[:i], existing[i+1:]...)...)
		}
	}
	return nil
}

// Populated returns an independent snapshot of all populated canonical fields.
func (e *Editor) Populated() map[Field][]Value { return cloneMap(e.fields) }

func (e *Editor) GetCustom(key string) ([]string, bool) {
	key, err := NormalizeCustomKey(key)
	if err != nil {
		return nil, false
	}
	values := e.custom[key]
	return append([]string(nil), values...), len(values) > 0
}

// SetCustom replaces the custom key's values. Custom keys are case-insensitive
// ASCII and are stored in a namespace independent of native field names.
func (e *Editor) SetCustom(key string, values ...string) error {
	key, err := NormalizeCustomKey(key)
	if err != nil {
		return err
	}
	for _, v := range values {
		if !validText(v) {
			return fmt.Errorf("%w: invalid custom text for %s", ErrInvalidValue, key)
		}
	}
	if len(values) == 0 {
		delete(e.custom, key)
		e.changes.Custom[key] = nil
		return nil
	}
	copied := append([]string(nil), values...)
	e.custom[key] = copied
	e.changes.Custom[key] = copied
	return nil
}

func (e *Editor) AppendCustom(key, value string) error {
	existing, _ := e.GetCustom(key)
	return e.SetCustom(key, append(existing, value)...)
}

func (e *Editor) RemoveCustom(key string) error { return e.SetCustom(key) }

func (e *Editor) RemoveCustomValue(key, value string) error {
	if _, err := NormalizeCustomKey(key); err != nil {
		return err
	}
	existing, _ := e.GetCustom(key)
	for i, v := range existing {
		if v == value {
			return e.SetCustom(key, append(existing[:i], existing[i+1:]...)...)
		}
	}
	return nil
}

// Custom returns an independent snapshot of the namespaced custom fields.
func (e *Editor) Custom() map[string][]string { return cloneMap(e.custom) }

func (e *Editor) ClearCustom() {
	for key := range e.custom {
		_ = e.RemoveCustom(key)
	}
}

// ClearAll removes editable canonical and custom values. When preserveUnknown
// is false, Save also removes unknown descriptive metadata where safe and may
// return ErrUnsupported if the container cannot distinguish it from structure.
func (e *Editor) ClearAll(preserveUnknown bool) {
	for field := range e.fields {
		if d, ok := Lookup(field); ok && d.Mutable {
			_ = e.RemoveField(field)
		}
	}
	e.ClearCustom()
	e.changes.ClearUnknown = e.changes.ClearUnknown || !preserveUnknown
}

// Save renders cumulative edits to w. It never closes w. Some adapters replace the contents of an *os.File. A
// failed write may leave partial output. Use a separate file and rename once
// saved when replacing a source.
func (e *Editor) Save(w io.Writer) error {
	if w == nil {
		return fmt.Errorf("%w: nil destination", ErrUnsupported)
	}
	if e.save == nil {
		return fmt.Errorf("%w: editor cannot save", ErrUnsupported)
	}
	return e.save(w, Changes{Fields: cloneMap(e.changes.Fields), Custom: cloneMap(e.changes.Custom), ClearUnknown: e.changes.ClearUnknown})
}

func SortedCustomKeys(custom map[string][]string) []string {
	keys := make([]string, 0, len(custom))
	for key := range custom {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
