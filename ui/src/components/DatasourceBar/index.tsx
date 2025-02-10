import { BaseSearchParams } from "@allape/gocrud";
import { CrudyButton } from "@allape/gocrud-react";
import { ICrudyButtonProps } from "@allape/gocrud-react/src/component/CrudyButton";
import NewCrudyButtonEventEmitter from "@allape/gocrud-react/src/component/CrudyButton/eventemitter.ts";
import { ILV } from "@allape/gocrud-react/src/helper/antd.tsx";
import { useLoading } from "@allape/use-loading";
import { DownOutlined } from "@ant-design/icons";
import {
  Button,
  Dropdown,
  Form,
  Input,
  MenuProps,
  Select,
  Space,
  Spin,
  Tag,
} from "antd";
import { ReactElement, useCallback, useEffect, useMemo, useState } from "react";
import { DatasourceCrudy } from "../../api/datasource.ts";
import { TagCrudy } from "../../api/tag.ts";
import IDatasource, { DatasourceTypes } from "../../model/datasource.ts";
import ITag from "../../model/tag.ts";
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

  const datasourceCrudyEmitter = useMemo(
    () => NewCrudyButtonEventEmitter<IDatasource>(),
    [],
  );

  const tagCrudyEmitter = useMemo(() => NewCrudyButtonEventEmitter<ITag>(), []);

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

  const datasourceColumns = useMemo<ICrudyButtonProps<IDatasource>["columns"]>(
    () => [
      {
        dataIndex: "id",
        title: "ID",
      },
      {
        dataIndex: "name",
        title: "Name",
      },
      {
        dataIndex: "type",
        title: "Type",
        render: (v) => DatasourceTypes.find((t) => t.value === v)?.label || v,
      },
      {
        dataIndex: "cwd",
        title: "Cwd",
        render: (v) => {
          try {
            const url = new URL(v);
            if (url.password) {
              url.password = "******";
            }
            if (url.username) {
              url.username = url.username.substring(0, 3) + "******";
            }
            return url.toString();
          } catch {
            // ignore
          }
          return v;
        },
      },
    ],
    [],
  );

  const tagColumns = useMemo<ICrudyButtonProps<ITag>["columns"]>(
    () => [
      {
        dataIndex: "id",
        title: "ID",
      },
      {
        dataIndex: "key",
        title: "Key",
      },
      {
        dataIndex: "name",
        title: "Name",
        render: (v, record) => (
          <Tag color={record.color}>
            <span style={{ color: record.color, filter: "invert(1)" }}>
              {v}
            </span>
          </Tag>
        ),
      },
    ],
    [],
  );

  useEffect(() => {
    setSelected(value);
  }, [value]);

  const items = useMemo<MenuProps["items"]>(
    () => [
      {
        key: "datasource",
        label: "Manage Datasource",
        onClick: () => datasourceCrudyEmitter.dispatchEvent("open"),
      },
      {
        key: "tag",
        label: "Manage Tag",
        onClick: () => tagCrudyEmitter.dispatchEvent("open"),
      },
    ],
    [datasourceCrudyEmitter, tagCrudyEmitter],
  );

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
        {/*<Button onClick={getSources}>*/}
        {/*  <ReloadOutlined />*/}
        {/*</Button>*/}
        <Dropdown menu={{ items }}>
          <Button type="link">
            <Space>
              More
              <DownOutlined />
            </Space>
          </Button>
        </Dropdown>
        <div className={styles.hidden}>
          <CrudyButton<IDatasource>
            name="Datasource"
            crudy={DatasourceCrudy}
            columns={datasourceColumns}
            pageable={false}
            afterSaved={getSources}
            emitter={datasourceCrudyEmitter}
            buttonProps={{
              type: "primary",
            }}
          >
            <Form.Item name="name" label="Name" rules={[{ required: true }]}>
              <Input placeholder="Name is required!" maxLength={200} />
            </Form.Item>
            <Form.Item name="type" label="Type" rules={[{ required: true }]}>
              <Select options={DatasourceTypes} placeholder="~" showSearch />
            </Form.Item>
            <Form.Item name="cwd" label="CWD" rules={[{ required: true }]}>
              <Input placeholder="CWD is required!" maxLength={3072} />
            </Form.Item>
          </CrudyButton>
          <CrudyButton<ITag>
            name="Tag"
            crudy={TagCrudy}
            columns={tagColumns}
            pageable={false}
            emitter={tagCrudyEmitter}
            buttonProps={{
              type: "primary",
            }}
          >
            <Form.Item name="name" label="Name" rules={[{ required: true }]}>
              <Input placeholder="Name is required!" maxLength={200} />
            </Form.Item>
            <Form.Item name="key" label="Key" rules={[{ required: true }]}>
              <Input placeholder="Key is required!" maxLength={200} />
            </Form.Item>
            <Form.Item name="color" label="Color">
              <Input type="color" />
            </Form.Item>
          </CrudyButton>
        </div>
      </div>
    </Spin>
  );
}
