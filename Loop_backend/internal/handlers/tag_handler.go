package handlers

import (
    "encoding/json"
    "net/http"
    "github.com/gorilla/mux"
    "strconv"
    "Loop_backend/internal/models"
    "Loop_backend/internal/services/tags"
)

type TagHandler struct {
    tagService services.TagService
}

func NewTagHandler(tagService services.TagService) *TagHandler {
    return &TagHandler{
        tagService: tagService,
    }
}

func (h *TagHandler) RegisterRoutes(router *mux.Router) {
    router.HandleFunc("/api/tags", h.CreateTag).Methods("POST")
    router.HandleFunc("/api/tags/{id:[0-9]+}", h.GetTagByID).Methods("GET")
    router.HandleFunc("/api/tags/name/{name}", h.GetTagByName).Methods("GET")
    router.HandleFunc("/api/tags/{id:[0-9]+}", h.UpdateTag).Methods("PUT")
    router.HandleFunc("/api/tags/relationship", h.CreateTagRelationship).Methods("POST")
    router.HandleFunc("/api/tags/{name}/related", h.GetRelatedTags).Methods("GET")
    router.HandleFunc("/api/tags/project/{projectID}", h.GetProjectTags).Methods("GET")
    router.HandleFunc("/api/tags/expertise", h.SetUserExpertise).Methods("POST")
    router.HandleFunc("/api/tags/expertise/{userID}", h.GetUserExpertise).Methods("GET")
    router.HandleFunc("/api/tags/{name}/experts", h.GetTagExperts).Methods("GET")
}

func (h *TagHandler) CreateTag(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Name     string    `json:"name"`
        Category string    `json:"category"`
        Vector   []float64 `json:"vector,omitempty"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    tag := &models.Tag{
        Name:     req.Name,
        Category: req.Category,
        Vector:   req.Vector,
    }

    if err := h.tagService.CreateTag(tag); err != nil {
        http.Error(w, "Failed to create tag", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(tag)
}

func (h *TagHandler) GetTagByID(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, ok := vars["id"]
    if !ok {
        http.Error(w, "Tag ID is required", http.StatusBadRequest)
        return
    }

    tag, err := h.tagService.GetTagByID(id)
    if err != nil {
        http.Error(w, "Tag not found", http.StatusNotFound)
        return
    }

    json.NewEncoder(w).Encode(tag)
}

func (h *TagHandler) GetTagByName(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    name, ok := vars["name"]
    if !ok {
        http.Error(w, "Tag name is required", http.StatusBadRequest)
        return
    }

    tag, err := h.tagService.GetTagByName(name)
    if err != nil {
        http.Error(w, "Tag not found", http.StatusNotFound)
        return
    }

    json.NewEncoder(w).Encode(tag)
}

func (h *TagHandler) UpdateTag(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, ok := vars["id"]
    if !ok {
        http.Error(w, "Tag ID is required", http.StatusBadRequest)
        return
    }

    var req struct {
        Name     string    `json:"name"`
        Category string    `json:"category"`
        Vector   []float64 `json:"vector,omitempty"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    tag := &models.Tag{
        ID:       id,
        Name:     req.Name,
        Category: req.Category,
        Vector:   req.Vector,
    }

    if err := h.tagService.UpdateTag(tag); err != nil {
        http.Error(w, "Failed to update tag", http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(tag)
}

func (h *TagHandler) CreateTagRelationship(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Tag1     string  `json:"tag1"`
        Tag2     string  `json:"tag2"`
        Strength float64 `json:"strength"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    if err := h.tagService.CreateTagRelationship(req.Tag1, req.Tag2, req.Strength); err != nil {
        http.Error(w, "Failed to create tag relationship", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
}

func (h *TagHandler) GetRelatedTags(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    name, ok := vars["name"]
    if !ok {
        http.Error(w, "Tag name is required", http.StatusBadRequest)
        return
    }

    minStrength := 0.5 // Default minimum strength
    if str := r.URL.Query().Get("min_strength"); str != "" {
        if s, err := strconv.ParseFloat(str, 64); err == nil {
            minStrength = s
        }
    }

    tags, err := h.tagService.GetRelatedTags(name, minStrength)
    if err != nil {
        http.Error(w, "Failed to get related tags", http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(tags)
}

func (h *TagHandler) GetProjectTags(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    projectID, ok := vars["projectID"]
    if !ok {
        http.Error(w, "Project ID is required", http.StatusBadRequest)
        return
    }

    tags, err := h.tagService.GetTagsByProject(projectID)
    if err != nil {
        http.Error(w, "Failed to get project tags", http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(tags)
}

func (h *TagHandler) SetUserExpertise(w http.ResponseWriter, r *http.Request) {
    var req struct {
        UserID  string `json:"user_id"`
        TagName string `json:"tag_name"`
        Level   string `json:"level"`
        Years   int    `json:"years"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    if err := h.tagService.SetUserTagExpertise(req.UserID, req.TagName, req.Level, req.Years); err != nil {
        http.Error(w, "Failed to set user expertise", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
}

func (h *TagHandler) GetUserExpertise(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    userID, ok := vars["userID"]
    if !ok {
        http.Error(w, "User ID is required", http.StatusBadRequest)
        return
    }

    expertise, err := h.tagService.GetUserExpertise(userID)
    if err != nil {
        http.Error(w, "Failed to get user expertise", http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(expertise)
}

func (h *TagHandler) GetTagExperts(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    name, ok := vars["name"]
    if !ok {
        http.Error(w, "Tag name is required", http.StatusBadRequest)
        return
    }

    experts, err := h.tagService.GetTagExperts(name)
    if err != nil {
        http.Error(w, "Failed to get tag experts", http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(experts)
}
