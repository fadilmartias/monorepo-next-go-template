"use client";
import { apiClient } from "@/lib/api-client";
import { Editor } from "@tinymce/tinymce-react";
import { useRef } from "react";
import Typography from "./typography";
import { Label } from "./ui/label";

const TinyEditor = ({
  onEditorChange,
  value,
  error,
  label,
  id,
  required,
  requiredSign = true,
  ...props
}: {
  onEditorChange: (content: string) => void;
  value?: string;
  error?: string;
  label?: string;
  id?: string;
  required?: boolean;
  requiredSign?: boolean;
}) => {
  const editorRef = useRef(null);
  return (
    <div className="space-y-2">
      {label && <Label htmlFor={id || label || ""}>{label}{required && requiredSign && <span className="text-destructive">*</span>}</Label>}
      <Editor
        tinymceScriptSrc={"/assets/lib/tinymce/tinymce.min.js"}
        onInit={(_evt, editor) => (editorRef.current = editor)}
        licenseKey="gpl"
        init={{
          menubar: false,
          plugins:
            "anchor autolink charmap codesample emoticons fullscreen image link lists media searchreplace table visualblocks wordcount preview code",
          toolbar:
            "fullscreen | undo redo | blocks fontfamily fontsize | bold italic underline strikethrough | link image media table | align lineheight | numlist bullist indent outdent | emoticons charmap | removeformat | code | preview",
          skin: "oxide-dark",
          images_upload_handler: async (blobInfo: any, progress: any) => {
            const formData = new FormData();
            formData.append("tinymce", blobInfo.blob(), blobInfo.filename());

            const res = await apiClient.put(
              "/v0/uploads",
              formData,
              {
                headers: {
                  "Content-Type": "multipart/form-data",
                },
                onUploadProgress: (e: any) => {
                  if (progress && e.lengthComputable) {
                    progress((e.loaded / e.total) * 100);
                  }
                },
              }
            );

            const imageUrl = res.data?.data?.file_url;
            if (!imageUrl) {
              throw new Error("Upload failed, no image URL returned");
            }

            return imageUrl;
          },
        }}
        onEditorChange={onEditorChange}
        value={value}
        id={id || label || ""}
        {...props}
      />
      {error && (
        <Typography variant="c1" className="text-destructive">
          {error}
        </Typography>
      )}
    </div>
  );
};

export default TinyEditor;
