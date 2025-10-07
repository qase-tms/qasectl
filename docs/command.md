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

## Example usage

```bash
qasectl testops run create --project <project_code> --token <token> --title <title> --description <description> --environment <environment> --milestone <milestone> --plan <plan> --cloud --browser <browser> --verbose
```

The `create` command has the following options:

- `--project`, `-p`: The project code where the test run will be created. Required.
- `--token`, `-t`: The API token to authenticate with the TestOps API. Required.
- `--title`: The name of the test run. Required.
- `--description`, `-d`: The description of the test run. Optional.
- `--environment`, `-e`: The environment where the test run will be executed. Optional.
- `--milestone`, `-m`: The milestone of the test run. Optional.
- `--plan`: The test plan of the test run. Optional.
- `--tags`: The tags of the test run. Optional.
- `--cloud`: Mark the test run as a cloud run. Optional. Must be used together with `--browser`.
- `--browser`: The browser for the cloud run. Optional. Must be used together with `--cloud`. Allowed values: `chromium`, `firefox`, `webkit`.
- `--output`, `-o`: The output path to save the test run ID. Optional. Default is `qase.env` in the current
  directory.
- `--verbose`, `-v`: Enable verbose mode. Optional.

The following example shows how to create a test run in the project with the code `PROJ`:

```bash
qasectl testops run create --project PROJ --token <token> --title "Test Run 1" --description "This is a test run" --environment "Production" --milestone "Milestone 1" --plan "Test Plan 1" --tags "tag1,tag2" --verbose
```

The following example shows how to create a cloud test run with a specific browser:

```bash
qasectl testops run create --project PROJ --token <token> --title "Cloud Test Run" --description "This is a cloud test run" --cloud --browser "chromium" --verbose
```

# Complete a test run

You can complete a test run by using the `complete` command. The `complete` command is used to complete a test run in
the specified project.

## Example usage

```bash
qasectl testops run complete --project <project_code> --token <token> --id <run_id> --verbose
```

The `complete` command has the following options:

- `--project`, `-p`: The project code where the test run will be completed. Required.
- `--token`, `-t`: The API token to authenticate with the TestOps API. Required.
- `--id`: The ID of the test run to complete. Required.
- `--verbose`, `-v`: Enable verbose mode. Optional.

The following example shows how to complete a test run with the ID `1` in the project with the code `PROJ`:

```bash
qasectl testops run complete --project PROJ --token <token> --id 1 --verbose
```

# Delete test runs

You can delete test runs by using the `delete` command. The `delete` command is used to delete test runs in the
specified project.

## Example usage

```bash
qasectl testops run delete --project <project_code> --token <token> --ids <run_id> --verbose
```

The `delete` command has the following options:

- `--project`, `-p`: The project code where the test runs will be deleted. Required.
- `--token`, `-t`: The API token to authenticate with the TestOps API. Required.
- `--ids`: The IDs of the test runs to delete. Optional if all doesn't set.
- `--all`: Delete all test runs in the project. Optional if ids doesn't set.
- `--start`, `-s`: The start date of the test runs to delete. Optional.
- `--end`, `-e`: The end date of the test runs to delete. Optional.
- `--verbose`, `-v`: Enable verbose mode. Optional.

The following example shows how to delete a test run with the ID `1` in the project with the code `PROJ`:

```bash
qasectl testops run delete --project PROJ --token <token> --ids 1 --verbose
```

The following example shows how to delete all test runs in the project with the code `PROJ`:

```bash
qasectl testops run delete --project PROJ --token <token> --all --verbose
```

The following example shows how to delete all test runs in the project with the code `PROJ` that were created between
`2022-01-01` and `2022-12-31`:

```bash
qasectl testops run delete --project PROJ --token <token> --all --start "2022-01-01" --end "2022-12-31" --verbose
```

# Upload test results

You can upload test results by using the `upload` command. The `upload` command is used to upload test results for a
test run in the specified project.

## Example usage

```bash
qasectl testops result upload --project <project_code> --token <token> --id <run_id> --format <format> --path <results_file> --batch <batch> --verbose
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
- `--replace-statuses`, `-r`: The statuses to replace. Optional. Pass like '{\"Passed\": \"Failed\"}' to replace all passed results with failed. Note: Use slugs of statuses.
- `--skip-params`: Skip parameters for the results. Optional.
- `--attachment-extensions`: Comma-separated list of file extensions to filter attachments. If not specified, all attachments will be uploaded. Optional.
- `--verbose`, `-v`: Enable verbose mode. Optional.

The following example shows how to upload test results in the JUnit format for a test run with the ID `1` in the project
with the code `PROJ`:

```bash
qasectl testops result upload --project PROJ --token <token> --id 1 --format junit --path /path/to/results.xml --verbose
```

The following example shows how to upload test results in the Qase format for a test run with the ID `1` in the project
with the code `PROJ`:

```bash
qasectl testops result upload --project PROJ --token <token> --id 1 --format qase --path /path/to/results --verbose
```

The following example shows how to upload test results in the Allure format for a test run with the ID `1` in the
project
with the code `PROJ`:

```bash
qasectl testops result upload --project PROJ --token <token> --id 1 --format allure --path /path/to/allure-results --verbose
```

The following example shows how to upload test results in the XCTest format for a test run with the ID `1` in the
project
with the code `PROJ`:

```bash
qasectl testops result upload --project PROJ --token <token> --id 1 --format xctest --steps user --path /path/to/xctest-results --verbose
```

The following example shows how to upload test results with filtered attachments (only PNG and JPG files) for a test run with the ID `1` in the project
with the code `PROJ`:

```bash
qasectl testops result upload --project PROJ --token <token> --id 1 --format allure --path /path/to/allure-results --attachment-extensions "png,jpg" --verbose
```

# Create an environment

You can create an environment by using the `create` command. The `create` command is used to create a new environment
in the specified project and save an environment slug to a file. You can specify the file path using the `--output`
option.
If the file path is not specified, the environment ID will be saved to `qase.env` in the current directory.

The file will contain the environment slug in the following format:

```text
QASE_ENVIRONMENT=production
```

You can use the environment slug in subsequent steps to specify the environment for a test run.
For exctract environment slug from file you can use command:

```bash
cat qase.env | grep QASE_ENVIRONMENT | cut -d'=' -f2
```

## Example usage

```bash
qasectl testops env create --project <project_code> --token <token> --title <title> --slug <slug> --description <description> --host <host> --verbose
```

The `create` command has the following options:

- `--project`, `-p`: The project code where the environment will be created. Required.
- `--token`, `-t`: The API token to authenticate with the TestOps API. Required.
- `--title` : The name of the environment. Required.
- `--slug`, `-s`: The slug of the environment. Required.
- `--description`, `-d`: The description of the environment. Optional.
- `--host` : The host of the environment. Optional.
- `--output`, `-o`: The output path to save the environment slug. Optional. Default is `qase.env` in the current
  directory.
- `--verbose`, `-v`: Enable verbose mode. Optional.

The following example shows how to create an environment in the project with the code `PROJ`:

```bash
qasectl testops env create --title 'New environment' --slug local --description 'This is an environment' --host app.server.com --project 'PRJ' --token 'TOKEN' --output 'env.env' --verbose
```

# Create a milestone

You can create a milestone by using the `create` command. The `create` command is used to create a new milestone in the
specified project and save a milestone ID to a file. You can specify the file path using the `--output` option. If the
file path is not specified, the milestone ID will be saved to `qase.env` in the current directory.

The file will contain the milestone ID in the following format:

```text
QASE_MILESTONE=123
```

You can use the milestone ID in subsequent steps to specify the milestone for a test run.
For exctract milestone ID from file you can use command:

```bash
cat qase.env | grep QASE_MILESTONE | cut -d'=' -f2
```

## Example usage

```bash
qasectl testops milestone create --project <project_code> --token <token> --title <title> --description <description> --status <status> --due-date <due_date> --verbose
```

The `create` command has the following options:

- `--project`, `-p`: The project code where the milestone will be created. Required.
- `--token`, `-t`: The API token to authenticate with the TestOps API. Required.
- `--title` : The name of the milestone. Required.
- `--description`, `-d`: The description of the milestone. Optional.
- `--status`, `-s`: The status of the milestone. Optional. Allow values: `active`, `completed`.
- `--due-date` : The due date of the milestone. Optional.
- `--output`, `-o`: The output path to save the milestone ID. Optional. Default is `qase.env` in the current directory.
- `--verbose`, `-v`: Enable verbose mode. Optional.

The following example shows how to create a milestone in the project with the code `PROJ`:

```bash
qasectl testops milestone create --project PROJ --token <token> --title "Milestone 1" --description "This is a milestone" --status active --due-date "2022-12-31" --verbose
```

# Get filtered results

You can get filtered results by using the `filter` command. The `filter` command is used to get filtered results for the given plan ID and framework. It prepares the filtered string for the given framework. For example, for playwright framework it will prepare the string like this: `(Qase ID: 1|2|3|...)`. You can use this string in your playwright tests to filter the tests by the given plan ID. You can specify the file path using the `--output` option. If the
file path is not specified, the filtered results will be saved to `qase.env` in the current directory.

The file will contain the filtered results in the following format:

```text
QASE_FILTERED_RESULTS=(Qase ID: 1|2|3|...)
```

You can use the filtered results in subsequent steps to filter the tests by the given plan ID.
For exctract filtered results from file you can use command:

```bash
cat qase.env | grep QASE_FILTERED_RESULTS | cut -d'=' -f2
```

## Example usage

```bash
qasectl testops filter --project <project_code> --token <token> --planID <plan_id> --framework <framework> --output <output_file> --verbose
```

The `filter` command has the following options:

- `--project`, `-p`: The project code where the filtered results will be saved. Required.
- `--token`, `-t`: The API token to authenticate with the TestOps API. Required.
- `--planID` : The ID of the test plan. Required.
- `--framework`, `-f`: The framework of the filtered results. Required. Allow values: `playwright`.
- `--output`, `-o`: The output path to save the filtered results. Optional. Default is `qase.env` in the current directory.
- `--verbose`, `-v`: Enable verbose mode. Optional.

The following example shows how to get filtered results for the plan with the ID `1` in the project with the code `PROJ` and save them to the file `qase.env`:

```bash
qasectl testops filter --project PROJ --token <token> --planID 1 --framework playwright --output qase.env --verbose
```

# Remove custom fields

You can remove custom fields by using the `remove` command. The `remove` command is used to remove custom fields in the
specified project.

## Example usage

```bash
qasectl testops field custom remove --project <project_code> --token <token> --id <field_id> --all --verbose
```

The `remove` command has the following options:

- `--token`, `-t`: The API token to authenticate with the TestOps API. Required.
- `--id`, `-i`: The ID of the custom field to remove. Optional if all doesn't set.
- `--all`: Remove all custom fields in the project. Optional if id doesn't set.
- `--verbose`, `-v`: Enable verbose mode. Optional.

The following example shows how to remove a custom field with the ID `1` in the project with the code `PROJ`:

```bash
qasectl testops field custom remove --project PROJ --token <token> --id 1 --verbose
```

The following example shows how to remove all custom fields in the project with the code `PROJ`:

```bash
qasectl testops field custom remove --project PROJ --token <token> --all --verbose
```
