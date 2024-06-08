package main

import (
    "html/template"
    "io"
    "log"
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
        state["Todos"] = append(state["Todos"], Todo{
            Description: description,
            Id: len(state["Todos"]) + 1,
            Done: false,
        })
        return c.Render(200, "todo", state["Todos"][len(state["Todos"])-1])
    })

    e.PUT("/todos/:id/done", func(c echo.Context) error {
        idStr := c.Param("id")
        id, _ := strconv.Atoi(idStr)
        for i, todo := range state["Todos"] {
            if id == todo.Id {
                log.Print("found")
                state["Todos"][i].Done = !state["Todos"][i].Done
            }
        }
        return c.Render(200, "index", state)
    })

    e.PUT("/todos/:id/description", func(c echo.Context) error {
        description := c.FormValue("description")
        idStr := c.Param("id")
        id, _ := strconv.Atoi(idStr)
        for i, todo := range state["Todos"] {
            if id == todo.Id {
                log.Print("found")
                state["Todos"][i].Description = description
            }
        }
        return c.Render(200, "index", state)
    })

    e.DELETE("/todos/:id", func(c echo.Context) error {
        idStr := c.Param("id")
        id, _ := strconv.Atoi(idStr)
        for i, todo := range state["Todos"] {
            if id == todo.Id {
                // Remove the todo from the list
                state["Todos"] = append(state["Todos"][:i], state["Todos"][i+1:]...)
                break
            }
        }
        return c.String(200, "")
    })

    e.Logger.Fatal(e.Start(":8080"))
}
