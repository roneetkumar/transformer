package primitive

import (
	"bytes"
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
func Transform(image io.Reader, shapes int, opts ...func() []string) (io.Reader, error) {

	in, err := ioutil.TempFile("", "in_")
	if err != nil {
		return nil, err
	}
	defer os.Remove(in.Name())

	out, err := ioutil.TempFile("", "out_")
	if err != nil {
		return nil, err
	}
	defer os.Remove(out.Name())

	_, err = io.Copy(in, image)
	if err != nil {
		return nil, err
	}

	std, err := primitive(in.Name(), out.Name(), shapes, ModeCombo)
	if err != nil {
		return nil, err
	}

	fmt.Println(std)

	b := bytes.NewBuffer(nil)

	_, err = io.Copy(b, out)
	if err != nil {
		return nil, err
	}

	return b, nil

}

func primitive(input string, output string, shapes int, mode Mode) (string, error) {
	argstr := fmt.Sprintf("-i %s -o %s -n %d -m %d", input, output, shapes, mode)
	cmd := exec.Command("primitive", strings.Fields(argstr)...)
	b, err := cmd.CombinedOutput()
	return string(b), err
}
