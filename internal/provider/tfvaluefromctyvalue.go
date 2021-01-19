package provider

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5/tftypes"
	"github.com/zclconf/go-cty/cty"
)

func tfValueFromCtyValue(val cty.Value) (*tftypes.Value, tftypes.Type, error) {
	typ := val.Type()
	switch {
	case typ.Equals(cty.String):
		v := tftypes.NewValue(tftypes.String, val.AsString())
		return &v, tftypes.String, nil
	case typ.Equals(cty.Number):
		v := tftypes.NewValue(tftypes.Number, val.AsBigFloat())
		return &v, tftypes.Number, nil
	case typ.Equals(cty.Bool):
		v := tftypes.NewValue(tftypes.Bool, val)
		return &v, tftypes.Bool, nil
	case typ.IsSetType():
		vals := make([]tftypes.Value, 0)
		for it := val.ElementIterator(); it.Next(); {
			_, ev := it.Element()
			v, _, err := tfValueFromCtyValue(ev)
			if err != nil {
				return nil, nil, err
			}
			vals = append(vals, *v)
		}
		t, err := tftypes.TypeFromElements(vals)
		if err != nil {
			return nil, nil, err
		}
		v := tftypes.NewValue(tftypes.Set{
			ElementType: t,
		}, vals)
		return &v, t, nil
	case typ.IsListType():
		vals := make([]tftypes.Value, 0)
		for it := val.ElementIterator(); it.Next(); {
			_, ev := it.Element()
			v, _, err := tfValueFromCtyValue(ev)
			if err != nil {
				return nil, nil, err
			}
			vals = append(vals, *v)
		}
		t, err := tftypes.TypeFromElements(vals)
		if err != nil {
			return nil, nil, err
		}
		v := tftypes.NewValue(tftypes.List{
			ElementType: t,
		}, vals)
		return &v, t, nil
	case typ.IsTupleType():
		typs := make([]tftypes.Type, 0)
		vals := make([]tftypes.Value, 0)
		for it := val.ElementIterator(); it.Next(); {
			_, ev := it.Element()
			v, t, err := tfValueFromCtyValue(ev)
			if err != nil {
				return nil, nil, err
			}
			typs = append(typs, t)
			vals = append(vals, *v)
		}
		t := tftypes.Tuple{
			ElementTypes: typs,
		}
		v := tftypes.NewValue(t, vals)
		return &v, t, nil
	case typ.IsMapType():
		vals := map[string]tftypes.Value{}
		for it := val.ElementIterator(); it.Next(); {
			k, ev := it.Element()
			rawK := k.AsString()
			v, _, err := tfValueFromCtyValue(ev)
			if err != nil {
				return nil, nil, err
			}
			vals[rawK] = *v
		}
		t := tftypes.Map{
			AttributeType: tftypes.String,
		}
		v := tftypes.NewValue(t, vals)
		return &v, t, nil
	case typ.IsObjectType():
		typs := make(map[string]tftypes.Type)
		vals := make(map[string]tftypes.Value)
		for it := val.ElementIterator(); it.Next(); {
			k, ev := it.Element()
			rawK := k.AsString()
			v, t, err := tfValueFromCtyValue(ev)
			if err != nil {
				return nil, nil, err
			}
			typs[rawK] = t
			vals[rawK] = *v
		}
		t := tftypes.Object{
			AttributeTypes: typs,
		}
		v := tftypes.NewValue(t, vals)
		return &v, t, nil
	default:
		return nil, nil, fmt.Errorf("unknown cty type %s", typ.GoString())
	}
}
