//go:generate stringer -type Type

package milestone

import (
	"fmt"
	"strings"
)

type Type int

const (
	unknown Type = iota

	InstanceCreated
	AuthenticationSucceededOnInstance
	ProjectCreated
	ApplicationCreated
	AuthenticationSucceededOnApplication
	InstanceDeleted
	CustomInstanceDomainConfigured

	IDPCreated
	IDPActivated
	IDPLinked
	IDPSignIn

	BrandingConfigured
	BrandingActivated

	SMTPConfigured

	B2BOrgCreated
	B2BProjectGranted
	B2BUserGranted

	typesCount
)

func AllTypes() []Type {
	types := make([]Type, typesCount-1)
	for i := Type(1); i < typesCount; i++ {
		types[i-1] = i
	}
	return types
}

func (t *Type) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, t.String())), nil
}

func (t *Type) UnmarshalJSON(data []byte) error {
	*t = typeFromString(strings.Trim(string(data), `"`))
	return nil
}

func typeFromString(t string) Type {
	idx := strings.Index(_Type_name, t)
	if idx <= 0 {
		return unknown
	}

	for i, typeIdx := range _Type_index {
		if int(typeIdx) != idx {
			continue
		}
		return Type(i)
	}
	return unknown
}
