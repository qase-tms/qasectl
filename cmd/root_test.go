package cmd

import (
	"testing"
)

func TestRootCommand_HasExpectedSubcommands(t *testing.T) {
	subcommands := make(map[string]bool)
	for _, cmd := range rootCmd.Commands() {
		subcommands[cmd.Name()] = true
	}

	expected := []string{"testops", "version"}
	for _, name := range expected {
		if !subcommands[name] {
			t.Errorf("rootCmd missing expected subcommand %q", name)
		}
	}
}

func TestRootCommand_VerboseFlag(t *testing.T) {
	flag := rootCmd.PersistentFlags().Lookup("verbose")
	if flag == nil {
		t.Fatal("verbose flag not registered on rootCmd")
	}

	if flag.DefValue != "false" {
		t.Errorf("verbose flag default = %q, want %q", flag.DefValue, "false")
	}

	if flag.Shorthand != "v" {
		t.Errorf("verbose flag shorthand = %q, want %q", flag.Shorthand, "v")
	}
}

func TestTestopsCommand_SubcommandRegistration(t *testing.T) {
	testopsCmd, _, err := rootCmd.Find([]string{"testops"})
	if err != nil {
		t.Fatalf("failed to find testops command: %v", err)
	}

	subcommands := make(map[string]bool)
	for _, cmd := range testopsCmd.Commands() {
		subcommands[cmd.Name()] = true
	}

	expected := []string{"run", "result", "env", "milestone", "filter", "field"}
	for _, name := range expected {
		if !subcommands[name] {
			t.Errorf("testops missing expected subcommand %q", name)
		}
	}
}

func TestTestopsCommand_PersistentFlags(t *testing.T) {
	testopsCmd, _, err := rootCmd.Find([]string{"testops"})
	if err != nil {
		t.Fatalf("failed to find testops command: %v", err)
	}

	if testopsCmd.PersistentFlags().Lookup("token") == nil {
		t.Error("testops missing persistent flag 'token'")
	}

	if testopsCmd.PersistentFlags().Lookup("project") == nil {
		t.Error("testops missing persistent flag 'project'")
	}
}

func TestTestopsCommand_RequiredFlags(t *testing.T) {
	testopsCmd, _, err := rootCmd.Find([]string{"testops"})
	if err != nil {
		t.Fatalf("failed to find testops command: %v", err)
	}

	tokenFlag := testopsCmd.PersistentFlags().Lookup("token")
	if tokenFlag == nil {
		t.Fatal("testops missing 'token' flag")
	}
	if _, ok := tokenFlag.Annotations["cobra_annotation_bash_completion_one_required_flag"]; !ok {
		t.Error("token flag is not marked as required")
	}

	projectFlag := testopsCmd.PersistentFlags().Lookup("project")
	if projectFlag == nil {
		t.Fatal("testops missing 'project' flag")
	}
	if _, ok := projectFlag.Annotations["cobra_annotation_bash_completion_one_required_flag"]; !ok {
		t.Error("project flag is not marked as required")
	}
}
