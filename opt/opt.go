package opt

// Custom types for parameters
type Field1 int
type Field2 int
type OutputFormat string
type EmptyString string

// Boolean flag types with constants
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

// Flags represents the configuration options for the join command
type Flags struct {
	Field1        Field1            // Join field for file 1 (-1)
	Field2        Field2            // Join field for file 2 (-2)
	OutputFormat  OutputFormat      // Output format specification (-o)
	EmptyString   EmptyString       // String to represent empty fields (-e)
	IgnoreCase    IgnoreCaseFlag    // Ignore case when comparing (-i)
	OuterJoin     OuterJoinFlag     // Produce all possible joins (-a)
	UnpairedFile1 UnpairedFile1Flag // Print unpairable lines from file 1 (-a 1)
	UnpairedFile2 UnpairedFile2Flag // Print unpairable lines from file 2 (-a 2)
	CheckOrder    CheckOrderFlag    // Check that input is sorted (--check-order)
}

// Configure methods for the opt system
func (f Field1) Configure(flags *Flags)            { flags.Field1 = f }
func (f Field2) Configure(flags *Flags)            { flags.Field2 = f }
func (o OutputFormat) Configure(flags *Flags)      { flags.OutputFormat = o }
func (e EmptyString) Configure(flags *Flags)       { flags.EmptyString = e }
func (i IgnoreCaseFlag) Configure(flags *Flags)    { flags.IgnoreCase = i }
func (o OuterJoinFlag) Configure(flags *Flags)     { flags.OuterJoin = o }
func (u UnpairedFile1Flag) Configure(flags *Flags) { flags.UnpairedFile1 = u }
func (u UnpairedFile2Flag) Configure(flags *Flags) { flags.UnpairedFile2 = u }
func (c CheckOrderFlag) Configure(flags *Flags)    { flags.CheckOrder = c }
