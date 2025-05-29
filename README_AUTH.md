# Better-Auth Integration with Go Gin Backend

This backend now integrates with better-auth for user authentication with advanced session synchronization. Here's how it works:

## How It Works

1. **Frontend Authentication**: Users authenticate through your frontend using better-auth
2. **Session Validation**: The Go backend validates sessions by making HTTP requests to your frontend's better-auth API
3. **Session Caching**: Backend caches validated sessions for 5 minutes to reduce API calls
4. **User Management**: The backend automatically creates users in its database when they first access protected endpoints
5. **Data Synchronization**: User data is automatically synced between frontend and backend
6. **Protected Routes**: All routes under `/api/*` require authentication

## Session Synchronization Features

### Caching
- Sessions are cached for 5 minutes to reduce validation calls
- Automatic cleanup of expired cache entries every 10 minutes
- Cache invalidation when sessions become invalid

### Real-time Sync
- Backend syncs user data from better-auth on each request
- Frontend can notify backend of session changes
- Session invalidation and refresh endpoints

### Data Consistency
- User names and other data automatically sync from better-auth
- Database updates when user information changes
- Graceful handling of sync failures

## Configuration

### Frontend URL
Update the `frontendURL` variable in `main.go` to match your frontend's URL:
```go
frontendURL = "http://localhost:3000" // Change this to your actual frontend URL
```

### Environment Variables
Set `NEXT_PUBLIC_BACKEND_URL` in your frontend `.env` file:
```bash
NEXT_PUBLIC_BACKEND_URL=http://localhost:8080
```

### CORS
The backend is configured to allow cross-origin requests from your frontend URL with credentials enabled.

## Authentication Flow

1. User logs in through your frontend (better-auth handles this)
2. Better-auth sets a session cookie (`better-auth.session_token`)
3. When making API requests to the backend, the browser automatically includes this cookie
4. The backend validates the session with your frontend's `/api/auth/session` endpoint
5. If valid, the backend either finds or creates the user in its database
6. User data is synced from better-auth if needed
7. The user object is made available in the Gin context for your handlers

## Endpoints

### Public Endpoints
- `GET /` - Basic health check
- `GET /health` - Health check with authentication status

### Session Management
- `POST /api/session/invalidate` - Invalidate a session in backend cache
- `POST /api/session/refresh` - Force session refresh

### Protected Endpoints (require authentication)
- `GET /api/*` - All API endpoints require authentication
- The user object is available in handlers via: `user := c.MustGet("user").(db.User)`

## Frontend Integration

### Using the Sync Utility

Import the sync utility in your frontend:
```typescript
import { BackendSync, useBackendSync } from '@/lib/backend-sync';
```

### In React Components
```typescript
import { useBackendSync } from '@/lib/backend-sync';

function MyComponent() {
  const { invalidateSession, refreshSession, checkAuthStatus } = useBackendSync();
  
  // Call when user logs out
  const handleLogout = async () => {
    await invalidateSession();
    // ... rest of logout logic
  };
  
  // Call when user data changes
  const handleProfileUpdate = async () => {
    await refreshSession();
  };
  
  // Check sync status
  const checkSync = async () => {
    const status = await checkAuthStatus();
    console.log('Backend auth status:', status);
  };
}
```

### Better-Auth Event Handlers

Add to your better-auth configuration:
```typescript
export const auth = betterAuth({
  // ... your existing config
  databaseHooks: {
    user: {
      update: {
        after: async (user) => {
          // Sync user changes with backend
          await BackendSync.refreshSession();
        },
      },
    },
    session: {
      delete: {
        after: async (session) => {
          // Notify backend of logout
          await BackendSync.invalidateSession();
        },
      },
    },
  },
});
```

## Testing Authentication & Sync

1. Start your frontend with better-auth running
2. Start the Go backend: `go run main.go`
3. Log in through your frontend
4. Test authentication status: `GET http://localhost:8080/health`
5. Test protected endpoint: `GET http://localhost:8080/api/acronyms`
6. Test session sync from frontend:
   ```typescript
   import { BackendSync } from '@/lib/backend-sync';
   
   // Check if backend recognizes the session
   const status = await BackendSync.checkAuthStatus();
   console.log('Sync status:', status);
   ```

## Performance Optimizations

### Session Caching
- Sessions are cached for 5 minutes to reduce API calls to frontend
- Cache automatically expires and cleans up unused entries
- Immediate cache invalidation on session errors

### Request Timeouts
- 10-second timeout on session validation requests
- Graceful fallback on network errors

### Background Cleanup
- Expired cache entries are cleaned up every 10 minutes
- No memory leaks from abandoned sessions

## Error Handling

The backend will return:
- `401 Unauthorized` - No session cookie or invalid session
- `500 Internal Server Error` - Database errors or other server issues

Session sync operations are designed to fail gracefully:
- Sync failures are logged but don't block requests
- Network errors don't affect core functionality
- Cache misses fall back to fresh validation

## User Management

### Automatic Creation
When a user first accesses a protected endpoint:
- Backend automatically creates a user record using better-auth session data
- Uses email and name from the session
- User is available for subsequent requests

### Data Synchronization
On each authenticated request:
- Backend checks if user data needs updating
- Syncs name and other fields from better-auth
- Updates are logged but don't block requests

## Security Considerations

1. **HTTPS**: In production, ensure both frontend and backend use HTTPS
2. **CORS**: Restrict CORS origins to your actual frontend domain
3. **Session Security**: Rely on better-auth's session security mechanisms
4. **Database**: Ensure your database connections are secure
5. **Cache Security**: Session cache is in-memory only and automatically expires
6. **Timeouts**: All external requests have timeouts to prevent hanging 