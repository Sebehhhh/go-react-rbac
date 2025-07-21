import { useAuthStore } from '../../hooks/useAuth';
import { UserCircle, LogOut } from 'lucide-react';
import { useNavigate } from 'react-router-dom';
import { showSuccessAlert, showConfirmAlert } from '../../utils/alerts';

const Header = () => {
  const { user, logout } = useAuthStore();
  const navigate = useNavigate();

  const handleLogout = async () => {
    const confirmed = await showConfirmAlert(
      'Logout Confirmation',
      'Are you sure you want to logout?'
    );
    
    if (confirmed) {
      logout();
      showSuccessAlert('Logged out successfully!');
      navigate('/login');
    }
  };

  return (
    <header className="flex items-center justify-between p-4 md:p-6 bg-white border-b">
      <div>
        {/* Can add breadcrumbs or page title here */}
      </div>
      <div className="flex items-center">
        <span className="mr-2 md:mr-4 text-sm md:text-base">Welcome, {user?.first_name}</span>
        <UserCircle className="w-6 h-6 md:w-8 md:h-8 text-gray-600" />
        <button onClick={handleLogout} className="ml-2 md:ml-4 text-gray-600 hover:text-red-500">
          <LogOut className="w-5 h-5 md:w-6 md:h-6" />
        </button>
      </div>
    </header>
  );
};

export default Header;
