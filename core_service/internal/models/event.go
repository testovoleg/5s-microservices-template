package models

type EventStatus string

const (
	EVENT_STATUS_UNKNOWN  EventStatus = "НЕТ_ДАННЫХ"
	EVENT_STATUS_CREATED  EventStatus = "СОЗДАН"
	EVENT_STATUS_FINISHED EventStatus = "ЗАВЕРШЕН"
	EVENT_STATUS_ERROR    EventStatus = "ОШИБКА"
)

type Event struct {
	// Идентификатор ивента
	Uuid string `json:"uuid,omitempty"`
	// Описание (например, текст ошибки)
	Description string `json:"description,omitempty"`
	// Время последнего обновления статуса
	Timestamp string `json:"timestamp,omitempty"`
	// Статус ивента
	Status EventStatus `json:"status,omitempty"`
}

type EventGetReq struct {
	// Идентификатор ивента
	EventUuid string `json:"event_uuid"`
}
