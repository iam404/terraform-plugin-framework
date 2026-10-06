// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package fwschemadata_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschemadata"
	"github.com/hashicorp/terraform-plugin-framework/internal/testing/testschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestDataNullifyCollectionBlocks(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		data          *fwschemadata.Data
		expected      *fwschemadata.Data
		expectedDiags diag.Diagnostics
	}{
		"list-attribute-unmodified": {
			data: &fwschemadata.Data{
				Description: fwschemadata.DataDescriptionConfiguration,
				Schema: testschema.Schema{
					Attributes: map[string]fwschema.Attribute{
						"list_attribute": testschema.Attribute{
							Optional: true,
							Type: types.ListType{
								ElemType: types.StringType,
							},
						},
					},
				},
				TerraformValue: tftypes.NewValue(
					tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"list_attribute": tftypes.List{
								ElementType: tftypes.String,
							},
						},
					},
					map[string]tftypes.Value{
						"list_attribute": tftypes.NewValue(
							tftypes.List{
								ElementType: tftypes.String,
							},
							[]tftypes.Value{}, // intentionally no elements
						),
					},
				),
			},
			expected: &fwschemadata.Data{
				Description: fwschemadata.DataDescriptionConfiguration,
				Schema: testschema.Schema{
					Attributes: map[string]fwschema.Attribute{
						"list_attribute": testschema.Attribute{
							Optional: true,
							Type: types.ListType{
								ElemType: types.StringType,
							},
						},
					},
				},
				TerraformValue: tftypes.NewValue(
					tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"list_attribute": tftypes.List{
								ElementType: tftypes.String,
							},
						},
					},
					map[string]tftypes.Value{
						"list_attribute": tftypes.NewValue(
							tftypes.List{
								ElementType: tftypes.String,
							},
							[]tftypes.Value{}, // intentionally no elements
						),
					},
				),
			},
		},
		"set-attribute-unmodified": {
			data: &fwschemadata.Data{
				Description: fwschemadata.DataDescriptionConfiguration,
				Schema: testschema.Schema{
					Attributes: map[string]fwschema.Attribute{
						"set_attribute": testschema.Attribute{
							Optional: true,
							Type: types.SetType{
								ElemType: types.StringType,
							},
						},
					},
				},
				TerraformValue: tftypes.NewValue(
					tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"set_attribute": tftypes.Set{
								ElementType: tftypes.String,
							},
						},
					},
					map[string]tftypes.Value{
						"set_attribute": tftypes.NewValue(
							tftypes.Set{
								ElementType: tftypes.String,
							},
							[]tftypes.Value{}, // intentionally no elements
						),
					},
				),
			},
			expected: &fwschemadata.Data{
				Description: fwschemadata.DataDescriptionConfiguration,
				Schema: testschema.Schema{
					Attributes: map[string]fwschema.Attribute{
						"set_attribute": testschema.Attribute{
							Optional: true,
							Type: types.SetType{
								ElemType: types.StringType,
							},
						},
					},
				},
				TerraformValue: tftypes.NewValue(
					tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"set_attribute": tftypes.Set{
								ElementType: tftypes.String,
							},
						},
					},
					map[string]tftypes.Value{
						"set_attribute": tftypes.NewValue(
							tftypes.Set{
								ElementType: tftypes.String,
							},
							[]tftypes.Value{}, // intentionally no elements
						),
					},
				),
			},
		},
		"list-block-null": {
			data: &fwschemadata.Data{
				Description: fwschemadata.DataDescriptionConfiguration,
				Schema: testschema.Schema{
					Blocks: map[string]fwschema.Block{
						"list_block": testschema.Block{
							NestedObject: testschema.NestedBlockObject{
								Attributes: map[string]fwschema.Attribute{
									"list_block_attribute": testschema.Attribute{
										Optional: true,
										Type:     types.StringType,
									},
								},
							},
							NestingMode: fwschema.BlockNestingModeList,
						},
					},
				},
				TerraformValue: tftypes.NewValue(
					tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"list_block": tftypes.List{
								ElementType: tftypes.Object{
									AttributeTypes: map[string]tftypes.Type{
										"list_block_attribute": tftypes.String,
									},
								},
							},
						},
					},
					map[string]tftypes.Value{
						"list_block": tftypes.NewValue(
							tftypes.List{
								ElementType: tftypes.Object{
									AttributeTypes: map[string]tftypes.Type{
										"list_block_attribute": tftypes.String,
									},
								},
							},
							nil,
						),
					},
				),
			},
			expected: &fwschemadata.Data{
				Description: fwschemadata.DataDescriptionConfiguration,
				Schema: testschema.Schema{
					Blocks: map[string]fwschema.Block{
						"list_block": testschema.Block{
							NestedObject: testschema.NestedBlockObject{
								Attributes: map[string]fwschema.Attribute{
									"list_block_attribute": testschema.Attribute{
										Optional: true,
										Type:     types.StringType,
									},
								},
							},
							NestingMode: fwschema.BlockNestingModeList,
						},
					},
				},
				TerraformValue: tftypes.NewValue(
					tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"list_block": tftypes.List{
								ElementType: tftypes.Object{
									AttributeTypes: map[string]tftypes.Type{
										"list_block_attribute": tftypes.String,
									},
								},
							},
						},
					},
					map[string]tftypes.Value{
						"list_block": tftypes.NewValue(
							tftypes.List{
								ElementType: tftypes.Object{
									AttributeTypes: map[string]tftypes.Type{
										"list_block_attribute": tftypes.String,
									},
								},
							},
							nil,
						),
					},
				),
			},
		},
		"list-block-unknown": {
			data: &fwschemadata.Data{
				Description: fwschemadata.DataDescriptionConfiguration,
				Schema: testschema.Schema{
					Blocks: map[string]fwschema.Block{
						"list_block": testschema.Block{
							NestedObject: testschema.NestedBlockObject{
								Attributes: map[string]fwschema.Attribute{
									"list_block_attribute": testschema.Attribute{
										Optional: true,
										Type:     types.StringType,
									},
								},
							},
							NestingMode: fwschema.BlockNestingModeList,
						},
					},
				},
				TerraformValue: tftypes.NewValue(
					tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"list_block": tftypes.List{
								ElementType: tftypes.Object{
									AttributeTypes: map[string]tftypes.Type{
										"list_block_attribute": tftypes.String,
									},
								},
							},
						},
					},
					map[string]tftypes.Value{
						"list_block": tftypes.NewValue(
							tftypes.List{
								ElementType: tftypes.Object{
									AttributeTypes: map[string]tftypes.Type{
										"list_block_attribute": tftypes.String,
									},
								},
							},
							tftypes.UnknownValue,
						),
					},
				),
			},
			expected: &fwschemadata.Data{
				Description: fwschemadata.DataDescriptionConfiguration,
				Schema: testschema.Schema{
					Blocks: map[string]fwschema.Block{
						"list_block": testschema.Block{
							NestedObject: testschema.NestedBlockObject{
								Attributes: map[string]fwschema.Attribute{
									"list_block_attribute": testschema.Attribute{
										Optional: true,
										Type:     types.StringType,
									},
								},
							},
							NestingMode: fwschema.BlockNestingModeList,
						},
					},
				},
				TerraformValue: tftypes.NewValue(
					tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"list_block": tftypes.List{
								ElementType: tftypes.Object{
									AttributeTypes: map[string]tftypes.Type{
										"list_block_attribute": tftypes.String,
									},
								},
							},
						},
					},
					map[string]tftypes.Value{
						"list_block": tftypes.NewValue(
							tftypes.List{
								ElementType: tftypes.Object{
									AttributeTypes: map[string]tftypes.Type{
										"list_block_attribute": tftypes.String,
									},
								},
							},
							tftypes.UnknownValue,
						),
					},
				),
			},
		},
		"list-block-elements": {
			data: &fwschemadata.Data{
				Description: fwschemadata.DataDescriptionConfiguration,
				Schema: testschema.Schema{
					Blocks: map[string]fwschema.Block{
						"list_block": testschema.Block{
							NestedObject: testschema.NestedBlockObject{
								Attributes: map[string]fwschema.Attribute{
									"list_block_attribute": testschema.Attribute{
										Optional: true,
										Type:     types.StringType,
									},
								},
							},
							NestingMode: fwschema.BlockNestingModeList,
						},
					},
				},
				TerraformValue: tftypes.NewValue(
					tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"list_block": tftypes.List{
								ElementType: tftypes.Object{
									AttributeTypes: map[string]tftypes.Type{
										"list_block_attribute": tftypes.String,
									},
								},
							},
						},
					},
					map[string]tftypes.Value{
						"list_block": tftypes.NewValue(
							tftypes.List{
								ElementType: tftypes.Object{
									AttributeTypes: map[string]tftypes.Type{
										"list_block_attribute": tftypes.String,
									},
								},
							},
							[]tftypes.Value{
								tftypes.NewValue(
									tftypes.Object{
										AttributeTypes: map[string]tftypes.Type{
											"list_block_attribute": tftypes.String,
										},
									},
									nil, // null value should not matter
								),
							},
						),
					},
				),
			},
			expected: &fwschemadata.Data{
				Description: fwschemadata.DataDescriptionConfiguration,
				Schema: testschema.Schema{
					Blocks: map[string]fwschema.Block{
						"list_block": testschema.Block{
							NestedObject: testschema.NestedBlockObject{
								Attributes: map[string]fwschema.Attribute{
									"list_block_attribute": testschema.Attribute{
										Optional: true,
										Type:     types.StringType,
									},
								},
							},
							NestingMode: fwschema.BlockNestingModeList,
						},
					},
				},
				TerraformValue: tftypes.NewValue(
					tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"list_block": tftypes.List{
								ElementType: tftypes.Object{
									AttributeTypes: map[string]tftypes.Type{
										"list_block_attribute": tftypes.String,
									},
								},
							},
						},
					},
					map[string]tftypes.Value{
						"list_block": tftypes.NewValue(
							tftypes.List{
								ElementType: tftypes.Object{
									AttributeTypes: map[string]tftypes.Type{
										"list_block_attribute": tftypes.String,
									},
								},
							},
							[]tftypes.Value{
								tftypes.NewValue(
									tftypes.Object{
										AttributeTypes: map[string]tftypes.Type{
											"list_block_attribute": tftypes.String,
										},
									},
									nil, // null value should not matter
								),
							},
						),
					},
				),
			},
		},
		"list-block-zero-elements": {
			data: &fwschemadata.Data{
				Description: fwschemadata.DataDescriptionConfiguration,
				Schema: testschema.Schema{
					Blocks: map[string]fwschema.Block{
						"list_block": testschema.Block{
							NestedObject: testschema.NestedBlockObject{
								Attributes: map[string]fwschema.Attribute{
									"list_block_attribute": testschema.Attribute{
										Optional: true,
										Type:     types.StringType,
									},
								},
							},
							NestingMode: fwschema.BlockNestingModeList,
						},
					},
				},
				TerraformValue: tftypes.NewValue(
					tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"list_block": tftypes.List{
								ElementType: tftypes.Object{
									AttributeTypes: map[string]tftypes.Type{
										"list_block_attribute": tftypes.String,
									},
								},
							},
						},
					},
					map[string]tftypes.Value{
						"list_block": tftypes.NewValue(
							tftypes.List{
								ElementType: tftypes.Object{
									AttributeTypes: map[string]tftypes.Type{
										"list_block_attribute": tftypes.String,
									},
								},
							},
							[]tftypes.Value{}, // intentionally no elements
						),
					},
				),
			},
			expected: &fwschemadata.Data{
				Description: fwschemadata.DataDescriptionConfiguration,
				Schema: testschema.Schema{
					Blocks: map[string]fwschema.Block{
						"list_block": testschema.Block{
							NestedObject: testschema.NestedBlockObject{
								Attributes: map[string]fwschema.Attribute{
									"list_block_attribute": testschema.Attribute{
										Optional: true,
										Type:     types.StringType,
									},
								},
							},
							NestingMode: fwschema.BlockNestingModeList,
						},
					},
				},
				TerraformValue: tftypes.NewValue(
					tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"list_block": tftypes.List{
								ElementType: tftypes.Object{
									AttributeTypes: map[string]tftypes.Type{
										"list_block_attribute": tftypes.String,
									},
								},
							},
						},
					},
					map[string]tftypes.Value{
						"list_block": tftypes.NewValue(
							tftypes.List{
								ElementType: tftypes.Object{
									AttributeTypes: map[string]tftypes.Type{
										"list_block_attribute": tftypes.String,
									},
								},
							},
							nil, // should be converted to null value
						),
					},
				),
			},
		},
		"set-block-null": {
			data: &fwschemadata.Data{
				Description: fwschemadata.DataDescriptionConfiguration,
				Schema: testschema.Schema{
					Blocks: map[string]fwschema.Block{
						"set_block": testschema.Block{
							NestedObject: testschema.NestedBlockObject{
								Attributes: map[string]fwschema.Attribute{
									"set_block_attribute": testschema.Attribute{
										Optional: true,
										Type:     types.StringType,
									},
								},
							},
							NestingMode: fwschema.BlockNestingModeSet,
						},
					},
				},
				TerraformValue: tftypes.NewValue(
					tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"set_block": tftypes.Set{
								ElementType: tftypes.Object{
									AttributeTypes: map[string]tftypes.Type{
										"set_block_attribute": tftypes.String,
									},
								},
							},
						},
					},
					map[string]tftypes.Value{
						"set_block": tftypes.NewValue(
							tftypes.Set{
								ElementType: tftypes.Object{
									AttributeTypes: map[string]tftypes.Type{
										"set_block_attribute": tftypes.String,
									},
								},
							},
							nil,
						),
					},
				),
			},
			expected: &fwschemadata.Data{
				Description: fwschemadata.DataDescriptionConfiguration,
				Schema: testschema.Schema{
					Blocks: map[string]fwschema.Block{
						"set_block": testschema.Block{
							NestedObject: testschema.NestedBlockObject{
								Attributes: map[string]fwschema.Attribute{
									"set_block_attribute": testschema.Attribute{
										Optional: true,
										Type:     types.StringType,
									},
								},
							},
							NestingMode: fwschema.BlockNestingModeSet,
						},
					},
				},
				TerraformValue: tftypes.NewValue(
					tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"set_block": tftypes.Set{
								ElementType: tftypes.Object{
									AttributeTypes: map[string]tftypes.Type{
										"set_block_attribute": tftypes.String,
									},
								},
							},
						},
					},
					map[string]tftypes.Value{
						"set_block": tftypes.NewValue(
							tftypes.Set{
								ElementType: tftypes.Object{
									AttributeTypes: map[string]tftypes.Type{
										"set_block_attribute": tftypes.String,
									},
								},
							},
							nil,
						),
					},
				),
			},
		},
		"set-block-unknown": {
			data: &fwschemadata.Data{
				Description: fwschemadata.DataDescriptionConfiguration,
				Schema: testschema.Schema{
					Blocks: map[string]fwschema.Block{
						"set_block": testschema.Block{
							NestedObject: testschema.NestedBlockObject{
								Attributes: map[string]fwschema.Attribute{
									"set_block_attribute": testschema.Attribute{
										Optional: true,
										Type:     types.StringType,
									},
								},
							},
							NestingMode: fwschema.BlockNestingModeSet,
						},
					},
				},
				TerraformValue: tftypes.NewValue(
					tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"set_block": tftypes.Set{
								ElementType: tftypes.Object{
									AttributeTypes: map[string]tftypes.Type{
										"set_block_attribute": tftypes.String,
									},
								},
							},
						},
					},
					map[string]tftypes.Value{
						"set_block": tftypes.NewValue(
							tftypes.Set{
								ElementType: tftypes.Object{
									AttributeTypes: map[string]tftypes.Type{
										"set_block_attribute": tftypes.String,
									},
								},
							},
							tftypes.UnknownValue,
						),
					},
				),
			},
			expected: &fwschemadata.Data{
				Description: fwschemadata.DataDescriptionConfiguration,
				Schema: testschema.Schema{
					Blocks: map[string]fwschema.Block{
						"set_block": testschema.Block{
							NestedObject: testschema.NestedBlockObject{
								Attributes: map[string]fwschema.Attribute{
									"set_block_attribute": testschema.Attribute{
										Optional: true,
										Type:     types.StringType,
									},
								},
							},
							NestingMode: fwschema.BlockNestingModeSet,
						},
					},
				},
				TerraformValue: tftypes.NewValue(
					tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"set_block": tftypes.Set{
								ElementType: tftypes.Object{
									AttributeTypes: map[string]tftypes.Type{
										"set_block_attribute": tftypes.String,
									},
								},
							},
						},
					},
					map[string]tftypes.Value{
						"set_block": tftypes.NewValue(
							tftypes.Set{
								ElementType: tftypes.Object{
									AttributeTypes: map[string]tftypes.Type{
										"set_block_attribute": tftypes.String,
									},
								},
							},
							tftypes.UnknownValue,
						),
					},
				),
			},
		},
		"set-block-elements": {
			data: &fwschemadata.Data{
				Description: fwschemadata.DataDescriptionConfiguration,
				Schema: testschema.Schema{
					Blocks: map[string]fwschema.Block{
						"set_block": testschema.Block{
							NestedObject: testschema.NestedBlockObject{
								Attributes: map[string]fwschema.Attribute{
									"set_block_attribute": testschema.Attribute{
										Optional: true,
										Type:     types.StringType,
									},
								},
							},
							NestingMode: fwschema.BlockNestingModeSet,
						},
					},
				},
				TerraformValue: tftypes.NewValue(
					tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"set_block": tftypes.Set{
								ElementType: tftypes.Object{
									AttributeTypes: map[string]tftypes.Type{
										"set_block_attribute": tftypes.String,
									},
								},
							},
						},
					},
					map[string]tftypes.Value{
						"set_block": tftypes.NewValue(
							tftypes.Set{
								ElementType: tftypes.Object{
									AttributeTypes: map[string]tftypes.Type{
										"set_block_attribute": tftypes.String,
									},
								},
							},
							[]tftypes.Value{
								tftypes.NewValue(
									tftypes.Object{
										AttributeTypes: map[string]tftypes.Type{
											"set_block_attribute": tftypes.String,
										},
									},
									nil, // null value should not matter
								),
							},
						),
					},
				),
			},
			expected: &fwschemadata.Data{
				Description: fwschemadata.DataDescriptionConfiguration,
				Schema: testschema.Schema{
					Blocks: map[string]fwschema.Block{
						"set_block": testschema.Block{
							NestedObject: testschema.NestedBlockObject{
								Attributes: map[string]fwschema.Attribute{
									"set_block_attribute": testschema.Attribute{
										Optional: true,
										Type:     types.StringType,
									},
								},
							},
							NestingMode: fwschema.BlockNestingModeSet,
						},
					},
				},
				TerraformValue: tftypes.NewValue(
					tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"set_block": tftypes.Set{
								ElementType: tftypes.Object{
									AttributeTypes: map[string]tftypes.Type{
										"set_block_attribute": tftypes.String,
									},
								},
							},
						},
					},
					map[string]tftypes.Value{
						"set_block": tftypes.NewValue(
							tftypes.Set{
								ElementType: tftypes.Object{
									AttributeTypes: map[string]tftypes.Type{
										"set_block_attribute": tftypes.String,
									},
								},
							},
							[]tftypes.Value{
								tftypes.NewValue(
									tftypes.Object{
										AttributeTypes: map[string]tftypes.Type{
											"set_block_attribute": tftypes.String,
										},
									},
									nil, // null value should not matter
								),
							},
						),
					},
				),
			},
		},
		"set-block-zero-elements": {
			data: &fwschemadata.Data{
				Description: fwschemadata.DataDescriptionConfiguration,
				Schema: testschema.Schema{
					Blocks: map[string]fwschema.Block{
						"set_block": testschema.Block{
							NestedObject: testschema.NestedBlockObject{
								Attributes: map[string]fwschema.Attribute{
									"set_block_attribute": testschema.Attribute{
										Optional: true,
										Type:     types.StringType,
									},
								},
							},
							NestingMode: fwschema.BlockNestingModeSet,
						},
					},
				},
				TerraformValue: tftypes.NewValue(
					tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"set_block": tftypes.Set{
								ElementType: tftypes.Object{
									AttributeTypes: map[string]tftypes.Type{
										"set_block_attribute": tftypes.String,
									},
								},
							},
						},
					},
					map[string]tftypes.Value{
						"set_block": tftypes.NewValue(
							tftypes.Set{
								ElementType: tftypes.Object{
									AttributeTypes: map[string]tftypes.Type{
										"set_block_attribute": tftypes.String,
									},
								},
							},
							[]tftypes.Value{}, // intentionally no elements
						),
					},
				),
			},
			expected: &fwschemadata.Data{
				Description: fwschemadata.DataDescriptionConfiguration,
				Schema: testschema.Schema{
					Blocks: map[string]fwschema.Block{
						"set_block": testschema.Block{
							NestedObject: testschema.NestedBlockObject{
								Attributes: map[string]fwschema.Attribute{
									"set_block_attribute": testschema.Attribute{
										Optional: true,
										Type:     types.StringType,
									},
								},
							},
							NestingMode: fwschema.BlockNestingModeSet,
						},
					},
				},
				TerraformValue: tftypes.NewValue(
					tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"set_block": tftypes.Set{
								ElementType: tftypes.Object{
									AttributeTypes: map[string]tftypes.Type{
										"set_block_attribute": tftypes.String,
									},
								},
							},
						},
					},
					map[string]tftypes.Value{
						"set_block": tftypes.NewValue(
							tftypes.Set{
								ElementType: tftypes.Object{
									AttributeTypes: map[string]tftypes.Type{
										"set_block_attribute": tftypes.String,
									},
								},
							},
							nil, // should be converted to null value
						),
					},
				),
			},
		},
		// Dynamic attributes that contain underlying list values should be skipped
		"dynamic-attribute-with-list-unmodified": {
			data: &fwschemadata.Data{
				Description: fwschemadata.DataDescriptionConfiguration,
				Schema: testschema.Schema{
					Attributes: map[string]fwschema.Attribute{
						"dynamic_attribute": testschema.Attribute{
							Optional: true,
							Type:     types.DynamicType,
						},
					},
				},
				TerraformValue: tftypes.NewValue(
					tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"dynamic_attribute": tftypes.DynamicPseudoType,
						},
					},
					map[string]tftypes.Value{
						"dynamic_attribute": tftypes.NewValue(
							tftypes.List{
								ElementType: tftypes.String,
							},
							[]tftypes.Value{
								tftypes.NewValue(tftypes.String, "hello"),
							},
						),
					},
				),
			},
			expected: &fwschemadata.Data{
				Description: fwschemadata.DataDescriptionConfiguration,
				Schema: testschema.Schema{
					Attributes: map[string]fwschema.Attribute{
						"dynamic_attribute": testschema.Attribute{
							Optional: true,
							Type:     types.DynamicType,
						},
					},
				},
				TerraformValue: tftypes.NewValue(
					tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"dynamic_attribute": tftypes.DynamicPseudoType,
						},
					},
					map[string]tftypes.Value{
						"dynamic_attribute": tftypes.NewValue(
							tftypes.List{
								ElementType: tftypes.String,
							},
							[]tftypes.Value{
								tftypes.NewValue(tftypes.String, "hello"),
							},
						),
					},
				),
			},
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			diags := testCase.data.NullifyCollectionBlocks(context.Background())

			if diff := cmp.Diff(diags, testCase.expectedDiags); diff != "" {
				t.Errorf("unexpected diagnostics difference: %s", diff)
			}

			if diff := cmp.Diff(testCase.data, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

// nestedBlocksFanout describes a schema of list nested blocks shaped like a
// large real-world schema (for example a Kubernetes pod template): a chain of
// wrapper blocks followed by a wide tree, 199 blocks over 8 levels.
var nestedBlocksFanout = []int{1, 1, 1, 4, 3, 3, 2, 1}

// nestedBlocks returns list nested blocks with the given fanout per level, and
// a value which sets one element in every block, with one known and one null
// attribute per element.
func nestedBlocks(fanout []int) (map[string]schema.Block, map[string]tftypes.Type, map[string]tftypes.Value) {
	blocks := map[string]schema.Block{}
	types := map[string]tftypes.Type{}
	values := map[string]tftypes.Value{}

	if len(fanout) == 0 {
		return blocks, types, values
	}

	for i := range fanout[0] {
		name := fmt.Sprintf("block_%d", i)
		childBlocks, childTypes, childValues := nestedBlocks(fanout[1:])

		blocks[name] = schema.ListNestedBlock{
			NestedObject: schema.NestedBlockObject{
				Attributes: map[string]schema.Attribute{
					"known": schema.StringAttribute{Optional: true},
					"null":  schema.StringAttribute{Optional: true},
				},
				Blocks: childBlocks,
			},
		}

		childTypes["known"] = tftypes.String
		childTypes["null"] = tftypes.String
		childValues["known"] = tftypes.NewValue(tftypes.String, "value")
		childValues["null"] = tftypes.NewValue(tftypes.String, nil)

		objectType := tftypes.Object{AttributeTypes: childTypes}
		listType := tftypes.List{ElementType: objectType}

		types[name] = listType
		values[name] = tftypes.NewValue(listType, []tftypes.Value{
			tftypes.NewValue(objectType, childValues),
		})
	}

	return blocks, types, values
}

func nestedBlocksData() *fwschemadata.Data {
	blocks, types, values := nestedBlocks(nestedBlocksFanout)

	return &fwschemadata.Data{
		Description:    fwschemadata.DataDescriptionConfiguration,
		Schema:         schema.Schema{Blocks: blocks},
		TerraformValue: tftypes.NewValue(tftypes.Object{AttributeTypes: types}, values),
	}
}

func BenchmarkDataNullifyCollectionBlocksNestedBlocks(b *testing.B) {
	ctx := context.Background()
	data := nestedBlocksData()
	value := data.TerraformValue

	b.ReportAllocs()

	for b.Loop() {
		data.TerraformValue = value

		if diags := data.NullifyCollectionBlocks(ctx); diags.HasError() {
			b.Fatalf("unexpected diagnostics: %v", diags)
		}
	}
}
