// Copyright 2017 Axel Springer SE
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package kv

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"github.com/docker/libkv/store"
)

// Transcode takes an interface and uses reflection
// to fill it with data from a kv.
func Transcode(s interface{}, prefix string, kv store.Store) error {
	config := &TranscoderConfig{
		Prefix:   prefix,
		KV:       kv,
		Metadata: nil,
		Result:   s,
	}

	transcoder, err := NewTranscoder(config)
	if err != nil {
		return err
	}

	return transcoder.Transcode()
}

// NewTranscoder returns a new transcoder for the given configuration.
// Once a transcoder has been returned, the same configuration must not be used
// again.
func NewTranscoder(config *TranscoderConfig) (*Transcoder, error) {
	val := reflect.ValueOf(config.Result)
	if val.Kind() != reflect.Ptr {
		return nil, errors.New("result must be a pointer")
	}

	val = val.Elem()
	if !val.CanAddr() {
		return nil, errors.New("result must be addressable (a pointer)")
	}

	if config.Metadata != nil {
		if config.Metadata.Keys == nil {
			config.Metadata.Keys = make([]string, 0)
		}

		if config.Metadata.Unused == nil {
			config.Metadata.Unused = make([]string, 0)
		}
	}

	if config.TagName == "" {
		config.TagName = "kvstructure"
	}

	result := &Transcoder{
		config: config,
	}

	return result, nil
}

// Transcode transcodes a given raw interface to a filled structure
func (t *Transcoder) Transcode() error {
	return t.transcode("", reflect.ValueOf(t.config.Result).Elem())
}

// transcode is doing the heavy lifting in the background
func (t *Transcoder) transcode(name string, val reflect.Value) error {
	var err error
	valKind := getKind(val)
	switch valKind {
	case reflect.String:
		err = t.transcodeString(name, val)
	case reflect.Bool:
		err = t.transcodeBool(name, val)
	case reflect.Int:
		err = t.transcodeInt(name, val)
	case reflect.Uint:
		err = t.transcodeUint(name, val)
	case reflect.Float32:
		err = t.transcodeFloat(name, val)
	case reflect.Struct:
		err = t.transcodeStruct(val)
	default:
		// we have to work on here for value to pointed to
		return fmt.Errorf("Unsupported type %s", valKind)
	}

	// should be nil
	return err
}

// transcodeBasic transcodes a basic type (bool, int, strinc, etc.)
// and eventually sets it to the retrieved value
func (t *Transcoder) transcodeBasic(val reflect.Value) error {
	return nil
}

// transcodeString
func (t *Transcoder) transcodeString(name string, val reflect.Value) error {
	kvPair, err := t.getKVPair(name)
	if err != nil {
		return err
	}
	kvVal := string(kvPair.Value)

	conv := true
	switch {
	case val.Kind() == reflect.String:
		val.SetString(kvVal)
	default:
		conv = false
	}

	// if conf was not successful
	if !conv {
		return err
	}

	return nil
}

// transcodeBool
func (t *Transcoder) transcodeBool(name string, val reflect.Value) error {
	kvPair, err := t.getKVPair(name)
	if err != nil {
		return err
	}
	kvVal := string(kvPair.Value)

	switch {
	case val.Kind() == reflect.Bool:
		conv, err := strconv.ParseBool(kvVal)
		if err != nil {
			return err
		}
		val.SetBool(conv)
	default:
		return fmt.Errorf("'%s' got unconvertible type '%s'", name, val.Type())
	}

	return nil
}

// transcodeInt
func (t *Transcoder) transcodeInt(name string, val reflect.Value) error {
	kvPair, err := t.getKVPair(name)
	if err != nil {
		return err
	}
	kvVal := string(kvPair.Value)

	switch {
	case val.Kind() == reflect.Int:
		conv, err := strconv.ParseInt(kvVal, 10, 64)
		if err != nil {
			return err
		}
		val.SetInt(conv)
	case val.Kind() == reflect.Uint:
		conv, err := strconv.ParseUint(kvVal, 10, 64)
		if err != nil {
			return err
		}
		val.SetUint(conv)
	default:
		return fmt.Errorf("'%s' got unconvertible type '%s'", name, val.Type())
	}

	return nil
}

// transcodeUint
func (t *Transcoder) transcodeUint(name string, val reflect.Value) error {
	kvPair, err := t.getKVPair(name)
	if err != nil {
		return err
	}
	kvVal := string(kvPair.Value)

	switch {
	case val.Kind() == reflect.Uint:
		conv, err := strconv.ParseUint(kvVal, 10, 64)
		if err != nil {
			return err
		}
		val.SetUint(conv)
	default:
		return fmt.Errorf("'%s' got unconvertible type '%s'", name, val.Type())
	}

	return nil
}

// transcodeFloat32
func (t *Transcoder) transcodeFloat(name string, val reflect.Value) error {
	kvPair, err := t.getKVPair(name)
	if err != nil {
		return err
	}
	kvVal := string(kvPair.Value)

	switch {
	case val.Kind() == reflect.Float32:
		conv, err := strconv.ParseFloat(kvVal, 64)
		if err != nil {
			return err
		}
		val.SetFloat(conv)
	default:
		return fmt.Errorf("'%s' got unconvertible type '%s'", name, val.Type())
	}

	return nil
}

// transcodeStruct
func (t *Transcoder) transcodeStruct(val reflect.Value) error {
	valInterface := reflect.Indirect(val)
	valType := valInterface.Type()

	var wg sync.WaitGroup
	wg.Add(valType.NumField())

	errors := make([]error, 0)

	// The slice will keep track of all structs we'll be transcoding.
	// There can be more structs, if we have embedded structs that are squashed.
	structs := make([]reflect.Value, 1, 5)
	structs[0] = val

	type field struct {
		field reflect.StructField
		val   reflect.Value
		json  bool
	}
	fields := []field{}
	for len(structs) > 0 { // could be easier
		structVal := structs[0]
		structs = structs[1:]
		// here we should do squashing

		for i := 0; i < valType.NumField(); i++ {
			fieldType := valType.Field(i)
			// json is somehow special
			// it is curated by golang json
			isJSON := false
			// fieldKind := fieldType.Type.Kind()

			tagParts := strings.Split(fieldType.Tag.Get(t.config.TagName), ",")
			for _, tag := range tagParts[1:] {
				// test here for squashing
				if tag == "json" {
					isJSON = true
				}
			}

			fields = append(fields, field{fieldType, structVal.Field(i), isJSON})
		}
	}

	for _, f := range fields {
		field, val, isJSON := f.field, f.val, f.json
		kv := field.Name

		tag := field.Tag.Get(t.config.TagName)
		tag = strings.SplitN(tag, ",", 2)[0]
		if tag != "" {
			kv = tag
		}

		if !val.CanSet() {
			wg.Done()

			continue
		}

		// we deal with
		if isJSON {
			// remove field from group
			wg.Done()

			if !val.CanAddr() {
				continue
			}

			kvPair, err := t.getKVPair(kv)
			if err != nil {
				errors = append(errors, err)
				continue
			}

			obj := reflect.New(field.Type).Interface()
			if err := json.Unmarshal(kvPair.Value, &obj); err != nil {
				errors = append(errors, err)
			}

			val.Set(reflect.ValueOf(obj).Elem())

			continue
		}

		go func() {
			defer wg.Done()
			if err := t.transcode(kv, val); err != nil {
				errors = append(errors, err)
			}
		}()
	}

	wg.Wait()

	return nil
}

func (t *Transcoder) getKVPair(key string) (*store.KVPair, error) {
	kvPair, err := t.config.KV.Get(trailingSlash(t.config.Prefix) + key)
	if err != nil {
		return nil, err
	}

	return kvPair, nil
}

// getKind is returning the kind of the reflected value
func getKind(val reflect.Value) reflect.Kind {
	kind := val.Kind()

	switch {
	case kind >= reflect.Int && kind <= reflect.Int64:
		return reflect.Int
	case kind >= reflect.Uint && kind <= reflect.Uint64:
		return reflect.Uint
	case kind >= reflect.Float32 && kind <= reflect.Float64:
		return reflect.Float32
	default:
		return kind
	}
}
