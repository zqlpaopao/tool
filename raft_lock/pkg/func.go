package pkg

import "context"

type CmdAbleFunc interface {
	Run(ctx context.Context, cmder Cmder) Cmder
}

type CmdAble func(cmd Cmder)

func (c CmdAble) Run(able Cmder) Cmder {
	var err error
	if err = able.Check(); err != nil {
		goto END
	}
	c(able)
END:
	able.SetError(err)
	return able
}
