package runn

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/chromedp/chromedp"
)

const cdpNewKey = "new"

const (
	cdpWindowWidth  = 1920
	cdpWindowHeight = 1080
)

type cdpRunner struct {
	name     string
	ctx      context.Context
	cancel   context.CancelFunc
	store    map[string]interface{}
	operator *operator
}

type CDPActions []CDPAction

type CDPAction struct {
	Fn   string
	Args map[string]interface{}
}

func newCDPRunner(name, remote string) (*cdpRunner, error) {
	if remote != cdpNewKey {
		return nil, errors.New("remote connect mode is planned, but not yet implemented")
	}

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.WindowSize(cdpWindowWidth, cdpWindowHeight),
	)
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	ctx, _ := chromedp.NewContext(allocCtx)
	return &cdpRunner{
		name:   name,
		ctx:    ctx,
		cancel: cancel,
		store:  map[string]interface{}{},
	}, nil
}

func (rnr *cdpRunner) Close() error {
	if rnr.cancel == nil {
		return nil
	}
	rnr.cancel()
	rnr.cancel = nil
	return nil
}

func (rnr *cdpRunner) Run(_ context.Context, cas CDPActions) error {
	rnr.operator.capturers.captureCDPStart(rnr.name)
	defer rnr.operator.capturers.captureCDPEnd(rnr.name)
	before := []chromedp.Action{
		chromedp.EmulateViewport(cdpWindowWidth, cdpWindowHeight),
	}
	if err := chromedp.Run(rnr.ctx, before...); err != nil {
		return err
	}
	for _, ca := range cas {
		rnr.operator.capturers.captureCDPAction(ca)
		a, err := rnr.evalAction(ca)
		if err != nil {
			return err
		}
		if err := chromedp.Run(rnr.ctx, a); err != nil {
			return err
		}
		fn, ok := CDPFnMap[ca.Fn]
		if !ok {
			return fmt.Errorf("invalid action: %v", ca)
		}
		ras := fn.Args.resArgs()
		if len(ras) > 0 {
			res := map[string]interface{}{}
			for _, arg := range ras {
				res[arg.key] = rnr.operator.store.latest()[arg.key]
			}
			// capture
			rnr.operator.capturers.captureCDPResponse(ca, res)
		}
	}

	// record
	r := map[string]interface{}{}
	for k, v := range rnr.store {
		switch vv := v.(type) {
		case *string:
			r[k] = *vv
		default:
			r[k] = vv
		}
	}
	rnr.operator.record(r)

	rnr.store = map[string]interface{}{} // clear

	return nil
}

func (rnr *cdpRunner) evalAction(ca CDPAction) (chromedp.Action, error) {
	fn, ok := CDPFnMap[ca.Fn]
	if !ok {
		return nil, fmt.Errorf("invalid action: %v", ca)
	}
	fv := reflect.ValueOf(fn.Fn)
	vs := []reflect.Value{}
	for i, a := range fn.Args {
		switch a.typ {
		case CDPArgTypeArg:
			v, ok := ca.Args[a.key]
			if !ok {
				return nil, fmt.Errorf("invalid action: %v", ca)
			}
			vs = append(vs, reflect.ValueOf(v))
		case CDPArgTypeRes, CDPArgTypeHiddenRes:
			k := a.key

			switch reflect.TypeOf(fn.Fn).In(i).Kind() {
			case reflect.Interface:
				// evaluate
				var v interface{}
				rnr.store[k] = &v
				vs = append(vs, reflect.ValueOf(&v))
			default:
				switch reflect.TypeOf(fn.Fn).In(i).Elem().Kind() {
				case reflect.String:
					var v string
					rnr.store[k] = &v
					vs = append(vs, reflect.ValueOf(&v))
				default:
					return nil, fmt.Errorf("invalid action: %v", ca)
				}
			}
		}
	}
	res := fv.Call(vs)
	return res[0].Interface().(chromedp.Action), nil
}
