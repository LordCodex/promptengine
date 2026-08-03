package manifest

import "strings"

func upsertTech(dst, src []TechnologyDefinition) []TechnologyDefinition {
	index := map[string]int{}
	for i, item := range dst {
		index[strings.ToLower(firstNonEmpty(item.ID, item.Framework, item.Language, item.Stack))] = i
	}
	for _, item := range src {
		key := strings.ToLower(firstNonEmpty(item.ID, item.Framework, item.Language, item.Stack))
		if i, ok := index[key]; ok {
			dst[i] = item
		} else {
			index[key] = len(dst)
			dst = append(dst, item)
		}
	}
	return dst
}

func upsertPlaybooks(dst, src []PlaybookDefinition) []PlaybookDefinition {
	index := map[string]int{}
	for i, item := range dst {
		index[item.ID] = i
	}
	for _, item := range src {
		if i, ok := index[item.ID]; ok {
			dst[i] = item
		} else {
			index[item.ID] = len(dst)
			dst = append(dst, item)
		}
	}
	return dst
}

func upsertWorkflows(dst, src []WorkflowDefinition) []WorkflowDefinition {
	index := map[string]int{}
	for i, item := range dst {
		index[item.ID] = i
	}
	for _, item := range src {
		if i, ok := index[item.ID]; ok {
			dst[i] = item
		} else {
			index[item.ID] = len(dst)
			dst = append(dst, item)
		}
	}
	return dst
}

func upsertPrompts(dst, src []PromptMapping) []PromptMapping {
	index := map[string]int{}
	for i, item := range dst {
		index[item.TaskType] = i
	}
	for _, item := range src {
		if i, ok := index[item.TaskType]; ok {
			dst[i] = item
		} else {
			index[item.TaskType] = len(dst)
			dst = append(dst, item)
		}
	}
	return dst
}

func upsertTemplates(dst, src []TemplateDefinition) []TemplateDefinition {
	index := map[string]int{}
	for i, item := range dst {
		index[strings.ToLower(item.Name)] = i
	}
	for _, item := range src {
		key := strings.ToLower(item.Name)
		if i, ok := index[key]; ok {
			dst[i] = item
		} else {
			index[key] = len(dst)
			dst = append(dst, item)
		}
	}
	return dst
}
