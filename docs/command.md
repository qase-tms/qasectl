# Create a test run

You can create a test run by using the `create` command. The `create` command is used to create a new test run in the
specified project and save a test run ID to a file. You can specify the file path using the `--output` option. If the
file path is not specified, the test run ID will be saved to `qase.env` in the current directory.

The file will contain the test run ID in the following format:

```text
QASE_TESTOPS_RUN_ID=123
```

You can use the test run ID in subsequent steps to upload test results for the test run.
For exctract test run ID from file you can use command:

```bash
cat qase.env | grep QASE_TESTOPS_RUN_ID | cut -d'=' -f2
```

## Example usage:

```bash
qli testops run create --project <project_code> --token <token> --title <title> --description <description> --environment <environment> --milestone <milestone> --plan <plan> --verbose
```

The `create` command has the following options:

- `--project`, `-p`: The project code where the test run will be created. Required.
- `--token`, `-t`: The API token to authenticate with the TestOps API. Required.
- `--title`: The name of the test run. Required.
- `--description`, `-d`: The description of the test run. Optional.
- `--environment`, `-e`: The environment where the test run will be executed. Optional.
- `--milestone`, `-m`: The milestone of the test run. Optional.
- `--plan`: The test plan of the test run. Optional.
- `--output`, `-o`: The output path to save the test run ID. Optional. Default is `qase.env` in the current
  directory.
- `--verbose`, `-v`: Enable verbose mode. Optional.

The following example shows how to create a test run in the project with the code `PROJ`:

```bash
qli testops run create --project PROJ --token <token> --title "Test Run 1" --description "This is a test run" --environment "Production" --milestone "Milestone 1" --plan "Test Plan 1" --verbose
```

# Complete a test run

You can complete a test run by using the `complete` command. The `complete` command is used to complete a test run in
the specified project.

## Example usage:

```bash
qli testops run complete --project <project_code> --token <token> --id <run_id> --verbose
```

The `complete` command has the following options:

- `--project`, `-p`: The project code where the test run will be completed. Required.
- `--token`, `-t`: The API token to authenticate with the TestOps API. Required.
- `--id`: The ID of the test run to complete. Required.
- `--verbose`, `-v`: Enable verbose mode. Optional.

The following example shows how to complete a test run with the ID `1` in the project with the code `PROJ`:

```bash
qli testops run complete --project PROJ --token <token> --id 1 --verbose
```

# Upload test results

You can upload test results by using the `upload` command. The `upload` command is used to upload test results for a
test run in the specified project.

## Example usage:

```bash
qli testops result upload --project <project_code> --token <token> --id <run_id> --format <format> --path <results_file> --batch <batch> --verbose
```

The `upload` command has the following options:

- `--project`, `-p`: The project code where the test results will be uploaded. Required.
- `--token`, `-t`: The API token to authenticate with the TestOps API. Required.
- `--id`: The ID of the test run to upload results for. Required if title doesn't set.
- `--title`: The title of the test results. Required if id doesn't set.
- `--description`, `-d`: The description of the test results. Optional.
- `--format`: The format of the test results file. Required. Allow values: `junit`, `qase`, `allure`, `xctest`.
- `--path`: The path to the test results file or folder. Required.
- `--steps`: The mode of upload steps for XCTest. Optional. Allow values: `all`, `user`.
- `--batch`: The batch number of the test results. Optional. Default is 200.
- `--suite`, `-s`: The suite name of the test results. Optional.
- `--verbose`, `-v`: Enable verbose mode. Optional.

The following example shows how to upload test results in the JUnit format for a test run with the ID `1` in the project
with the code `PROJ`:

```bash
qli testops result upload --project PROJ --token <token> --id 1 --format junit --path /path/to/results.xml --verbose
```

The following example shows how to upload test results in the Qase format for a test run with the ID `1` in the project
with the code `PROJ`:

```bash
qli testops result upload --project PROJ --token <token> --id 1 --format qase --path /path/to/results.json --verbose
```

The following example shows how to upload test results in the Allure format for a test run with the ID `1` in the
project
with the code `PROJ`:

```bash
qli testops result upload --project PROJ --token <token> --id 1 --format allure --path /path/to/allure-results --verbose
```

The following example shows how to upload test results in the XCTest format for a test run with the ID `1` in the
project
with the code `PROJ`:

```bash
qli testops result upload --project PROJ --token <token> --id 1 --format xctest --steps user --path /path/to/xctest-results --verbose
```
