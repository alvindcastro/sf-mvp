package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"sf-mvp/internal/eval"
)

const (
	exitPass        = 0
	exitGateFailure = 1
	exitUsageError  = 2
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("evalops-gate", flag.ContinueOnError)
	flags.SetOutput(stderr)
	inputPath := flags.String("input", "", "optional shared Promptfoo/EvalOps result JSON file")
	summaryPath := flags.String("summary", "", "optional Markdown summary path; defaults to GITHUB_STEP_SUMMARY when set")
	warningOnly := flags.String("warning-only", "", "comma-separated scorer names that should warn instead of blocking")
	minSeverity := flags.Float64("min-severity", 1, "minimum severity accuracy")
	minCitation := flags.Float64("min-citation", 1, "minimum citation coverage")
	minRecommendation := flags.Float64("min-recommendation", 1, "minimum recommendation accuracy")
	if err := flags.Parse(args); err != nil {
		return exitUsageError
	}

	outputs, err := releaseGateOutputs(*inputPath)
	if err != nil {
		fmt.Fprintf(stderr, "%v\n", err)
		return exitUsageError
	}

	config := eval.DefaultReleaseGateConfig()
	config.MinSeverityAccuracy = *minSeverity
	config.MinCitationCoverage = *minCitation
	config.MinRecommendationAccuracy = *minRecommendation
	config.WarningOnlyScorers = splitCSV(*warningOnly)

	result, err := eval.EvaluateReleaseGate(outputs, config)
	if err != nil {
		fmt.Fprintf(stderr, "evaluate release gate: %v\n", err)
		return exitUsageError
	}
	summary := eval.ReleaseGateMarkdownSummary(result)

	path := firstNonEmpty(*summaryPath, os.Getenv("GITHUB_STEP_SUMMARY"))
	if path != "" {
		if err := os.WriteFile(path, []byte(summary), 0o644); err != nil {
			fmt.Fprintf(stderr, "write summary %q: %v\n", path, err)
			return exitUsageError
		}
	}

	if result.Passed {
		if len(result.Warnings) > 0 {
			fmt.Fprintf(stdout, "EvalOps release gate passed with %d warning(s)\n", len(result.Warnings))
		} else {
			fmt.Fprintln(stdout, "EvalOps release gate passed")
		}
		return exitPass
	}
	fmt.Fprintf(stdout, "EvalOps release gate failed with %d blocking failure(s)\n", len(result.Failures))
	if path == "" {
		fmt.Fprint(stdout, summary)
	}
	return exitGateFailure
}

func releaseGateOutputs(inputPath string) ([]eval.PromptfooOutput, error) {
	if strings.TrimSpace(inputPath) == "" {
		report := eval.Run(eval.GoldenCases(), eval.DefaultThresholds())
		return eval.PromptfooOutputsFromReport(report), nil
	}

	data, err := os.ReadFile(inputPath)
	if err != nil {
		return nil, fmt.Errorf("read eval results %q: %w", inputPath, err)
	}
	results, err := eval.ImportSharedResultsJSON(data)
	if err != nil {
		return nil, fmt.Errorf("import eval results %q: %w", inputPath, err)
	}
	return eval.PromptfooOutputsFromReport(eval.Report{Cases: results}), nil
}

func splitCSV(value string) []string {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	parts := strings.Split(value, ",")
	values := make([]string, 0, len(parts))
	for _, part := range parts {
		if strings.TrimSpace(part) != "" {
			values = append(values, strings.TrimSpace(part))
		}
	}
	return values
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}
