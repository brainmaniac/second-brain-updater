# Second-Brain Updater

## Overview
The Second-Brain Updater is a Go application designed to automate the creation of a daily schedule based on a to-do list. It utilizes the OpenAI API to generate an organized, priority-based schedule in org-mode format. The application reads a to-do list, formats a prompt for OpenAI, and processes the response to create a daily schedule.

## Requirements
- Go programming environment
- OpenAI API key
- `.env` file with necessary environment variables

## Installation
1. Clone the repository to your local machine.
2. Navigate to the project directory.
3. Rename `.env.example` to `.env` and ensure the variables in the file are properly set with your data.
4. Run `go build` to compile the application.

## Usage
Execute the compiled binary to generate the daily schedule. The schedule is based on the to-do list specified in the `.env` file and is written to a file named with the current date in org-mode format.

## Functionality
- `initConfig()`: Initializes configuration variables from the environment.
- `readTodoList()`: Reads the to-do list from a file.
- `addPrePrompt()`: Formats the OpenAI API prompt with the to-do list.
- `writeDailySchedule()`: Writes the generated schedule to a file.
- `callOpenAI()`: Sends a request to the OpenAI API with the prompt.
- `extractContent()`: Extracts the schedule content from the API response.

## Customization
- Modify the `.env` file to point to your to-do list location.
- Adjust the `prePrompt` string in `initConfig()` for different prompt structures.
- Change the `temperature` parameter in `main()` to control the creativity of the OpenAI response.

## Notes
- Ensure that your OpenAI API key has sufficient permissions and quota.
- The format of the to-do list should be compatible with the expected org-mode format.

## Troubleshooting
- Check the `.env` file for correct paths and file names.
- Ensure your OpenAI API key is valid and active.
- Verify the format of the to-do list for compatibility.

## Contributing
Contributions to the project are welcome. Please follow the standard GitHub pull request process.
