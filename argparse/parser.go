package argparse

import "github.com/jessevdk/go-flags"

type Options struct {
	BgHidden bool `long:"nobg" description:"Disable background"`
}

func Parse(argstr []string) (*Options, error) {
	opts := &Options{}
	_, err := flags.ParseArgs(opts, argstr)

	if err != nil {
		return nil, err
	}

	return opts, nil
}
