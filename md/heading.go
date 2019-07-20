package md

import "io"

type Heading struct {
	Level   int
	Content Node
}

func (x *Heading) Render(w io.Writer) error {
	for i := 0; i < x.Level; i++ {
		if _, err := w.Write([]byte("#")); err != nil {
			return err
		}
	}

	if _, err := w.Write([]byte(" ")); err != nil {
		return err
	}

	if err := x.Content.Render(w); err != nil {
		return err
	}

	if _, err := w.Write([]byte("\n\n")); err != nil {
		return err
	}
	return nil
}
