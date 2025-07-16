#!/bin/bash

# Конфигурация тестов
TEST_INPUT="tests/test_input.txt"
EXPECTED_OUTPUT="tests/test_input_expected.txt"
ACTUAL_OUTPUT="tests/test_input_sorted.txt"

TESTS=(
    # Название теста         Флаги        Ожидаемый вывод (не используется)
    "basic_sort"            ""           ""
    "numeric_sort"          "-n -k 2"    ""
    "month_sort"            "-M -k 3"    ""  # Оригинальный sort использует -M для месяцев
    "human_sort"            "-h -k 4"    ""
    "reverse_sort"          "-r"         ""
    "unique_sort"           "-u"         ""
    "combined_sort"         "-n -r -k 2" ""
    "complex_sort"          "-k 3 -M -u" ""
)

# Функция сравнения результатов
run_test() {
    local test_name=$1
    local flags=$2
    local expected_file=$3
    
    echo "Running test: $test_name"
    echo "Command: ./mysort $flags $TEST_INPUT"
    
    # Удаляем предыдущие результаты
    rm -f "$EXPECTED_OUTPUT" "$ACTUAL_OUTPUT"
    
    # 1. Запускаем оригинальный sort и сохраняем результат
    echo "Original sort command: sort $flags $TEST_INPUT"
    sort $flags "$TEST_INPUT" > "$EXPECTED_OUTPUT"
    
    # 2. Запускаем вашу утилиту
    ./mysort $flags "$TEST_INPUT" >> "$ACTUAL_OUTPUT"
    # mv "${TEST_INPUT%.txt}_sorted.txt" "$ACTUAL_OUTPUT"
    
    # 3. Сравниваем результаты
    if diff -u "$EXPECTED_OUTPUT" "$ACTUAL_OUTPUT" > "diff_${test_name}.txt"; then
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