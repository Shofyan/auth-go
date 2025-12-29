# Postman Collection

## 📦 Import Instructions

1. Open Postman
2. Click "Import" button
3. Select `Auth-Service.postman_collection.json`
4. The collection will be imported with all endpoints

## 🔧 Setup Environment Variables

The collection uses these environment variables (automatically set by test scripts):

| Variable | Description | Auto-set |
|----------|-------------|----------|
| `base_url` | API base URL | ✅ (default: http://localhost:8080) |
| `access_token` | JWT access token | ✅ (from login) |
| `refresh_token` | Refresh token | ✅ (from login) |
| `user_email` | User email | ✅ (from register) |
| `admin_access_token` | Admin access token | ✅ (from admin login) |
| `admin_refresh_token` | Admin refresh token | ✅ (from admin login) |

## 🚀 Quick Start

### 1. Health Check
Test if the service is running:
```
GET /health
```

### 2. Register User
Create a new user account:
```json
POST /api/v1/auth/register
{
  "email": "user@example.com",
  "password": "SecurePass123!"
}
```

### 3. Login
Login and automatically save tokens:
```json
POST /api/v1/auth/login
{
  "email": "user@example.com",
  "password": "SecurePass123!"
}
```
✅ Access token and refresh token are automatically saved to environment variables

### 4. Get Profile
Get user profile (uses saved access token):
```
GET /api/v1/auth/profile
Authorization: Bearer {{access_token}}
```

### 5. Refresh Token
Get new tokens (token rotation):
```json
POST /api/v1/auth/refresh
{
  "refresh_token": "{{refresh_token}}"
}
```
✅ New tokens are automatically saved

### 6. Logout
Revoke all tokens:
```
POST /api/v1/auth/logout
Authorization: Bearer {{access_token}}
```

## 🔐 Testing RBAC

### Create Admin User

1. **Register admin user**:
```json
POST /api/v1/auth/register
{
  "email": "admin@example.com",
  "password": "AdminPass123!"
}
```

2. **Update user role in database** (manual step):
```sql
UPDATE users 
SET roles = ARRAY['admin'] 
WHERE email = 'admin@example.com';
```

3. **Login as admin**:
```json
POST /api/v1/auth/login
{
  "email": "admin@example.com",
  "password": "AdminPass123!"
}
```

4. **Access admin endpoint**:
```
GET /api/v1/admin/users
Authorization: Bearer {{admin_access_token}}
```

## 📝 Test Scenarios

### Scenario 1: Complete Authentication Flow
1. Register → Login → Get Profile → Logout
2. All steps should succeed

### Scenario 2: Token Expiration
1. Login (get tokens)
2. Wait 15+ minutes (access token expires)
3. Try Get Profile → Should fail with 401
4. Use Refresh Token → Get new access token
5. Try Get Profile again → Should succeed

### Scenario 3: Token Rotation
1. Login (get refresh_token_1)
2. Refresh Token (get refresh_token_2)
3. Refresh Token again (get refresh_token_3)
4. Try to use refresh_token_1 → Should fail (revoked)

### Scenario 4: Token Reuse Detection (Security)
1. Login (get refresh_token_1)
2. Refresh Token (get refresh_token_2, token_1 revoked)
3. Try to use refresh_token_1 → Security breach detected
4. All tokens in family should be revoked
5. Must login again

### Scenario 5: RBAC Testing
1. Login as regular user
2. Try Admin endpoint → Should fail with 403 Forbidden
3. Login as admin user
4. Try Admin endpoint → Should succeed with 200

## 🛠️ Tips

- **Auto-save tokens**: Login and Refresh requests automatically save tokens to environment
- **Bearer token**: Protected endpoints automatically use `{{access_token}}`
- **Test scripts**: Check the "Tests" tab in each request to see automation
- **Environment**: Create a new environment to save variables persistently

## 🐛 Troubleshooting

### 401 Unauthorized
- Access token may be expired → Use Refresh Token
- Token not saved → Re-run Login request
- Check Authorization header format: `Bearer <token>`

### 403 Forbidden
- Insufficient permissions (RBAC)
- Update user role in database
- Re-login after role change

### 500 Internal Server Error
- Check if database is running
- Check if migrations are applied
- Check server logs

## 📊 Expected Response Codes

| Endpoint | Success Code | Error Codes |
|----------|-------------|-------------|
| Register | 201 Created | 400, 409 |
| Login | 200 OK | 401, 403 |
| Refresh | 200 OK | 401 |
| Profile | 200 OK | 401 |
| Logout | 200 OK | 401 |
| Admin | 200 OK | 401, 403 |
