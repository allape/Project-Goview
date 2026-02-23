import cls from "classnames";
import { PropsWithChildren, ReactElement, ReactNode } from "react";
import { NO_PREVIEW } from "../../config";
import styles from "./style.module.scss";

export interface IFileProps {
  name: ReactNode;
  alt?: string;
  hidden?: boolean;
  center?: boolean;
  cover?: string;
  className?: string;
  onClick?: () => void;
}

export default function File({
  children,
  name,
  alt,
  hidden,
  center,
  cover,
  className,
  onClick,
}: PropsWithChildren<IFileProps>): ReactElement {
  return (
    <div className={cls(styles.file, hidden && styles.hidden, className)}>
      <div
        className={cls(styles.cover, onClick && styles.clickable)}
        onClick={onClick}
      >
        <img src={cover || NO_PREVIEW} alt={alt} loading="lazy" />
      </div>
      <div className={cls(styles.name, center && styles.center)}>{name}</div>
      {children && <div className={styles.actions}>{children}</div>}
    </div>
  );
}
