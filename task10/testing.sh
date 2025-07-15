#!/bin/bash

# Конфигурация тестов
TEST_INPUT="tests/test_input.txt"

TESTS=(
    # Название теста         Флаги        Ожидаемый вывод
    "basic_sort"            ""           "tests/expected_basic.txt"
    "numeric_sort"          "-n -k 2"    "tests/expected_numeric.txt"
    "month_sort"            "-m -k 3"    "tests/expected_month.txt"
    "human_sort"            "-h -k 4"    "tests/expected_human.txt"
    "reverse_sort"          "-r"         "tests/expected_reverse.txt"
    "unique_sort"           "-u"         "tests/expected_unique.txt"
    "combined_sort"         "-n -r -k 2" "tests/expected_combined.txt"
    "complex_sort"          "-k 3 -m -u" "tests/expected_complex.txt"
)

# Функция сравнения результатов
run_test() {
    local test_name=$1
    local flags=$2
    local expected_file=$3
    
    echo "Running test: $test_name"
    echo "Command: ./mysort $flags $TEST_INPUT"
    
    # Удаляем предыдущий результат, если есть
    rm -f "${TEST_INPUT%.txt}_sorted.txt"
    
    # Выполняем команду
    ./mysort $flags "$TEST_INPUT"
    
    # Сравниваем с ожидаемым результатом
    if diff -u "$expected_file" "${TEST_INPUT%.txt}_sorted.txt" > "diff_${test_name}.txt"; then
        echo "✅ PASSED: $test_name"
        rm "diff_${test_name}.txt"
        return 0
    else
        echo "❌ FAILED: $test_name"
        echo "Differences:"
        cat "diff_${test_name}.txt"
        rm "diff_${test_name}.txt"
        return 1
    fi
}

# Основная функция
main() {
    if [ ! -f "$TEST_INPUT" ]; then
        echo "Error: Test input file $TEST_INPUT not found!"
        echo "Please create it manually with following content:"
        echo
        cat <<EOF
apple 5 Jan 100K
banana 3 Feb 200M
apple 5 Mar 300G
cherry 1 Apr 400
date 12 May 500T
fig 7 Jun 600K
grape 8 Jul 700
banana 9 Aug 800M
kiwi 11 Sep 900G
lemon 4 Oct 1000
orange 2 Nov 1100T
pear 10 Dec 1200K
EOF
        echo
        exit 1
    fi
    
    total_tests=0
    passed_tests=0
    
    for ((i=0; i<${#TESTS[@]}; i+=3)); do
        ((total_tests++))
        if run_test "${TESTS[i]}" "${TESTS[i+1]}" "${TESTS[i+2]}"; then
            ((passed_tests++))
        fi
        echo
    done
    
    # Итоговая статистика
    echo "Test results:"
    echo "✅ $passed_tests passed"
    echo "❌ $((total_tests - passed_tests)) failed"
    
    # Возвращаем код ошибки, если есть проваленные тесты
    if [ $passed_tests -ne $total_tests ]; then
        exit 1
    fi
}

main