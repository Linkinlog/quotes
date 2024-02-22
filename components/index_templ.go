// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.543
package components

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

func Index(contents templ.Component) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<!doctype html><html lang=\"en\"><head><meta charset=\"UTF-8\"><link rel=\"icon\" type=\"image/svg+xml\" href=\"/assets/wisePup.webp\"><meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\"><meta name=\"theme-color\" media=\"(prefers-color-scheme: dark)\" content=\"black\"><meta name=\"description\" content=\"Wise Pup Quotes\"><script src=\"https://unpkg.com/htmx.org@1.9.10\" integrity=\"sha384-D1Kt99CQMDuVetoL1lrYwg5t+9QdHe7NLX/SoJYkXDFfX37iInKRy5xLSi8nO7UC\" crossorigin=\"anonymous\"></script>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = ga().Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<title>Wise Pup Quotes</title></head><body><main id=\"app\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = Nav().Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<h1>Wise Pup</h1><br><img id=\"wisePup\" src=\"/assets/wisePup.webp\" alt=\"wise pup\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = contents.Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</main></body></html><style>\n        :root {\n          font-family: Inter, system-ui, Avenir, Helvetica, Arial, sans-serif;\n          line-height: 1.5;\n          font-weight: 400;\n\n          color-scheme: dark;\n          color: rgba(255, 255, 255, 0.87);\n          background-color: #242424;\n\n          font-synthesis: none;\n          text-rendering: optimizeLegibility;\n          -webkit-font-smoothing: antialiased;\n          -moz-osx-font-smoothing: grayscale;\n        }\n\n        nav ul {\n          display: flex;\n          justify-content: space-between;\n          list-style: none;\n          margin: 0;\n          margin-bottom: 2rem;\n          padding: 0;\n        }\n\n        a {\n          font-weight: 500;\n          color: #C05CFF\n          cursor: pointer;\n          text-decoration: inherit;\n        }\n\n        a:hover {\n          color: #535bf2;\n        }\n\n        body {\n          margin: 0;\n          display: flex;\n          place-items: center;\n        }\n\n        h1 {\n          line-height: 1.1;\n        }\n\n        #app {\n          width: 80vw;\n          height: 100%;\n          margin: 0 auto;\n          padding: 1rem;\n          text-align: center;\n        }\n\n        button {\n          border-radius: 8px;\n          border: 1px solid transparent;\n          padding: 0.6em 1.2em;\n          margin: 0.5rem;\n          font-size: 1em;\n          font-weight: 500;\n          font-family: inherit;\n          background-color: #1a1a1a;\n          cursor: pointer;\n          transition: border-color 0.25s;\n        }\n        button:hover {\n          border-color: #646cff;\n        }\n        button:focus,\n        button:focus-visible {\n          outline: 4px auto -webkit-focus-ring-color;\n        }\n\n        .quote {\n          font-size: 1.5em;\n          font-style: italic;\n          color: rgba(255, 255, 255, 0.87);\n          margin: 2em 0;\n          border: 1px solid rgba(255, 255, 255, 0.2);\n          border-radius: 8px;\n        }\n\n        #add-quote-form div {\n            margin: 7px;\n        }\n\n        alert-success {\n          display: block;\n          padding: 1em;\n          margin: 1em 0;\n          border-radius: 8px;\n          border-color: #00e676;\n          color: #00e676;\n        }\n\n        p {\n          font-size: 1rem;\n        }\n        footer {\n          font-size: 0.8em;\n          color: rgba(255, 255, 255, 0.6);\n        }\n        img#wisePup {\n          height: 300px;\n          width: 200px;\n        }\n        @media screen and (max-width: 767px) {\n            img#wisePup {\n                height: 157px;\n                width: 105px;\n            }\n        }\n    </style>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}
