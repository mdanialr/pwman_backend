package response

// WithMsg option to add given msg to success response as `message` field.
func WithMsg(msg string) SuccessOption {
	return func(s *appSuccess) {
		s.Message = msg
	}
}

// WithData option to add the given data to success response as `data` field.
func WithData(data any) SuccessOption {
	return func(s *appSuccess) {
		s.Data = data
	}
}

// WithMeta option to add the given meta to success response as `meta` field.
func WithMeta(meta any) SuccessOption {
	return func(s *appSuccess) {
		s.Meta = meta
	}
}
