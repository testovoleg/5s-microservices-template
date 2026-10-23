package utils

// UuidPtr constrains T to a pointer type that exposes its own uuid via GetUuid.
type UuidPtr[T any] interface {
	*T
	GetUuid() string
}

// MapByUuid indexes a slice of uuid-identified pointers by their uuid, skipping nils.
func MapByUuid[T any, PT UuidPtr[T]](items []PT) map[string]PT {
	m := make(map[string]PT, len(items))
	for _, v := range items {
		if v != nil {
			m[v.GetUuid()] = v
		}
	}
	return m
}

// maxPageSize is the largest page size core_service accepts for a list request.
const maxPageSize = 100

// PaginationFor picks a page size large enough to return every match for count items in a
// single page (capped at maxPageSize), and how many pages are needed to cover all of them
// once count grows past that cap. Without this, list requests fall back to core_service's
// default page size, silently truncating results once there are more than a few ids.
func PaginationFor(count int) (size, pages int) {
	size = count + 1
	if size > maxPageSize {
		size = maxPageSize
	}
	pages = (count + maxPageSize - 1) / maxPageSize
	if pages < 1 {
		pages = 1
	}
	return size, pages
}
