package primitive

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

// Mode type
type Mode int

// Mode supported by premitive packages
const (
	ModeCombo Mode = iota
	ModeTriangle
	ModeRect
	ModeEllipse
	ModeCircle
	ModeRotatedRect
	ModeBeziers
	ModeRotatedEllipse
	ModePolygon
)

//WithMode func
func WithMode(mode Mode) func() []string {
	return func() []string {
		return []string{"-m", fmt.Sprintf("%d", mode)}
	}
}

//Transform func
func Transform(image io.Reader, ext string, shapes int, opts ...func() []string) (io.Reader, error) {

	var args []string
	for _, opt := range opts {
		args = append(args, opt()...)
	}

	in, err := tempFile("in_", ext)
	if err != nil {
		return nil, errors.New("Primitive : failed to create temp input file")
	}
	defer os.Remove(in.Name())

	out, err := tempFile("out_", ext)
	if err != nil {
		return nil, errors.New("Primitive : failed to create temp output file")
	}
	defer os.Remove(out.Name())

	_, err = io.Copy(in, image)
	if err != nil {
		return nil, errors.New("Primitive : failed to copy image into temp input file")
	}

	std, err := primitive(in.Name(), out.Name(), shapes, args...)
	if err != nil {
		return nil, fmt.Errorf("Primitive : failed to run the primitive cmd, stdcombo=%s", std)
	}

	fmt.Println(std)

	b := bytes.NewBuffer(nil)

	_, err = io.Copy(b, out)
	if err != nil {
		return nil, errors.New("Primitive : failed to copy output file into byte buffer")
	}

	return b, nil

}

func primitive(input string, output string, shapes int, args ...string) (string, error) {
	argstr := fmt.Sprintf("-i %s -o %s -n %d", input, output, shapes)

	args = append(strings.Fields(argstr), args...)
	cmd := exec.Command("primitive", args...)
	b, err := cmd.CombinedOutput()
	return string(b), err
}

func tempFile(prefix string, ext string) (*os.File, error) {
	in, err := ioutil.TempFile("", prefix)
	if err != nil {
		return nil, errors.New("Primitive : failed to create temp input file")
	}
	defer os.Remove(in.Name())

	return os.Create(fmt.Sprintf("%s.%s", in.Name(), ext))
}
