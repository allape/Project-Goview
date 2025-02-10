import { ThemeProvider } from "@allape/gocrud-react";
import { ReactElement, useState } from "react";
import DatasourceBar from "./components/DatasourceBar";
import Explorer from "./components/Explorer";
import IDatasource from "./model/datasource.ts";
import styles from "./style.module.scss";

export default function App(): ReactElement {
  const [id, setId] = useState<IDatasource["id"] | undefined>();
  return (
    <ThemeProvider>
      <div className={styles.wrapper}>
        <DatasourceBar value={id} onChange={setId} />
        <Explorer value={id} />
      </div>
    </ThemeProvider>
  );
}
