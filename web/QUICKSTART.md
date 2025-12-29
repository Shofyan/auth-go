# HTMX Web UI - Quick Start Guide

## 🚀 Getting Started

### 1. Start the Server

Using Go:
```bash
go run cmd/server/main.go
```

Using Docker Compose:
```bash
docker-compose up
```

The server will start on `http://localhost:8080`

### 2. Access the Web UI

Open your browser and navigate to:
- **Home**: http://localhost:8080/
- **Login**: http://localhost:8080/web/login
- **Register**: http://localhost:8080/web/register

## 📝 User Flow

### Registration
1. Go to http://localhost:8080/web/register
2. Enter your email and password (minimum 8 characters)
3. Click "Create Account"
4. You'll be redirected to the login page

### Login
1. Go to http://localhost:8080/web/login
2. Enter your registered email and password
3. Click "Sign In"
4. Upon successful login, you'll be redirected to the dashboard
5. Access and refresh tokens are automatically stored in browser localStorage

### Dashboard
- View your user information
- Access quick actions
- Navigate to profile page
- Logout functionality

### Profile
- View detailed user information
- User ID, email, and assigned roles
- All data loaded dynamically with HTMX

## 🔧 How It Works

### Token Management
- **Login**: Access token (15 min) and refresh token (7 days) are stored in localStorage
- **API Calls**: Access token is automatically added to Authorization header
- **Expired Token**: If access token expires, you'll be redirected to login
- **Refresh**: Use the refresh button on dashboard to get new tokens

### HTMX Features
- Forms submit via AJAX (no page reload)
- Profile data loads dynamically
- Real-time error/success messages
- Smooth transitions and loading states

## 🎨 Pages Overview

### Public Pages (No Authentication Required)
- `/` - Home (redirects to login)
- `/web/login` - Login form
- `/web/register` - Registration form

### Protected Pages (Authentication Required)
- `/web/dashboard` - Main dashboard
- `/web/profile` - User profile
- `/web/profile-data` - Profile data endpoint (used by HTMX)

## 🔐 Testing the API

### 1. Register a new user
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'
```

### 2. Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'
```

### 3. Access protected route
```bash
curl http://localhost:8080/api/v1/auth/profile \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

## 🎯 Features Demonstrated

✅ User registration with validation  
✅ User login with JWT tokens  
✅ Protected routes with authentication  
✅ Profile viewing  
✅ Token refresh mechanism  
✅ Logout functionality  
✅ HTMX dynamic updates  
✅ Error handling and messages  
✅ Responsive design  

## 🐛 Troubleshooting

### "Cannot connect to server"
- Ensure the server is running on port 8080
- Check if another service is using port 8080

### "Invalid credentials"
- Make sure you registered the account first
- Check email and password are correct

### "Unauthorized" on protected pages
- Clear localStorage and login again
- Check if access token has expired

### Templates not loading
- Ensure `web/templates/` directory exists
- Check all template files are present

## 📚 Next Steps

1. **Customize Styling**: Edit CSS in [layout.html](templates/layout.html)
2. **Add Features**: Create new templates and handlers
3. **Enhance Security**: Implement CSRF protection, httpOnly cookies
4. **Add Validation**: Client-side form validation with HTMX
5. **Error Pages**: Create 404, 500 error pages

## 🔗 Related Documentation

- [Main README](../README.md) - Full project documentation
- [Postman Collection](../postman/README.md) - API testing with Postman
- [HTMX Documentation](https://htmx.org/) - Learn more about HTMX
