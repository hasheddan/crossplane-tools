/*
Copyright 2019 The Crossplane Authors.

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

// Package fields defines and matches common struct fields.
package fields

import (
	"go/types"
	"strings"
)

// Field names.
const (
	NameTypeMeta                     = "TypeMeta"
	NameObjectMeta                   = "ObjectMeta"
	NameListMeta                     = "ListMeta"
	NameSpec                         = "Spec"
	NameSpecTemplate                 = "SpecTemplate"
	NameStatus                       = "Status"
	NameResourceSpec                 = "ResourceSpec"
	NameResourceStatus               = "ResourceStatus"
	NameResourceClaimSpec            = "ResourceClaimSpec"
	NameNonPortableClassSpecTemplate = "NonPortableClassSpecTemplate"
	NamePortableClass                = "PortableClass"
	NameItems                        = "Items"
)

// Field type suffixes.
const (
	TypeSuffixTypeMeta                     = "k8s.io/apimachinery/pkg/apis/meta/v1.TypeMeta"
	TypeSuffixObjectMeta                   = "k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"
	TypeSuffixListMeta                     = "k8s.io/apimachinery/pkg/apis/meta/v1.ListMeta"
	TypeSuffixSpec                         = NameSpec
	TypeSuffixSpecTemplate                 = NameSpecTemplate
	TypeSuffixStatus                       = NameStatus
	TypeSuffixResourceSpec                 = "github.com/crossplaneio/crossplane-runtime/apis/core/v1alpha1.ResourceSpec"
	TypeSuffixResourceStatus               = "github.com/crossplaneio/crossplane-runtime/apis/core/v1alpha1.ResourceStatus"
	TypeSuffixResourceClaimSpec            = "github.com/crossplaneio/crossplane-runtime/apis/core/v1alpha1.ResourceClaimSpec"
	TypeSuffixResourceClaimStatus          = "github.com/crossplaneio/crossplane-runtime/apis/core/v1alpha1.ResourceClaimStatus"
	TypeSuffixNonPortableClassSpecTemplate = "github.com/crossplaneio/crossplane-runtime/apis/core/v1alpha1.NonPortableClassSpecTemplate"
	TypeSuffixPortableClass                = "github.com/crossplaneio/crossplane-runtime/apis/core/v1alpha1.PortableClass"
)

func matches(s *types.Struct, m Matcher) bool {
	for i := 0; i < s.NumFields(); i++ {
		if m(s.Field(i)) {
			return true
		}
	}
	return false
}

// A Matcher is a function that returns true if the supplied Var (assumed to be
// a struct field) matches its requirements.
type Matcher func(f *types.Var) bool

// And chains the original Matcher o with a new Matcher m.
func (o Matcher) And(m Matcher) Matcher {
	return func(f *types.Var) bool {
		return o(f) && m(f)
	}
}

// Has returns true if the supplied Object's underlying type is struct (or a
// slice or map of struct), and it matches all of the supplied field Matchers.
func Has(o types.Object, m ...Matcher) bool {
	s := findStruct(o)
	if s == nil {
		return false
	}
	for _, matcher := range m {
		if !matches(s, matcher) {
			return false
		}
	}
	return true
}

func findStruct(o types.Object) *types.Struct {
	switch t := o.Type().Underlying().(type) {
	case *types.Struct:
		return t
	case *types.Slice:
		s, ok := t.Elem().Underlying().(*types.Struct)
		if !ok {
			return nil
		}
		return s
	case *types.Map:
		s, ok := t.Elem().Underlying().(*types.Struct)
		if !ok {
			return nil
		}
		return s
	}
	return nil
}

// IsEmbedded returns a Matcher that returns true if the supplied field is
// embedded.
func IsEmbedded() Matcher {
	return func(f *types.Var) bool {
		return f.Embedded()
	}
}

// IsSlice returns a Matcher that returns true if the supplied field is a slice.
func IsSlice() Matcher {
	return func(f *types.Var) bool {
		_, ok := f.Type().Underlying().(*types.Slice)
		return ok
	}
}

// IsNamed returns a Matcher that returns true if the supplied field has the
// supplied name.
func IsNamed(name string) Matcher {
	return func(f *types.Var) bool {
		if !f.IsField() {
			return false
		}
		return f.Name() == name
	}
}

// IsTypeNamed returns a Matcher that returns true if the supplied field has the
// supplied type name suffix and name.
func IsTypeNamed(typeNameSuffix, name string) Matcher {
	return func(f *types.Var) bool {
		if !IsNamed(name)(f) {
			return false
		}
		return strings.HasSuffix(f.Type().String(), typeNameSuffix)
	}
}

// HasFieldThat returns a Matcher that returns true if the supplied field is a
// struct that matches the supplied field matchers.
func HasFieldThat(m ...Matcher) Matcher {
	return func(f *types.Var) bool {
		return Has(f, m...)
	}
}

// IsTypeMeta returns a Matcher that returns true if the supplied field appears
// to be Kubernetes type metadata.
func IsTypeMeta() Matcher { return IsTypeNamed(TypeSuffixTypeMeta, NameTypeMeta) }

// IsObjectMeta returns a Matcher that returns true if the supplied field
// appears to be Kubernetes object metadata.
func IsObjectMeta() Matcher { return IsTypeNamed(TypeSuffixObjectMeta, NameObjectMeta) }

// IsListMeta returns a Matcher that returns true if the supplied field appears
// to be Kubernetes list metadata.
func IsListMeta() Matcher { return IsTypeNamed(TypeSuffixListMeta, NameListMeta) }

// IsSpec returns a Matcher that returns true if the supplied field appears to
// be a Kubernetes resource spec.
func IsSpec() Matcher { return IsTypeNamed(NameSpec, TypeSuffixSpec) }

// IsSpecTemplate returns a Matcher that returns true if the supplied field
// appears to be a Crossplane resource class spec template.
func IsSpecTemplate() Matcher { return IsTypeNamed(NameSpecTemplate, TypeSuffixSpecTemplate) }

// IsStatus returns a Matcher that returns true if the supplied field appears to
// be a Kubernetes resource status.
func IsStatus() Matcher { return IsTypeNamed(NameStatus, TypeSuffixStatus) }

// IsResourceSpec returns a Matcher that returns true if the supplied field
// appears to be a Crossplane managed resource spec.
func IsResourceSpec() Matcher { return IsTypeNamed(TypeSuffixResourceSpec, NameResourceSpec) }

// IsResourceStatus returns a Matcher that returns true if the supplied field
// appears to be a Crossplane managed resource status.
func IsResourceStatus() Matcher { return IsTypeNamed(TypeSuffixResourceStatus, NameResourceStatus) }

// IsResourceClaimSpec returns a Matcher that returns true if the supplied field
// appears to be a Crossplane resource claim spec.
func IsResourceClaimSpec() Matcher {
	return IsTypeNamed(TypeSuffixResourceClaimSpec, NameResourceClaimSpec)
}

// IsResourceClaimStatus returns a Matcher that returns true if the supplied
// field appears to be a Crossplane resource claim status.
func IsResourceClaimStatus() Matcher {
	return IsTypeNamed(TypeSuffixResourceClaimStatus, NameStatus)
}

// IsNonPortableClassSpecTemplate returns a Matcher that returns true if the
// supplied field appears to be a Crossplane non-portable resource class spec
// template.
func IsNonPortableClassSpecTemplate() Matcher {
	return IsTypeNamed(TypeSuffixNonPortableClassSpecTemplate, NameNonPortableClassSpecTemplate)
}

// IsPortableClass returns a Matcher that returns true if the supplied field
// appears to be a Crossplane portable resource class.
func IsPortableClass() Matcher {
	return IsTypeNamed(TypeSuffixPortableClass, NamePortableClass)
}

// IsItems returns a Matcher that returns true if the supplied field appears to
// be the Items of a Kubernetes list.
func IsItems() Matcher {
	return IsNamed(NameItems)
}
