# clockify-mcp

An [MCP](https://modelcontextprotocol.io) server that exposes the [Clockify](https://clockify.me) API as tools for AI assistants.

Written in Go, it speaks the Streamable HTTP transport and covers workspaces, time entries, projects, tasks and tags.

## Running

The server needs two environment variables:

| Variable | Required | Description |
| --- | --- | --- |
| `CLOCKIFY_API_KEY` | yes | Your Clockify API key (Clockify → Profile settings → API) |
| `MCP_TOKEN` | yes | Bearer token clients must present to reach the server |
| `PORT` | no | Listen port, defaults to `8080` |

### Docker

```sh
docker run --rm -p 8080:8080 \
  -e CLOCKIFY_API_KEY=your-clockify-key \
  -e MCP_TOKEN=your-secret-token \
  ghcr.io/mdallastella/clockify-mcp:latest
```

### From source

```sh
go build -o clockify-mcp ./main.go
CLOCKIFY_API_KEY=... MCP_TOKEN=... ./clockify-mcp
```

The MCP endpoint is served at `/mcp`.

## Connecting a client

Every request must carry `Authorization: Bearer <MCP_TOKEN>`; anything else gets a 401.

```json
{
  "mcp": {
    "clockify": {
      "type": "remote",
      "url": "http://127.0.0.1:8080/mcp",
      "enabled": true,
      "headers": {
        "Authorization": "Bearer your-secret-token"
      }
    }
  }
}
```

## Tools

**Workspaces** — `list_workspaces`, `get_workspace`

**Time entries** — `create_time_entry`, `get_time_entry`, `update_time_entry`, `delete_time_entry`, `list_time_entries`, `get_in_progress_time_entries`, `stop_timer`, `create_time_entry_for_user`, `duplicate_time_entry`, `delete_user_time_entries`, `mark_time_entries_invoiced`

**Projects** — `list_projects`, `create_project`, `get_project`, `update_project`, `delete_project`, `update_project_estimate`, `add_users_to_project`, `create_project_from_template`, `update_project_template_flag`, `set_project_user_hourly_rate`, `set_project_user_cost_rate`

**Tasks** — `list_tasks`, `create_task`, `get_task`, `update_task`, `delete_task`, `set_task_hourly_rate`, `set_task_cost_rate`

**Tags** — `list_tags`, `create_tag`, `get_tag`, `update_tag`, `delete_tag`

Omit `end` on `create_time_entry` to start a running timer. Note that `delete_user_time_entries` wipes every time entry for a user and cannot be undone.

## Notes

The bearer token is the only thing standing between the internet and your Clockify account — use a long random value and put the server behind TLS if you expose it beyond localhost.

## License

[MIT](LICENSE)
