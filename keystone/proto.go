package keystone

import (
	"github.com/kubex/definitions-go/k4"
	"github.com/kubex/keystone-go/proto"
)

func EntityFromProto(p *proto.EntityResponse) *Entity {
	e := &Entity{}
	e.ID = k4.IDFromString(p.EntityId)
	e.Schema = GetScopedKey(p.Schema)
	e.Properties = make(map[string]Property)
	for _, prop := range p.Properties {
		t := prop.GetValue().GetTime().AsTime()
		e.Properties[prop.Name] = Property{
			Name:           prop.GetName(),
			Type:           GetPropertyType(prop.GetValue().GetType()),
			Classification: GetClassification(prop.GetClassification()),
			Text:           prop.GetValue().GetText(),
			Secret:         prop.GetValue().GetSecureText(),
			Int:            prop.GetValue().GetInt(),
			Bool:           prop.GetValue().GetBool(),
			Float:          float64(prop.GetValue().GetFloat()),
			Time:           &t,
			indexed:        prop.Indexed,
			lookup:         prop.Lookup,
		}
	}
	return e
}
