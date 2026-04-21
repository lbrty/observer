import {
  FileArchiveIcon,
  FileAudioIcon,
  FileCsvIcon,
  FileDashedIcon,
  FileDocIcon,
  FileImageIcon,
  FilePdfIcon,
  FilePngIcon,
  FilePptIcon,
  FileSvgIcon,
  FileTextIcon,
  FileVideoIcon,
  FileXlsIcon,
} from "@/components/icons";
import type { Icon } from "@/components/icons";

export function formatBytes(bytes: number): string {
  if (bytes === 0) return "0 B";
  const k = 1024;
  const sizes = ["B", "KB", "MB", "GB"];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return `${(bytes / Math.pow(k, i)).toFixed(1)} ${sizes[i]}`;
}

export function mimeIcon(mime: string): Icon {
  if (mime === "application/pdf") return FilePdfIcon;
  if (mime.startsWith("image/png")) return FilePngIcon;
  if (mime.startsWith("image/svg")) return FileSvgIcon;
  if (mime === "image/avif" || mime === "image/heif" || mime === "image/heic")
    return FileDashedIcon;
  if (mime.startsWith("image/")) return FileImageIcon;
  if (mime.startsWith("video/")) return FileVideoIcon;
  if (mime.startsWith("audio/")) return FileAudioIcon;
  if (mime === "text/csv") return FileCsvIcon;
  if (mime.startsWith("text/")) return FileTextIcon;
  if (
    mime === "application/msword" ||
    mime === "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
  )
    return FileDocIcon;
  if (
    mime === "application/vnd.ms-excel" ||
    mime === "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
  )
    return FileXlsIcon;
  if (
    mime === "application/vnd.ms-powerpoint" ||
    mime === "application/vnd.openxmlformats-officedocument.presentationml.presentation"
  )
    return FilePptIcon;
  if (
    mime === "application/zip" ||
    mime === "application/x-rar-compressed" ||
    mime === "application/gzip"
  )
    return FileArchiveIcon;
  return FileDashedIcon;
}
