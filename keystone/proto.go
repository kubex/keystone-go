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
		t := prop.GetTime().AsTime()
		e.Properties[prop.Name] = Property{
			Name:           prop.GetName(),
			Type:           GetPropertyType(prop.GetType()),
			Classification: GetClassification(prop.GetClassification()),
			Text:           prop.GetText(),
			Secret:         prop.GetSecureText(),
			Int:            prop.GetInt(),
			Bool:           prop.GetBool(),
			Float:          float64(prop.GetFloat()),
			Time:           &t,
			indexed:        prop.Indexed,
			lookup:         prop.Lookup,
		}
	}
	return e
}
