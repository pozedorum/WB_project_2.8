#!/bin/bash

TEST_PATH="tests/test_file.txt"
TEMP_DIR=$(mktemp -d)
MYGREP_OUT="$TEMP_DIR/mygrep.out"
GREP_OUT="$TEMP_DIR/grep.out"

# Создаём тестовый файл
# cat > "$TEST_PATH" <<EOF
# first line
# second line
# TEST LINE
# third line
# fourth TEST line
# fifth line
# TEST LINE AGAIN
# sixth line
# EOF

# Функция для выполнения теста
run_test() {
    local test_name=$1
    local test_description=$2
    shift 2
    
    echo "Running test: $test_name"
    echo "Description: $test_description"
    echo "Command: $@"
    
    # Запускаем mygrep и сохраняем вывод
    "$@" > "$MYGREP_OUT" 2>&1
    
    # Запускаем оригинальный grep с теми же аргументами
    if [[ $1 == "./mygrep" ]]; then
        shift
        grep "${@}" > "$GREP_OUT" 2>&1
    else
        grep "${@:2}" > "$GREP_OUT" 2>&1
    fi
    
    # Сравниваем вывод
    if diff -u "$MYGREP_OUT" "$GREP_OUT"; then
        echo "✅ Test PASSED"
    else
        echo "❌ Test FAILED"
    fi
    echo "--------------------------------------"
}

# 1. Базовый тест (флаг: без флагов)
run_test "Basic search" "Поиск без флагов, проверка основного функционала" \
    ./mygrep "TEST" "$TEST_PATH"

# 2. Игнорирование регистра (флаг: -i)
run_test "Case insensitive" "Поиск с игнорированием регистра (-i)" \
    ./mygrep -i "test" "$TEST_PATH"

# 3. Инвертированный поиск (флаг: -v)
run_test "Inverted search" "Инвертированный поиск, вывод несовпадающих строк (-v)" \
    ./mygrep -v "TEST" "$TEST_PATH"

# 4. Контекст после совпадения (флаг: -A)
run_test "Context after match" "Вывод 1 строки после совпадения (-A 1)" \
    ./mygrep -A 1 "TEST" "$TEST_PATH"

# 5. Контекст до совпадения (флаг: -B)
run_test "Context before match" "Вывод 1 строки до совпадения (-B 1)" \
    ./mygrep -B 1 "TEST" "$TEST_PATH"

# 6. Контекст вокруг совпадения (флаг: -C)
run_test "Context around match" "Вывод 1 строки до и после совпадения (-C 1)" \
    ./mygrep -C 1 "TEST" "$TEST_PATH"

# 7. Подсчёт строк (флаг: -c)
run_test "Count matches" "Только подсчёт количества совпадений (-c)" \
    ./mygrep -c "TEST" "$TEST_PATH"

# 8. Номера строк (флаг: -n)
run_test "Line numbers" "Вывод номеров строк с совпадениями (-n)" \
    ./mygrep -n "TEST" "$TEST_PATH"

# 9. Фиксированная строка (флаг: -F)
run_test "Fixed string" "Поиск фиксированной строки, не как regex (-F)" \
    ./mygrep -F "TEST" "$TEST_PATH"

# 10. Комбинированные флаги (-i, -n, -A)
run_test "Combined flags 1" "Комбинация: игнорирование регистра + номера строк + контекст (-i -n -A 1)" \
    ./mygrep -i -n -A 1 "test" "$TEST_PATH"

# 11. Комбинированные флаги (-in, -B) - объединённые флаги
run_test "Combined flags 2" "Комбинация: объединённые флаги (-in) + контекст (-B 2)" \
    ./mygrep -in -B 2 "test" "$TEST_PATH"

# 12. Комбинированные флаги (-iv, -C)
run_test "Combined flags 3" "Комбинация: игнорирование регистра + инвертирование + контекст (-iv -C 1)" \
    ./mygrep -iv -C 1 "test" "$TEST_PATH"

# 13. Комбинированные флаги (-nF)
run_test "Combined flags 4" "Комбинация: номера строк + фиксированная строка (-nF)" \
    ./mygrep -nF "TEST" "$TEST_PATH"

# Удаляем временные файлы
rm -rf "$TEMP_DIR"
#rm -f "$TEST_PATH"