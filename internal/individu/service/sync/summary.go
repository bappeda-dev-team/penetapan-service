package sync

import "github.com/bappeda-dev-team/penetapan-service/internal/individu/web"

type SummaryCounter struct {
	Rekin     int
	Indikator int
	Target    int
	Renaksi   int
}

func (s *SummaryCounter) AddRekin(n int) {
	s.Rekin += n
}

func (s *SummaryCounter) AddIndikator(n int) {
	s.Indikator += n
}

func (s *SummaryCounter) AddTarget(n int) {
	s.Target += n
}

func (s *SummaryCounter) AddRenaksi(n int) {
	s.Renaksi += n
}

func (s SummaryCounter) Response() web.SyncPenetapanSummary {
	var rekin *int
	if s.Rekin > 0 {
		rekin = &s.Rekin
	}
	return web.SyncPenetapanSummary{
		Rekin:     rekin,
		Indikator: s.Indikator,
		Target:    s.Target,
		Renaksi:   s.Renaksi,
	}

}
