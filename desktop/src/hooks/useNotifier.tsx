import { createContext, useContext, useEffect, useState } from "react";
import { Notification } from "../types/notification";

const NotifierContext = createContext<
  | {
      notification: Notification | null;
      setNotification: (notification: Notification | null) => void;
    }
  | undefined
>(undefined);

export function NotifierProvider({ children }: { children: React.ReactNode }) {
  const [notification, setNotification] = useState<Notification | null>(null);

  useEffect(() => {
    if (notification?.autoClose) {
      setTimeout(() => {
        setNotification(null);
      }, notification.autoCloseDelay || 10000);
    }
  }, [notification]);

  return (
    <NotifierContext.Provider value={{ notification, setNotification }}>
      {children}
    </NotifierContext.Provider>
  );
}

export function useNotifier() {
  const context = useContext(NotifierContext);
  if (context === undefined) {
    throw new Error("useNotifier must be used within a NotifierProvider");
  }
  return context;
}
