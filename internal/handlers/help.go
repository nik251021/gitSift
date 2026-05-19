package handlers

import "fmt"

const helpText = `
goSift — твой AI помощник для работы с репозиториями Git.

Использование:
  goSift <команда> [аргументы]

Команды:
  setAPIkey      Сохранить API ключ в .env файл
  getAPIkey      Получить текущий API ключ из .env файла
  help           Показать эту справку

Примеры:
  goSift setAPIkey "sk-123..."
  goSift getAPIkey
  goSift help
`

func HelpHandler(args []string) error {
	fmt.Print(helpText)
	return nil
}
