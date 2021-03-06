/*
Copyright 2017 Caicloud Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package service

import (
	"context"
	"io"
	"mime/multipart"
	"reflect"
	"strings"

	"github.com/caicloud/nirvana/definition"
)

// ParameterGenerator is used to generate object for a parameter.
type ParameterGenerator interface {
	// Source returns the source generated by current generator.
	Source() definition.Source
	// Validate validates whether defaultValue and target type is valid.
	Validate(name string, defaultValue interface{}, target reflect.Type) error
	// Generate generates an object by data from value container.
	Generate(ctx context.Context, vc ValueContainer, consumers []Consumer, name string, target reflect.Type) (interface{}, error)
}

var generators = map[definition.Source]ParameterGenerator{
	definition.Path:   &PathParameterGenerator{},
	definition.Query:  &QueryParameterGenerator{},
	definition.Header: &HeaderParameterGenerator{},
	definition.Form:   &FormParameterGenerator{},
	definition.File:   &FileParameterGenerator{},
	definition.Body:   &BodyParameterGenerator{},
	definition.Auto:   &AutoParameterGenerator{},
	definition.Prefab: &PrefabParameterGenerator{},
}

// ParameterGeneratorFor gets a parameter generator for specified source.
func ParameterGeneratorFor(source definition.Source) ParameterGenerator {
	return generators[source]
}

// RegisterParameterGenerator register a generator.
func RegisterParameterGenerator(generator ParameterGenerator) error {
	generators[generator.Source()] = generator
	return nil
}

func assignable(defaultValue interface{}, target reflect.Type) error {
	if defaultValue == nil {
		return nil
	}
	value := reflect.ValueOf(defaultValue)
	if !value.Type().AssignableTo(target) {
		return unassignableType.Error(value.Type(), target)
	}
	return nil
}

func convertible(target reflect.Type) error {
	c := ConverterFor(target)
	if c == nil {
		return noConverter.Error(target)
	}
	return nil
}

// PathParameterGenerator is used to generate object by value from path
type PathParameterGenerator struct{}

func (g *PathParameterGenerator) Source() definition.Source { return definition.Path }
func (g *PathParameterGenerator) Validate(name string, defaultValue interface{}, target reflect.Type) error {
	if name == "" {
		return noName.Error(g.Source())
	}
	if err := assignable(defaultValue, target); err != nil {
		return err
	}
	if err := convertible(target); err != nil {
		return err
	}
	return nil
}
func (g *PathParameterGenerator) Generate(ctx context.Context, vc ValueContainer, consumers []Consumer,
	name string, target reflect.Type) (interface{}, error) {
	data, ok := vc.Path(name)
	if !ok || len(data) <= 0 {
		return nil, nil
	}
	if converter := ConverterFor(target); converter != nil {
		return converter(ctx, []string{data})
	}
	return nil, nil
}

// QueryParameterGenerator is used to generate object by value from query string.
type QueryParameterGenerator struct{}

func (g *QueryParameterGenerator) Source() definition.Source { return definition.Query }
func (g *QueryParameterGenerator) Validate(name string, defaultValue interface{}, target reflect.Type) error {
	if name == "" {
		return noName.Error(g.Source())
	}
	if err := assignable(defaultValue, target); err != nil {
		return err
	}
	if err := convertible(target); err != nil {
		return err
	}
	return nil
}
func (g *QueryParameterGenerator) Generate(ctx context.Context, vc ValueContainer, consumers []Consumer,
	name string, target reflect.Type) (interface{}, error) {
	data, ok := vc.Query(name)
	if !ok || len(data) <= 0 {
		return nil, nil
	}
	if converter := ConverterFor(target); converter != nil {
		return converter(ctx, data)
	}
	return nil, nil
}

// HeaderParameterGenerator is used to generate object by value from request header.
type HeaderParameterGenerator struct{}

func (g *HeaderParameterGenerator) Source() definition.Source { return definition.Header }
func (g *HeaderParameterGenerator) Validate(name string, defaultValue interface{}, target reflect.Type) error {
	if name == "" {
		return noName.Error(g.Source())
	}
	if err := assignable(defaultValue, target); err != nil {
		return err
	}
	if err := convertible(target); err != nil {
		return err
	}
	return nil
}

func (g *HeaderParameterGenerator) Generate(ctx context.Context, vc ValueContainer, consumers []Consumer,
	name string, target reflect.Type) (interface{}, error) {
	data, ok := vc.Header(name)
	if !ok || len(data) <= 0 {
		return nil, nil
	}
	if converter := ConverterFor(target); converter != nil {
		return converter(ctx, data)
	}
	return nil, nil
}

// FormParameterGenerator is used to generate object by value from request form.
type FormParameterGenerator struct{}

func (g *FormParameterGenerator) Source() definition.Source { return definition.Form }
func (g *FormParameterGenerator) Validate(name string, defaultValue interface{}, target reflect.Type) error {
	if name == "" {
		return noName.Error(g.Source())
	}
	if err := assignable(defaultValue, target); err != nil {
		return err
	}
	if err := convertible(target); err != nil {
		return err
	}
	return nil
}

func (g *FormParameterGenerator) Generate(ctx context.Context, vc ValueContainer, consumers []Consumer,
	name string, target reflect.Type) (interface{}, error) {
	data, ok := vc.Form(name)
	if !ok || len(data) <= 0 {
		return nil, nil
	}
	if converter := ConverterFor(target); converter != nil {
		return converter(ctx, data)
	}
	return nil, nil
}

type repeatableCloserForFile struct {
	multipart.File
	closed bool
}

func (c *repeatableCloserForFile) Close() error {
	if c.closed {
		return nil
	}
	err := c.File.Close()
	if err != nil {
		return err
	}
	c.closed = true
	return nil
}

// FileParameterGenerator is used to generate file reader by value from request form file.
type FileParameterGenerator struct {
}

func (g *FileParameterGenerator) Source() definition.Source { return definition.File }
func (g *FileParameterGenerator) Validate(name string, defaultValue interface{}, target reflect.Type) error {
	if name == "" {
		return noName.Error(g.Source())
	}
	err := assignable(defaultValue, target)
	if err != nil {
		return err
	}
	if !reflect.TypeOf((*multipart.File)(nil)).Elem().AssignableTo(target) {
		return unassignableType.Error("multipart.File", target)
	}
	return nil
}

func (g *FileParameterGenerator) Generate(ctx context.Context, vc ValueContainer, consumers []Consumer,
	name string, target reflect.Type) (interface{}, error) {
	file, ok := vc.File(name)
	if !ok {
		return nil, nil
	}
	return repeatableCloserForFile{file, false}, nil
}

// BodyParameterGenerator is used to generate object or body reader by value from request body.
type BodyParameterGenerator struct{}

func (g *BodyParameterGenerator) Source() definition.Source { return definition.Body }
func (g *BodyParameterGenerator) Validate(name string, defaultValue interface{}, target reflect.Type) error {
	err := assignable(defaultValue, target)
	if err != nil {
		return err
	}
	kind := target.Kind()
	switch {
	case kind == reflect.String:
	case kind == reflect.Slice:
	case kind == reflect.Struct:
	case kind == reflect.Ptr && target.Elem().Kind() == reflect.Struct:
	case kind == reflect.Interface && reflect.TypeOf((*io.ReadCloser)(nil)).Elem().AssignableTo(target):
	default:
		return invalidBodyType.Error(target)
	}
	return nil
}

type repeatableCloserForBody struct {
	io.ReadCloser
	closed bool
}

func (c *repeatableCloserForBody) Close() error {
	if c.closed {
		return nil
	}
	err := c.ReadCloser.Close()
	if err != nil {
		return err
	}
	c.closed = true
	return nil
}

func (g *BodyParameterGenerator) Generate(ctx context.Context, vc ValueContainer, consumers []Consumer,
	name string, target reflect.Type) (interface{}, error) {
	reader, contentType, ok := vc.Body()
	if !ok {
		return nil, nil
	}
	kind := target.Kind()
	if kind == reflect.Interface && reflect.TypeOf((*io.ReadCloser)(nil)).Elem().AssignableTo(target) {
		return &repeatableCloserForBody{reader, false}, nil
	}
	var consumer Consumer
	for _, c := range consumers {
		if c.ContentType() == contentType {
			consumer = c
			break
		}
	}
	if consumer == nil {
		return nil, nil
	}
	var value reflect.Value
	// Create a pointer to target.
	// Value is a pointer and it points to nil.
	// consumer should fill it.
	switch {
	case kind == reflect.String || kind == reflect.Slice || kind == reflect.Struct:
		value = reflect.New(target)
	case kind == reflect.Ptr && target.Elem().Kind() == reflect.Struct:
		value = reflect.New(target.Elem())
	default:
		return nil, nil
	}
	if err := consumer.Consume(reader, value.Interface()); err != nil {
		return nil, err
	}
	if kind == reflect.String || kind == reflect.Slice || kind == reflect.Struct {
		return value.Elem().Interface(), nil
	}
	return value.Interface(), nil
}

// PrefabParameterGenerator is used to generate object by prefabs.
type PrefabParameterGenerator struct{}

func (g *PrefabParameterGenerator) Source() definition.Source { return definition.Prefab }
func (g *PrefabParameterGenerator) Validate(name string, defaultValue interface{}, target reflect.Type) error {
	if name == "" {
		return noName.Error(g.Source())
	}
	err := assignable(defaultValue, target)
	if err != nil {
		return err
	}
	prefab := PrefabFor(name)
	if prefab == nil {
		return noPrefab.Error(name)
	}
	if !prefab.Type().AssignableTo(target) {
		return unassignableType.Error(prefab.Type(), target)
	}
	return nil
}

func (g *PrefabParameterGenerator) Generate(ctx context.Context, vc ValueContainer, consumers []Consumer,
	name string, target reflect.Type) (interface{}, error) {
	prefab := PrefabFor(name)
	if prefab == nil {
		return nil, noPrefab.Error(name)
	}
	return prefab.Make(ctx)
}

// AutoParameterGenerator generates an object from a struct type. The fields in a struct can have tag.
// Tag name is "source". Its value format is "Source,Name".
//
// ex.
// type Example struct {
//     Start       int    `source:"Query,start"`
//     ContentType string `source:"Header,Content-Type"`
// }
type AutoParameterGenerator struct{}

type autoTagParams map[paramsKey]string
type paramsKey string

const (
	keyDefault paramsKey = "default"
)

func (params autoTagParams) get(key paramsKey) string {
	return params[key]
}

func (params autoTagParams) set(key paramsKey, value string) {
	params[key] = value
}

func (g *AutoParameterGenerator) Source() definition.Source { return definition.Auto }
func (g *AutoParameterGenerator) Validate(name string, defaultValue interface{}, target reflect.Type) error {
	err := assignable(defaultValue, target)
	if err != nil {
		return err
	}
	if target.Kind() != reflect.Struct && !(target.Kind() == reflect.Ptr && target.Elem().Kind() == reflect.Struct) {
		return invalidAutoParameter.Error(target)
	}
	f := func(index []int, field reflect.StructField) error {
		source, name, params, err := g.split(field.Tag.Get("source"))
		if err != nil {
			return err
		}
		generator := ParameterGeneratorFor(definition.Source(source))
		if generator == nil {
			return noParameterGenerator.Error(source)
		}

		var value interface{}
		defaultValue := params.get(keyDefault)
		if defaultValue != "" {
			if c := ConverterFor(field.Type); c != nil {
				var err error
				value, err = c(context.Background(), []string{defaultValue})
				if err != nil {
					return err
				}
			}
		}

		return generator.Validate(name, value, field.Type)
	}
	if target.Kind() == reflect.Struct {
		err = g.enum([]int{}, target, f)
	} else {
		err = g.enum([]int{}, target.Elem(), f)
	}
	return err
}

func (g *AutoParameterGenerator) split(tag string) (source, name string, atp autoTagParams, err error) {
	atp = make(autoTagParams)
	result := strings.Split(tag, ",")

	length := len(result)

	if length < 1 {
		return "", "", nil, invalidFieldTag.Error(tag)
	}

	if length >= 1 {
		source = strings.Title(strings.ToLower(strings.TrimSpace(result[0])))
	}
	if length >= 2 {
		name = strings.TrimSpace(result[1])
	}
	if length >= 3 {
		params := result[2:]

		for _, param := range params {
			keyValue := strings.Split(param, "=")
			if len(keyValue) == 2 {
				key := paramsKey(strings.TrimSpace(keyValue[0]))
				value := strings.TrimSpace(keyValue[1])
				if key == keyDefault {
					atp.set(key, value)
				}
			}
		}
	}

	return
}

func (g *AutoParameterGenerator) Generate(ctx context.Context, vc ValueContainer, consumers []Consumer, name string, target reflect.Type) (interface{}, error) {
	var result reflect.Value
	var value reflect.Value
	if target.Kind() == reflect.Struct {
		result = reflect.New(target).Elem()
		value = result
	} else {
		result = reflect.New(target.Elem())
		value = result.Elem()
	}
	if err := g.generate(ctx, vc, consumers, value); err != nil {
		return nil, err
	}
	return result.Interface(), nil
}

func (g *AutoParameterGenerator) generate(ctx context.Context, vc ValueContainer, consumers []Consumer, value reflect.Value) error {
	f := func(index []int, field reflect.StructField) error {
		source, name, params, err := g.split(field.Tag.Get("source"))
		if err != nil {
			return err
		}
		generator := ParameterGeneratorFor(definition.Source(source))
		if generator == nil {
			return noParameterGenerator.Error(source)
		}
		ins, err := generator.Generate(ctx, vc, consumers, name, field.Type)
		if err != nil {
			return err
		}

		defaultValue := params.get(keyDefault)
		if ins == nil && defaultValue != "" {
			if c := ConverterFor(field.Type); c != nil {
				// After passing the validation phase, here will never return an error
				ins, _ = c(ctx, []string{defaultValue}) // #nosec
			}
		}

		if ins != nil {
			value.FieldByIndex(index).Set(reflect.ValueOf(ins))
		}

		return nil
	}
	return g.enum([]int{}, value.Type(), f)
}

func (g *AutoParameterGenerator) enum(index []int, typ reflect.Type, f func(index []int, field reflect.StructField) error) error {
	var err error
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if field.Tag.Get("source") != "" {
			err = f(append(index, i), field)
		} else if field.Type.Kind() == reflect.Struct {
			err = g.enum(append(index, i), field.Type, f)
		}
		if err != nil {
			return err
		}
	}
	return nil
}
