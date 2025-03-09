CREATE OR REPLACE FUNCTION create_project(
    p_title TEXT,
    p_description TEXT,
    p_introduction TEXT,
    p_owner_id INT,
    p_tags TEXT[],
    p_sections JSONB
) RETURNS INT AS $$
DECLARE
    new_project_id INT;
BEGIN
    -- Insert the main project and get the project_id
    INSERT INTO projects (title, description, introduction, owner_id)
    VALUES (p_title, p_description, p_introduction, p_owner_id)
    RETURNING project_id INTO new_project_id;

    -- Insert tags associated with the project
    INSERT INTO project_tags (project_id, tag_description)
    SELECT new_project_id, unnest(p_tags);

    -- Insert sections associated with the project using JSONB
    INSERT INTO project_sections (project_id, title, body, section_number)
    SELECT
        new_project_id,
        section->>'title',
        section->>'body',
        (section->>'section_number')::int
    FROM jsonb_array_elements(p_sections) AS section;

    -- Return the new project_id
    RETURN new_project_id;
END;
$$ LANGUAGE plpgsql;