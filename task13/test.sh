#!/bin/bash

TEST_FILE="./tests/test_file.txt"
TEMP_DIR=$(mktemp -d)
MYCUT_OUT="$TEMP_DIR/mycut.out"
CUT_OUT="$TEMP_DIR/cut.out"

# Функция для выполнения теста
run_test() {
    local test_name=$1
    local test_description=$2
    local mycut_cmd=$3
    local cut_cmd=$4
    
    echo "Running test: $test_name"
    echo "Description: $test_description"
    echo "Command: $mycut_cmd"
    
    # Запускаем mycut и сохраняем вывод
    eval "$mycut_cmd" > "$MYCUT_OUT" 2>&1
    
    # Запускаем оригинальный cut
    eval "$cut_cmd" > "$CUT_OUT" 2>&1
    
    # Сравниваем вывод
    if diff -u "$MYCUT_OUT" "$CUT_OUT"; then
        echo "✅ Test PASSED"
    else
        echo "❌ Test FAILED"
        echo "--- Output diff ---"
        diff -y --suppress-common-lines "$MYCUT_OUT" "$CUT_OUT"
    fi
    echo "--------------------------------------"
}

# 1. Базовый тест (выбор полей)
run_test "Basic field selection" "Выбор отдельных полей" \
    "./mycut -f 2,4 -d ':' $TEST_FILE" \
    "cut -f 2,4 -d ':' $TEST_FILE"

# 2. Диапазон полей
run_test "Field range" "Выбор диапазона полей" \
    "./mycut -f 2-4 -d ':' $TEST_FILE" \
    "cut -f 2-4 -d ':' $TEST_FILE"

# 3. Обратный порядок полей
run_test "Reverse order" "Поля в обратном порядке" \
    "./mycut -f 5,3,1 -d ':' $TEST_FILE" \
    "cut -f 5,3,1 -d ':' $TEST_FILE"

# 4. Пустые поля
run_test "Empty fields" "Обработка пустых полей" \
    "./mycut -f 1,3,5 -d ':' $TEST_FILE" \
    "cut -f 1,3,5 -d ':' $TEST_FILE"

# 5. Флаг -s (только строки с разделителем)
run_test "Separated flag" "Только строки с разделителем (-s)" \
    "./mycut -f 2 -d ':' -s $TEST_FILE" \
    "cut -f 2 -d ':' -s $TEST_FILE"

# 6. Несуществующие поля
run_test "Non-existent fields" "Запрос несуществующих полей" \
    "./mycut -f 10 -d ':' $TEST_FILE" \
    "cut -f 10 -d ':' $TEST_FILE"

# 7. Разные разделители
run_test "Different delimiter" "Использование другого разделителя" \
    "./mycut -f 2 -d 'e' $TEST_FILE" \
    "cut -f 2 -d 'e' $TEST_FILE"

# 8. Комбинированные флаги
run_test "Combined flags" "Комбинация -s и диапазона полей" \
    "./mycut -f 2-3 -d ':' -s $TEST_FILE" \
    "cut -f 2-3 -d ':' -s $TEST_FILE"

# 9. Повторяющиеся поля
run_test "Duplicate fields" "Повторяющиеся номера полей" \
    "./mycut -f 2,2,2 -d ':' $TEST_FILE" \
    "cut -f 2,2,2 -d ':' $TEST_FILE"

# Удаляем временные файлы
rm -rf "$TEMP_DIR"