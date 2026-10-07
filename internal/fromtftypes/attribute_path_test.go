// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package fromtftypes_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/fromtftypes"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema"
	"github.com/hashicorp/terraform-plugin-framework/internal/testing/testschema"
	"github.com/hashicorp/terraform-plugin-framework/internal/testing/testtypes"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// nestedSchema contains each kind of schema element a path can walk through.
var nestedSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"dynamic": schema.DynamicAttribute{Optional: true},
		"list_nested": schema.ListNestedAttribute{
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"string": schema.StringAttribute{Optional: true},
				},
			},
			Optional: true,
		},
		"map_nested": schema.MapNestedAttribute{
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"string": schema.StringAttribute{Optional: true},
				},
			},
			Optional: true,
		},
		"object": schema.ObjectAttribute{
			AttributeTypes: map[string]attr.Type{
				"list": types.ListType{ElemType: types.StringType},
			},
			Optional: true,
		},
		"set_nested": schema.SetNestedAttribute{
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"string": schema.StringAttribute{Optional: true},
				},
			},
			Optional: true,
		},
		"single_nested": schema.SingleNestedAttribute{
			Attributes: map[string]schema.Attribute{
				"string": schema.StringAttribute{Optional: true},
			},
			Optional: true,
		},
	},
	Blocks: map[string]schema.Block{
		"list_block": schema.ListNestedBlock{
			NestedObject: schema.NestedBlockObject{
				Attributes: map[string]schema.Attribute{
					"set": schema.SetAttribute{
						ElementType: types.StringType,
						Optional:    true,
					},
					"string": schema.StringAttribute{Optional: true},
				},
				Blocks: map[string]schema.Block{
					"list_block": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"string": schema.StringAttribute{Optional: true},
							},
						},
					},
				},
			},
		},
		"set_block": schema.SetNestedBlock{
			NestedObject: schema.NestedBlockObject{
				Attributes: map[string]schema.Attribute{
					"string": schema.StringAttribute{Optional: true},
				},
			},
		},
		"single_block": schema.SingleNestedBlock{
			Attributes: map[string]schema.Attribute{
				"string": schema.StringAttribute{Optional: true},
			},
		},
	},
}

func TestAttributePath(t *testing.T) {
	t.Parallel()

	type testCase struct {
		tfType        *tftypes.AttributePath
		schema        fwschema.Schema
		expected      path.Path
		expectedDiags diag.Diagnostics
	}

	testCases := map[string]testCase{
		"nil": {
			tfType:   nil,
			expected: path.Empty(),
		},
		"empty": {
			tfType:   tftypes.NewAttributePath(),
			expected: path.Empty(),
		},
		"AttributeName": {
			tfType: tftypes.NewAttributePath().WithAttributeName("test"),
			schema: testschema.Schema{
				Attributes: map[string]fwschema.Attribute{
					"test": testschema.Attribute{
						Type: types.StringType,
					},
				},
			},
			expected: path.Root("test"),
		},
		"AttributeName-nonexistent-attribute": {
			tfType: tftypes.NewAttributePath().WithAttributeName("test"),
			schema: testschema.Schema{
				Attributes: map[string]fwschema.Attribute{
					"not-test": testschema.Attribute{
						Type: testtypes.StringType{},
					},
				},
			},
			expected: path.Empty(),
			expectedDiags: diag.Diagnostics{
				diag.NewErrorDiagnostic(
					"Unable to Convert Attribute Path",
					"An unexpected error occurred while trying to convert an attribute path. "+
						"This is an error in terraform-plugin-framework used by the provider. "+
						"Please report the following to the provider developers.\n\n"+
						"Attribute Path: AttributeName(\"test\")\n"+
						"Original Error: AttributeName(\"test\") still remains in the path: could not find attribute or block \"test\" in schema",
				),
			},
		},
		"AttributeName-ElementKeyInt": {
			tfType: tftypes.NewAttributePath().WithAttributeName("test").WithElementKeyInt(1),
			schema: testschema.Schema{
				Attributes: map[string]fwschema.Attribute{
					"test": testschema.Attribute{
						Type: types.ListType{
							ElemType: types.StringType,
						},
					},
				},
			},
			expected: path.Root("test").AtListIndex(1),
		},
		"AttributeName-ElementKeyValue": {
			tfType: tftypes.NewAttributePath().WithAttributeName("test").WithElementKeyValue(tftypes.NewValue(tftypes.String, "test-value")),
			schema: testschema.Schema{
				Attributes: map[string]fwschema.Attribute{
					"test": testschema.Attribute{
						Type: types.SetType{
							ElemType: types.StringType,
						},
					},
				},
			},
			expected: path.Root("test").AtSetValue(types.StringValue("test-value")),
		},
		"AttributeName-ElementKeyValue-value-conversion-error": {
			tfType: tftypes.NewAttributePath().WithAttributeName("test").WithElementKeyValue(tftypes.NewValue(tftypes.String, "test-value")),
			schema: testschema.Schema{
				Attributes: map[string]fwschema.Attribute{
					"test": testschema.Attribute{
						Type: types.SetType{
							ElemType: testtypes.InvalidType{},
						},
					},
				},
			},
			expected: path.Empty(),
			expectedDiags: diag.Diagnostics{
				diag.NewErrorDiagnostic(
					"Unable to Convert Attribute Path",
					"An unexpected error occurred while trying to convert an attribute path. "+
						"This is either an error in terraform-plugin-framework or a custom attribute type used by the provider. "+
						"Please report the following to the provider developers.\n\n"+
						"Attribute Path: AttributeName(\"test\").ElementKeyValue(tftypes.String<\"test-value\">)\n"+
						"Original Error: unable to create PathStepElementKeyValue from tftypes.Value: unable to convert tftypes.Value (tftypes.String<\"test-value\">) to attr.Value: intentional ValueFromTerraform error",
				),
			},
		},
		"ListNestedBlock-ElementKeyInt-AttributeName": {
			tfType:   tftypes.NewAttributePath().WithAttributeName("list_block").WithElementKeyInt(0).WithAttributeName("string"),
			schema:   nestedSchema,
			expected: path.Root("list_block").AtListIndex(0).AtName("string"),
		},
		"ListNestedBlock-nested-ListNestedBlock": {
			tfType: tftypes.NewAttributePath().
				WithAttributeName("list_block").WithElementKeyInt(1).
				WithAttributeName("list_block").WithElementKeyInt(2).
				WithAttributeName("string"),
			schema:   nestedSchema,
			expected: path.Root("list_block").AtListIndex(1).AtName("list_block").AtListIndex(2).AtName("string"),
		},
		"ListNestedBlock-SetAttribute-ElementKeyValue": {
			tfType: tftypes.NewAttributePath().
				WithAttributeName("list_block").WithElementKeyInt(0).
				WithAttributeName("set").WithElementKeyValue(tftypes.NewValue(tftypes.String, "test-value")),
			schema:   nestedSchema,
			expected: path.Root("list_block").AtListIndex(0).AtName("set").AtSetValue(types.StringValue("test-value")),
		},
		"ListNestedBlock-ElementKeyString": {
			tfType:   tftypes.NewAttributePath().WithAttributeName("list_block").WithElementKeyString("test"),
			schema:   nestedSchema,
			expected: path.Empty(),
			expectedDiags: diag.Diagnostics{
				diag.NewErrorDiagnostic(
					"Unable to Convert Attribute Path",
					"An unexpected error occurred while trying to convert an attribute path. "+
						"This is an error in terraform-plugin-framework used by the provider. "+
						"Please report the following to the provider developers.\n\n"+
						"Attribute Path: AttributeName(\"list_block\").ElementKeyString(\"test\")\n"+
						"Original Error: ElementKeyString(\"test\") still remains in the path: cannot apply step tftypes.ElementKeyString to ListNestedBlock",
				),
			},
		},
		"ListNestedBlock-ElementKeyInt-nonexistent-attribute": {
			tfType:   tftypes.NewAttributePath().WithAttributeName("list_block").WithElementKeyInt(0).WithAttributeName("nonexistent"),
			schema:   nestedSchema,
			expected: path.Empty(),
			expectedDiags: diag.Diagnostics{
				diag.NewErrorDiagnostic(
					"Unable to Convert Attribute Path",
					"An unexpected error occurred while trying to convert an attribute path. "+
						"This is an error in terraform-plugin-framework used by the provider. "+
						"Please report the following to the provider developers.\n\n"+
						"Attribute Path: AttributeName(\"list_block\").ElementKeyInt(0).AttributeName(\"nonexistent\")\n"+
						"Original Error: AttributeName(\"nonexistent\") still remains in the path: no attribute or block \"nonexistent\" on NestedBlockObject",
				),
			},
		},
		"SetNestedBlock-ElementKeyValue": {
			tfType: tftypes.NewAttributePath().WithAttributeName("set_block").WithElementKeyValue(tftypes.NewValue(
				tftypes.Object{AttributeTypes: map[string]tftypes.Type{"string": tftypes.String}},
				map[string]tftypes.Value{"string": tftypes.NewValue(tftypes.String, "test-value")},
			)),
			schema: nestedSchema,
			expected: path.Root("set_block").AtSetValue(types.ObjectValueMust(
				map[string]attr.Type{"string": types.StringType},
				map[string]attr.Value{"string": types.StringValue("test-value")},
			)),
		},
		"SetNestedBlock-ElementKeyValue-AttributeName": {
			tfType: tftypes.NewAttributePath().WithAttributeName("set_block").WithElementKeyValue(tftypes.NewValue(
				tftypes.Object{AttributeTypes: map[string]tftypes.Type{"string": tftypes.String}},
				map[string]tftypes.Value{"string": tftypes.NewValue(tftypes.String, "test-value")},
			)).WithAttributeName("string"),
			schema: nestedSchema,
			expected: path.Root("set_block").AtSetValue(types.ObjectValueMust(
				map[string]attr.Type{"string": types.StringType},
				map[string]attr.Value{"string": types.StringValue("test-value")},
			)).AtName("string"),
		},
		"SingleNestedBlock-AttributeName": {
			tfType:   tftypes.NewAttributePath().WithAttributeName("single_block").WithAttributeName("string"),
			schema:   nestedSchema,
			expected: path.Root("single_block").AtName("string"),
		},
		"ListNestedAttribute-ElementKeyInt-AttributeName": {
			tfType:   tftypes.NewAttributePath().WithAttributeName("list_nested").WithElementKeyInt(0).WithAttributeName("string"),
			schema:   nestedSchema,
			expected: path.Root("list_nested").AtListIndex(0).AtName("string"),
		},
		"MapNestedAttribute-ElementKeyString-AttributeName": {
			tfType:   tftypes.NewAttributePath().WithAttributeName("map_nested").WithElementKeyString("key").WithAttributeName("string"),
			schema:   nestedSchema,
			expected: path.Root("map_nested").AtMapKey("key").AtName("string"),
		},
		"SetNestedAttribute-ElementKeyValue-AttributeName": {
			tfType: tftypes.NewAttributePath().WithAttributeName("set_nested").WithElementKeyValue(tftypes.NewValue(
				tftypes.Object{AttributeTypes: map[string]tftypes.Type{"string": tftypes.String}},
				map[string]tftypes.Value{"string": tftypes.NewValue(tftypes.String, "test-value")},
			)).WithAttributeName("string"),
			schema: nestedSchema,
			expected: path.Root("set_nested").AtSetValue(types.ObjectValueMust(
				map[string]attr.Type{"string": types.StringType},
				map[string]attr.Value{"string": types.StringValue("test-value")},
			)).AtName("string"),
		},
		"SingleNestedAttribute-AttributeName": {
			tfType:   tftypes.NewAttributePath().WithAttributeName("single_nested").WithAttributeName("string"),
			schema:   nestedSchema,
			expected: path.Root("single_nested").AtName("string"),
		},
		"ObjectAttribute-AttributeName-ElementKeyInt": {
			tfType:   tftypes.NewAttributePath().WithAttributeName("object").WithAttributeName("list").WithElementKeyInt(0),
			schema:   nestedSchema,
			expected: path.Root("object").AtName("list").AtListIndex(0),
		},
		"DynamicAttribute": {
			tfType:   tftypes.NewAttributePath().WithAttributeName("dynamic"),
			schema:   nestedSchema,
			expected: path.Root("dynamic"),
		},
		"DynamicAttribute-ElementKeyInt": {
			tfType:   tftypes.NewAttributePath().WithAttributeName("dynamic").WithElementKeyInt(0),
			schema:   nestedSchema,
			expected: path.Empty(),
			expectedDiags: diag.Diagnostics{
				diag.NewErrorDiagnostic(
					"Unable to Convert Attribute Path",
					"An unexpected error occurred while trying to convert an attribute path. "+
						"This is an error in terraform-plugin-framework used by the provider. "+
						"Please report the following to the provider developers.\n\n"+
						"Attribute Path: AttributeName(\"dynamic\").ElementKeyInt(0)\n"+
						"Original Error: path leads to element or attribute nested in a schema.DynamicAttribute",
				),
			},
		},
		"ElementKeyInt": {
			tfType: tftypes.NewAttributePath().WithElementKeyInt(1),
			schema: testschema.Schema{
				Attributes: map[string]fwschema.Attribute{
					"test": testschema.Attribute{
						Type: testtypes.StringType{},
					},
				},
			},
			expected: path.Empty(),
			expectedDiags: diag.Diagnostics{
				diag.NewErrorDiagnostic(
					"Unable to Convert Attribute Path",
					"An unexpected error occurred while trying to convert an attribute path. "+
						"This is an error in terraform-plugin-framework used by the provider. "+
						"Please report the following to the provider developers.\n\n"+
						"Attribute Path: ElementKeyInt(1)\n"+
						"Original Error: ElementKeyInt(1) still remains in the path: cannot apply AttributePathStep tftypes.ElementKeyInt to schema",
				),
			},
		},
		"ElementKeyString": {
			tfType: tftypes.NewAttributePath().WithElementKeyString("test"),
			schema: testschema.Schema{
				Attributes: map[string]fwschema.Attribute{
					"test": testschema.Attribute{
						Type: testtypes.StringType{},
					},
				},
			},
			expected: path.Empty(),
			expectedDiags: diag.Diagnostics{
				diag.NewErrorDiagnostic(
					"Unable to Convert Attribute Path",
					"An unexpected error occurred while trying to convert an attribute path. "+
						"This is an error in terraform-plugin-framework used by the provider. "+
						"Please report the following to the provider developers.\n\n"+
						"Attribute Path: ElementKeyString(\"test\")\n"+
						"Original Error: ElementKeyString(\"test\") still remains in the path: cannot apply AttributePathStep tftypes.ElementKeyString to schema",
				),
			},
		},
		"ElementKeyValue": {
			tfType: tftypes.NewAttributePath().WithElementKeyValue(tftypes.NewValue(tftypes.String, "test-value")),
			schema: testschema.Schema{
				Attributes: map[string]fwschema.Attribute{
					"test": testschema.Attribute{
						Type: testtypes.StringType{},
					},
				},
			},
			expected: path.Empty(),
			expectedDiags: diag.Diagnostics{
				diag.NewErrorDiagnostic(
					"Unable to Convert Attribute Path",
					"An unexpected error occurred while trying to convert an attribute path. "+
						"This is an error in terraform-plugin-framework used by the provider. "+
						"Please report the following to the provider developers.\n\n"+
						"Attribute Path: ElementKeyValue(tftypes.String<\"test-value\">)\n"+
						"Original Error: ElementKeyValue(tftypes.String<\"test-value\">) still remains in the path: cannot apply AttributePathStep tftypes.ElementKeyValue to schema",
				),
			},
		},
	}

	objectTypes := map[string]attr.Type{"string": types.StringType}
	tfObjectType := tftypes.Object{AttributeTypes: map[string]tftypes.Type{"string": tftypes.String}}
	customSetSchema := schema.Schema{
		Blocks: map[string]schema.Block{
			"set_block": schema.SetNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"string": schema.StringAttribute{Optional: true},
					},
					CustomType: testtypes.ObjectTypeWithSemanticEquals{
						ObjectType: types.ObjectType{AttrTypes: objectTypes},
					},
				},
			},
		},
	}

	// Set element paths must preserve custom values, including unknown and null
	// objects and objects whose children are unknown.
	for name, value := range map[string]struct {
		raw      any
		expected types.Object
	}{
		"known": {
			raw:      map[string]tftypes.Value{"string": tftypes.NewValue(tftypes.String, "test-value")},
			expected: types.ObjectValueMust(objectTypes, map[string]attr.Value{"string": types.StringValue("test-value")}),
		},
		"null": {
			expected: types.ObjectNull(objectTypes),
		},
		"unknown": {
			raw:      tftypes.UnknownValue,
			expected: types.ObjectUnknown(objectTypes),
		},
		"unknown-child": {
			raw:      map[string]tftypes.Value{"string": tftypes.NewValue(tftypes.String, tftypes.UnknownValue)},
			expected: types.ObjectValueMust(objectTypes, map[string]attr.Value{"string": types.StringUnknown()}),
		},
	} {
		for _, custom := range []bool{false, true} {
			s := nestedSchema
			var expected attr.Value = value.expected

			if custom {
				s = customSetSchema
				expected = testtypes.ObjectValueWithSemanticEquals{ObjectValue: value.expected}
			}

			testCases[fmt.Sprintf("SetNestedBlock-%s-custom-%t", name, custom)] = testCase{
				tfType: tftypes.NewAttributePath().WithAttributeName("set_block").
					WithElementKeyValue(tftypes.NewValue(tfObjectType, value.raw)).WithAttributeName("string"),
				schema:   s,
				expected: path.Root("set_block").AtSetValue(expected).AtName("string"),
			}
		}
	}

	invalidObject := testtypes.ObjectTypeWithSemanticEquals{
		ObjectType: types.ObjectType{AttrTypes: map[string]attr.Type{"string": testtypes.InvalidType{}}},
	}
	invalidSetSchema := schema.Schema{
		Blocks: map[string]schema.Block{
			"set_block": schema.SetNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{"string": schema.StringAttribute{Optional: true}},
					CustomType: invalidObject,
				},
			},
		},
	}
	invalidSetValue := tftypes.NewValue(tfObjectType, map[string]tftypes.Value{"string": tftypes.NewValue(tftypes.String, "test-value")})
	invalidSetPath := tftypes.NewAttributePath().WithAttributeName("set_block").WithElementKeyValue(invalidSetValue)
	testCases["SetNestedBlock-custom-value-conversion-error"] = testCase{
		tfType:   invalidSetPath.WithAttributeName("string"),
		schema:   invalidSetSchema,
		expected: path.Empty(),
		expectedDiags: diag.Diagnostics{
			diag.NewErrorDiagnostic(
				"Unable to Convert Attribute Path",
				"An unexpected error occurred while trying to convert an attribute path. "+
					"This is either an error in terraform-plugin-framework or a custom attribute type used by the provider. "+
					"Please report the following to the provider developers.\n\n"+
					fmt.Sprintf("Attribute Path: %s\n", invalidSetPath)+
					fmt.Sprintf("Original Error: unable to create PathStepElementKeyValue from tftypes.Value: unable to convert tftypes.Value (%s) to attr.Value: intentional ValueFromTerraform error", invalidSetValue),
			),
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, diags := fromtftypes.AttributePath(context.Background(), testCase.tfType, testCase.schema)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}

			if diff := cmp.Diff(diags, testCase.expectedDiags); diff != "" {
				for _, d := range diags {
					t.Logf("diag: %s", d.Detail())
				}
				t.Errorf("unexpected diagnostics difference: %s", diff)
			}
		})
	}
}

// nestedBlocks returns list nested blocks with the given fanout per level.
func nestedBlocks(fanout []int) map[string]schema.Block {
	blocks := map[string]schema.Block{}

	if len(fanout) == 0 {
		return blocks
	}

	for i := range fanout[0] {
		blocks[fmt.Sprintf("block_%d", i)] = schema.ListNestedBlock{
			NestedObject: schema.NestedBlockObject{
				Attributes: map[string]schema.Attribute{
					"attr": schema.StringAttribute{Optional: true},
				},
				Blocks: nestedBlocks(fanout[1:]),
			},
		}
	}

	return blocks
}

func BenchmarkAttributePathNestedBlocks(b *testing.B) {
	ctx := context.Background()

	// 199 blocks over 8 levels, shaped like a Kubernetes pod template.
	fanout := []int{1, 1, 1, 4, 3, 3, 2, 1}
	s := schema.Schema{Blocks: nestedBlocks(fanout)}

	// Deepest attribute: block_0[0].block_0[0]...block_0[0].attr
	tfPath := tftypes.NewAttributePath()
	for range fanout {
		tfPath = tfPath.WithAttributeName("block_0").WithElementKeyInt(0)
	}
	tfPath = tfPath.WithAttributeName("attr")

	b.ReportAllocs()

	for b.Loop() {
		if _, diags := fromtftypes.AttributePath(ctx, tfPath, s); diags.HasError() {
			b.Fatalf("unexpected diagnostics: %v", diags)
		}
	}
}

func BenchmarkAttributePathSmallSchemas(b *testing.B) {
	ctx := context.Background()
	smallSchema := schema.Schema{
		Attributes: map[string]schema.Attribute{
			"string": schema.StringAttribute{Optional: true},
			"set":    schema.SetAttribute{Optional: true, ElementType: types.StringType},
		},
	}

	for name, benchmark := range map[string]struct {
		schema schema.Schema
		path   *tftypes.AttributePath
	}{
		"shallow": {
			schema: smallSchema,
			path:   tftypes.NewAttributePath().WithAttributeName("string"),
		},
		"set-element": {
			schema: smallSchema,
			path:   tftypes.NewAttributePath().WithAttributeName("set").WithElementKeyValue(tftypes.NewValue(tftypes.String, "value")),
		},
		"set-block-child": {
			schema: nestedSchema,
			path: tftypes.NewAttributePath().WithAttributeName("set_block").WithElementKeyValue(tftypes.NewValue(
				tftypes.Object{AttributeTypes: map[string]tftypes.Type{"string": tftypes.String}},
				map[string]tftypes.Value{"string": tftypes.NewValue(tftypes.String, "value")},
			)).WithAttributeName("string"),
		},
	} {
		b.Run(name, func(b *testing.B) {
			b.ReportAllocs()

			for b.Loop() {
				if _, diags := fromtftypes.AttributePath(ctx, benchmark.path, benchmark.schema); diags.HasError() {
					b.Fatalf("unexpected diagnostics: %v", diags)
				}
			}
		})
	}
}
