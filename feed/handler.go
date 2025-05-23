package feed

type Handler struct {
	feeds []Feed
}

func (h *Handler) Add(f Feed) {
	h.feeds = append(h.feeds, f)

	f.handle()
}

func (h *Handler) Failing() bool {
	for _, f := range h.feeds {
		if f.failing() {
			return true
		}
	}

	return false
}

func (h *Handler) Failed() bool {
	for _, f := range h.feeds {
		if f.failed() {
			return true
		}
	}

	return false
}

func NewHandler() *Handler {
	return &Handler{}
}
