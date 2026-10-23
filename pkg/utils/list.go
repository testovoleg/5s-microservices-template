package utils

// ListUuids collects unique non-empty uuids from a slice of pointers.
func ListUuids[T any](items []*T, uuid func(*T) string) []string {
	uuids := []string{}
	for _, item := range items {
		if item != nil {
			if id := uuid(item); id != "" {
				uuids = AppendUnique(uuids, id)
			}
		}
	}
	return uuids
}

// EnrichList replaces items in list with matching entries from incoming (by uuid).
func EnrichList[T any](list []*T, incoming []*T, uuid func(*T) string) {
	if len(incoming) == 0 {
		return
	}
	m := make(map[string]*T, len(incoming))
	for _, item := range incoming {
		if item != nil {
			m[uuid(item)] = item
		}
	}
	for i, obj := range list {
		if obj == nil {
			continue
		}
		if item, ok := m[uuid(obj)]; ok {
			list[i] = item
		}
	}
}

// IncompleteUuids returns uuids of items that satisfy the incomplete predicate.
func IncompleteUuids[T any](items []*T, uuid func(*T) string, incomplete func(*T) bool) []string {
	uuids := []string{}
	for _, item := range items {
		if item == nil {
			continue
		}
		if id := uuid(item); id != "" && incomplete(item) {
			uuids = AppendUnique(uuids, id)
		}
	}
	return uuids
}
