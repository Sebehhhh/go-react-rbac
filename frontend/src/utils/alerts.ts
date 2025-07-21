import Swal from 'sweetalert2';

export const showSuccessAlert = (message: string) => {
  Swal.fire({
    icon: 'success',
    title: 'Success',
    text: message,
    showConfirmButton: false,
    timer: 1500,
  });
};

export const showErrorAlert = (message: string) => {
  Swal.fire({
    icon: 'error',
    title: 'Error',
    text: message,
    showConfirmButton: true,
  });
};

export const showConfirmAlert = async (title: string, text: string) => {
  const result = await Swal.fire({
    title: title,
    text: text,
    icon: 'warning',
    showCancelButton: true,
    confirmButtonColor: '#3085d6',
    cancelButtonColor: '#d33',
    confirmButtonText: 'Yes, proceed!',
  });
  return result.isConfirmed;
};
