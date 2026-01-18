# 04_API_IMPLEMENTATION_GUIDE.md

## Layering
- Controller: validate request -> call service -> envelope response
- Request: parse + basic validation
- Service: business rules (validate/exceptions/gating/events/audit)
- Repository: DB only
- Resource: response mapping

## Error codes
- Use docs/ERR_CODES.md

## center scoping
- never trust client center_id alone; verify admin_user.center_id == path center_id
