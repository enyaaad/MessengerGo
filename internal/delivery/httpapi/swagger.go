package httpapi

import (
	"net/http"
)

const openapiYAML = `openapi: 3.0.3
info:
  title: Chats API
  version: 1.0.0
servers:
  - url: http://localhost:8080

paths:
  /chats:
    post:
      summary: Create chat
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/CreateChatRequest'
      responses:
        '201':
          description: Created chat
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Chat'
        '400':
          $ref: '#/components/responses/BadRequest'

  /chats/{id}:
    get:
      summary: Get chat with last N messages
      parameters:
        - $ref: '#/components/parameters/ChatID'
        - name: limit
          in: query
          required: false
          schema:
            type: integer
            default: 20
            minimum: 1
            maximum: 100
      responses:
        '200':
          description: Chat and messages
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/GetChatResponse'
        '404':
          $ref: '#/components/responses/NotFound'
        '400':
          $ref: '#/components/responses/BadRequest'
    delete:
      summary: Delete chat (cascade delete messages)
      parameters:
        - $ref: '#/components/parameters/ChatID'
      responses:
        '204':
          description: No Content
        '404':
          $ref: '#/components/responses/NotFound'

  /chats/{id}/messages:
    post:
      summary: Create message in chat
      parameters:
        - $ref: '#/components/parameters/ChatID'
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/CreateMessageRequest'
      responses:
        '201':
          description: Created message
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Message'
        '404':
          $ref: '#/components/responses/NotFound'
        '400':
          $ref: '#/components/responses/BadRequest'

components:
  parameters:
    ChatID:
      name: id
      in: path
      required: true
      schema:
        type: integer
        minimum: 1

  responses:
    BadRequest:
      description: Bad Request
      content:
        application/json:
          schema:
            $ref: '#/components/schemas/ErrorResponse'
    NotFound:
      description: Not Found
      content:
        application/json:
          schema:
            $ref: '#/components/schemas/ErrorResponse'

  schemas:
    ErrorResponse:
      type: object
      required: [error]
      properties:
        error:
          type: string

    CreateChatRequest:
      type: object
      required: [title]
      properties:
        title:
          type: string
          minLength: 1
          maxLength: 200

    CreateMessageRequest:
      type: object
      required: [text]
      properties:
        text:
          type: string
          minLength: 1
          maxLength: 5000

    Chat:
      type: object
      required: [id, title, created_at]
      properties:
        id:
          type: integer
        title:
          type: string
        created_at:
          type: string
          format: date-time

    Message:
      type: object
      required: [id, chat_id, text, created_at]
      properties:
        id:
          type: integer
        chat_id:
          type: integer
        text:
          type: string
        created_at:
          type: string
          format: date-time

    GetChatResponse:
      type: object
      required: [chat, messages]
      properties:
        chat:
          $ref: '#/components/schemas/Chat'
        messages:
          type: array
          items:
            $ref: '#/components/schemas/Message'
`

func (h *Handlers) OpenAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/yaml; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(openapiYAML))
}

func (h *Handlers) SwaggerUI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	_, _ = w.Write([]byte(`<!doctype html>
<html>
  <head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <title>messengerTest — Swagger</title>
    <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css" />
  </head>
  <body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
    <script>
      window.onload = () => {
        SwaggerUIBundle({
          url: '/openapi.yaml',
          dom_id: '#swagger-ui',
          deepLinking: true,
          presets: [SwaggerUIBundle.presets.apis]
        });
      };
    </script>
  </body>
</html>`))
}
