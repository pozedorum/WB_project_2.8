# Тестирование утилиты Wget (краткая версия)

## Основные тесты

1. **Простая страница**  
   `./mywget http://example.com 0`  
   - Должен создать `example.com/index.html`
   - Быстро завершиться (1-3 сек)

2. **Страница с картинками**  
   `./mywget https://www.techinsider.ru/science/1573409-neveroyatnye-fakty-ob-utkonosah-vy-porazites 1`  
   - Должен скачать HTML + изображения/CSS/JS
   - Игнорировать битые ссылки (favicon.ico)
   (скачивает статью об утконосах, страницы оглавлений по темам, скачал ещё одну статью про кошек, как я понимаю там прямая ссылка , всё лежит в папке science)
   Статься об утконосах имеет ссылку на саму себя
   `./mywget https://www.techinsider.ru/science/1573409-neveroyatnye-fakty-ob-utkonosah-vy-porazites 1 | grep 1573409-neveroyatnye-fakty-ob-utkonosah-vy-porazites`


3. **Статья с Хабра**
   `./mywget https://habr.com/ru/articles/934286/ 1`
   Большинство информационных ресурсов содержит тэги, которые не получается обработать, но вот с Хабром всё хорошо получается
