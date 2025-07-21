# 🔐 Go React RBAC System

A modern, full-stack Role-Based Access Control (RBAC) system built with Go and React. This boilerplate provides a complete authentication and authorization solution with a clean, scalable architecture.

## 🚀 Features

- **Complete RBAC Implementation**: Users, Roles, Permissions with fine-grained access control
- **JWT Authentication**: Secure token-based authentication with refresh tokens
- **Real-time Dashboard**: Statistics, analytics, and activity monitoring
- **User Management**: Complete CRUD operations with role assignments
- **Role Management**: Dynamic role creation with permission assignments
- **Password Management**: Secure password updates and reset functionality
- **Activity Logging**: Track all user actions and system events
- **Responsive UI**: Mobile-first design with modern components
- **Type Safety**: Full TypeScript support on frontend
- **Database Seeding**: Smart seeding system with duplicate prevention
- **API Documentation**: Well-structured REST API endpoints

## 🛠️ Tech Stack

### Backend
- **[Go](https://golang.org/)** `1.21+` - Modern programming language
- **[Fiber](https://gofiber.io/)** `v2.52.0` - Express-inspired web framework
- **[GORM](https://gorm.io/)** `v1.25.5` - ORM for database operations
- **[MySQL](https://mysql.com/)** `8.0+` - Relational database
- **[JWT](https://github.com/golang-jwt/jwt)** `v5.2.0` - JSON Web Token implementation
- **[Bcrypt](https://golang.org/x/crypto/bcrypt)** - Password hashing
- **[Validator](https://github.com/go-playground/validator)** `v10.16.0` - Struct validation

### Frontend
- **[React](https://reactjs.org/)** `18.2.0` - UI library
- **[TypeScript](https://www.typescriptlang.org/)** `5.2.2` - Type safety
- **[Vite](https://vitejs.dev/)** `5.0.8` - Build tool and dev server
- **[Zustand](https://zustand-demo.pmnd.rs/)** `4.4.6` - State management
- **[React Hook Form](https://react-hook-form.com/)** `7.48.2` - Form handling
- **[Zod](https://zod.dev/)** `3.22.4` - Schema validation
- **[Tailwind CSS](https://tailwindcss.com/)** `3.3.6` - Utility-first styling
- **[Headless UI](https://headlessui.com/)** `1.7.17` - Accessible UI components
- **[Lucide React](https://lucide.dev/)** `0.294.0` - Icon library
- **[Recharts](https://recharts.org/)** `2.8.0` - Charts and analytics
- **[SweetAlert2](https://sweetalert2.github.io/)** `11.10.1` - Beautiful alerts
- **[Axios](https://axios-http.com/)** `1.6.2` - HTTP client

### Development Tools
- **[Air](https://github.com/cosmtrek/air)** - Hot reload for Go
- **[ESLint](https://eslint.org/)** `8.55.0` - Code linting
- **[Prettier](https://prettier.io/)** - Code formatting

## 📋 Prerequisites

- **Go** 1.21 or higher
- **Node.js** 18.0 or higher
- **MySQL** 8.0 or higher
- **Git** for version control

## 🚀 Quick Start

### 1. Clone the Repository
```bash
git clone https://github.com/Sebehhhh/go-react-rbac.git
cd go-react-rbac
```

### 2. Backend Setup

#### Install Dependencies
```bash
cd backend
go mod download
```

#### Environment Configuration
Create `.env` file in the backend directory:
```env
# Database Configuration
DB_HOST=localhost
DB_PORT=3306
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_NAME=rbac_system

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
JWT_REFRESH_SECRET=your-super-secret-refresh-key-change-this-too

# Server Configuration
PORT=8080
CORS_ALLOWED_ORIGINS=http://localhost:3000

# System Configuration
DEFAULT_ADMIN_EMAIL=admin@example.com
DEFAULT_ADMIN_PASSWORD=admin123456
```

#### Database Setup
Create MySQL database:
```sql
CREATE DATABASE rbac_system CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

#### Run Backend
```bash
# Development with hot reload (install Air first)
go install github.com/cosmtrek/air@latest
air

# Or run directly
go run cmd/server/main.go
```

### 3. Frontend Setup

#### Install Dependencies
```bash
cd frontend
npm install
```

#### Environment Configuration
Create `.env` file in the frontend directory:
```env
VITE_API_URL=http://localhost:8080/api
```

#### Run Frontend
```bash
npm run dev
```

### 4. Access the Application

- **Frontend**: http://localhost:3000
- **Backend API**: http://localhost:8080
- **Health Check**: http://localhost:8080/health

### 5. Default Login Credentials

```
Email: admin@example.com
Password: admin123456
```

## 📁 Project Structure

```
go-react-rbac/
├── backend/
│   ├── cmd/server/           # Application entry point
│   ├── internal/
│   │   ├── config/           # Configuration management
│   │   ├── database/         # Database connection & migrations
│   │   ├── handlers/         # HTTP request handlers
│   │   ├── middleware/       # Custom middleware
│   │   ├── models/           # Database models & DTOs
│   │   ├── services/         # Business logic layer
│   │   └── utils/            # Utility functions
│   ├── .env.example          # Environment variables template
│   ├── go.mod                # Go dependencies
│   └── go.sum                # Go dependencies checksum
└── frontend/
    ├── public/               # Static assets
    ├── src/
    │   ├── components/       # Reusable UI components
    │   ├── guards/           # Route protection
    │   ├── hooks/            # Custom React hooks
    │   ├── pages/            # Page components
    │   ├── services/         # API service layer
    │   ├── types/            # TypeScript type definitions
    │   └── utils/            # Utility functions
    ├── .env.example          # Environment variables template
    ├── package.json          # Node dependencies
    └── vite.config.ts        # Vite configuration
```

## 🔧 API Endpoints

### Authentication
- `POST /api/auth/login` - User login
- `POST /api/auth/register` - User registration
- `POST /api/auth/refresh` - Refresh JWT token
- `POST /api/auth/logout` - User logout
- `POST /api/auth/forgot-password` - Request password reset
- `POST /api/auth/reset-password` - Reset password

### Users
- `GET /api/users` - List users (paginated)
- `POST /api/users` - Create user
- `GET /api/users/:id` - Get user details
- `PUT /api/users/:id` - Update user
- `DELETE /api/users/:id` - Delete user
- `PUT /api/users/:id/password` - Update user password
- `PUT /api/users/:id/activate` - Activate user
- `PUT /api/users/:id/deactivate` - Deactivate user

### Roles
- `GET /api/roles` - List roles
- `POST /api/roles` - Create role
- `GET /api/roles/:id` - Get role details
- `PUT /api/roles/:id` - Update role
- `DELETE /api/roles/:id` - Delete role
- `PUT /api/roles/:id/permissions` - Assign permissions to role
- `GET /api/roles/:id/permissions` - Get role permissions

### Permissions
- `GET /api/permissions` - List all permissions

### Profile
- `GET /api/profile` - Get current user profile
- `PUT /api/profile` - Update profile
- `PUT /api/profile/password` - Change password

### Dashboard
- `GET /api/dashboard/stats` - Get dashboard statistics
- `GET /api/dashboard/role-distribution` - Get role distribution
- `GET /api/dashboard/recent-activity` - Get recent activity logs
- `GET /api/dashboard/user-analytics` - Get user analytics

## 🔐 Default Roles & Permissions

The system comes with 4 pre-configured roles:

### Super Admin
- Full system access
- All CRUD operations on users, roles, and permissions
- System administration capabilities

### Admin
- User management (create, read, update, delete)
- Role viewing and basic management
- Dashboard and activity log access

### Manager
- User viewing and limited editing
- Role viewing
- Dashboard access

### User
- Profile management only
- Basic dashboard access

## 🛡️ Security Features

- **Password Hashing**: Bcrypt with salt
- **JWT Tokens**: Secure authentication with refresh tokens
- **CORS Protection**: Configurable CORS settings
- **Input Validation**: Comprehensive request validation
- **SQL Injection Prevention**: GORM ORM with parameterized queries
- **Rate Limiting**: Built-in Fiber rate limiting
- **Activity Logging**: All user actions are tracked
- **Role-Based Access**: Fine-grained permission system

## 🧪 Development

### Backend Development
```bash
cd backend

# Run with hot reload
air

# Run tests
go test ./...

# Build for production
go build -o main cmd/server/main.go
```

### Frontend Development
```bash
cd frontend

# Development server
npm run dev

# Build for production
npm run build

# Preview production build
npm run preview

# Lint code
npm run lint
```

## 📦 Production Deployment

### Backend
```bash
# Build binary
go build -o rbac-server cmd/server/main.go

# Run in production
./rbac-server
```

### Frontend
```bash
# Build for production
npm run build

# Serve static files with any web server
# The build output will be in the dist/ directory
```

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Fiber](https://gofiber.io/) for the excellent Go web framework
- [React](https://reactjs.org/) team for the amazing UI library
- [Tailwind CSS](https://tailwindcss.com/) for the utility-first CSS framework
- All other open-source libraries that made this project possible

## 📞 Support

If you have any questions or need help getting started:

- Open an issue on [GitHub](https://github.com/Sebehhhh/go-react-rbac/issues)
- Check the documentation in the codebase
- Review the example environment files

---

Made with ❤️ by [Sebeh](https://github.com/Sebehhhh)