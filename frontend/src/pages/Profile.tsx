import { useAuthStore } from '../hooks/useAuth';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import api from '../services/api';
import { showSuccessAlert, showErrorAlert } from '../utils/alerts';

// Schema for profile update
const profileSchema = z.object({
  first_name: z.string().min(1, 'First name is required'),
  last_name: z.string().min(1, 'Last name is required'),
  email: z.string().email('Invalid email address'),
  username: z.string().min(3, 'Username must be at least 3 characters'),
});

type ProfileFormInputs = z.infer<typeof profileSchema>;

// Schema for password change
const passwordSchema = z.object({
  current_password: z.string().min(8, 'Current password must be at least 8 characters'),
  new_password: z.string().min(8, 'New password must be at least 8 characters'),
  confirm_password: z.string().min(8, 'Confirm password must be at least 8 characters'),
}).refine((data) => data.new_password === data.confirm_password, {
  message: "Passwords don't match",
  path: ["confirm_password"],
});

type PasswordFormInputs = z.infer<typeof passwordSchema>;

const ProfilePage = () => {
  const { user, setUser } = useAuthStore();
  

  const { register: profileRegister, handleSubmit: handleProfileSubmit, formState: { errors: profileErrors, isSubmitting: isProfileSubmitting } } = useForm<ProfileFormInputs>({
    resolver: zodResolver(profileSchema),
    defaultValues: {
      first_name: user?.first_name || '',
      last_name: user?.last_name || '',
      email: user?.email || '',
      username: user?.username || '',
    },
  });

  const { register: passwordRegister, handleSubmit: handlePasswordSubmit, formState: { errors: passwordErrors, isSubmitting: isPasswordSubmitting }, reset: resetPasswordForm } = useForm<PasswordFormInputs>({
    resolver: zodResolver(passwordSchema),
  });

  const onProfileSubmit = async (data: ProfileFormInputs) => {
    try {
      const response = await api.put('/profile', data);
      setUser(response.data.data); // Update user in store
      showSuccessAlert('Profile updated successfully!');
    } catch (error: any) {
      showErrorAlert(error.response?.data?.message || 'Failed to update profile.');
    }
  };

  const onPasswordSubmit = async (data: PasswordFormInputs) => {
    try {
      await api.put('/profile/password', {
        current_password: data.current_password,
        new_password: data.new_password,
      });
      showSuccessAlert('Password updated successfully!');
      resetPasswordForm();
    } catch (error: any) {
      showErrorAlert(error.response?.data?.message || 'Failed to change password.');
    }
  };

  return (
    <div className="p-6 bg-gray-100 min-h-screen">
      <h1 className="text-3xl font-bold text-gray-800 mb-6">My Profile</h1>

      {/* Profile Information Section */}
      <div className="bg-white p-6 rounded-lg shadow-md mb-8">
        <h2 className="text-xl font-semibold text-gray-800 mb-4">Profile Information</h2>
        
        <form onSubmit={handleProfileSubmit(onProfileSubmit)} className="space-y-4">
          <div>
            <label htmlFor="first_name" className="block text-sm font-medium text-gray-700">First Name</label>
            <input id="first_name" {...profileRegister('first_name')} className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-300 focus:ring focus:ring-indigo-200 focus:ring-opacity-50" />
            {profileErrors.first_name && <p className="mt-1 text-sm text-red-600">{profileErrors.first_name.message}</p>}
          </div>
          <div>
            <label htmlFor="last_name" className="block text-sm font-medium text-gray-700">Last Name</label>
            <input id="last_name" {...profileRegister('last_name')} className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-300 focus:ring focus:ring-indigo-200 focus:ring-opacity-50" />
            {profileErrors.last_name && <p className="mt-1 text-sm text-red-600">{profileErrors.last_name.message}</p>}
          </div>
          <div>
            <label htmlFor="email" className="block text-sm font-medium text-gray-700">Email</label>
            <input id="email" type="email" {...profileRegister('email')} className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-300 focus:ring focus:ring-indigo-200 focus:ring-opacity-50" />
            {profileErrors.email && <p className="mt-1 text-sm text-red-600">{profileErrors.email.message}</p>}
          </div>
          <div>
            <label htmlFor="username" className="block text-sm font-medium text-gray-700">Username</label>
            <input id="username" {...profileRegister('username')} className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-300 focus:ring focus:ring-indigo-200 focus:ring-opacity-50" />
            {profileErrors.username && <p className="mt-1 text-sm text-red-600">{profileErrors.username.message}</p>}
          </div>
          <button
            type="submit"
            disabled={isProfileSubmitting}
            className="px-4 py-2 bg-indigo-600 text-white font-semibold rounded-md shadow-sm hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {isProfileSubmitting ? 'Saving...' : 'Save Profile'}
          </button>
        </form>
      </div>

      {/* Change Password Section */}
      <div className="bg-white p-6 rounded-lg shadow-md">
        <h2 className="text-xl font-semibold text-gray-800 mb-4">Change Password</h2>
        
        <form onSubmit={handlePasswordSubmit(onPasswordSubmit)} className="space-y-4">
          <div>
            <label htmlFor="current_password" className="block text-sm font-medium text-gray-700">Current Password</label>
            <input id="current_password" type="password" {...passwordRegister('current_password')} className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-300 focus:ring focus:ring-indigo-200 focus:ring-opacity-50" />
            {passwordErrors.current_password && <p className="mt-1 text-sm text-red-600">{passwordErrors.current_password.message}</p>}
          </div>
          <div>
            <label htmlFor="new_password" className="block text-sm font-medium text-gray-700">New Password</label>
            <input id="new_password" type="password" {...passwordRegister('new_password')} className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-300 focus:ring focus:ring-indigo-200 focus:ring-opacity-50" />
            {passwordErrors.new_password && <p className="mt-1 text-sm text-red-600">{passwordErrors.new_password.message}</p>}
          </div>
          <div>
            <label htmlFor="confirm_password" className="block text-sm font-medium text-gray-700">Confirm New Password</label>
            <input id="confirm_password" type="password" {...passwordRegister('confirm_password')} className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-300 focus:ring focus:ring-indigo-200 focus:ring-opacity-50" />
            {passwordErrors.confirm_password && <p className="mt-1 text-sm text-red-600">{passwordErrors.confirm_password.message}</p>}
          </div>
          <button
            type="submit"
            disabled={isPasswordSubmitting}
            className="px-4 py-2 bg-indigo-600 text-white font-semibold rounded-md shadow-sm hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {isPasswordSubmitting ? 'Changing...' : 'Change Password'}
          </button>
        </form>
      </div>
    </div>
  );
};

export default ProfilePage;
