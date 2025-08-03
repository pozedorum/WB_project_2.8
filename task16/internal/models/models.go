package models

import "net/url"

type Resource struct {
	URL         *url.URL    // Оригинальный URL
	LocalPath   string      // Локальный путь (например "css/styles.css")
	ContentType string      // MIME-тип (например "text/html")
	Content     []byte      // Скачанные данные
	Links       []*Resource // Ссылки на связанные ресурсы (для HTML)
	IsHTML      bool        // Флаг HTML-страницы
}

// Сюда можно поместить уровень рекурсии
