export type Notification = {
  title: string;
  message: string;
  type: NotificationType;
  autoClose?: boolean;
  autoCloseDelay?: number;
};

export enum NotificationType {
  Error = "error",
  Success = "success",
  Info = "info",
  Warning = "warning",
}

export function colorFromType(type: NotificationType) {
  switch (type) {
    case NotificationType.Error:
      return "red";
    case NotificationType.Success:
      return "green";
    case NotificationType.Info:
      return "blue";
    case NotificationType.Warning:
      return "yellow";
  }
}
