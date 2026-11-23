package configuration

import (
	"flag"
	"fmt"
	"strings"
)

type Options struct {
	inputs    []string
	output    string
	id        string
	name      string
	alternate string
	self      string
}

func Read(args []string) (*Options, error) {
	flags := flag.NewFlagSet("", flag.ContinueOnError)
	output := flags.String("output", "-", "The destination to write the output feed to, or - for stdout")
	id := flags.String("id", "", "The ID for the processed feed, or use the current ID of the feed if blank")
	name := flags.String("name", "", "The name for the processed feed, or use the current name of the feed if blank")
	alternate := flags.String("alternate", "", "The URL for the alternate representation of this feed, or use the current if left blank")
	self := flags.String("self", "", "The URL for this feed, or nothing if left blank")

	if err := flags.Parse(args); err != nil {
		return nil, err
	}

	inputs := flags.Args()
	if len(inputs) == 0 {
		return nil, fmt.Errorf("No input files given")
	}

	return &Options{
		inputs:    inputs,
		output:    *output,
		id:        *id,
		name:      *name,
		alternate: *alternate,
		self:      *self,
	}, nil
}

func (options *Options) Inputs() []string {
	return options.inputs
}

func (options *Options) Output() string {
	return options.output
}

func (options *Options) ID() string {
	return options.id
}

func (options *Options) Name() string {
	return options.name
}

func (options *Options) Alternate() string {
	return options.alternate
}

func (options *Options) Self() string {
	return options.self
}

func (options *Options) String() string {
	return fmt.Sprintf(
		"Configuration{output: %s, id: %s, name: %s, alternate: %s, inputs: %s}",
		options.output,
		options.id,
		options.name,
		options.alternate,
		strings.Join(options.inputs, ", "),
	)
}
