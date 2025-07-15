# Исходный тестовый файл
TEST_INPUT := tests/test_input.txt
TEST_OUTPUT := $(basename $(TEST_INPUT))ed.txt

# Список всех тестов
TESTS := \
	basic \
	numeric \
	month \
	human \
	reverse \
	unique \
	combined \
	complex

# Флаги для каждого теста
basic_FLAGS := 
numeric_FLAGS := -n -k 2
month_FLAGS := -m -k 3
human_FLAGS := -h -k 4
reverse_FLAGS := -r
unique_FLAGS := -u
combined_FLAGS := -n -r -k 2
complex_FLAGS := -k 3 -m -u

# Ожидаемые файлы результатов
EXPECTED_FILES := $(addprefix tests/expected_, $(addsuffix .txt, $(TESTS)))

.PHONY: all clean $(TESTS)

all: mysort
	./testing.sh
	@echo "Все тесты выполнены"

small_test: mysort
	./mysort -r tests/small_test.txt

mysort: clean
	go build -o mysort ./cmd/main.go



clean:
	rm -rf mysort
	rm -f $(TEST_OUTPUT)







# Общее правило для запуска тестов
$(TESTS): %: $(TEST_INPUT)
	@echo "Запуск теста $@ с флагами: $($@_FLAGS)"
	@rm -f $(TEST_OUTPUT)
	@./mysort $($@_FLAGS) $(TEST_INPUT)
	@echo "Результат сохранен в $(TEST_OUTPUT)"
	@echo "Ожидаемый результат в tets/expected_$@.txt"
#	@echo "Для сравнения выполните:"
#   @echo "diff -u expected_$@.txt $(TEST_OUTPUT) || meld expected_$@.txt $(TEST_OUTPUT)"

# Проверка наличия тестового файла
$(TEST_INPUT):
	@echo "Тестовый файл $(TEST_INPUT) не найден!"
	@echo "Создайте его вручную со следующим содержимом:"
	@echo "----------------------------------------"
	@echo "apple 5 Jan 100K"
	@echo "banana 3 Feb 200M"
	@echo "apple 5 Mar 300G"
	@echo "cherry 1 Apr 400"
	@echo "date 12 May 500T"
	@echo "fig 7 Jun 600K"
	@echo "grape 8 Jul 700"
	@echo "banana 9 Aug 800M"
	@echo "kiwi 11 Sep 900G"
	@echo "lemon 4 Oct 1000"
	@echo "orange 2 Nov 1100T"
	@echo "pear 10 Dec 1200K"
	@echo "----------------------------------------"
	@exit 1

	