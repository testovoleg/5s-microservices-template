package dto

type EventStatusDto string

const (
	EVENT_STATUS_UNKNOWN  EventStatusDto = "НЕТ_ДАННЫХ"
	EVENT_STATUS_CREATED  EventStatusDto = "СОЗДАН"
	EVENT_STATUS_FINISHED EventStatusDto = "ЗАВЕРШЕН"
	EVENT_STATUS_ERROR    EventStatusDto = "ОШИБКА"
)

type EventDto struct {
	// Идентификатор ивента
	Uuid string `json:"uuid,omitempty"`
	// Описание (например, текст ошибки)
	Description string `json:"description,omitempty"`
	// Время последнего обновления статуса
	Timestamp string `json:"timestamp,omitempty"`
	// Статус ивента
	Status EventStatusDto `json:"status,omitempty"`
}

type EventGetReqDto struct {
	// Идентификатор компании
	CompanyUuid string `query:"company_uuid" swaggerignore:"true"`
	// Идентификатор ивента
	EventUuid string `query:"event_uuid" swaggerignore:"true"`
}
