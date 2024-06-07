package main

import (
    "html/template"
    "io"
    "log"
    "strconv"

    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
)

type Todo struct {
    Description string
    Id          int
    Done        bool
}

type Templates struct {
    templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
    return t.templates.ExecuteTemplate(w, name, data)
}

func newTemplate() *Templates {
    return &Templates{
        templates: template.Must(template.ParseGlob("templates/*.html")),
    }
}

func main() {
    e := echo.New()
    e.Use(middleware.Logger())
    e.Use(middleware.Recover())

    e.Renderer = newTemplate()

    todos := map[string][]Todo{
        "Todos": {
            {Description: "Learn Go", Id: 1, Done: false},
            {Description: "Build a web app", Id: 2, Done: false},
            {Description: "Profit", Id: 3, Done: false},
        },
    }
    e.GET("/", func(c echo.Context) error {
        return c.Render(200, "index", todos)
    })

    e.POST("/todos", func(c echo.Context) error {
        description := c.FormValue("description")
        todos["Todos"] = append(todos["Todos"], Todo{Description: description, Id: len(todos["Todos"]) + 1, Done: false})
        return c.Render(200, "todo", todos["Todos"][len(todos["Todos"])-1])
    })

    e.POST("/todos/:id", func(c echo.Context) error {
        idStr := c.Param("id")
        id, _ := strconv.Atoi(idStr)
        for i, todo := range todos["Todos"] {
            if id == todo.Id {
                log.Print("found")
                todos["Todos"][i].Done = !todos["Todos"][i].Done
            }
        }
        return c.Render(200, "index", todos)
    })

    e.Logger.Fatal(e.Start(":8080"))
}
