package gologging

import "regexp"

var (
	getFunctionInfoRgx = regexp.MustCompile(`(((.*/)|^)\w+)\.?(\(.*\))?(.*?)$`)
	lambdaNameRgx      = regexp.MustCompile(`^((\w+)\.)?func([0-9]+)`)
	tagRgx             = regexp.MustCompile(`^(([[:lower:]^:]{2,10}): )?(.+)`)
)
