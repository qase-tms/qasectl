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
- `--api-host`: The Qase API host. Optional. Default is `api.qase.io`. Use it to point the CLI at a Qase Enterprise
  instance, for example `api.qase.example.com`. Can also be set via the `QASE_TESTOPS_API_HOST` environment variable.
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
- `--api-host`: The Qase API host. Optional. Default is `api.qase.io`. Use it to point the CLI at a Qase Enterprise
  instance, for example `api.qase.example.com`. Can also be set via the `QASE_TESTOPS_API_HOST` environment variable.
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
- `--api-host`: The Qase API host. Optional. Default is `api.qase.io`. Use it to point the CLI at a Qase Enterprise
  instance, for example `api.qase.example.com`. Can also be set via the `QASE_TESTOPS_API_HOST` environment variable.
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
- `--api-host`: The Qase API host. Optional. Default is `api.qase.io`. Use it to point the CLI at a Qase Enterprise
  instance, for example `api.qase.example.com`. Can also be set via the `QASE_TESTOPS_API_HOST` environment variable.
- `--id`: The ID of the test run to upload results for. Required if title doesn't set.
- `--title`: The title of the test results. Required if id doesn't set.
- `--description`, `-d`: The description of the test results. Optional.
- `--format`: The format of the test results file. Required. Allow values: `junit`, `qase`, `allure`, `xctest`, `maestro`.
- `--path`: The path to the test results file or folder. Required. For `maestro` it is the folder of one run.
- `--steps`: The mode of upload steps for XCTest and Maestro. Optional. Allow values: `all`, `user`. `user` keeps only the steps
  written in the test code and drops the internal steps XCUITest records for every action, which makes large UI test
  runs much smaller. For Maestro, `all` also shows the service commands (`defineVariables`, `applyConfiguration`),
  while `user` (default) hides them.
- `--step-attachments`: Which XCTest steps keep their attachments. Optional. Allow values: `failed`, `all`, `none`.
  Default is `failed`. XCUITest takes a screenshot and records several debug files for every action, so `all` can
  add thousands of files to a run. The attachments of the failure itself are always uploaded.
- `--batch`: The batch number of the test results. Optional. Default is 200. A batch is also split when its results
  hold more than 50,000 steps in total, nested steps included: the Qase API rejects a larger request.
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

## Maestro

The `maestro` format reads the debug output of a Maestro 2.x run; Maestro 1.x is not supported. Run the flows with
the debug output and the JUnit report in one folder:

```bash
maestro test --debug-output <run-dir> --flatten-debug-output --format junit --output <run-dir>/report.xml <flows>
qasectl testops result upload --project PROJ --token <token> --id 1 --format maestro --path <run-dir> --verbose
```

`--path` is the folder of ONE run. If it holds several runs, the command fails with an error. Without the JUnit
report the device names are not set. The report is searched only inside `--path`, so pass the run folder (not a
single flow folder) to get the `device` parameter.

Each flow becomes one test result, and its commands become steps. Set the Qase attributes in the flow `properties`.
Maestro accepts only string values, so several IDs are a comma-separated string:

```yaml
appId: com.example.app
name: Search products
tags:
  - smoke
properties:
  qase.id: "123, 456"
  qase.suite: "Search.Smoke"
  severity: critical
---
- launchApp
```

- `qase.id`: the IDs of the Qase test cases the result belongs to. The `PRJ-123` form is accepted too (the number
  after the last `-` is used). Elements that cannot be parsed are skipped with a warning. With several IDs, Qase
  records a result for each case.
- `qase.suite`: the suite of the test case. Levels are separated by a dot.
- Any other key becomes a result field with the same name; Qase applies it only if the project has a system or custom
  field with that name. `qase.*` keys are not sent as fields. The flow `tags` become the tags of the test case.

Without `qase.id`, Qase matches the result to a test case by its signature: `qase.suite`, the flow file name, the flow
name and the `device` parameter. Renaming the flow or running it on another device creates a new test case.

How a flow is reported:

- The result title is the flow `name`, or the flow file name when `name` is not set.
- Steps keep the nesting of `runFlow`, so a failure inside a subflow is shown on the step where it happened, and its
  error becomes the result message.
- A retried command is shown as `Attempt N` steps, a repeated one as `Iteration N` steps. The result gets the
  `is_flaky` field only when the flow finally passes after a failed retry attempt.
- A failed `optional: true` command is a passed step with the comment `Optional command failed`.
- A flow that did not finish (Maestro or the device stopped mid-run) is failed with the message `Flow did not finish`.
  A flow with an empty or broken `commands.json` is skipped with a warning.
- With several shards, the shard number is written as the thread `shard-N`. The device name from the JUnit report is
  written to the `device` parameter; `--skip-params` drops it.
- Screenshots Maestro takes on failure and with `takeScreenshot` are attached to their steps. Screen hierarchy dumps
  are not uploaded.
- `maestro.log` and the device logs are attached only to failed results. Files over 32 MB are gzipped; a file that is
  still over 32 MB after gzip is skipped with a warning. `--attachment-extensions` matches the gzipped log as `gz`.

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
- `--api-host`: The Qase API host. Optional. Default is `api.qase.io`. Use it to point the CLI at a Qase Enterprise
  instance, for example `api.qase.example.com`. Can also be set via the `QASE_TESTOPS_API_HOST` environment variable.
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
- `--api-host`: The Qase API host. Optional. Default is `api.qase.io`. Use it to point the CLI at a Qase Enterprise
  instance, for example `api.qase.example.com`. Can also be set via the `QASE_TESTOPS_API_HOST` environment variable.
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
- `--api-host`: The Qase API host. Optional. Default is `api.qase.io`. Use it to point the CLI at a Qase Enterprise
  instance, for example `api.qase.example.com`. Can also be set via the `QASE_TESTOPS_API_HOST` environment variable.
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
- `--api-host`: The Qase API host. Optional. Default is `api.qase.io`. Use it to point the CLI at a Qase Enterprise
  instance, for example `api.qase.example.com`. Can also be set via the `QASE_TESTOPS_API_HOST` environment variable.
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
