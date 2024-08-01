package messages

const (
	MsgErrInternal       = "Внутренняя ошибка"
	MsgErrUnderCon       = "Функционал находится в разработке"
	MsgErrInvalidCommand = "Неправильная команда"
	MsgErrEmptyList      = "Пустой результат"
	MsgErrBadRequest     = "Неправильный запрос"

	MsgErrFoodNotFound = "Еда не найдена в базе данных"
	MsgErrFoodIsUsed   = "Еда уже используется в журнале приема пищи или бандле"

	MsgErrBundleDepBundleNotFound  = "Зависимый бандл не найден в базе данных"
	MsgErrBundleDepFoodNotFound    = "Зависимая еда не найдена в базе данных"
	MsgErrBundleDepBundleRecursive = "Зависимый бандл не может быть рекурсивным"
	MsgErrBundleNotFound           = "Бандл не найден"
	MsgErrBundleIsUsed             = "Бандл уже используется в другом бандле"

	MsgErrUserSettingsNotFound = "Не найдены пользовательские настройки"

	MsgErrJournalCopy = "Не пустое назначение копирования"
	MsgJournalCopied  = "Скопировано записей: %d"

	MsgOK = "OK"
)
