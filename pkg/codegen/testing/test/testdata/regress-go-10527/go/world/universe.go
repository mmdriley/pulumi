// Code generated by test DO NOT EDIT.
// *** WARNING: Do not edit by hand unless you're certain you know what you are doing! ***

package world

import (
	"context"
	"reflect"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Universe struct {
	pulumi.CustomResourceState
}

// NewUniverse registers a new resource with the given unique name, arguments, and options.
func NewUniverse(ctx *pulumi.Context,
	name string, args *UniverseArgs, opts ...pulumi.ResourceOption) (*Universe, error) {
	if args == nil {
		args = &UniverseArgs{}
	}

	var resource Universe
	err := ctx.RegisterResource("world::Universe", name, args, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// GetUniverse gets an existing Universe resource's state with the given name, ID, and optional
// state properties that are used to uniquely qualify the lookup (nil if not required).
func GetUniverse(ctx *pulumi.Context,
	name string, id pulumi.IDInput, state *UniverseState, opts ...pulumi.ResourceOption) (*Universe, error) {
	var resource Universe
	err := ctx.ReadResource("world::Universe", name, id, state, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// Input properties used for looking up and filtering Universe resources.
type universeState struct {
}

type UniverseState struct {
}

func (UniverseState) ElementType() reflect.Type {
	return reflect.TypeOf((*universeState)(nil)).Elem()
}

type universeArgs struct {
	Worlds map[string]World `pulumi:"worlds"`
}

// The set of arguments for constructing a Universe resource.
type UniverseArgs struct {
	Worlds WorldMapInput
}

func (UniverseArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*universeArgs)(nil)).Elem()
}

type UniverseInput interface {
	pulumi.Input

	ToUniverseOutput() UniverseOutput
	ToUniverseOutputWithContext(ctx context.Context) UniverseOutput
}

func (*Universe) ElementType() reflect.Type {
	return reflect.TypeOf((**Universe)(nil)).Elem()
}

func (i *Universe) ToUniverseOutput() UniverseOutput {
	return i.ToUniverseOutputWithContext(context.Background())
}

func (i *Universe) ToUniverseOutputWithContext(ctx context.Context) UniverseOutput {
	return pulumi.ToOutputWithContext(ctx, i).(UniverseOutput)
}

type UniverseOutput struct{ *pulumi.OutputState }

func (UniverseOutput) ElementType() reflect.Type {
	return reflect.TypeOf((**Universe)(nil)).Elem()
}

func (o UniverseOutput) ToUniverseOutput() UniverseOutput {
	return o
}

func (o UniverseOutput) ToUniverseOutputWithContext(ctx context.Context) UniverseOutput {
	return o
}

func init() {
	pulumi.RegisterInputType(reflect.TypeOf((*UniverseInput)(nil)).Elem(), &Universe{})
	pulumi.RegisterOutputType(UniverseOutput{})
}
