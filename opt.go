package command

type Field1 int
type Field2 int
type OutputFormat string
type EmptyString string

type IgnoreCaseFlag bool

const (
	IgnoreCase    IgnoreCaseFlag = true
	CaseSensitive IgnoreCaseFlag = false
)

type OuterJoinFlag bool

const (
	OuterJoin OuterJoinFlag = true
	InnerJoin OuterJoinFlag = false
)

type UnpairedFile1Flag bool

const (
	UnpairedFile1   UnpairedFile1Flag = true
	NoUnpairedFile1 UnpairedFile1Flag = false
)

type UnpairedFile2Flag bool

const (
	UnpairedFile2   UnpairedFile2Flag = true
	NoUnpairedFile2 UnpairedFile2Flag = false
)

type CheckOrderFlag bool

const (
	CheckOrder   CheckOrderFlag = true
	NoCheckOrder CheckOrderFlag = false
)

type flags struct {
	Field1        Field1
	Field2        Field2
	OutputFormat  OutputFormat
	EmptyString   EmptyString
	IgnoreCase    IgnoreCaseFlag
	OuterJoin     OuterJoinFlag
	UnpairedFile1 UnpairedFile1Flag
	UnpairedFile2 UnpairedFile2Flag
	CheckOrder    CheckOrderFlag
}

func (f Field1) Configure(flags *flags)            { flags.Field1 = f }
func (f Field2) Configure(flags *flags)            { flags.Field2 = f }
func (o OutputFormat) Configure(flags *flags)      { flags.OutputFormat = o }
func (e EmptyString) Configure(flags *flags)       { flags.EmptyString = e }
func (i IgnoreCaseFlag) Configure(flags *flags)    { flags.IgnoreCase = i }
func (o OuterJoinFlag) Configure(flags *flags)     { flags.OuterJoin = o }
func (u UnpairedFile1Flag) Configure(flags *flags) { flags.UnpairedFile1 = u }
func (u UnpairedFile2Flag) Configure(flags *flags) { flags.UnpairedFile2 = u }
func (c CheckOrderFlag) Configure(flags *flags)    { flags.CheckOrder = c }
