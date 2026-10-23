package utils

import "encoding/json"

// ChunkByJSONSize greedily groups items into chunks whose combined JSON-encoded size
// stays under maxBytes (Kafka message.max.bytes, HTTP body limits, etc.).
// maxBytes <= 0 disables splitting. An item larger than maxBytes is still emitted
// as its own chunk.
func ChunkByJSONSize[T any](items []T, maxBytes int) [][]T {
	if len(items) == 0 {
		return [][]T{items}
	}
	if maxBytes <= 0 {
		return [][]T{items}
	}

	const jsonOverheadPerItem = 2 // ',' separator + margin

	chunks := make([][]T, 0, 1)
	current := make([]T, 0)
	currentSize := 0
	for _, item := range items {
		itemBytes, err := json.Marshal(item)
		itemSize := len(itemBytes) + jsonOverheadPerItem
		if err != nil {
			itemSize = maxBytes
		}
		if len(current) > 0 && currentSize+itemSize > maxBytes {
			chunks = append(chunks, current)
			current = make([]T, 0)
			currentSize = 0
		}
		current = append(current, item)
		currentSize += itemSize
	}
	if len(current) > 0 {
		chunks = append(chunks, current)
	}
	return chunks
}
