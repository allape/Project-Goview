import { BaseSearchParams } from "@allape/gocrud";
import { ILV } from "@allape/gocrud-react/src/helper/antd.tsx";
import { useLoading } from "@allape/use-loading";
import { Select, Spin } from "antd";
import { ReactElement, useCallback, useEffect, useState } from "react";
import { DatasourceCrudy } from "../../api/datasource.ts";
import IDatasource from "../../model/datasource.ts";
import PrimaryDropdown from "../PrimaryDropdown";
import styles from "./style.module.scss";

export interface IDatasourceBarProps {
  value?: IDatasource["id"];
  onChange?: (value: IDatasource["id"]) => void;
}

export default function DatasourceBar({
  value,
  onChange,
}: IDatasourceBarProps): ReactElement {
  const { loading, execute } = useLoading();

  const [selected, setSelected] = useState<IDatasource["id"] | undefined>();
  const [options, setOptions] = useState<ILV<IDatasource["id"]>[]>([]);

  const getSources = useCallback(async () => {
    await execute(async () => {
      const sources = await DatasourceCrudy.all({
        ...BaseSearchParams,
      });
      setOptions(
        sources.map((d) => ({
          value: d.id,
          label: d.name,
        })),
      );
    });
  }, [execute]);

  useEffect(() => {
    getSources().then();
  }, [getSources]);

  useEffect(() => {
    setSelected(value);
  }, [value]);

  return (
    <Spin spinning={loading}>
      <div className={styles.wrapper}>
        <Select
          className={styles.selector}
          placeholder="Select a datasource to preview"
          showSearch
          allowClear
          optionFilterProp="label"
          options={options}
          onChange={(v) => {
            setSelected(v);
            onChange?.(v);
          }}
          value={selected}
        />
        <PrimaryDropdown afterSaved={getSources} />
      </div>
    </Spin>
  );
}
