# Code Snippet
A web application called Code Snippet, which lets people paste and share snippets of text — a bit like Pastebin or
GitHub’s Gists. This is the final code from version 1 of the [Let's Go e-book by Alex Edwards](https://lets-go.alexedwards.net/). It's a full-stack Go web application called "Code Snippet" that lets users CRUD text snippets (similar to GitHub gists).

### Features

- Authentication. Users can register and sign in.
- Protected endpoints. Only signed-in users can create snippets.
- RESTful routing.
- Middleware.
- MySQL database.
- SSL/TLS web server using HTTP 2.0.
- Generated HTML via Golang templates.
- CRSF protection.

### Development

##### `go run cmd/web/*`

Starts the local web server with HTTPS on port 4000 ([https://localhost:4000](https://localhost:4000))
