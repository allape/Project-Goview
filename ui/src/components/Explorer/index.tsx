import { Flex } from "@allape/gocrud-react";
import { SERVER_URL } from "@allape/gocrud-react/src/config";
import { useLoading, useProxy } from "@allape/use-loading";
import { ReloadOutlined } from "@ant-design/icons";
import { Button, Empty, Input, Spin } from "antd";
import cls from "classnames";
import { ReactElement, useCallback, useEffect, useState } from "react";
import { readDir } from "../../api/datasource.ts";
import { generatePreview } from "../../api/preview.ts";
import IDatasource, { IFileInfo } from "../../model/datasource.ts";
import styles from "./style.module.scss";

const PreviewableSuffix: string[] = [
  "jpg",
  "jpeg",
  "png",
  "gif",
  "bmp",
  "webp",
  "svg",
  "mp4",
  "raw",
  "arw",
  "webm",
];

export interface IModifiedFileInfo extends IFileInfo {
  url: string;
}

interface IFileProps {
  dummy?: boolean;
  file: IModifiedFileInfo | string;
  onClick?: (file: IModifiedFileInfo | string) => void;
}

const NO_PREVIEW = `${SERVER_URL}/preview/no-preview`;

function File({ dummy, file, onClick }: IFileProps): ReactElement {
  if (typeof file === "string") {
    return (
      <div
        className={cls(styles.file, dummy && styles.dummy)}
        onClick={() => onClick?.(file)}
      >
        <div className={styles.preview}>
          <img src={NO_PREVIEW} alt={file} />
        </div>
        <div className={cls(styles.name, styles.center)}>{file}</div>
      </div>
    );
  }

  return (
    <div className={styles.file} onClick={() => onClick?.(file)}>
      <div className={styles.preview}>
        <img src={file.isDir ? NO_PREVIEW : file.url} alt={file.name} />
      </div>
      <div className={styles.name}>{file.name}</div>
    </div>
  );
}

export interface IExplorerProps {
  value?: IDatasource["id"];
}

export default function Explorer({
  value: valueFromProps,
}: IExplorerProps): ReactElement {
  const { loading, execute } = useLoading();

  const [value, setValue] = useState<IExplorerProps["value"] | undefined>();
  const [cwd, cwdRef, setCwd] = useProxy<string>("");
  const [files, filesRef, setFiles] = useProxy<IModifiedFileInfo[]>([]);

  const [dummyFiles, setDummyFiles] = useState<string[]>([]);

  const reload = useCallback(
    (value: IDatasource["id"] | undefined, cwd: string) => {
      setFiles([]);
      if (!value) {
        return;
      }

      execute(async () => {
        const files = await readDir(value, cwd || "/");
        files.sort((a, b) =>
          a.isDir === b.isDir ? a.name.localeCompare(b.name) : a.isDir ? -1 : 1,
        );
        setFiles(
          files.map((file) => ({
            ...file,
            url: `${SERVER_URL}/preview/${value}${cwd}/${encodeURIComponent(file.name)}`,
          })),
        );
      }).then();
    },
    [execute, setFiles],
  );

  const handleReload = useCallback(() => {
    reload(value, cwd);
  }, [reload, cwd, value]);

  useEffect(() => {
    handleReload();
  }, [handleReload]);

  useEffect(() => {
    setCwd("");
    setValue(valueFromProps);
  }, [setCwd, valueFromProps]);

  const handleClick = useCallback(
    (file: IModifiedFileInfo | string) => {
      if (typeof file === "string") {
        if (file == "..") {
          const cwd = cwdRef.current;
          const parts = cwd.split("/");
          parts.pop();
          location.hash = parts.join("/");
        }
        return;
      }

      if (file.isDir) {
        // setCwd((cwd) => `${cwd}/${encodeURIComponent(file.name)}`);
        location.hash = `${cwdRef.current}/${encodeURIComponent(file.name)}`;
      } else {
        window.open(
          `${SERVER_URL}/datasource/fetch/${value}${cwdRef.current}/${encodeURIComponent(file.name)}`,
        );
      }
    },
    [cwdRef, value],
  );

  const handleGenerate = useCallback(async () => {
    if (!value) {
      return;
    }
    await execute(async () => {
      for (const file of filesRef.current) {
        if (file.isDir) {
          continue;
        }

        const ext = file.name.split(".").pop();
        if (!ext || !PreviewableSuffix.includes(ext.toLowerCase())) {
          continue;
        }

        await generatePreview(
          value,
          `${cwdRef.current}/${encodeURIComponent(file.name)}`,
        );
      }
      reload(value, cwdRef.current);
    });
  }, [cwdRef, execute, filesRef, reload, value]);

  useEffect(() => {
    const handleHashChange = (e: HashChangeEvent) => {
      e.preventDefault();
      const url = new URL(e.newURL);
      setCwd(url.hash.slice(1));
    };
    window.addEventListener("hashchange", handleHashChange);
    return () => {
      window.removeEventListener("hashchange", handleHashChange);
    };
  }, [setCwd]);

  useEffect(() => {
    const handleResize = () => {
      const width = window.innerWidth > 1400 ? 1400 : window.innerWidth;
      const count = Math.floor(width / 300) + 1;
      setDummyFiles(new Array(count).fill(0).map((_, i) => `_dummy_${i}`));
    };
    handleResize();
    window.addEventListener("resize", handleResize);
    return () => {
      window.removeEventListener("resize", handleResize);
    };
  }, []);

  return (
    <Spin spinning={loading}>
      <div className={styles.wrapper}>
        <Flex justifyContent="stretch">
          <Input value={cwd} readOnly placeholder="CWD" />
          <Button onClick={handleReload}>
            <ReloadOutlined />
          </Button>
          <Button onClick={handleGenerate}>Generate Preview</Button>
        </Flex>
        <div className={styles.files}>
          {cwd && <File file=".." onClick={handleClick} />}
          {files.map((file) => (
            <File key={file.name} file={file} onClick={handleClick} />
          ))}
          {files.length === 0 && !cwd ? (
            <Empty className={styles.emtpy} />
          ) : undefined}
          {dummyFiles.map((file) => (
            <File key={file} dummy file={file} />
          ))}
        </div>
      </div>
    </Spin>
  );
}
