#!/bin/bash
# Comparative GShell Test Script

make build
GSHELL_PATH="./GShell"
TEST_DIR="./shell_test"
mkdir -p "$TEST_DIR"

TEMP_GSH_OUT="$TEST_DIR/gshell_output.tmp"
TEMP_STD_OUT="$TEST_DIR/std_output.tmp"

print_test_header() {
    echo -e "\n[TEST $1] $2"
    echo "Command: $3"
}

run_comparison() {
    local test_num=$1
    local description=$2
    local command=$3
    
    print_test_header "$test_num" "$description" "$command"
    
    # Run in GShell
    echo -e "=== GShell Output ==="
    echo -e "$command" | "$GSHELL_PATH" 2>&1 | tee "$TEMP_GSH_OUT"
    local gsh_status=${PIPESTATUS[0]}
    
    # Run in standard shell
    echo -e "\n=== Standard Shell Output ==="
    eval "$command" 2>&1 | tee "$TEMP_STD_OUT"
    local std_status=${PIPESTATUS[0]}
}

echo "Starting Comparative GShell Test Suite - $(date)"
echo "--------------------------------------"

# 1. Built-in commands
run_comparison "1.1" "Current directory" "pwd"
run_comparison "1.2" "Change directory" "cd /tmp && pwd && cd ~ && pwd"
run_comparison "1.3" "Echo commands" 'echo "Hello World" && echo "Home: $HOME"'

# 2. External commands
run_comparison "2.1" "List files" "ls -l"
run_comparison "2.2" "Find executable" "which bash"

# 3. Pipes
run_comparison "3.1" "Simple pipe" "ls | wc -l"
run_comparison "3.2" "Text transformation" 'echo "test" | tr a-z A-Z'

# 4. Conditional execution
run_comparison "4.1" "Logical AND" 'true && echo "Success"'
run_comparison "4.2" "Logical OR" 'false || echo "Failed"'

# 5. Redirections
run_comparison "5.1" "Output redirection" 'echo "Test" > test_output.txt && cat test_output.txt'

# 6. Environment variables
run_comparison "6.1" "Standard variables" 'echo "PATH: $PATH" && echo "HOME: $HOME"'

# Cleanup
rm -f "$TEMP_GSH_OUT" "$TEMP_STD_OUT" "test_output.txt"

echo -e "\n--------------------------------------"
