import { CrudyButton, searchable } from "@allape/gocrud-react";
import { ICrudyButtonProps } from "@allape/gocrud-react/src/component/CrudyButton";
import NewCrudyButtonEventEmitter from "@allape/gocrud-react/src/component/CrudyButton/eventemitter.ts";
import { EllipsisCell } from "@allape/gocrud-react/src/helper/antd.tsx";
import { DownOutlined } from "@ant-design/icons";
import {
  Avatar,
  Button,
  Dropdown,
  Form,
  Input,
  MenuProps,
  Select,
  Space,
  Tag,
} from "antd";
import { ReactElement, useMemo, useState } from "react";
import { DatasourceCrudy } from "../../api/datasource.ts";
import { getPreviewURLByKey, PreviewCrudy } from "../../api/preview.ts";
import { TagCrudy } from "../../api/tag.ts";
import IDatasource, { DatasourceTypes } from "../../model/datasource.ts";
import IPreview, { IPreviewSearchParams } from "../../model/preview.ts";
import ITag from "../../model/tag.ts";

export interface IPrimaryDropdownProps {
  afterSaved?: () => void;
}

export default function PrimaryDropdown({
  afterSaved,
}: IPrimaryDropdownProps): ReactElement {
  const [previewSearchParams, setPreviewSearchParams] =
    useState<IPreviewSearchParams>({});

  const datasourceCrudyEmitter = useMemo(
    () => NewCrudyButtonEventEmitter<IDatasource>(),
    [],
  );

  const tagCrudyEmitter = useMemo(() => NewCrudyButtonEventEmitter<ITag>(), []);

  const previewCrudyEmitter = useMemo(
    () => NewCrudyButtonEventEmitter<IPreview>(),
    [],
  );

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
      {
        key: "preview",
        label: "Manage Preview",
        onClick: () => previewCrudyEmitter.dispatchEvent("open"),
      },
    ],
    [datasourceCrudyEmitter, previewCrudyEmitter, tagCrudyEmitter],
  );

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

  const previewColumns = useMemo<ICrudyButtonProps<IPreview>["columns"]>(
    () => [
      {
        dataIndex: "id",
        title: "ID",
      },
      {
        dataIndex: "cover",
        title: "Cover",
        render: (_, record) => {
          const url = getPreviewURLByKey(record.key);
          return (
            <Avatar
              size={64}
              shape="square"
              style={{ cursor: "pointer" }}
              src={url}
              onClick={() => window.open(url)}
              data-url={url}
            />
          );
        },
      },
      {
        dataIndex: "key",
        title: "Key",
        render: (v, record) => (
          <div>
            <div>{v}</div>
            <div>{record.digest}</div>
          </div>
        ),
        filtered: !!previewSearchParams["key"],
        ...searchable("key", "Key", (dataIndex, value) =>
          setPreviewSearchParams((old) => ({
            ...old,
            [dataIndex]: value,
          })),
        ),
      },
      {
        dataIndex: "ffprobeInfo",
        title: "FFProbe Info",
        render: EllipsisCell(),
        filtered: !!previewSearchParams["ffprobeInfo"],
        ...searchable("ffprobeInfo", "FFProbe Info", (dataIndex, value) =>
          setPreviewSearchParams((old) => ({
            ...old,
            [dataIndex]: value,
          })),
        ),
      },
    ],
    [previewSearchParams],
  );

  return (
    <>
      <Dropdown menu={{ items }}>
        <Button type="link">
          <Space>
            More
            <DownOutlined />
          </Space>
        </Button>
      </Dropdown>
      <div style={{ display: "none" }}>
        <CrudyButton<IDatasource>
          name="Datasource"
          crudy={DatasourceCrudy}
          columns={datasourceColumns}
          pageable={false}
          afterSaved={afterSaved}
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
        <CrudyButton<IPreview>
          name="Preview"
          crudy={PreviewCrudy}
          columns={previewColumns}
          emitter={previewCrudyEmitter}
          buttonProps={{
            type: "primary",
          }}
          searchParams={previewSearchParams}
          creatable={false}
          editable={false}
          deletable={false}
        />
      </div>
    </>
  );
}
