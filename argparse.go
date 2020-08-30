package main

import "github.com/jessevdk/go-flags"

type Options struct {
	BgHidden bool `long:"nobg" description:"Disable background"`
}

func ArgParse(argstr []string) (*Options, error) {
	opts := &Options{}
	_, err := flags.ParseArgs(opts, argstr)

	if err != nil {
		return nil, err
	}

	return opts, nil
}
