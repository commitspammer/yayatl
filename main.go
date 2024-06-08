package main

import (
    "html/template"
    "io"
    "strconv"
    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
)

type Templates struct {
    templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
    t.templates = template.Must(template.ParseGlob("templates/*.html")) // COMMENT THIS LINE IN PROD!
    return t.templates.ExecuteTemplate(w, name, data)
}

func newTemplate() *Templates {
    return &Templates{
        templates: template.Must(template.ParseGlob("templates/*.html")),
    }
}

type Todo struct {
    Description string
    Id          int
    Done        bool
}

func main() {
    e := echo.New()
    e.Use(middleware.Logger())
    e.Use(middleware.Recover())
    e.Renderer = newTemplate()
    e.File("/favicon.ico", "assets/favicon.ico")

    next_id := 0
    state := map[string][]Todo{
        "Todos": {
            { Description: "Learn HTMX",     Id: 1, Done: false },
            { Description: "Build a Go app", Id: 2, Done: false },
            { Description: "Profit",         Id: 3, Done: false },
        },
    }

    e.GET("/", func(c echo.Context) error {
        return c.Render(200, "index", state)
    })

    e.POST("/todos", func(c echo.Context) error {
        description := c.FormValue("description")
        todo := Todo{
            Description: description,
            Id: next_id,
            Done: false,
        }
        state["Todos"] = append(state["Todos"], todo)
        next_id += 1
        return c.Render(200, "todo", todo)
    })

    e.PUT("/todos/:id/done", func(c echo.Context) error {
        idStr := c.Param("id")
        id, _ := strconv.Atoi(idStr)
        for i, todo := range state["Todos"] {
            if id == todo.Id {
                state["Todos"][i].Done = !state["Todos"][i].Done
                return c.Render(200, "todo", state["Todos"][i])
            }
        }
        return c.Render(404, "error", state) //error.html no exist
    })

    e.PUT("/todos/:id/description", func(c echo.Context) error {
        idStr := c.Param("id")
        id, _ := strconv.Atoi(idStr)
        description := c.FormValue("description")
        for i, todo := range state["Todos"] {
            if id == todo.Id {
                state["Todos"][i].Description = description
                return c.Render(200, "todo", state["Todos"][i])
            }
        }
        return c.Render(404, "error", state) //error.html no exist
    })

    e.DELETE("/todos/:id", func(c echo.Context) error {
        idStr := c.Param("id")
        id, _ := strconv.Atoi(idStr)
        for i, todo := range state["Todos"] {
            if id == todo.Id {
                state["Todos"] = append(state["Todos"][:i], state["Todos"][i+1:]...)
                break
            }
        }
        return c.String(200, "")
    })

    e.GET("/todos/:id/edit", func(c echo.Context) error {
        idStr := c.Param("id")
        id, _ := strconv.Atoi(idStr)
        for _, todo := range state["Todos"] {
            if id == todo.Id {
                return c.Render(200, "edit-todo", todo)
            }
        }
        return c.Render(404, "error", state) //error.html no exist
    })

    e.Logger.Fatal(e.Start(":8080"))
}
