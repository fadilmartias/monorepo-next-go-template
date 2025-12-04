import Swal from "sweetalert2";

export const shadSwal = Swal.mixin({
  customClass: {
    popup: "shadcn",
    confirmButton: "swal2-confirm",
    cancelButton: "swal2-cancel",
    title: "swal2-title",
    htmlContainer: "swal2-html-container",
    actions: "swal2-actions",
  },
  buttonsStyling: false,
});
