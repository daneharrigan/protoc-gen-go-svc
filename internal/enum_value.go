package internal

import (
	"google.golang.org/protobuf/compiler/protogen"

	"github.com/dane/protoc-gen-go-svc/internal/options"
)

type EnumValue struct {
	IsLatest     bool
	IsDeprecated bool
	IsPrivate    bool
	Name         string
	Next         *EnumValue
	Private      *EnumValue
	Receive      []*EnumValue
}

// NewEnumValue creates a `EnumValue`. An error will be returned if the
// enum value cannot be created for any reason.
func NewEnumValue(f *Field, value *protogen.EnumValue) (*EnumValue, error) {
	v := &EnumValue{
		IsLatest:     f.IsLatest,
		IsPrivate:    f.IsPrivate,
		IsDeprecated: f.IsDeprecated,
		Name:         value.GoIdent.GoName,
	}

	// Private enum values are the last in the service chain.
	if v.IsPrivate {
		return v, nil
	}

	valueName := options.EnumValueName(value)
	receiveNames := options.ReceiveEnumValueNames(value)
	var ok bool

	// Messages of the latest service or deprecated messages read/write directly
	// to the private service.
	if v.IsLatest || v.IsDeprecated {
		v.Private, ok = f.Private.EnumValueByName[valueName]
		if !ok {
			return nil, NewErrEnumValueNotFound(valueName, f.Private)
		}

		// An enum value can have many receiveNames because multiple values
		// of a later version may map to a single enum in the service being
		// constructed.
		for _, name := range receiveNames {
			pv, ok := f.Private.EnumValueByName[name]
			if !ok {
				return nil, NewErrEnumValueNotFound(name, f.Private)
			}

			v.Receive = append(v.Receive, pv)
		}

		return v, nil
	}

	v.Next, ok = f.Next.EnumValueByName[valueName]
	if !ok {
		return nil, NewErrEnumValueNotFound(valueName, f.Next)
	}

	v.Private = v.Next.Private

	for _, name := range receiveNames {
		nv, ok := f.Next.EnumValueByName[name]
		if !ok {
			return nil, NewErrEnumValueNotFound(name, f.Next)
		}

		v.Receive = append(v.Receive, nv)
	}

	return v, nil
}
