package services

func MergeTags(userTags, generatedTags []string) []string {
    tagSet := make(map[string]struct{})

    // Add user-provided tags
    for _, tag := range userTags {
        tagSet[tag] = struct{}{}
    }

    // Add generated tags
    for _, tag := range generatedTags {
        tagSet[tag] = struct{}{}
    }

    // Convert map keys back to a slice
    mergedTags := make([]string, 0, len(tagSet))
    for tag := range tagSet {
        mergedTags = append(mergedTags, tag)
    }

    return mergedTags
}
