import { Flex } from "@allape/gocrud-react";
import { useLoading, useProxy } from "@allape/use-loading";
import {
  ExportOutlined,
  FullscreenExitOutlined,
  ReloadOutlined,
} from "@ant-design/icons";
import { Button, Empty, Input, Spin } from "antd";
import { ReactElement, useCallback, useEffect, useState } from "react";
import { getFileURLFromDatasource, readDir } from "../../api/datasource.ts";
import {
  generatePreview,
  getPreviewURLByDatasource,
} from "../../api/preview.ts";
import IDatasource, { IFileInfo } from "../../model/datasource.ts";
import File from "../File";
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
            url: getPreviewURLByDatasource(
              value,
              `${cwd}/${encodeURIComponent(file.name)}`,
            ),
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
          getFileURLFromDatasource(
            value!,
            `${cwdRef.current}/${encodeURIComponent(file.name)}`,
          ),
        );
      }
    },
    [cwdRef, value],
  );

  const handleGenerate = useCallback(
    async (file?: IModifiedFileInfo) => {
      if (!value) {
        return;
      }
      await execute(async () => {
        const noFilter = !!file;

        let files = filesRef.current;
        if (file) {
          files = [file];
        }

        for (const file of files) {
          if (file.isDir) {
            continue;
          }

          const ext = file.name.split(".").pop();
          if (
            !noFilter &&
            (!ext || !PreviewableSuffix.includes(ext.toLowerCase()))
          ) {
            continue;
          }

          await generatePreview(
            value,
            `${cwdRef.current}/${encodeURIComponent(file.name)}`,
          );
        }
        reload(value, cwdRef.current);
      });
    },
    [cwdRef, execute, filesRef, reload, value],
  );

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
          <Input value={decodeURIComponent(cwd)} readOnly placeholder="CWD" />
          <Button onClick={handleReload}>
            <ReloadOutlined />
          </Button>
          <Button onClick={() => handleGenerate()}>Generate Preview</Button>
        </Flex>
        <div className={styles.files}>
          {cwd && <File name=".." center onClick={() => handleClick("..")} />}
          {files.map((file) => (
            <File
              key={file.name}
              name={`${file.isDir ? `📁` : "📃"} ${file.name}`}
              cover={file.isDir ? undefined : file.url}
              onClick={file.isDir ? () => handleClick(file) : undefined}
            >
              {!file.isDir ? (
                <>
                  <Button type="link" onClick={() => handleGenerate(file)}>
                    <FullscreenExitOutlined />
                  </Button>
                  <Button type="link" onClick={() => handleClick(file)}>
                    <ExportOutlined />
                  </Button>
                </>
              ) : (
                <></>
              )}
            </File>
          ))}
          {files.length === 0 && !cwd ? (
            <Empty className={styles.emtpy} />
          ) : undefined}
          {dummyFiles.map((file) => (
            <File key={file} hidden name={file} />
          ))}
        </div>
      </div>
    </Spin>
  );
}
