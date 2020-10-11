package frontmatter

import (
	"bufio"
	"bytes"
	"io"
)

type parser struct {
	reader *bufio.Reader
	output *bytes.Buffer

	read  int
	start int
	end   int
}

func newParser(r io.Reader) *parser {
	return &parser{
		reader: bufio.NewReader(r),
		output: bytes.NewBuffer(nil),
	}
}

func (p *parser) parse(v interface{}, formats []*Format) ([]byte, error) {
	// Detect format.
	f, err := p.detect(formats)
	if err != nil {
		return nil, err
	}

	// Extract front matter.
	if f != nil {
		if err = p.extract(f, v); err != nil {
			return nil, err
		}
	}

	// Read remaining data.
	if _, err := p.output.ReadFrom(p.reader); err != nil {
		return nil, err
	}

	return p.output.Bytes()[p.end:], nil
}

func (p *parser) detect(formats []*Format) (*Format, error) {
	for {
		read := p.read

		line, atEOF, err := p.readLine()
		if err != nil || atEOF {
			return nil, err
		}
		if line == "" {
			continue
		}

		for _, f := range formats {
			if f.Start == line {
				if !f.UnmarshalDelims {
					read = p.read
				}

				p.start = read
				return f, nil
			}
		}

		return nil, nil
	}
}

func (p *parser) extract(f *Format, v interface{}) error {
	for {
		read := p.read

		line, atEOF, err := p.readLine()
		if err != nil {
			return err
		}

	CheckLine:
		if line != f.End {
			if atEOF {
				return nil
			}
			continue
		}
		if f.RequiresNewLine {
			if line, atEOF, err = p.readLine(); err != nil {
				return err
			}
			if line != "" {
				goto CheckLine
			}
		}
		if f.UnmarshalDelims {
			read = p.read
		}

		if err := f.Unmarshal(p.output.Bytes()[p.start:read], v); err != nil {
			return err
		}

		p.end = p.read
		return nil
	}
}

func (p *parser) readLine() (string, bool, error) {
	line, err := p.reader.ReadBytes('\n')

	atEOF := err == io.EOF
	if err != nil && !atEOF {
		return "", false, err
	}

	p.read += len(line)
	_, err = p.output.Write(line)
	return string(bytes.TrimSpace(line)), atEOF, err
}
