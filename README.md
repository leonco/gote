gote
====

gote is a simple but powerful text templating library for go.  The syntax is similar to Google's ctemplate library, and emphasizes the separation of logic from presentation.

Example
==========
```
t, err := gote.Parse([]byte("Hello {{WORLD:h}}"))
if err != nil {
	panic(err)
}

dict := gote.NewTemplateDictionary()
dict.Put("WORLD", "<h1>world</h1>")
t.Render(dict, os.Stdout)
```

Features
==========
* Enforces a strict separation of "view" from application logic.
* Based loosely on the syntax and behavior of Google's ctemplate library.
* Templates do not require a compilation step. Templates are parsed at runtime.
* Native support for some types of content escaping (JavaScript, XML, HTML, URLs).
