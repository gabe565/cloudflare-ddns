package output

import "context"

type Format string

const (
	FormatANSI     Format = "ansi"
	FormatMarkdown Format = "markdown"
)

type ctxKey uint8

const formatKey ctxKey = iota

func NewContext(ctx context.Context, format Format) context.Context {
	return context.WithValue(ctx, formatKey, format)
}

func FromContext(ctx context.Context) (Format, bool) {
	format, ok := ctx.Value(formatKey).(Format)
	return format, ok
}
