// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package fromtftypes

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// AttributePath returns the path.Path equivalent of a *tftypes.AttributePath.
func AttributePath(ctx context.Context, tfType *tftypes.AttributePath, schema fwschema.Schema) (path.Path, diag.Diagnostics) {
	if fwPath, ok := attributePathSingleWalk(ctx, tfType, schema); ok {
		return fwPath, nil
	}

	return attributePathPerPrefix(ctx, tfType, schema)
}

// attributePathSingleWalk converts the path while walking the schema once.
// The framework type of a step is only needed for set element steps, so it is
// not resolved for other steps. Resolving the type of a block rebuilds the
// type of all of its nested blocks, which made the conversion of every path
// prefix expensive for deeply nested block schemas.
//
// It returns false when the path cannot be converted this way, in which case
// attributePathPerPrefix produces the result and any diagnostics.
func attributePathSingleWalk(ctx context.Context, tfType *tftypes.AttributePath, schema fwschema.Schema) (path.Path, bool) {
	fwPath := path.Empty()
	steps := tfType.Steps()

	var current any = schema

	for stepIndex, step := range steps {
		stepper, ok := current.(tftypes.AttributePathStepper)

		if !ok {
			return path.Empty(), false
		}

		next, err := stepper.ApplyTerraform5AttributePathStep(step)

		if err != nil {
			return path.Empty(), false
		}

		// Only continue with the schema elements which
		// fwschema.SchemaTypeAtTerraformPath can resolve a type for.
		switch next.(type) {
		case attr.Type, fwschema.Attribute, fwschema.Block, fwschema.NestedAttributeObject,
			fwschema.NestedBlockObject, fwschema.UnderlyingAttributes:
		default:
			return path.Empty(), false
		}

		current = next

		switch step := step.(type) {
		case tftypes.AttributeName:
			fwPath = fwPath.AtName(string(step))
		case tftypes.ElementKeyInt:
			fwPath = fwPath.AtListIndex(int(step))
		case tftypes.ElementKeyString:
			fwPath = fwPath.AtMapKey(string(step))
		case tftypes.ElementKeyValue:
			attrType, err := schema.TypeAtTerraformPath(ctx, tftypes.NewAttributePathWithSteps(steps[:stepIndex+1]))

			if err != nil {
				return path.Empty(), false
			}

			attrValue, err := Value(ctx, tftypes.Value(step), attrType)

			if err != nil {
				return path.Empty(), false
			}

			fwPath = fwPath.AtSetValue(attrValue)
		default:
			return path.Empty(), false
		}
	}

	return fwPath, true
}

// attributePathPerPrefix converts the path by resolving the framework type of
// every path prefix.
func attributePathPerPrefix(ctx context.Context, tfType *tftypes.AttributePath, schema fwschema.Schema) (path.Path, diag.Diagnostics) {
	fwPath := path.Empty()

	for tfTypeStepIndex, tfTypeStep := range tfType.Steps() {
		currentTfTypeSteps := tfType.Steps()[:tfTypeStepIndex+1]
		currentTfTypePath := tftypes.NewAttributePathWithSteps(currentTfTypeSteps)
		attrType, err := schema.TypeAtTerraformPath(ctx, currentTfTypePath)

		if err != nil {
			return path.Empty(), diag.Diagnostics{
				diag.NewErrorDiagnostic(
					"Unable to Convert Attribute Path",
					"An unexpected error occurred while trying to convert an attribute path. "+
						"This is an error in terraform-plugin-framework used by the provider. "+
						"Please report the following to the provider developers.\n\n"+
						// Since this is an error with the attribute path
						// conversion, we cannot return a protocol path-based
						// diagnostic. Returning a framework human-readable
						// representation seems like the next best thing to do.
						fmt.Sprintf("Attribute Path: %s\n", currentTfTypePath.String())+
						fmt.Sprintf("Original Error: %s", err),
				),
			}
		}

		fwStep, err := AttributePathStep(ctx, tfTypeStep, attrType)

		if err != nil {
			return path.Empty(), diag.Diagnostics{
				diag.NewErrorDiagnostic(
					"Unable to Convert Attribute Path",
					"An unexpected error occurred while trying to convert an attribute path. "+
						"This is either an error in terraform-plugin-framework or a custom attribute type used by the provider. "+
						"Please report the following to the provider developers.\n\n"+
						// Since this is an error with the attribute path
						// conversion, we cannot return a protocol path-based
						// diagnostic. Returning a framework human-readable
						// representation seems like the next best thing to do.
						fmt.Sprintf("Attribute Path: %s\n", currentTfTypePath.String())+
						fmt.Sprintf("Original Error: %s", err),
				),
			}
		}

		// In lieu of creating a path.NewPathFromSteps function, this path
		// building logic is inlined to not expand the path package API.
		switch fwStep := fwStep.(type) {
		case path.PathStepAttributeName:
			fwPath = fwPath.AtName(string(fwStep))
		case path.PathStepElementKeyInt:
			fwPath = fwPath.AtListIndex(int(fwStep))
		case path.PathStepElementKeyString:
			fwPath = fwPath.AtMapKey(string(fwStep))
		case path.PathStepElementKeyValue:
			fwPath = fwPath.AtSetValue(fwStep.Value)
		default:
			return fwPath, diag.Diagnostics{
				diag.NewErrorDiagnostic(
					"Unable to Convert Attribute Path",
					"An unexpected error occurred while trying to convert an attribute path. "+
						"This is an error in terraform-plugin-framework used by the provider. "+
						"Please report the following to the provider developers.\n\n"+
						// Since this is an error with the attribute path
						// conversion, we cannot return a protocol path-based
						// diagnostic. Returning a framework human-readable
						// representation seems like the next best thing to do.
						fmt.Sprintf("Attribute Path: %s\n", currentTfTypePath.String())+
						fmt.Sprintf("Original Error: unknown path.PathStep type: %#v", fwStep),
				),
			}
		}
	}

	return fwPath, nil
}
