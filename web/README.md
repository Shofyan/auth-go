# Web UI with HTMX

This directory contains the HTMX-based web interface for the authentication service.

## Features

- **Login Page** (`/web/login`) - User authentication with email and password
- **Registration Page** (`/web/register`) - New user registration
- **Dashboard** (`/web/dashboard`) - Main user dashboard after login
- **Profile Page** (`/web/profile`) - View user profile information

## Technology Stack

- **HTMX** - Modern HTML-over-the-wire framework for dynamic interactions
- **Go Templates** - Server-side rendering
- **Vanilla JavaScript** - Token management and authentication flow

## How It Works

### Authentication Flow

1. **Login/Register**: Users enter credentials on the web form
2. **Token Storage**: Access and refresh tokens are stored in browser's localStorage
3. **Protected Pages**: Dashboard and profile pages check for valid tokens
4. **API Calls**: HTMX makes requests to REST API endpoints with Bearer token
5. **Auto-redirect**: 401 responses automatically redirect to login

### HTMX Features Used

- **hx-post**: Send POST requests to API endpoints
- **hx-get**: Fetch data dynamically (profile information)
- **hx-target**: Specify where to insert response
- **hx-swap**: Control how content is swapped
- **hx-trigger**: Define when requests should fire
- **hx-on::after-request**: Handle response events

### Template Structure

```
web/templates/
├── layout.html      # Base layout with styles and HTMX script
├── login.html       # Login form
├── register.html    # Registration form
├── dashboard.html   # Main dashboard
└── profile.html     # Profile page
```

## Usage

1. Start the server:
```bash
go run cmd/server/main.go
```

2. Open browser and navigate to:
   - Home: `http://localhost:8080/`
   - Login: `http://localhost:8080/web/login`
   - Register: `http://localhost:8080/web/register`

3. Register a new account or login with existing credentials

4. Access protected pages:
   - Dashboard: `http://localhost:8080/web/dashboard`
   - Profile: `http://localhost:8080/web/profile`

## API Integration

The web UI integrates with the following REST API endpoints:

### Public Endpoints
- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/refresh` - Refresh access token

### Protected Endpoints (require Authorization header)
- `GET /api/v1/auth/profile` - Get user profile
- `POST /api/v1/auth/logout` - User logout
- `GET /web/profile-data` - Get profile data as HTML fragment

## Styling

The UI uses embedded CSS with:
- Clean, modern gradient design
- Responsive layout
- Form validation
- Loading states
- Error/success messages

## Security

- Tokens stored in localStorage (consider httpOnly cookies for production)
- Authorization header sent with all protected requests
- Automatic logout on 401 responses
- Client-side token validation

## Extending the UI

To add new pages:

1. Create a new template in `web/templates/`
2. Add a handler method in `internal/interface/http/handler/web_handler.go`
3. Register the route in `internal/interface/http/router.go`

Example:
```go
// In web_handler.go
func (h *WebHandler) ServeNewPage(w http.ResponseWriter, r *http.Request) {
    data := map[string]interface{}{
        "Title": "New Page",
    }
    h.templates.ExecuteTemplate(w, "layout.html", data)
}

// In router.go
mux.HandleFunc("/web/newpage", rt.webHandler.ServeNewPage)
```
