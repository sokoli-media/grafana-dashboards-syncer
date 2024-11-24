package syncer

type Dashboard struct {
	filename  string
	dashboard string
}

func (d Dashboard) Equals(o Dashboard) bool {
	return d.filename == o.filename && d.dashboard == o.dashboard
}
