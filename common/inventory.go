package common

import (
	"bytes"
	"encoding/json"
)

type InventoryRecord struct {
	ID     string            `json:"id"`
	Region string            `json:"region"`
	Images map[string]string `json:"images"`
	Tags   map[string]string `json:"tags,omitempty"`
}

func EncodeInventoryRecord(clientId, region string, images map[string]string) ([]byte, error) {
	inventoryRecord := InventoryRecord{
		ID:     clientId,
		Region: region,
		Images: images,
	}

	recordBytes, err := json.Marshal(inventoryRecord)
	if err != nil {
		return nil, err
	}

	return recordBytes, nil
}

func DecodeInventoryRecord(msgData interface{}) (*InventoryRecord, error) {
	inventoryRecord := &InventoryRecord{}
	err := json.NewDecoder(bytes.NewBuffer([]byte(msgData.(string)))).Decode(inventoryRecord)
	if err != nil {
		return nil, err
	}

	return inventoryRecord, nil
}
