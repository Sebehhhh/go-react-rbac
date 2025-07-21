import { Dialog, Transition } from '@headlessui/react';
import { Fragment } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';

const passwordSchema = z.object({
  current_password: z.string().min(8, 'Current password must be at least 8 characters'),
  new_password: z.string().min(8, 'New password must be at least 8 characters'),
  confirm_password: z.string().min(8, 'Confirm password must be at least 8 characters'),
}).refine((data) => data.new_password === data.confirm_password, {
  message: 'Passwords do not match',
  path: ['confirm_password'],
});

type PasswordFormInputs = z.infer<typeof passwordSchema>;

interface PasswordUpdateModalProps {
  isOpen: boolean;
  onClose: () => void;
  userId: number;
  userFullName: string;
  onUpdatePassword: (userId: number, data: PasswordFormInputs) => void;
}

const PasswordUpdateModal: React.FC<PasswordUpdateModalProps> = ({ 
  isOpen, 
  onClose, 
  userId, 
  userFullName, 
  onUpdatePassword 
}) => {
  const { register, handleSubmit, formState: { errors }, reset } = useForm<PasswordFormInputs>({
    resolver: zodResolver(passwordSchema),
  });

  const onSubmit = (data: PasswordFormInputs) => {
    onUpdatePassword(userId, data);
    reset();
    onClose();
  };

  const handleClose = () => {
    reset();
    onClose();
  };

  return (
    <Transition appear show={isOpen} as={Fragment}>
      <Dialog as="div" className="relative z-10" onClose={handleClose}>
        <Transition.Child as={Fragment} enter="ease-out duration-300" enterFrom="opacity-0" enterTo="opacity-100" leave="ease-in duration-200" leaveFrom="opacity-100" leaveTo="opacity-0">
          <div className="fixed inset-0 bg-black bg-opacity-25" />
        </Transition.Child>
        <div className="fixed inset-0 overflow-y-auto">
          <div className="flex items-center justify-center min-h-full p-4 text-center">
            <Transition.Child as={Fragment} enter="ease-out duration-300" enterFrom="opacity-0 scale-95" enterTo="opacity-100 scale-100" leave="ease-in duration-200" leaveFrom="opacity-100 scale-100" leaveTo="opacity-0 scale-95">
              <Dialog.Panel className="w-full max-w-md p-6 overflow-hidden text-left align-middle transition-all transform bg-white shadow-xl rounded-2xl">
                <Dialog.Title as="h3" className="text-lg font-medium leading-6 text-gray-900">
                  Update Password for {userFullName}
                </Dialog.Title>
                <form onSubmit={handleSubmit(onSubmit)} className="mt-4 space-y-4">
                  <div>
                    <label htmlFor="current_password" className="block text-sm font-medium text-gray-700">
                      Current Password
                    </label>
                    <input 
                      id="current_password" 
                      type="password" 
                      {...register('current_password')} 
                      className="w-full px-3 py-2 mt-1 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500" 
                    />
                    {errors.current_password && <p className="text-sm text-red-500">{errors.current_password.message}</p>}
                  </div>
                  <div>
                    <label htmlFor="new_password" className="block text-sm font-medium text-gray-700">
                      New Password
                    </label>
                    <input 
                      id="new_password" 
                      type="password" 
                      {...register('new_password')} 
                      className="w-full px-3 py-2 mt-1 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500" 
                    />
                    {errors.new_password && <p className="text-sm text-red-500">{errors.new_password.message}</p>}
                  </div>
                  <div>
                    <label htmlFor="confirm_password" className="block text-sm font-medium text-gray-700">
                      Confirm New Password
                    </label>
                    <input 
                      id="confirm_password" 
                      type="password" 
                      {...register('confirm_password')} 
                      className="w-full px-3 py-2 mt-1 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500" 
                    />
                    {errors.confirm_password && <p className="text-sm text-red-500">{errors.confirm_password.message}</p>}
                  </div>
                  <div className="mt-6 flex justify-end space-x-2">
                    <button 
                      type="button" 
                      onClick={handleClose} 
                      className="px-4 py-2 text-sm font-medium text-gray-700 bg-gray-100 border border-transparent rounded-md hover:bg-gray-200 focus:outline-none focus:ring-2 focus:ring-gray-500"
                    >
                      Cancel
                    </button>
                    <button 
                      type="submit" 
                      className="px-4 py-2 text-sm font-medium text-white bg-indigo-600 border border-transparent rounded-md hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500"
                    >
                      Update Password
                    </button>
                  </div>
                </form>
              </Dialog.Panel>
            </Transition.Child>
          </div>
        </div>
      </Dialog>
    </Transition>
  );
};

export default PasswordUpdateModal;