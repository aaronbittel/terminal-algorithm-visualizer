package component

import (
	"fmt"
	"strings"
	"sync"
	"time"
	term_utils "tui/internal/term-utils"
)

type Statusbar struct {
	instrCh chan<- string

	rows int
	cols int

	message string
	buf     *Buf

	mutex *sync.Mutex
}

func NewStatusbar(instrCh chan<- string, rows, cols int, buf *Buf) *Statusbar {
	return &Statusbar{
		instrCh: instrCh,
		buf:     buf,
		rows:    rows,
		cols:    cols,
		mutex:   &sync.Mutex{},
	}
}

func (s Statusbar) Pos() (row, col int) {
	return s.rows - 1, 1
}

func (s Statusbar) Chan() chan<- string {
	return s.instrCh
}

func (s Statusbar) String() string {
	var b strings.Builder
	b.WriteString(term_utils.Lightgray)

	b.WriteString(strings.Repeat(term_utils.Underscore, s.cols))
	b.WriteString(term_utils.ResetCode)
	return b.String()
}

func (s *Statusbar) Set(updated string) {
	var b strings.Builder
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.message == updated {
		return
	}

	b.WriteString(term_utils.MoveCur(s.rows, 1))
	b.WriteString(updated)
	diff := term_utils.StringLen(s.message) - term_utils.StringLen(updated)
	if diff > 0 {
		b.WriteString(strings.Repeat(" ", diff))
	}
	s.message = updated

	s.instrCh <- b.String()
}

func (s *Statusbar) After(updated string, duration time.Duration) {
	go func() {
		s.Set(updated)
		<-time.After(duration)
		// if status was updated in the meantime dont delete the content
		if s.message == updated {
			s.instrCh <- term_utils.ClearLineInst(s.rows, 1)
		}
	}()
}

func Error(message string) string {
	return fmt.Sprintf("%s%s%s", term_utils.BoldRed, message, term_utils.ResetCode)
}

func Info(message string) string {
	return fmt.Sprintf("%s%s%s", term_utils.BoldCyan, message, term_utils.ResetCode)
}
