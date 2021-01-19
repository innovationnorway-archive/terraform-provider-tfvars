package provider

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5/tftypes"
)

type dataFile struct {
	p *provider
}

func newDataFile(p *provider) (*dataFile, error) {
	if p == nil {
		return nil, fmt.Errorf("a provider is required")
	}

	return &dataFile{
		p: p,
	}, nil
}

func (d *dataFile) Schema(context.Context) *tfprotov5.Schema {
	return &tfprotov5.Schema{
		Block: &tfprotov5.SchemaBlock{
			Description:     "Use this data source to read Terraform variable definitions (.tfvars) files.",
			DescriptionKind: tfprotov5.StringKindMarkdown,
			Attributes: []*tfprotov5.SchemaAttribute{
				{
					Name:       "id",
					Computed:   true,
					Deprecated: true,
					Description: "This attribute is only present for some compatibility issues and should not be used. It " +
						"will be removed in a future version.",
					DescriptionKind: tfprotov5.StringKindMarkdown,
					Type:            tftypes.String,
				},
				{
					Name:            "filename",
					Required:        true,
					Description:     "The path to the variable definitions (`.tfvars`) file.",
					DescriptionKind: tfprotov5.StringKindMarkdown,
					Type:            tftypes.String,
				},
				{
					Name:            "variables",
					Description:     "An object where the top-level arguments in the variable definitions (`.tfvars`) file are attributes.",
					Computed:        true,
					DescriptionKind: tfprotov5.StringKindMarkdown,
					Type:            tftypes.DynamicPseudoType,
				},
			},
		},
	}
}

func (d *dataFile) Validate(ctx context.Context, config map[string]tftypes.Value) ([]*tfprotov5.Diagnostic, error) {
	return nil, nil
}

func (d *dataFile) Read(ctx context.Context, config map[string]tftypes.Value) (map[string]tftypes.Value, []*tfprotov5.Diagnostic, error) {
	var (
		filename string
		diags    []*tfprotov5.Diagnostic
	)
	err := config["filename"].As(&filename)
	if err != nil {
		return nil, nil, err
	}

	src, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, nil, err
	}

	f, hclDiags := hclsyntax.ParseConfig(src, filename, hcl.Pos{Line: 1, Column: 1})
	if f == nil || f.Body == nil {
		return nil, nil, fmt.Errorf(hclDiags.Error())
	}

	attrs, hclDiags := f.Body.JustAttributes()
	if len(hclDiags) != 0 {
		for _, hclDiag := range hclDiags {
			diags = append(diags, &tfprotov5.Diagnostic{
				Severity: tfprotov5.DiagnosticSeverity(hclDiag.Severity),
				Summary:  hclDiag.Summary,
				Detail:   hclDiag.Detail,
			})
		}

		return nil, diags, nil
	}

	types := map[string]tftypes.Type{}
	values := map[string]tftypes.Value{}
	for name, attr := range attrs {
		val, hclDiags := attr.Expr.Value(nil)
		if len(hclDiags) != 0 {
			for _, hclDiag := range hclDiags {
				diags = append(diags, &tfprotov5.Diagnostic{
					Severity: tfprotov5.DiagnosticSeverity(hclDiag.Severity),
					Summary:  hclDiag.Summary,
					Detail:   hclDiag.Detail,
				})
			}

			return nil, diags, nil
		}

		v, t, err := tfValueFromCtyValue(val)
		if err != nil {
			return nil, nil, err
		}

		types[name] = t
		values[name] = *v
	}

	return map[string]tftypes.Value{
		"id":       config["filename"],
		"filename": config["filename"],
		"variables": tftypes.NewValue(tftypes.Object{
			AttributeTypes: types,
		}, values),
	}, nil, nil
}
