package messages

const (
	MsgErrInternal       = "Внутренняя ошибка"
	MsgErrUnderCon       = "Функционал находится в разработке"
	MsgErrInvalidCommand = "Неправильная команда"
	MsgErrEmptyList      = "Пустой результат"

	MsgErrFoodNotFound = "Еда не найдена в базе данных"
	MsgErrFoodIsUsed   = "Еда уже используется в дневнике приема пищи"

	MsgErrUserSettingsNotFound = "Не найдены пользовательские настройки"

	MsgErrJournalNotStartOfWeek = "Дата не является началом недели"
	MsgErrJournalCopy           = "Не пустое назначение копирования"
	MsgJournalCopied            = "Скопировано записей: %d"

	MsgOK = "OK"
)

const (
	MsgClassError   = "danger"
	MsgClassWarning = "warning"
	MsgClassInfo    = "primary"
)
