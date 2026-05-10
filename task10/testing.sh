#!/bin/bash

set -u

TEST_INPUT="tests/test_input.txt"
OUTPUT_DIR="tests/output"

TESTS=(
    "basic_sort"            ""           "tests/expected_basic.txt"
    "numeric_sort"          "-n -k 2"    "tests/expected_numeric.txt"
    "month_sort"            "-M -k 3"    "tests/expected_month.txt"
    "human_sort"            "-h -k 4"    "tests/expected_human.txt"
    "reverse_sort"          "-r"         "tests/expected_reverse.txt"
    "unique_sort"           "-u"         "tests/expected_unique.txt"
    "combined_sort"         "-n -r -k 2" "tests/expected_combined.txt"
    "complex_sort"          "-k 3 -M -u" "tests/expected_complex.txt"
)

run_test() {
    local test_name=$1
    local flags=$2
    local expected_file=$3

    local output_file="${OUTPUT_DIR}/${test_name}.out"
    local diff_file="${OUTPUT_DIR}/${test_name}.diff"

    echo "Running test: $test_name"

    rm -f "$output_file" "$diff_file"

    # shellcheck disable=SC2206
    local flags_array=($flags)

    echo "Command: ./mysort ${flags_array[*]} $TEST_INPUT > $output_file"

    if ! ./mysort "${flags_array[@]}" "$TEST_INPUT" > "$output_file"; then
        echo "❌ FAILED: $test_name"
        echo "Program exited with error"
        return 1
    fi

    if diff -u "$expected_file" "$output_file" > "$diff_file"; then
        echo "✅ PASSED: $test_name"
        rm -f "$diff_file"
        return 0
    fi

    echo "❌ FAILED: $test_name"
    echo "Differences:"
    cat "$diff_file"
    return 1
}

main() {
    if [ ! -x ./mysort ]; then
        echo "Error: ./mysort not found or not executable"
        echo "Build it first, for example:"
        echo "go build -o mysort ./cmd/mysort"
        exit 1
    fi

    if [ ! -f "$TEST_INPUT" ]; then
        echo "Error: Test input file $TEST_INPUT not found!"
        exit 1
    fi

    mkdir -p "$OUTPUT_DIR"

    local total_tests=0
    local passed_tests=0

    for ((i=0; i<${#TESTS[@]}; i+=3)); do
        ((total_tests++))

        if run_test "${TESTS[i]}" "${TESTS[i+1]}" "${TESTS[i+2]}"; then
            ((passed_tests++))
        fi

        echo
    done

    echo "Test results:"
    echo "✅ $passed_tests passed"
    echo "❌ $((total_tests - passed_tests)) failed"

    if [ "$passed_tests" -ne "$total_tests" ]; then
        exit 1
    fi
}

main