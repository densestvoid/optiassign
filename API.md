# OptiAssign API Documentation

## Overview

OptiAssign provides a RESTful API for group-based, randomized, prioritized, snaking-draft item assignment.

## Authentication

### Google OAuth 2.0
- **Endpoint**: `/auth/google`
- **Method**: GET
- **Description**: Initiates Google OAuth flow
- **Response**: Redirects to Google OAuth consent screen

### OAuth Callback
- **Endpoint**: `/auth/google/callback`
- **Method**: GET
- **Description**: Handles OAuth callback
- **Response**: Redirects to dashboard on success

### Logout
- **Endpoint**: `/auth/logout`
- **Method**: POST
- **Description**: Logs out the current user
- **Response**: Redirects to login page

## Groups

### List Groups
- **Endpoint**: `/groups`
- **Method**: GET
- **Auth**: Required
- **Description**: Returns list of user's groups
- **Response**: HTML page with groups list

### Create Group
- **Endpoint**: `/groups`
- **Method**: POST
- **Auth**: Required
- **Description**: Creates a new group
- **Body**: Form data with group details
- **Response**: HTMX partial HTML

### View Group
- **Endpoint**: `/groups/{id}`
- **Method**: GET
- **Auth**: Required
- **Description**: Shows group details and management interface
- **Response**: HTML page with group details

### Execute Assignment
- **Endpoint**: `/groups/{id}/execute`
- **Method**: POST
- **Auth**: Required
- **Description**: Manually triggers assignment execution
- **Response**: HTMX partial HTML with results

### Assignment Results
- **Endpoint**: `/groups/{id}/results`
- **Method**: GET
- **Auth**: Required
- **Description**: Shows complete assignment results
- **Response**: HTML page with results

## Participants

### Participant Access
- **Endpoint**: `/participant/{token}`
- **Method**: GET
- **Auth**: None (token-based)
- **Description**: Shows participant access page
- **Response**: HTML page with group and item information

### Priority Form
- **Endpoint**: `/participant/{token}/priorities`
- **Method**: GET
- **Auth**: None (token-based)
- **Description**: Shows priority submission form
- **Response**: HTMX partial HTML with form

### Submit Priorities
- **Endpoint**: `/participant/{token}/priorities`
- **Method**: POST
- **Auth**: None (token-based)
- **Description**: Submits participant priorities
- **Body**: Form data with item rankings
- **Response**: HTMX partial HTML with success/error

### Assignment Status
- **Endpoint**: `/participant/{token}/status`
- **Method**: GET
- **Auth**: None (token-based)
- **Description**: Returns assignment status
- **Response**: JSON with status information

### Participant Results
- **Endpoint**: `/participant/{token}/results`
- **Method**: GET
- **Auth**: None (token-based)
- **Description**: Shows participant's assigned items
- **Response**: HTML page with results

## Health Checks

### Health Check
- **Endpoint**: `/health`
- **Method**: GET
- **Auth**: None
- **Description**: Returns application health status
- **Response**: JSON with health information

### Readiness Check
- **Endpoint**: `/ready`
- **Method**: GET
- **Auth**: None
- **Description**: Returns application readiness status
- **Response**: JSON with readiness information

### Liveness Check
- **Endpoint**: `/live`
- **Method**: GET
- **Auth**: None
- **Description**: Returns application liveness status
- **Response**: JSON with liveness information

## Error Handling

### HTTP Status Codes
- **200**: Success
- **400**: Bad Request (validation errors)
- **401**: Unauthorized (authentication required)
- **403**: Forbidden (insufficient permissions)
- **404**: Not Found
- **429**: Too Many Requests (rate limit exceeded)
- **500**: Internal Server Error

### Error Response Format
```json
{
  "error": "Error message",
  "details": "Additional error details",
  "field": "Field name (for validation errors)"
}
```

## Rate Limiting

- **Limit**: 100 requests per minute per IP (production)
- **Headers**: 
  - `X-RateLimit-Limit`: Request limit
  - `X-RateLimit-Remaining`: Remaining requests
  - `X-RateLimit-Reset`: Reset time

## Security

### CSRF Protection
- All POST requests require CSRF token
- Token included in forms and headers
- Validates token on server-side

### Input Validation
- All inputs validated server-side
- SQL injection prevention with prepared statements
- XSS protection with proper escaping

### Session Security
- Secure session cookies
- Session timeout configuration
- Session invalidation on logout

## Examples

### Creating a Group
```bash
curl -X POST http://localhost:8080/groups \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "name=My Group&rule=equal"
```

### Submitting Priorities
```bash
curl -X POST http://localhost:8080/participant/abc123/priorities \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "priority_1=1&priority_2=2&priority_3=3"
```

### Checking Assignment Status
```bash
curl http://localhost:8080/participant/abc123/status
```

## Webhooks (Future Enhancement)

### Assignment Complete
- **Event**: `assignment.complete`
- **Payload**: Assignment results and participant information
- **Endpoint**: Configurable webhook URL

### Participant Joined
- **Event**: `participant.joined`
- **Payload**: Participant and group information
- **Endpoint**: Configurable webhook URL