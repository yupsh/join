package join

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	yup "github.com/yupsh/framework"
	"github.com/yupsh/framework/opt"

	localopt "github.com/yupsh/join/opt"
)

// Flags represents the configuration options for the join command
type Flags = localopt.Flags

// Command implementation
type command opt.Inputs[string, Flags]

// Join creates a new join command with the given parameters
func Join(parameters ...any) yup.Command {
	cmd := command(opt.Args[string, Flags](parameters...))
	// Set defaults
	if cmd.Flags.Field1 == 0 {
		cmd.Flags.Field1 = 1 // Default to first field
	}
	if cmd.Flags.Field2 == 0 {
		cmd.Flags.Field2 = 1 // Default to first field
	}
	if cmd.Flags.EmptyString == "" {
		cmd.Flags.EmptyString = ""
	}
	return cmd
}

func (c command) Execute(ctx context.Context, stdin io.Reader, stdout, stderr io.Writer) error {
	// Check for cancellation before starting
	if err := yup.CheckContextCancellation(ctx); err != nil {
		return err
	}

	if len(c.Positional) < 2 {
		fmt.Fprintln(stderr, "join: need at least 2 files")
		return fmt.Errorf("need at least 2 files")
	}

	file1Name := c.Positional[0]
	file2Name := c.Positional[1]

	// Open files
	var file1, file2 io.ReadCloser
	var err error

	if file1Name == "-" {
		file1 = io.NopCloser(stdin)
	} else {
		file1, err = os.Open(file1Name)
		if err != nil {
			fmt.Fprintf(stderr, "join: %s: %v\n", file1Name, err)
			return err
		}
		defer file1.Close()
	}

	// Check for cancellation after opening first file
	if err := yup.CheckContextCancellation(ctx); err != nil {
		return err
	}

	if file2Name == "-" {
		file2 = io.NopCloser(stdin)
	} else {
		file2, err = os.Open(file2Name)
		if err != nil {
			fmt.Fprintf(stderr, "join: %s: %v\n", file2Name, err)
			return err
		}
		defer file2.Close()
	}

	return c.performJoin(ctx, file1, file2, stdout, stderr)
}

func (c command) performJoin(ctx context.Context, file1, file2 io.Reader, output, stderr io.Writer) error {
	lines1, err := c.readLines(ctx, file1)
	if err != nil {
		return err
	}

	// Check for cancellation after reading first file
	if err := yup.CheckContextCancellation(ctx); err != nil {
		return err
	}

	lines2, err := c.readLines(ctx, file2)
	if err != nil {
		return err
	}

	// Check for cancellation after reading second file
	if err := yup.CheckContextCancellation(ctx); err != nil {
		return err
	}

	// Parse lines into records
	records1 := c.parseRecords(ctx, lines1)
	records2 := c.parseRecords(ctx, lines2)

	// Check for cancellation after parsing
	if err := yup.CheckContextCancellation(ctx); err != nil {
		return err
	}

	// Create indexes on join fields
	index1 := c.createIndex(ctx, records1, int(c.Flags.Field1)-1)
	index2 := c.createIndex(ctx, records2, int(c.Flags.Field2)-1)

	// Check for cancellation after indexing
	if err := yup.CheckContextCancellation(ctx); err != nil {
		return err
	}

	// Perform join
	joined := make(map[string]bool)
	joinCount := 0

	// Inner join
	for key, recs1 := range index1 {
		// Check for cancellation periodically during join (every 100 keys for efficiency)
		joinCount++
		if joinCount%100 == 0 {
			if err := yup.CheckContextCancellation(ctx); err != nil {
				return err
			}
		}

		if recs2, exists := index2[c.normalizeKey(key)]; exists {
			for _, rec1 := range recs1 {
				for _, rec2 := range recs2 {
					c.outputJoinedRecord(rec1, rec2, output)
					joined[c.recordKey(rec1)] = true
					joined[c.recordKey(rec2)] = true
				}
			}
		}
	}

	// Check for cancellation before processing unmatched records
	if err := yup.CheckContextCancellation(ctx); err != nil {
		return err
	}

	// Output unmatched records if requested
	if bool(c.Flags.UnpairedFile1) || bool(c.Flags.OuterJoin) {
		for i, rec := range records1 {
			// Check for cancellation periodically (every 1000 records for efficiency)
			if i%1000 == 0 {
				if err := yup.CheckContextCancellation(ctx); err != nil {
					return err
				}
			}
			if !joined[c.recordKey(rec)] {
				c.outputUnpairedRecord(rec, true, output)
			}
		}
	}

	if bool(c.Flags.UnpairedFile2) || bool(c.Flags.OuterJoin) {
		for i, rec := range records2 {
			// Check for cancellation periodically (every 1000 records for efficiency)
			if i%1000 == 0 {
				if err := yup.CheckContextCancellation(ctx); err != nil {
					return err
				}
			}
			if !joined[c.recordKey(rec)] {
				c.outputUnpairedRecord(rec, false, output)
			}
		}
	}

	return nil
}

func (c command) readLines(ctx context.Context, reader io.Reader) ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(reader)

	for yup.ScanWithContext(ctx, scanner) {
		lines = append(lines, scanner.Text())
	}

	// Check if context was cancelled
	if err := yup.CheckContextCancellation(ctx); err != nil {
		return lines, err
	}

	return lines, scanner.Err()
}

func (c command) parseRecords(ctx context.Context, lines []string) [][]string {
	var records [][]string
	for i, line := range lines {
		// Check for cancellation periodically (every 1000 lines for efficiency)
		if i%1000 == 0 {
			if err := yup.CheckContextCancellation(ctx); err != nil {
				return records // Return partial results on cancellation
			}
		}
		fields := strings.Fields(line)
		if len(fields) > 0 {
			records = append(records, fields)
		}
	}
	return records
}

func (c command) createIndex(ctx context.Context, records [][]string, fieldIndex int) map[string][][]string {
	index := make(map[string][][]string)

	for i, record := range records {
		// Check for cancellation periodically (every 1000 records for efficiency)
		if i%1000 == 0 {
			if err := yup.CheckContextCancellation(ctx); err != nil {
				return index // Return partial index on cancellation
			}
		}
		var key string
		if fieldIndex < len(record) {
			key = record[fieldIndex]
		}
		normalizedKey := c.normalizeKey(key)
		index[normalizedKey] = append(index[normalizedKey], record)
	}

	return index
}

func (c command) normalizeKey(key string) string {
	if bool(c.Flags.IgnoreCase) {
		return strings.ToLower(key)
	}
	return key
}

func (c command) recordKey(record []string) string {
	return strings.Join(record, "\t")
}

func (c command) outputJoinedRecord(rec1, rec2 []string, output io.Writer) {
	// Default join: join field, remaining fields from file1, remaining fields from file2
	joinField := ""
	if int(c.Flags.Field1)-1 < len(rec1) {
		joinField = rec1[int(c.Flags.Field1)-1]
	}

	var result []string
	result = append(result, joinField)

	// Add remaining fields from file1 (excluding join field)
	for i, field := range rec1 {
		if i != int(c.Flags.Field1)-1 {
			result = append(result, field)
		}
	}

	// Add remaining fields from file2 (excluding join field)
	for i, field := range rec2 {
		if i != int(c.Flags.Field2)-1 {
			result = append(result, field)
		}
	}

	fmt.Fprintln(output, strings.Join(result, " "))
}

func (c command) outputUnpairedRecord(record []string, isFile1 bool, output io.Writer) {
	fmt.Fprintln(output, strings.Join(record, " "))
}

func (c command) String() string {
	return fmt.Sprintf("join %v", c.Positional)
}
