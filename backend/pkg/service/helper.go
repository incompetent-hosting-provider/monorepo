package service

import db_presets "incompetent-hosting-provider/backend/pkg/db/tables/presets"

func serializePresetListResponse(presets []db_presets.PresetTable) PresetListResponse {
	var serializedPresets []preset
	for _, p := range presets {
		serializedPresets = append(serializedPresets, preset{
			Name:        p.Name,
			PresetId:    p.PresetId,
			Description: p.Description,
		})
	}
	return PresetListResponse{
		Presets: serializedPresets,
	}
}
