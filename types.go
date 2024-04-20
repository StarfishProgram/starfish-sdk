package starfish_sdk

type Result[D any] struct {
	Code ICode
	Data D
}

func (r *Result[D]) IsOk() bool {
	return r.Code == nil
}

func (r *Result[D]) IsFaild() bool {
	return !r.IsOk()
}
