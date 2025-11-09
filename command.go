package command

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	gloo "github.com/gloo-foo/framework"
)

type command gloo.Inputs[string, flags]

func Join(parameters ...any) gloo.Command {
	cmd := command(gloo.Initialize[string, flags](parameters...))
	if cmd.Flags.Field1 == 0 {
		cmd.Flags.Field1 = 1
	}
	if cmd.Flags.Field2 == 0 {
		cmd.Flags.Field2 = 1
	}
	if cmd.Flags.EmptyString == "" {
		cmd.Flags.EmptyString = ""
	}
	return cmd
}

func (p command) Executor() gloo.CommandExecutor {
	return func(ctx context.Context, stdin io.Reader, stdout, stderr io.Writer) error {
		// Need two file paths
		if len(p.Positional) < 2 {
			_, _ = fmt.Fprintf(stderr, "join: missing operand\n")
			return fmt.Errorf("join requires two files")
		}

		file1Path := p.Positional[0]
		file2Path := p.Positional[1]

		// Read both files
		lines1, err := readFileLines(file1Path)
		if err != nil {
			_, _ = fmt.Fprintf(stderr, "join: %s: %v\n", file1Path, err)
			return err
		}

		lines2, err := readFileLines(file2Path)
		if err != nil {
			_, _ = fmt.Fprintf(stderr, "join: %s: %v\n", file2Path, err)
			return err
		}

		// Use whitespace as separator (standard join behavior)
		sep := " "

		// Get field indices (1-indexed)
		field1 := int(p.Flags.Field1) - 1
		field2 := int(p.Flags.Field2) - 1

		// Build index of file2 by join field
		file2Map := make(map[string][]string)
		for _, line := range lines2 {
			var fields []string
			if sep == " " {
				fields = strings.Fields(line)
			} else {
				fields = strings.Split(line, sep)
			}

			if field2 >= 0 && field2 < len(fields) {
				joinKey := fields[field2]
				file2Map[joinKey] = append(file2Map[joinKey], line)
			}
		}

		// Process file1 and join with file2
		for _, line := range lines1 {
			var fields []string
			if sep == " " {
				fields = strings.Fields(line)
			} else {
				fields = strings.Split(line, sep)
			}

			if field1 >= 0 && field1 < len(fields) {
				joinKey := fields[field1]

				// Find matching lines in file2
				if matches, found := file2Map[joinKey]; found {
					for _, match := range matches {
						var matchFields []string
						if sep == " " {
							matchFields = strings.Fields(match)
						} else {
							matchFields = strings.Split(match, sep)
						}

						// Output join key and fields from both files
						output := joinKey

						// Add remaining fields from file1
						for i, f := range fields {
							if i != field1 {
								output += sep + f
							}
						}

						// Add remaining fields from file2
						for i, f := range matchFields {
							if i != field2 {
								output += sep + f
							}
						}

						_, _ = fmt.Fprintln(stdout, output)
					}
				} else if bool(p.Flags.UnpairedFile1) || bool(p.Flags.OuterJoin) {
					// Output unmatched line from file1
					_, _ = fmt.Fprintln(stdout, line)
				}
			}
		}

		// Output unmatched lines from file2 if requested
		if bool(p.Flags.UnpairedFile2) || bool(p.Flags.OuterJoin) {
			seen := make(map[string]bool)
			for _, line := range lines1 {
				var fields []string
				if sep == " " {
					fields = strings.Fields(line)
				} else {
					fields = strings.Split(line, sep)
				}
				if field1 >= 0 && field1 < len(fields) {
					seen[fields[field1]] = true
				}
			}

			for _, line := range lines2 {
				var fields []string
				if sep == " " {
					fields = strings.Fields(line)
				} else {
					fields = strings.Split(line, sep)
				}
				if field2 >= 0 && field2 < len(fields) {
					if !seen[fields[field2]] {
						_, _ = fmt.Fprintln(stdout, line)
					}
				}
			}
		}

		return nil
	}
}

// readFileLines reads all lines from a file
func readFileLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}
