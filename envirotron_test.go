package envirotron_test

import (
	"os"
	"strconv"
	"strings"
	"testing"

	env "github.com/jhunt/go-envirotron"
)

type Env map[string]string

func set(e Env) {
	for name, val := range e {
		os.Setenv(name, val)
	}
}

type Shallow struct {
	Name    string `env:"NAME"`
	Ignored string

	hidden string `env:"SECRET"` /* noop */
}

type Deep struct {
	Family string `env:"FAMILY"`
	Nested Shallow
}

type Values struct {
	Bool    bool    `env:"SOME_BOOL"`
	Int     int     `env:"SOME_INT"`
	Int8    int8    `env:"SOME_INT_8"`
	Int16   int16   `env:"SOME_INT_16"`
	Int32   int32   `env:"SOME_INT_32"`
	Int64   int64   `env:"SOME_INT_64"`
	Uint    uint    `env:"SOME_UINT"`
	Uint8   uint8   `env:"SOME_UINT_8"`
	Uint16  uint16  `env:"SOME_UINT_16"`
	Uint32  uint32  `env:"SOME_UINT_32"`
	Uint64  uint64  `env:"SOME_UINT_64"`
	Float32 float32 `env:"SOME_FLOAT_32"`
	Float64 float64 `env:"SOME_FLOAT_64"`
}

type doubler int

func (d *doubler) UnmarshalEnv(raw string) {
	i, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		panic(err)
	}
	*d = (doubler)(i * 2)
}

type Callback struct {
	Doubler doubler `env:"DOUBLE_ME"`
}

func TestEnvirotronShallow(t *testing.T) {
	is := func(got, expect, message string) {
		if got != expect {
			t.Errorf("%s failed - got '%s', expected '%s'\n", message, got, expect)
		}
	}

	set(Env{
		"NAME":    "overridden name",
		"SECRET":  "overridden secret",
		"IGNORED": "overridden ignored (BAD!)",
	})

	shallow := Shallow{
		Name:    "initial name",
		hidden:  "initial hidden",
		Ignored: "initial ignored",
	}
	is(shallow.Name, "initial name", "initial name is set before testing")
	is(shallow.Ignored, "initial ignored", "initial ignored is set before testing")
	is(shallow.hidden, "initial hidden", "initial hidden is set before testing")

	env.Override(&shallow)
	is(shallow.Name, "overridden name", "Name is overridden from NAME env var")
	is(shallow.Ignored, "initial ignored", "initial ignored is still set")
	is(shallow.hidden, "initial hidden", "hidden fields cannot be overridden")
}

func TestEnvirotronNested(t *testing.T) {
	is := func(got, expect, message string) {
		if got != expect {
			t.Errorf("%s failed - got '%s', expected '%s'\n", message, got, expect)
		}
	}

	set(Env{
		"FAMILY":  "overridden family",
		"NAME":    "overridden name",
		"SECRET":  "overridden secret",
		"IGNORED": "overridden ignored (BAD!)",
	})

	deep := Deep{
		Family: "initial family",
		Nested: Shallow{
			Name:    "initial name",
			hidden:  "initial hidden",
			Ignored: "initial ignored",
		},
	}
	is(deep.Family, "initial family", "initial family is set before testing")
	is(deep.Nested.Name, "initial name", "initial name is set before testing")
	is(deep.Nested.Ignored, "initial ignored", "initial ignored is set before testing")
	is(deep.Nested.hidden, "initial hidden", "initial hidden is set before testing")

	env.Override(&deep)
	is(deep.Family, "overridden family", "Family is overridden from FAMILY env var")
	is(deep.Nested.Name, "overridden name", "Name is overridden from NAME env var")
	is(deep.Nested.Ignored, "initial ignored", "initial ignored is still set")
	is(deep.Nested.hidden, "initial hidden", "hidden fields cannot be overridden")
}

func TestEnvirotronValues(t *testing.T) {
	ok := func(good bool, got, expect interface{}, message string) {
		if !good {
			t.Errorf("%s failed - got [%v], expected [%v]\n", message, got, expect)
		}
	}

	var (
		b   bool    = true
		u   uint    = 42
		u8  uint8   = 255
		u16 uint16  = 65535
		u32 uint32  = 4294967295
		u64 uint64  = 18446744073709551615
		i   int     = -42
		i8  int8    = 127
		i16 int16   = 32767
		i32 int32   = 2147483647
		i64 int64   = 9223372036854775807
		f32 float32 = 1.2345
		f64 float64 = 123456789.123456789123456789
	)

	set(Env{
		"SOME_BOOL":     "y",
		"SOME_UINT":     "42",
		"SOME_UINT_8":   "255",
		"SOME_UINT_16":  "65535",
		"SOME_UINT_32":  "4294967295",
		"SOME_UINT_64":  "18446744073709551615",
		"SOME_INT":      "-42",
		"SOME_INT_8":    "127",
		"SOME_INT_16":   "32767",
		"SOME_INT_32":   "2147483647",
		"SOME_INT_64":   "9223372036854775807",
		"SOME_FLOAT_32": "1.2345",
		"SOME_FLOAT_64": "123456789.123456789123456789",
	})

	values := Values{
		Bool:    false,
		Uint:    1,
		Uint8:   2,
		Uint16:  3,
		Uint32:  4,
		Uint64:  5,
		Int:     6,
		Int8:    7,
		Int16:   8,
		Int32:   9,
		Int64:   10,
		Float32: 11.12,
		Float64: 13.14,
	}
	ok(values.Bool == false, values.Bool, false, "initial Bool is set before testing")
	ok(values.Uint == 1, values.Uint, 1, "initial Uint is set before testing")
	ok(values.Uint8 == 2, values.Uint8, 2, "initial Uint8 is set before testing")
	ok(values.Uint16 == 3, values.Uint16, 3, "initial Uint16 is set before testing")
	ok(values.Uint32 == 4, values.Uint32, 4, "initial Uint32 is set before testing")
	ok(values.Uint64 == 5, values.Uint64, 5, "initial Uint64 is set before testing")
	ok(values.Int == 6, values.Int, 6, "initial Int is set before testing")
	ok(values.Int8 == 7, values.Int8, 7, "initial Int8 is set before testing")
	ok(values.Int16 == 8, values.Int16, 8, "initial Int16 is set before testing")
	ok(values.Int32 == 9, values.Int32, 9, "initial Int32 is set before testing")
	ok(values.Int64 == 10, values.Int64, 10, "initial Int64 is set before testing")
	ok(values.Float32 == 11.12, values.Float32, 11.12, "initial Float32 is set before testing")
	ok(values.Float64 == 13.14, values.Float64, 13.14, "initial Float64 is set before testing")

	env.Override(&values)

	ok(values.Bool == b, values.Bool, b, "initial Bool is overridden from env")
	ok(values.Uint == u, values.Uint, u, "initial Uint is overridden from env")
	ok(values.Uint8 == u8, values.Uint8, u8, "initial Uint8 is overridden from env")
	ok(values.Uint16 == u16, values.Uint16, u16, "initial Uint16 is overridden from env")
	ok(values.Uint32 == u32, values.Uint32, u32, "initial Uint32 is overridden from env")
	ok(values.Uint64 == u64, values.Uint64, u64, "initial Uint64 is overridden from env")
	ok(values.Int == i, values.Int, i, "initial Int is overridden from env")
	ok(values.Int8 == i8, values.Int8, i8, "initial Int8 is overridden from env")
	ok(values.Int16 == i16, values.Int16, i16, "initial Int16 is overridden from env")
	ok(values.Int32 == i32, values.Int32, i32, "initial Int32 is overridden from env")
	ok(values.Int64 == i64, values.Int64, i64, "initial Int64 is overridden from env")
	ok(values.Float32 == f32, values.Float32, f32, "initial Float32 is overridden from env")
	ok(values.Float64 == f64, values.Float64, f64, "initial Float64 is overridden from env")
}

func TestEnvirotronBools(t *testing.T) {
	is := func(got, expect bool, test, message string) {
		if got != expect {
			t.Errorf("%s ['%s' test] failed - got [%v], expected [%v]\n", message, test, got, expect)
		}
	}

	trues := strings.Split("Y y Yes YES yES true TrUe 1", " ")
	falses := strings.Split("N n No NO nO false fALse 0", " ")
	check := Values{}

	for _, yes := range trues {
		check.Bool = false
		set(Env{"SOME_BOOL": yes})
		is(check.Bool, false, yes, "structure bool is initially false, before override")
		env.Override(&check)
		is(check.Bool, true, yes, "structure bool is overridden to be true")
	}

	for _, no := range falses {
		check.Bool = true
		set(Env{"SOME_BOOL": no})
		is(check.Bool, true, no, "structure bool is initially true, before override")
		env.Override(&check)
		is(check.Bool, false, no, "structure bool is overridden to be false")
	}
}

func TestEnvirotronCallback(t *testing.T) {
	is := func(got doubler, expect int, message string) {
		if got != (doubler)(expect) {
			t.Errorf("%s failed - got '%d', expected '%d'\n", message, got, expect)
		}
	}

	set(Env{
		"DOUBLE_ME": "42",
	})

	callback := Callback{
		Doubler: 1,
	}
	is(callback.Doubler, 1, "initial doubler value is set before testing")

	env.Override(&callback)
	is(callback.Doubler, 84, "doubler value is overridden from DOUBLE_ME env var")
}
